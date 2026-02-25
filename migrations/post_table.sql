-- 0. CƯỠNG CHẾ BẢNG MÃ UNICODE CHO TOÀN BỘ LUỒNG KẾT NỐI
SET NAMES utf8mb4 COLLATE utf8mb4_unicode_ci;
SET CHARACTER SET utf8mb4;

USE ahihi_db;

-- 1. ĐỊNH NGHĨA LƯỢC ĐỒ (SCHEMA)
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
    
    INDEX idx_status_created_at (status, created_at DESC),
    INDEX idx_created_at (created_at DESC),
    FULLTEXT INDEX idx_fts_search (title, description, content)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;