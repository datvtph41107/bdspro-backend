package providers

import (
	"context"
	"encoding/json"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
)

type RedisProvider struct {
	client *redis.Client
	ctx    context.Context
}

func NewRedisProvider() *RedisProvider {
	rdb := redis.NewClient(&redis.Options{
		Addr:     viper.GetString("redis.host"), // Địa chỉ Redis (ví dụ: "localhost:6379")
		Password: viper.GetString("redis.pass"),
	})
	return &RedisProvider{
		client: rdb,
		ctx:    context.Background(),
	}
}

func GetValueWithRedis[T any](r *RedisProvider, ctx *gin.Context, key string) (*T, error) {
	var result T
	val, err := r.client.Get(ctx, key).Result()
	if err != nil {
		return &result, err
	}

	err = json.Unmarshal([]byte(val), &result)
	return &result, err
}
