package usecase

import (
	"context"
	"fmt"
	"time"

	"Test2/internal/domain"
)

type postUseCase struct {
	postRepo       domain.PostRepository
	cache          domain.CacheRepository
	contextTimeout time.Duration
}

// NewPostUseCase khởi tạo PostUseCase với Dependency Injection
func NewPostUseCase(
	repo domain.PostRepository,
	cache domain.CacheRepository,
	timeout time.Duration,
) domain.PostUseCase {
	return &postUseCase{
		postRepo:       repo,
		cache:          cache,
		contextTimeout: timeout,
	}
}

// Helper: Xóa các khóa cache liên quan đến danh sách bài viết để đảm bảo tính nhất quán
func (pu *postUseCase) invalidatePostListCache(ctx context.Context) {
	// Trong hệ thống thực tế có nhiều trang, cần cơ chế xóa theo pattern (SCAN) hoặc cấu trúc Hash.
	// Minh họa xóa cache danh sách trang 1 mặc định.
	//_ = pu.cache.Delete(ctx, "posts:list:page:1:size:10")
}

// Helper: Xóa cache của một bài viết cụ thể
func (pu *postUseCase) invalidateSinglePostCache(ctx context.Context, id int64) {
	cacheKey := fmt.Sprintf("post:detail:%d", id)
	_ = pu.cache.Delete(ctx, cacheKey)
}

func (pu *postUseCase) Fetch(ctx context.Context, page int64, pageSize int64) ([]domain.Post, error) {
	c, cancel := context.WithTimeout(ctx, pu.contextTimeout)
	defer cancel()

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	cacheKey := fmt.Sprintf("posts:list:page:%d:size:%d", page, pageSize)

	// Kiểm tra Cache Hit
	if cachedPosts, found := pu.cache.Get(c, cacheKey); found {
		return cachedPosts, nil
	}

	// Cache Miss -> Gọi MySQL
	offset := (page - 1) * pageSize
	posts, err := pu.postRepo.Fetch(c, pageSize, offset)
	if err != nil {
		return nil, err
	}

	// Ghi vào Cache với TTL = 5 phút
	_ = pu.cache.Set(c, cacheKey, posts, 5*time.Minute)

	return posts, nil
}

func (pu *postUseCase) Store(ctx context.Context, p *domain.Post) error {
	c, cancel := context.WithTimeout(ctx, pu.contextTimeout)
	defer cancel()

	if p.Status == "" {
		p.Status = domain.StatusDraft
	}
	now := time.Now()
	p.CreatedAt = now
	p.UpdateDate = now

	err := pu.postRepo.Store(c, p)
	if err == nil {
		// Dữ liệu mới thay đổi danh sách -> Xóa cache danh sách
		pu.invalidatePostListCache(c)
	}
	return err
}

func (pu *postUseCase) Delete(ctx context.Context, id int64) error {
	c, cancel := context.WithTimeout(ctx, pu.contextTimeout)
	defer cancel()

	_, err := pu.postRepo.GetByID(c, id)
	if err != nil {
		return err
	}

	err = pu.postRepo.Delete(c, id)
	if err == nil {
		pu.invalidatePostListCache(c)
		pu.invalidateSinglePostCache(c, id)
	}
	return err
}

func (pu *postUseCase) GetByID(ctx context.Context, id int64) (*domain.Post, error) {
	c, cancel := context.WithTimeout(ctx, pu.contextTimeout)
	defer cancel()

	cacheKey := fmt.Sprintf("post:detail:%d", id)

	// Lấy mảng từ cache, giả sử phần tử đầu tiên là kết quả cần tìm
	if cachedData, found := pu.cache.Get(c, cacheKey); found && len(cachedData) > 0 {
		return &cachedData[0], nil
	}

	post, err := pu.postRepo.GetByID(c, id)
	if err != nil {
		return nil, err
	}

	_ = pu.cache.Set(c, cacheKey, []domain.Post{*post}, 10*time.Minute)

	return post, nil
}

func (pu *postUseCase) Update(ctx context.Context, p *domain.Post) error {
	c, cancel := context.WithTimeout(ctx, pu.contextTimeout)
	defer cancel()

	p.UpdateDate = time.Now()
	err := pu.postRepo.Update(c, p)
	if err == nil {
		pu.invalidatePostListCache(c)
		pu.invalidateSinglePostCache(c, p.ID)
	}
	return err
}

func (pu *postUseCase) Search(ctx context.Context, keyword string, page int64, pageSize int64) ([]domain.Post, error) {
	c, cancel := context.WithTimeout(ctx, pu.contextTimeout)
	defer cancel()

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	cacheKey := fmt.Sprintf("posts:search:%s:page:%d:size:%d", keyword, page, pageSize)

	if cachedPosts, found := pu.cache.Get(c, cacheKey); found {
		return cachedPosts, nil
	}

	offset := (page - 1) * pageSize
	posts, err := pu.postRepo.Search(c, keyword, pageSize, offset)
	if err != nil {
		return nil, err
	}

	_ = pu.cache.Set(c, cacheKey, posts, 3*time.Minute)

	return posts, nil
}
