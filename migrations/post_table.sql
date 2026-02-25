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

-- 2. TẠO STORED PROCEDURE SINH DỮ LIỆU CHIA LÔ (BATCHING)
DELIMITER $$
CREATE PROCEDURE SeedMockData()
BEGIN
    DECLARE i INT DEFAULT 0;
    
    -- Tạm thời tắt các kiểm tra ràng buộc để tối ưu Disk I/O
    SET autocommit = 0;
    SET unique_checks = 0;
    SET foreign_key_checks = 0;

    WHILE i < 10 DO
        START TRANSACTION;
        -- Chèn 100.000 dòng mỗi lô (sử dụng 5 bảng lai chéo)
        INSERT INTO posts (title, description, content, status, publish_date)
        SELECT 
            CONCAT('Tieu de bai viet ', a.N + b.N * 10 + c.N * 100 + d.N * 1000 + e.N * 10000 + (i * 100000)),
            'Mo ta mau thuc te duoc tao tu dong.',
            'Noi dung chi tiet duoc tao tu dong.',
            IF((a.N + b.N) % 5 = 0, 'Published', 'Draft'), 
            IF((a.N + b.N) % 5 = 0, CURRENT_TIMESTAMP - INTERVAL ((a.N + b.N) % 365) DAY, NULL)
        FROM 
            (SELECT 0 AS N UNION ALL SELECT 1 UNION ALL SELECT 2 UNION ALL SELECT 3 UNION ALL SELECT 4 UNION ALL SELECT 5 UNION ALL SELECT 6 UNION ALL SELECT 7 UNION ALL SELECT 8 UNION ALL SELECT 9) a
            CROSS JOIN (SELECT 0 AS N UNION ALL SELECT 1 UNION ALL SELECT 2 UNION ALL SELECT 3 UNION ALL SELECT 4 UNION ALL SELECT 5 UNION ALL SELECT 6 UNION ALL SELECT 7 UNION ALL SELECT 8 UNION ALL SELECT 9) b
            CROSS JOIN (SELECT 0 AS N UNION ALL SELECT 1 UNION ALL SELECT 2 UNION ALL SELECT 3 UNION ALL SELECT 4 UNION ALL SELECT 5 UNION ALL SELECT 6 UNION ALL SELECT 7 UNION ALL SELECT 8 UNION ALL SELECT 9) c
            CROSS JOIN (SELECT 0 AS N UNION ALL SELECT 1 UNION ALL SELECT 2 UNION ALL SELECT 3 UNION ALL SELECT 4 UNION ALL SELECT 5 UNION ALL SELECT 6 UNION ALL SELECT 7 UNION ALL SELECT 8 UNION ALL SELECT 9) d
            CROSS JOIN (SELECT 0 AS N UNION ALL SELECT 1 UNION ALL SELECT 2 UNION ALL SELECT 3 UNION ALL SELECT 4 UNION ALL SELECT 5 UNION ALL SELECT 6 UNION ALL SELECT 7 UNION ALL SELECT 8 UNION ALL SELECT 9) e;
        
        -- Giải phóng bộ nhớ đệm Buffer Pool và ghi Redo Log
        COMMIT;
        SET i = i + 1;
    END WHILE;

    -- Khôi phục cấu hình
    SET unique_checks = 1;
    SET foreign_key_checks = 1;
    SET autocommit = 1;
END$$
DELIMITER ;

-- 3. GỌI THỦ TỤC THỰC THI
CALL SeedMockData();

-- 4. XÓA THỦ TỤC SAU KHI DÙNG (Dọn dẹp rác)
DROP PROCEDURE SeedMockData;