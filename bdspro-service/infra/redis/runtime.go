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

// =====================================================================
// CHỈ CÁC LỆNH NGUYÊN THỦY - KHÔNG VIẾT HÀM RIPÊNG
// =====================================================================

// ---------- STRING OPERATIONS ----------
func (r *RedisClient) Get(ctx context.Context, key string) *redis.StringCmd {
	return r.client.Get(ctx, key)
}

func (r *RedisClient) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd {
	return r.client.Set(ctx, key, value, expiration)
}

func (r *RedisClient) Del(ctx context.Context, keys ...string) *redis.IntCmd {
	return r.client.Del(ctx, keys...)
}

// ---------- HASH OPERATIONS ----------
func (r *RedisClient) HSet(ctx context.Context, key string, values ...interface{}) *redis.IntCmd {
	return r.client.HSet(ctx, key, values...)
}

func (r *RedisClient) HGet(ctx context.Context, key, field string) *redis.StringCmd {
	return r.client.HGet(ctx, key, field)
}

func (r *RedisClient) HGetAll(ctx context.Context, key string) *redis.MapStringStringCmd {
	return r.client.HGetAll(ctx, key)
}

func (r *RedisClient) HMGet(
	ctx context.Context,
	key string,
	fields ...string,
) *redis.SliceCmd {
	return r.client.HMGet(ctx, key, fields...)
}

func (r *RedisClient) HMSet(ctx context.Context, key string, values ...interface{}) *redis.BoolCmd {
	return r.client.HMSet(ctx, key, values...)
}

func (r *RedisClient) HDel(ctx context.Context, key string, fields ...string) *redis.IntCmd {
	return r.client.HDel(ctx, key, fields...)
}

func (r *RedisClient) ZRangeByScoreWithScores(ctx context.Context, key string, opt *redis.ZRangeBy) *redis.ZSliceCmd {
	return r.client.ZRangeByScoreWithScores(ctx, key, opt)
}

func (r *RedisClient) ZRevRangeByScoreWithScores(
	ctx context.Context,
	key string,
	opt *redis.ZRangeBy,
) *redis.ZSliceCmd {
	return r.client.ZRevRangeByScoreWithScores(ctx, key, opt)
}

func (r *RedisClient) ExpireNX(
	ctx context.Context,
	key string,
	expiration time.Duration,
) *redis.BoolCmd {
	return r.client.ExpireNX(ctx, key, expiration)
}

// ---------- SET OPERATIONS ----------
func (r *RedisClient) SAdd(ctx context.Context, key string, members ...interface{}) *redis.IntCmd {
	return r.client.SAdd(ctx, key, members...)
}

func (r *RedisClient) SMembers(ctx context.Context, key string) *redis.StringSliceCmd {
	return r.client.SMembers(ctx, key)
}

func (r *RedisClient) SRem(ctx context.Context, key string, members ...interface{}) *redis.IntCmd {
	return r.client.SRem(ctx, key, members...)
}

// ---------- SORTED SET OPERATIONS ----------
func (r *RedisClient) ZAdd(ctx context.Context, key string, members ...redis.Z) *redis.IntCmd {
	return r.client.ZAdd(ctx, key, members...)
}

func (r *RedisClient) ZRangeByScore(ctx context.Context, key string, opt *redis.ZRangeBy) *redis.StringSliceCmd {
	return r.client.ZRangeByScore(ctx, key, opt)
}

func (r *RedisClient) ZRem(ctx context.Context, key string, members ...interface{}) *redis.IntCmd {
	return r.client.ZRem(ctx, key, members...)
}

// ---------- KEY OPERATIONS ----------
func (r *RedisClient) Expire(ctx context.Context, key string, expiration time.Duration) *redis.BoolCmd {
	return r.client.Expire(ctx, key, expiration)
}

func (r *RedisClient) TTL(ctx context.Context, key string) *redis.DurationCmd {
	return r.client.TTL(ctx, key)
}

func (r *RedisClient) Keys(ctx context.Context, pattern string) *redis.StringSliceCmd {
	return r.client.Keys(ctx, pattern)
}

// ---------- PIPELINE ----------
func (r *RedisClient) Pipeline() redis.Pipeliner {
	return r.client.Pipeline()
}
