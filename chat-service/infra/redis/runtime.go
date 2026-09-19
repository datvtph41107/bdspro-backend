package redis

import (
	"context"
	"fmt"
	"os"

	configs "chat/config"
	"common/logging"

	"github.com/go-redis/redis/v8"
)

var updateTypingScript = redis.NewScript(`
local current = tonumber(redis.call('GET', KEYS[1]) or '0')
if ARGV[1] == '1' then
  current = redis.call('INCR', KEYS[1])
  redis.call('EXPIRE', KEYS[1], ARGV[2])
  return current
end
if current <= 1 then
  redis.call('DEL', KEYS[1])
  return 0
end
current = redis.call('DECR', KEYS[1])
redis.call('EXPIRE', KEYS[1], ARGV[2])
return current
`)

type RedisClient struct {
	client *redis.Client
}

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

func (r *RedisClient) UpdateTyping(ctx context.Context, conversationID, userID uint64, typing bool) (int64, error) {
	key := fmt.Sprintf("chat:typing:%d:%d", conversationID, userID)
	value := "0"
	if typing {
		value = "1"
	}
	return updateTypingScript.Run(ctx, r.client, []string{key}, value, 15).Int64()
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

func (r *RedisClient) Close() error {
	return r.client.Close()
}

func GetRedisChannel() string {
	return "ws_channel"
}
