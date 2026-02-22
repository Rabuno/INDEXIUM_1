package redis

import (
	"context"
	"encoding/json"
	"time"

	"Test2/internal/domain"
	redisclient "github.com/redis/go-redis/v9"
)

type redisCacheRepo struct {
	client *redisclient.Client
}

func NewRedisCacheRepository(client *redisclient.Client) domain.CacheRepository {
	return &redisCacheRepo{client: client}
}

func (r *redisCacheRepo) Get(ctx context.Context, key string) ([]domain.Post, bool) {
	val, err := r.client.Get(ctx, key).Result()
	if err != nil {
		return nil, false // Cache miss hoặc lỗi kết nối
	}

	var posts []domain.Post
	err = json.Unmarshal([]byte(val), &posts)
	if err != nil {
		return nil, false // Lỗi parse JSON
	}

	return posts, true
}

func (r *redisCacheRepo) Set(ctx context.Context, key string, value []domain.Post, ttl time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return r.client.Set(ctx, key, data, ttl).Err()
}

func (r *redisCacheRepo) Delete(ctx context.Context, key string) error {
	return r.client.Del(ctx, key).Err()
}
