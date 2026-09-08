package providers

import (
	"encoding/json"

	redis_cli "bdspro/infra/redis"

	"github.com/gin-gonic/gin"
)

// RedisProvider is a compatibility adapter over the process-owned Redis client.
// It must never construct or own a second Redis connection.
type RedisProvider struct {
	client *redis_cli.RedisClient
}

func NewRedisProvider(client *redis_cli.RedisClient) *RedisProvider {
	return &RedisProvider{client: client}
}

func GetValueWithRedis[T any](r *RedisProvider, ctx *gin.Context, key string) (*T, error) {
	var result T
	client, err := r.client.Client(ctx)
	if err != nil {
		return &result, err
	}
	val, err := client.Get(ctx, key).Result()
	if err != nil {
		return &result, err
	}

	err = json.Unmarshal([]byte(val), &result)
	return &result, err
}
