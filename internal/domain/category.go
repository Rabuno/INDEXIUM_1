package domain

import (
	"context"
	"time"
)

// --- ENUMS & CONSTANTS ---
const (
	CategoryStatusActive   = "Active"
	CategoryStatusInactive = "Inactive"
)

// --- ENTITIES ---

// Category đại diện cho danh mục trong hệ thống

type Category struct {
	ID          int64     `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Thumbnail   string    `json:"thumbnail"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// --- INTERFACES (PORTS) ---

// CategoryRepository định nghĩa các hành vi tương tác với dữ liệu (Output Port)
// Lớp Repository (MySQL) sẽ phải implement interface này.

type CategoryRepository interface {
	Fetch(ctx context.Context, limit int64, offset int64) ([]Category, error)
	GetByID(ctx context.Context, id int64) (*Category, error)
	Store(ctx context.Context, c *Category) error
	Update(ctx context.Context, c *Category) error
	Delete(ctx context.Context, id int64) error
}

type CategoryUseCase interface {
	Fetch(ctx context.Context, page int64, pageSize int64) ([]Category, error)
	GetByID(ctx context.Context, id int64) (*Category, error)
	Store(ctx context.Context, c *Category) error
	Update(ctx context.Context, c *Category) error
	Delete(ctx context.Context, id int64) error
}
