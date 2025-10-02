package providers

import (
	"context"
	"time"
	"tqd/internal/interface/provider"

	_redis "common/redis"
)

type CacheProvider struct {
	*_redis.RedisService
}

func NewCacheProvider(client *_redis.RedisService) provider.CacheProvider {
	return &CacheProvider{
		RedisService: client,
	}
}

func (p *CacheProvider) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	// Return nil for now, implement later when needed
	return nil
}

func (p *CacheProvider) Get(ctx context.Context, key string, dest interface{}) error {
	// Return nil for now, implement later when needed
	return nil
}

func (p *CacheProvider) Delete(ctx context.Context, key string) error {
	// Return nil for now, implement later when needed
	return nil
}

func (p *CacheProvider) Exists(ctx context.Context, key string) (bool, error) {
	// Return false for now, implement later when needed
	return false, nil
}

func (p *CacheProvider) SetNX(ctx context.Context, key string, value interface{}, expiration time.Duration) (bool, error) {
	// Return true for now, implement later when needed
	return true, nil
}

func (p *CacheProvider) Increment(ctx context.Context, key string, delta int64) (int64, error) {
	// Return delta for now, implement later when needed
	return delta, nil
}

func (p *CacheProvider) Expire(ctx context.Context, key string, expiration time.Duration) error {
	// Return nil for now, implement later when needed
	return nil
}

func (p *CacheProvider) GetTTL(ctx context.Context, key string) (time.Duration, error) {
	// Return 0 for now, implement later when needed
	return 0, nil
}
