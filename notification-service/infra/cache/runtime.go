package cache

import (
	"context"

	"notification/config"

	"github.com/go-redis/redis/v8"
	"github.com/hyperledger/fabric/common/flogging"
)

type RedisClient struct {
	client *redis.Client
	logger *flogging.FabricLogger
}

func NewRedisClient() *RedisClient {
	logger := flogging.MustGetLogger("redis")
	client := redis.NewClient(&redis.Options{
		Addr:     config.AppProperties.Redis.Host,
		Password: config.AppProperties.Redis.Password,
		DB:       0,
	})

	_, err := client.Ping(client.Context()).Result()
	if err != nil {
		logger.Fatalf("Failed to connect to Redis: %v", err)
	}

	logger.Info("Connected to Redis")
	return &RedisClient{client: client, logger: logger}
}

// Publish message to a Redis channel
func (r *RedisClient) Publish(ctx context.Context, channel string, message string) error {
	err := r.client.Publish(ctx, channel, message).Err()
	if err != nil {
		r.logger.Errorf("Failed to publish message to Redis channel %s: %v", channel, err)
	}
	return err
}

// Get lấy giá trị từ Redis key
func (r *RedisClient) Get(ctx context.Context, key string) (string, error) {
	result, err := r.client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return "", nil // Key không tồn tại
		}
		r.logger.Errorf("Failed to get key %s from Redis: %v", key, err)
		return "", err
	}
	return result, nil
}
