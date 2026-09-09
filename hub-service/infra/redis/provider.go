package redis

import (
	common_redis "common/redis"
	"context"
	"fmt"
	"hub/config"
	providers "hub/internal/provider"

	"github.com/redis/go-redis/v9"
)

type RedisProvider struct {
	*common_redis.RedisService
}

func NewRedisService(runtime config.Runtime) (*common_redis.RedisService, func(), error) {
	service, err := common_redis.Open(common_redis.Config{
		Address:  runtime.Redis.Address,
		Password: runtime.Redis.Password,
		DB:       runtime.Redis.DB,
	})
	if err != nil {
		return nil, nil, fmt.Errorf("open hub redis: %w", err)
	}

	cleanup := func() {
		_ = service.Close()
	}
	return service, cleanup, nil
}

func NewRedisProvider(service *common_redis.RedisService) providers.CacheProvider {
	return &RedisProvider{RedisService: service}
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

func (r *RedisProvider) Get(c context.Context, key string) (string, error) {
	client, err := r.Instance()
	if err != nil {
		return "", err
	}
	result, err := client.Get(c, key).Result()
	if err == redis.Nil {
		return "", nil
	}
	return result, err
}

func (r *RedisProvider) Set(c context.Context, key string, value string) error {
	client, err := r.Instance()
	if err != nil {
		return err
	}
	return client.Set(c, key, value, 0).Err()
}
