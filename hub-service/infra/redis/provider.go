package redis

import (
	_redis "common/redis"
	"context"
	providers "hub/internal/provider"
)

type RedisProvider struct {
	*_redis.RedisService
}

func NewRedisProvider() providers.CacheProvider {
	return &RedisProvider{
		RedisService: _redis.NewRedisService(),
	}
}

func (r *RedisProvider) LRange(c context.Context, key string, start, stop int64) ([]string, error) {
	client, err := r.Instance()
	if err != nil {
		return nil, err
	}
	return client.LRange(c, key, start, stop).Result()
}

func (r *RedisProvider) LTrim(c context.Context, key string, start, stop int64) error {
	client, err := r.Instance()
	if err != nil {
		return err
	}
	return client.LTrim(c, key, start, stop).Err()
}

func (r *RedisProvider) RPush(c context.Context, key string, id uint64) error {
	client, err := r.Instance()
	if err != nil {
		return err
	}
	return client.RPush(c, key, id).Err()
}
