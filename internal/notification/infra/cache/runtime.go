package cache

import (
	"context"
	"errors"

	sharedredis "common/redis"

	"github.com/redis/go-redis/v9"
)

// RedisClient is a Notification-specific adapter over the one Redis resource
// owned by the process. It never opens or closes a connection itself.
type RedisClient struct{ service *sharedredis.RedisService }

func NewRedisClient(service *sharedredis.RedisService) *RedisClient {
	return &RedisClient{service: service}
}

func (r *RedisClient) Publish(ctx context.Context, channel, message string) error {
	if r == nil || r.service == nil || r.service.Client == nil {
		return errors.New("notification redis unavailable")
	}
	return r.service.Client.Publish(ctx, channel, message).Err()
}

func (r *RedisClient) Get(ctx context.Context, key string) (string, error) {
	if r == nil || r.service == nil || r.service.Client == nil {
		return "", errors.New("notification redis unavailable")
	}
	value, err := r.service.Client.Get(ctx, key).Result()
	if errors.Is(err, redis.Nil) {
		return "", nil
	}
	return value, err
}
