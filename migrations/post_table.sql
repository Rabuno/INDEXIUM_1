USE ahihi_db;

CREATE TABLE IF NOT EXISTS posts (
    id INT AUTO_INCREMENT PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    content LONGTEXT,
    thumbnail VARCHAR(512),
    status VARCHAR(50) DEFAULT 'Draft',
    publish_date DATETIME NULL,
    update_date DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    
    -- TỐI ƯU: Chỉ cần 1 index này là đủ cho cả Filter và Sort
    -- Nó giúp MySQL: "Lọc status xong là có ngay danh sách đã sắp xếp theo ngày"
    INDEX idx_status_created_at (status, created_at DESC)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;