package redis

import (
	"common/logging"
	"context"
	"fmt"
	"log/slog"

	"relay/config"

	"github.com/go-redis/redis/v8"
)

type RedisClient struct {
	client *redis.Client
}

func NewRedisClient() (*RedisClient, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     config.AppProperties.Redis.Host,
		Password: config.AppProperties.Redis.Password,
		DB:       0,
	})

	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		_ = client.Close()
		return nil, fmt.Errorf("connect Relay Redis: %w", err)
	}

	slog.Info(
		"Relay Redis connected",
		slog.String("component", "redis"),
	)
	return &RedisClient{client: client}, nil
}

func (r *RedisClient) Close() error {
	if r == nil || r.client == nil {
		return nil
	}
	return r.client.Close()
}

// Push message to list
func (r *RedisClient) PushMessage(ctx context.Context, key string, value []byte) error {
	err := r.client.LPush(ctx, key, value).Err()
	if err != nil {
		logging.WithComponent(ctx, "redis").Error(
			"push Relay message to Redis",
			slog.String("redis.key", key),
			slog.Any("error", err),
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

// Set key-value pair in Redis
func (r *RedisClient) Set(ctx context.Context, key string, value string) error {
	err := r.client.Set(ctx, key, value, 0).Err()
	if err != nil {
		logging.WithComponent(ctx, "redis").Error(
			"set Relay Redis key",
			slog.String("redis.key", key),
			slog.Any("error", err),
		)
	}
	return err
}

// Get value from Redis
func (r *RedisClient) Get(ctx context.Context, key string) (string, error) {
	val, err := r.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return "", nil
	}
	if err != nil {
		logging.WithComponent(ctx, "redis").Error(
			"get Relay Redis key",
			slog.String("redis.key", key),
			slog.Any("error", err),
		)
		return "", err
	}
	return val, nil
}
