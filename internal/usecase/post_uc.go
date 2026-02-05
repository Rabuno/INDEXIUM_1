package usecase

import (
	"context"
	"time"

	"Test2/internal/domain"
)

type postUseCase struct {
	postRepo       domain.PostRepository
	contextTimeout time.Duration // Timeout cho mỗi request để tránh treo hệ thống
}

// NewPostUseCase khởi tạo
func NewPostUseCase(repo domain.PostRepository, timeout time.Duration) domain.PostUseCase {
	return &postUseCase{
		postRepo:       repo,
		contextTimeout: timeout,
	}
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

	offset := (page - 1) * pageSize
	return pu.postRepo.Fetch(c, pageSize, offset)
}

func (pu *postUseCase) Store(ctx context.Context, p *domain.Post) error {
	c, cancel := context.WithTimeout(ctx, pu.contextTimeout)
	defer cancel()

	// Logic nghiệp vụ: Set default status
	if p.Status == "" {
		p.Status = domain.StatusDraft
	}

	// Logic nghiệp vụ: Set timestamps
	now := time.Now()
	p.CreatedAt = now
	p.UpdateDate = now

	return pu.postRepo.Store(c, p)
}

func (pu *postUseCase) Delete(ctx context.Context, id int64) error {
	c, cancel := context.WithTimeout(ctx, pu.contextTimeout)
	defer cancel()

	// Logic mở rộng: Kiểm tra xem bài viết có tồn tại không trước khi xóa
	_, err := pu.postRepo.GetByID(c, id)
	if err != nil {
		return err // Trả về lỗi nếu không tìm thấy
	}

	return pu.postRepo.Delete(c, id)
}

func (pu *postUseCase) GetByID(ctx context.Context, id int64) (*domain.Post, error) {
	c, cancel := context.WithTimeout(ctx, pu.contextTimeout)
	defer cancel()
	return pu.postRepo.GetByID(c, id)
}

func (pu *postUseCase) Update(ctx context.Context, p *domain.Post) error {
	c, cancel := context.WithTimeout(ctx, pu.contextTimeout)
	defer cancel()

	p.UpdateDate = time.Now()
	return pu.postRepo.Update(c, p)
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

	offset := (page - 1) * pageSize
	return pu.postRepo.Fetch(c, pageSize, offset)
}
