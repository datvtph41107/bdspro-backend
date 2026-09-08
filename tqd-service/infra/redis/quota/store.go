package redisquota

import (
	"errors"
	"time"

	"github.com/redis/go-redis/v9"
)

const defaultKeepAfterPeriod = 48 * time.Hour

/**
 * Store giữ quota runtime trong Redis.
 *
 * Lua scripts được dùng để Reserve/Commit/Cancel không bị race giữa nhiều request.
 */
type Store struct {
	client *redis.Client
}

func NewStore(client *redis.Client) *Store {
	return &Store{client: client}
}

func (s *Store) validate() error {
	if s == nil || s.client == nil {
		return errors.New("redis quota store is not configured")
	}
	return nil
}

func keepTTL(periodEnd time.Time, now time.Time) time.Duration {
	if periodEnd.IsZero() {
		return 0
	}
	if periodEnd.After(now) {
		return periodEnd.Sub(now) + defaultKeepAfterPeriod
	}
	return defaultKeepAfterPeriod
}
