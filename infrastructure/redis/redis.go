package redis

import (
	"context"
	"log"

	"Test2/config"
	redisclient "github.com/redis/go-redis/v9"
)

var Ctx = context.Background()
var Client *redisclient.Client

// Sửa chữ ký hàm để nhận tham số cấu hình
func InitRedis(cfg *config.Config) {
	Client = redisclient.NewClient(&redisclient.Options{
		Addr:     cfg.GetRedisAddr(), // Sử dụng địa chỉ động từ config
		Password: "",                 // Có thể mở rộng để lấy Password từ config nếu cần
		DB:       0,
	})

	// Kiểm tra kết nối
	_, err := Client.Ping(Ctx).Result()
	if err != nil {
		log.Fatalf("Redis connection failed at %s: %v", cfg.GetRedisAddr(), err)
	}

	log.Printf("Redis connected successfully at %s", cfg.GetRedisAddr())
}
