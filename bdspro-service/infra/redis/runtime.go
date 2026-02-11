// infra/client/redis/client.go
package redis_cli

import (
	"context"
	"fmt"
	"time"

	configs "bdspro/config"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

// RedisClient wraps Redis connection
type RedisClient struct {
	client *redis.Client
	logger *zap.Logger
}

// NewClient creates a new Redis client
func NewClient() (*RedisClient, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     configs.AppProperties.Redis.Host,
		Password: configs.AppProperties.Redis.Password,
	})

	// Test connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	return &RedisClient{
		client: client,
	}, nil
}

// Client returns the underlying Redis client
func (r *RedisClient) Client() *redis.Client {
	return r.client
}

// Close closes the Redis connection
func (r *RedisClient) Close() error {
	return r.client.Close()
}

// HealthCheck checks Redis health
func (r *RedisClient) HealthCheck(ctx context.Context) error {
	return r.client.Ping(ctx).Err()
}

// Stats returns Redis statistics
func (r *RedisClient) Stats(ctx context.Context) (*redis.PoolStats, error) {
	stats := r.client.PoolStats()
	return stats, nil
}

// ============= Wrapper Methods for Redis Operations =============

// Get wraps redis.Client.Get
func (r *RedisClient) Get(ctx context.Context, key string) *redis.StringCmd {
	return r.client.Get(ctx, key)
}

// Set wraps redis.Client.Set
func (r *RedisClient) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd {
	return r.client.Set(ctx, key, value, expiration)
}

// SAdd wraps redis.Client.SAdd
func (r *RedisClient) SAdd(ctx context.Context, key string, members ...interface{}) *redis.IntCmd {
	return r.client.SAdd(ctx, key, members...)
}

// SMembers wraps redis.Client.SMembers
func (r *RedisClient) SMembers(ctx context.Context, key string) *redis.StringSliceCmd {
	return r.client.SMembers(ctx, key)
}

// Del wraps redis.Client.Del
func (r *RedisClient) Del(ctx context.Context, keys ...string) *redis.IntCmd {
	return r.client.Del(ctx, keys...)
}

// Expire wraps redis.Client.Expire
func (r *RedisClient) Expire(ctx context.Context, key string, expiration time.Duration) *redis.BoolCmd {
	return r.client.Expire(ctx, key, expiration)
}

// Pipeline wraps redis.Client.Pipeline
func (r *RedisClient) Pipeline() redis.Pipeliner {
	return r.client.Pipeline()
}

// SCard wraps redis.Client.SCard
func (r *RedisClient) SCard(ctx context.Context, key string) *redis.IntCmd {
	return r.client.SCard(ctx, key)
}

// TTL wraps redis.Client.TTL
func (r *RedisClient) TTL(ctx context.Context, key string) *redis.DurationCmd {
	return r.client.TTL(ctx, key)
}
