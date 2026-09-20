package _redis

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

const (
	// TileSessionAESPrefix — ss:k:{sessionK} → sessionEncryptKey (AES base64).
	TileSessionAESPrefix = "ss:k:"
	TileSessionTTL       = 12 * time.Hour
)

func tileSessionAESKey(sessionK uint64) string {
	return TileSessionAESPrefix + strconv.FormatUint(sessionK, 10)
}

func (r *RedisService) tileClient() (*redis.Client, error) {
	if r == nil || r.Client == nil {
		return nil, fmt.Errorf("redis client not initialized")
	}
	if _, err := r.Client.Ping(r.Ctx).Result(); err != nil {
		return nil, fmt.Errorf("redis ping: %w", err)
	}
	return r.Client, nil
}

// SaveTileSession lưu sessionEncryptKey theo sessionK (TTL theo expiresAt).
func (r *RedisService) SaveTileSession(sessionK uint64, sessionEncryptKey string, expiresAt time.Time) error {
	if sessionK == 0 || sessionEncryptKey == "" {
		return fmt.Errorf("sessionK and sessionEncryptKey are required")
	}
	client, err := r.tileClient()
	if err != nil {
		return err
	}
	ttl := time.Until(expiresAt)
	if ttl <= 0 {
		ttl = TileSessionTTL
	}
	kRedis := tileSessionAESKey(sessionK)
	slog.Info(fmt.Sprintf("save tile session: %s, %s, %v", kRedis, sessionEncryptKey, sessionK))
	if err = client.Set(r.Ctx, kRedis, sessionEncryptKey, ttl).Err(); err != nil {
		return fmt.Errorf("save tile session: %w", err)
	}
	return nil
}

// GetTileSessionEncryptKey lấy sessionEncryptKey theo sessionK (ss:k:{sessionK}).
func (r *RedisService) GetTileSessionEncryptKey(sessionK uint64) (string, error) {
	if sessionK == 0 {
		return "", redis.Nil
	}
	client, err := r.tileClient()
	if err != nil {
		return "", err
	}
	val, err := client.Get(r.Ctx, tileSessionAESKey(sessionK)).Result()
	if errors.Is(err, redis.Nil) {
		return "", redis.Nil
	}
	return val, err
}

// DeleteTileSession xóa session theo sessionK.
func (r *RedisService) DeleteTileSession(ctx context.Context, sessionK uint64) error {
	if sessionK == 0 {
		return nil
	}
	client, err := r.tileClient()
	if err != nil {
		return err
	}
	c := r.Ctx
	if ctx != nil {
		c = ctx
	}
	return client.Del(c, tileSessionAESKey(sessionK)).Err()
}
