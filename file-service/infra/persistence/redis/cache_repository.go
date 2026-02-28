// file/internal/infrastructure/persistence/redis/cache_repository.go
package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"file/internal/ports"

	"github.com/redis/go-redis/v9"
)

type cacheRepository struct {
	client *redis.Client
	prefix string
}

func NewCacheRepository(client *redis.Client) ports.CacheRepository {
	return &cacheRepository{
		client: client,
		prefix: "file-service:",
	}
}

func (r *cacheRepository) Get(ctx context.Context, key string) ([]byte, error) {
	fullKey := r.prefix + key
	data, err := r.client.Get(ctx, fullKey).Bytes()
	if err == redis.Nil {
		return nil, nil
	}
	return data, err
}

func (r *cacheRepository) Set(ctx context.Context, key string, value []byte, ttl time.Duration) error {
	fullKey := r.prefix + key
	return r.client.Set(ctx, fullKey, value, ttl).Err()
}

func (r *cacheRepository) Delete(ctx context.Context, key string) error {
	fullKey := r.prefix + key
	return r.client.Del(ctx, fullKey).Err()
}

func (r *cacheRepository) Exists(ctx context.Context, key string) (bool, error) {
	fullKey := r.prefix + key
	result, err := r.client.Exists(ctx, fullKey).Result()
	return result > 0, err
}

func (r *cacheRepository) Incr(ctx context.Context, key string) (int64, error) {
	fullKey := r.prefix + key
	return r.client.Incr(ctx, fullKey).Result()
}

func (r *cacheRepository) SetNX(ctx context.Context, key string, value []byte, ttl time.Duration) (bool, error) {
	fullKey := r.prefix + key
	return r.client.SetNX(ctx, fullKey, value, ttl).Result()
}

func (r *cacheRepository) GetJSON(ctx context.Context, key string, v interface{}) error {
	data, err := r.Get(ctx, key)
	if err != nil || data == nil {
		return err
	}
	return json.Unmarshal(data, v)
}

func (r *cacheRepository) SetJSON(ctx context.Context, key string, v interface{}, ttl time.Duration) error {
	data, err := json.Marshal(v)
	if err != nil {
		return err
	}
	return r.Set(ctx, key, data, ttl)
}

func (r *cacheRepository) HSet(ctx context.Context, key string, values map[string]interface{}) error {
	fullKey := r.prefix + key
	return r.client.HSet(ctx, fullKey, values).Err()
}

func (r *cacheRepository) HGet(ctx context.Context, key, field string) (string, error) {
	fullKey := r.prefix + key
	return r.client.HGet(ctx, fullKey, field).Result()
}

// Key generators
func (r *cacheRepository) FileKey(fileID string) string {
	return fmt.Sprintf("file:%s", fileID)
}

func (r *cacheRepository) FileAccessKey(fileID string, userID uint64) string {
	return fmt.Sprintf("access:%s:%d", fileID, userID)
}

func (r *cacheRepository) SignatureKey(signature string) string {
	return fmt.Sprintf("signature:%s", signature)
}

func (r *cacheRepository) UploadSessionKey(sessionID string) string {
	return fmt.Sprintf("upload:session:%s", sessionID)
}

func (r *cacheRepository) RateLimitKey(ip string, endpoint string) string {
	return fmt.Sprintf("rate:%s:%s", ip, endpoint)
}
