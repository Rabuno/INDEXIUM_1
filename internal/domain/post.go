package domain

import (
	"context"
	"time"
)

// --- ENUMS & CONSTANTS ---
// Sử dụng hằng số để tránh magic string trong code
const (
	StatusDraft     = "Draft"
	StatusPending   = "Pending"
	StatusPublished = "Published"
	StatusDeleted   = "Deleted" // Key cho tính năng Soft Delete
)

// --- ENTITIES ---

// Post đại diện cho bài viết trong hệ thống
type Post struct {
	ID          int64     `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Content     string    `json:"content"`
	Thumbnail   string    `json:"thumbnail"`
	Status      string    `json:"status"`
	UpdateDate  time.Time `json:"update_date"`
	CreatedAt   time.Time `json:"created_at"`
}

// --- INTERFACES (PORTS) ---

// PostRepository định nghĩa các hành vi tương tác với dữ liệu (Output Port)
// Lớp Repository (MySQL) sẽ phải implement interface này.
type PostRepository interface {
	// Fetch lấy danh sách bài viết có phân trang
	Fetch(ctx context.Context, limit int64, offset int64) ([]Post, error)
	// GetByID lấy chi tiết một bài viết
	GetByID(ctx context.Context, id int64) (*Post, error)
	// Store tạo mới một bài viết
	Store(ctx context.Context, p *Post) error
	// Update cập nhật thông tin bài viết
	Update(ctx context.Context, p *Post) error
	// Delete thực hiện xóa mềm (Soft Delete)
	Delete(ctx context.Context, id int64) error
	// Search tìm kiếm bài viết theo từ khóa với phân trang
	Search(ctx context.Context, keyword string, limit int64, offset int64) ([]Post, error)
}

// PostUseCase định nghĩa các logic nghiệp vụ (Input Port)
// Lớp Delivery (Gin Handler) sẽ gọi interface này.
type PostUseCase interface {
	Fetch(ctx context.Context, page int64, pageSize int64) ([]Post, error)
	GetByID(ctx context.Context, id int64) (*Post, error)
	Store(ctx context.Context, p *Post) error
	Update(ctx context.Context, p *Post) error
	Delete(ctx context.Context, id int64) error
	Search(ctx context.Context, keyword string, page int64, pageSize int64) ([]Post, error)
}
