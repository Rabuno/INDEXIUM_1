package domain

import (
	"context"
	"time"
)

type CacheRepository interface {
	Get(ctx context.Context, key string) ([]Post, bool)
	Set(ctx context.Context, key string, value []Post, ttl time.Duration) error
	Delete(ctx context.Context, key string) error // Bổ sung cơ chế Invalidation
}
