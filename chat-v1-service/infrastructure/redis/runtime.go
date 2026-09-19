package redis

import (
	"context"
	"fmt"
	"os"

	configs "chat/config"
	_provider "chat/internal/interface/provider"
	"common/logging"

	"github.com/go-redis/redis/v8"
)

// @bind: infrastructure/redis
type RedisClient struct {
	client *redis.Client
}

var _ _provider.ITypingCounter = (*RedisClient)(nil)

func NewRedisClient() *RedisClient {
	client := redis.NewClient(&redis.Options{
		Addr:     configs.AppProperties.Redis.Host,
		Password: configs.AppProperties.Redis.Password,
		DB:       0,
	})
	logger := logging.WithComponent(client.Context(), "redis")

	_, err := client.Ping(client.Context()).Result()
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to connect to Redis: %v", err))
		os.Exit(1)
	}

	logger.Info("Connected to Redis")
	return &RedisClient{client: client}
}

// Push message to list
func (r *RedisClient) PushMessage(ctx context.Context, key string, value []byte) error {
	err := r.client.LPush(ctx, key, value).Err()
	if err != nil {
		logging.WithComponent(ctx, "redis").Error(
			fmt.Sprintf("Failed to push message to Redis: %v", err),
		)
	}
	return err
}

// Get and delete messages from list
func (r *RedisClient) GetAndDeleteMessages(ctx context.Context, key string, maxReads int) ([]string, error) {
	msgs, err := r.client.LRange(ctx, key, 0, -1).Result()
	if err != nil {
		return nil, err
	}

	if len(msgs) == 0 {
		return nil, nil
	}

	// Track read count in Hash
	for _, msg := range msgs {
		readCountKey := "read_count:" + key + ":" + msg
		count, _ := r.client.Incr(ctx, readCountKey).Result()

		if int(count) >= maxReads {
			r.client.LRem(ctx, key, 0, msg) // Remove message
			r.client.Del(ctx, readCountKey) // Delete read counter
		}
	}

	return msgs, nil
}

// Publish message to a Redis channel
func (r *RedisClient) Publish(ctx context.Context, channel string, message string) error {
	return r.client.Publish(ctx, channel, message).Err()
}

// AdjustTypingRef INCR khi đang gõ, DECR khi dừng; key theo conversation + user, DECR không âm.
func (r *RedisClient) AdjustTypingRef(ctx context.Context, conversationID, userID uint64, increment bool) (int64, error) {
	key := fmt.Sprintf("typing:%d", conversationID)
	if increment {
		return r.client.Incr(ctx, key).Result()
	}
	val, err := r.client.Decr(ctx, key).Result()
	if err != nil {
		return 0, err
	}
	if val < 0 {
		if err := r.client.Set(ctx, key, 0, 0).Err(); err != nil {
			return 0, err
		}
		return 0, nil
	}
	return val, nil
}

// Subscribe to a Redis channel
func (r *RedisClient) Subscribe(ctx context.Context, channel string) *redis.PubSub {
	return r.client.Subscribe(ctx, channel)
}

// Add message to Redis Stream
func (r *RedisClient) AddToStream(ctx context.Context, stream string, values map[string]interface{}) error {
	_, err := r.client.XAdd(ctx, &redis.XAddArgs{
		Stream: stream,
		Values: values,
	}).Result()
	return err
}

// Read messages from Redis Stream
func (r *RedisClient) ReadStream(ctx context.Context, stream, consumerGroup, consumerName string, count int) ([]redis.XMessage, error) {
	res, err := r.client.XReadGroup(ctx, &redis.XReadGroupArgs{
		Group:    consumerGroup,
		Consumer: consumerName,
		Streams:  []string{stream, ">"},
		Count:    int64(count),
		Block:    0,
	}).Result()

	if err != nil {
		return nil, err
	}

	var messages []redis.XMessage
	for _, stream := range res {
		messages = append(messages, stream.Messages...)
	}

	return messages, nil
}
