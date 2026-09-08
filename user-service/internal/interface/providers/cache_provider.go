package providers

import (
	"context"
	"time"
)

type CacheProvider interface {
	SaveToken(c context.Context, userID uint64, token string, expirationInSeconds uint64) error
	GetToken(c context.Context, userID uint64) (string, error)
	IsTokenValid(c context.Context, userID uint64, token string) (bool, error)
	DeleteToken(c context.Context, userID uint64) error
	Increment(c context.Context, key string) error
	SaveAppState(c context.Context, userID uint64, state string) error

	MGet(c context.Context, keys ...string) ([]interface{}, error)
	Set(c context.Context, key string, value interface{}, expiration time.Duration) error
	Get(c context.Context, key string) (string, error)
	Delete(c context.Context, key string) error
	DeleteAllKeysWithPrefix(c context.Context, prefix string) error
	Expire(c context.Context, key string, expiration time.Duration) error
	IncrementBy(c context.Context, key string, delta int64) error
	GetAsInt(c context.Context, key string) (int, error)
}
