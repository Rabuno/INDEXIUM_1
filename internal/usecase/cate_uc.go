package usecase

import (
	"context"
	"time"

	"Test2/internal/domain"
)

type cateUseCase struct {
	cateRepo       domain.CategoryRepository
	contextTimeout time.Duration
}

func NewCateUseCase(repo domain.CategoryRepository, timeout time.Duration) domain.CategoryUseCase {
	return &cateUseCase{
		cateRepo:       repo,
		contextTimeout: timeout,
	}
}

func (cu *cateUseCase) Fetch(ctx context.Context, page int64, pageSize int64) ([]domain.Category, error) {
	c, cancel := context.WithTimeout(ctx, cu.contextTimeout)
	defer cancel()

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize
	return cu.cateRepo.Fetch(c, pageSize, offset)
}

func (cu *cateUseCase) Store(ctx context.Context, c *domain.Category) error {
	p, cancel := context.WithTimeout(ctx, cu.contextTimeout)
	defer cancel()

	if c.Status == "" {
		c.Status = domain.CategoryStatusActive
	}

	now := time.Now()
	c.CreatedAt = now
	c.UpdatedAt = now

	return cu.cateRepo.Store(p, c)
}

func (cu *cateUseCase) Delete(ctx context.Context, id int64) error {
	p, cancel := context.WithTimeout(ctx, cu.contextTimeout)
	defer cancel()

	_, err := cu.cateRepo.GetByID(p, id)
	if err != nil {
		return err
	}

	return cu.cateRepo.Delete(p, id)
}

func (cu *cateUseCase) GetByID(ctx context.Context, id int64) (*domain.Category, error) {
	p, cancel := context.WithTimeout(ctx, cu.contextTimeout)
	defer cancel()

	return cu.cateRepo.GetByID(p, id)
}

func (cu *cateUseCase) Update(ctx context.Context, c *domain.Category) error {
	p, cancel := context.WithTimeout(ctx, cu.contextTimeout)
	defer cancel()

	c.UpdatedAt = time.Now()
	return cu.cateRepo.Update(p, c)
}
