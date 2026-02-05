package kernel

import (
	"context"
	"core/kernel/container"
	"time"

	"github.com/redis/go-redis/v9"
)

func Redis() *redis.Client {
	rs, err := container.ResolveTyped[*redis.Client](appC, "redis")
	if err != nil {
		panic(err)
	}
	return rs
}

func CacheSet(ctx context.Context, key string, val any, ttl time.Duration) error {
	return Redis().Set(ctx, key, val, ttl).Err()
}

func CacheGet(ctx context.Context, key string) (string, error) {
	return Redis().Get(ctx, key).Result()
}
