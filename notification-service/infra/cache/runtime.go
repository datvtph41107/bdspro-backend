package cache

import (
	"context"
	"sync"

	"notification/config"

	"github.com/go-redis/redis/v8"
	"github.com/hyperledger/fabric/common/flogging"
)

type RedisClient struct {
	client    *redis.Client
	logger    *flogging.FabricLogger
	mu        sync.RWMutex
	connected bool
}

func NewRedisClient() *RedisClient {
	logger := flogging.MustGetLogger("redis")
	client := redis.NewClient(&redis.Options{
		Addr:     config.AppProperties.Redis.Host,
		Password: config.AppProperties.Redis.Password,
		DB:       0,
	})

	rc := &RedisClient{
		client:    client,
		logger:    logger,
		connected: false,
	}

	// Thử kết nối nhưng không crash nếu lỗi
	_, err := client.Ping(client.Context()).Result()
	if err != nil {
		logger.Warnf("Failed to connect to Redis on startup: %v. App will continue running and retry on next operation.", err)
		rc.connected = false
	} else {
		logger.Info("Connected to Redis")
		rc.connected = true
	}

	return rc
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

// Publish message to a Redis channel
func (r *RedisClient) Publish(ctx context.Context, channel string, message string) error {
	// Thử kết nối lại nếu chưa kết nối
	if err := r.ensureConnection(ctx); err != nil {
		r.logger.Warnf("Cannot publish to Redis channel %s: connection unavailable", channel)
		return err
	}

	err := r.client.Publish(ctx, channel, message).Err()
	if err != nil {
		// Nếu lỗi khi publish, đánh dấu là mất kết nối
		r.mu.Lock()
		r.connected = false
		r.mu.Unlock()
		r.logger.Errorf("Failed to publish message to Redis channel %s: %v", channel, err)
	}
	return err
}

// Get lấy giá trị từ Redis key
func (r *RedisClient) Get(ctx context.Context, key string) (string, error) {
	// Thử kết nối lại nếu chưa kết nối
	if err := r.ensureConnection(ctx); err != nil {
		r.logger.Warnf("Cannot get key %s from Redis: connection unavailable", key)
		return "", err
	}

	result, err := r.client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return "", nil // Key không tồn tại
		}
		// Nếu lỗi khi get, đánh dấu là mất kết nối
		r.mu.Lock()
		r.connected = false
		r.mu.Unlock()
		r.logger.Errorf("Failed to get key %s from Redis: %v", key, err)
		return "", err
	}
	return result, nil
}
