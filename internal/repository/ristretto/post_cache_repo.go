package ristretto

import (
	"context"
	"time"

	"Test2/internal/domain"
	ristretto "github.com/dgraph-io/ristretto/v2"
)

type ristrettoCacheRepo struct {
	cache *ristretto.Cache[string, []domain.Post]
}

func NewRistrettoCacheRepository(cache *ristretto.Cache[string, []domain.Post]) domain.CacheRepository {
	return &ristrettoCacheRepo{cache: cache}
}

func (r *ristrettoCacheRepo) Get(ctx context.Context, key string) ([]domain.Post, bool) {
	value, found := r.cache.Get(key)
	if !found {
		return nil, false // Cache miss
	}
	return value, true
}

func (r *ristrettoCacheRepo) Set(ctx context.Context, key string, value []domain.Post, ttl time.Duration) error {
	r.cache.SetWithTTL(key, value, 1, ttl) // Sử dụng chi phí cố định cho mỗi mục và TTL
	return nil
}

func (r *ristrettoCacheRepo) Delete(ctx context.Context, key string) error {
	r.cache.Del(key)
	return nil
}
