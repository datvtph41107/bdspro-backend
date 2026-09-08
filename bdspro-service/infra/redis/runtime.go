package redis_cli

import (
	"context"
	"sync"
	"time"

	configs "bdspro/config"

	"github.com/hyperledger/fabric/common/flogging"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

// RedisClient wraps Redis connection
type RedisClient struct {
	client    *redis.Client
	logger    *flogging.FabricLogger
	connected bool
	mu        sync.RWMutex
}

// NewClient creates a new Redis client
func NewClient() (*RedisClient, error) {
	logger := flogging.MustGetLogger("redis")
	client := redis.NewClient(&redis.Options{
		Addr:     configs.AppProperties.Redis.Host,
		Password: configs.AppProperties.Redis.Password,
		DB:       0,
	})

	rc := &RedisClient{
		client:    client,
		logger:    logger,
		connected: false,
	}

	// Thử kết nối nhưng không crash nếu lỗi
	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		logger.Warn("Failed to connect to Redis on startup: %v. App will continue running and retry on next operation.", zap.Error(err))
		rc.connected = false
	} else {
		logger.Info("Connected to Redis")
		rc.connected = true
	}

	return rc, nil
}

// ensureConnection kiểm tra và kết nối lại nếu cần
func (r *RedisClient) ensureConnection(ctx context.Context) error {
	r.mu.RLock()
	connected := r.connected
	r.mu.RUnlock()

	if connected {
		// Kiểm tra kết nối còn hoạt động không
		_, err := r.client.Ping(ctx).Result()
		if err == nil {
			return nil
		}
		// Kết nối bị mất, đánh dấu là chưa kết nối
		r.mu.Lock()
		r.connected = false
		r.mu.Unlock()
		r.logger.Warnf("Redis connection lost: %v", err)
	}

	// Thử kết nối lại
	r.mu.Lock()
	defer r.mu.Unlock()

	// Double check sau khi lock
	if r.connected {
		return nil
	}

	_, err := r.client.Ping(ctx).Result()
	if err != nil {
		r.logger.Warnf("Failed to reconnect to Redis: %v", err)
		return err
	}

	r.connected = true
	r.logger.Info("Reconnected to Redis successfully")
	return nil
}

// Client returns the underlying Redis client
func (r *RedisClient) Client(ctx context.Context) (*redis.Client, error) {
	if err := r.ensureConnection(ctx); err != nil {
		r.logger.Warnf("Cannot get Redis client: connection unavailable")
		return nil, err
	}
	return r.client, nil
}

// Close closes the Redis connection
func (r *RedisClient) Close() error {
	client, err := r.Client(context.Background())
	if err != nil {
		return err
	}
	return client.Close()
}

// =====================================================================
// CHỈ CÁC LỆNH NGUYÊN THỦY - KHÔNG VIẾT HÀM RIPÊNG
// =====================================================================

// ---------- STRING OPERATIONS ----------
func (r *RedisClient) Get(ctx context.Context, key string) *redis.StringCmd {
	client, err := r.Client(ctx)
	if err != nil {
		return nil
	}
	return client.Get(ctx, key)
}

func (r *RedisClient) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd {
	client, err := r.Client(ctx)
	if err != nil {
		return nil
	}
	return client.Set(ctx, key, value, expiration)
}

func (r *RedisClient) Del(ctx context.Context, keys ...string) *redis.IntCmd {
	client, err := r.Client(ctx)
	if err != nil {
		return nil
	}
	return client.Del(ctx, keys...)
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
	client, err := r.Client(ctx)
	if err != nil {
		return nil
	}
	return client.Expire(ctx, key, expiration)
}

func (r *RedisClient) TTL(ctx context.Context, key string) *redis.DurationCmd {
	client, err := r.Client(ctx)
	if err != nil {
		return nil
	}
	return client.TTL(ctx, key)
}

func (r *RedisClient) Keys(ctx context.Context, pattern string) *redis.StringSliceCmd {
	return r.client.Keys(ctx, pattern)
}

// ---------- PIPELINE ----------
func (r *RedisClient) Pipeline() redis.Pipeliner {
	return r.client.Pipeline()
}
