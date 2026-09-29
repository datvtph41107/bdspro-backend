package tilesession

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	commonredis "common/redis"

	"github.com/redis/go-redis/v9"
)

const (
	keyPrefix = "ss:k:"

	// TTL is the canonical lifetime for a tile encryption session.
	TTL = 12 * time.Hour
)

// Store owns the shared tile-session Redis keyspace and persistence semantics.
// Redis connection lifecycle remains owned by common/redis and the process root.
type Store struct {
	redis *commonredis.RedisService
}

func NewStore(service *commonredis.RedisService) *Store {
	return &Store{redis: service}
}

func (s *Store) Save(ctx context.Context, sessionK uint64, sessionEncryptKey string, expiresAt time.Time) error {
	if sessionK == 0 || sessionEncryptKey == "" {
		return errors.New("sessionK and sessionEncryptKey are required")
	}
	client, err := s.client()
	if err != nil {
		return err
	}
	if ctx == nil {
		ctx = context.Background()
	}
	if err := client.Set(ctx, key(sessionK), sessionEncryptKey, expirationTTL(expiresAt, time.Now())).Err(); err != nil {
		return fmt.Errorf("save tile session: %w", err)
	}
	return nil
}

func (s *Store) Get(ctx context.Context, sessionK uint64) (string, error) {
	if sessionK == 0 {
		return "", redis.Nil
	}
	client, err := s.client()
	if err != nil {
		return "", err
	}
	if ctx == nil {
		ctx = context.Background()
	}
	value, err := client.Get(ctx, key(sessionK)).Result()
	if errors.Is(err, redis.Nil) {
		return "", redis.Nil
	}
	if err != nil {
		return "", fmt.Errorf("get tile session: %w", err)
	}
	return value, nil
}

func (s *Store) Delete(ctx context.Context, sessionK uint64) error {
	if sessionK == 0 {
		return nil
	}
	client, err := s.client()
	if err != nil {
		return err
	}
	if ctx == nil {
		ctx = context.Background()
	}
	if err := client.Del(ctx, key(sessionK)).Err(); err != nil {
		return fmt.Errorf("delete tile session: %w", err)
	}
	return nil
}

func (s *Store) client() (*redis.Client, error) {
	if s == nil || s.redis == nil {
		return nil, errors.New("tile session redis unavailable")
	}
	client, err := s.redis.Instance()
	if err != nil {
		return nil, fmt.Errorf("tile session redis: %w", err)
	}
	return client, nil
}

func key(sessionK uint64) string {
	return keyPrefix + strconv.FormatUint(sessionK, 10)
}

func expirationTTL(expiresAt, now time.Time) time.Duration {
	ttl := expiresAt.Sub(now)
	if ttl <= 0 {
		return TTL
	}
	return ttl
}
