package provider

import (
	"context"
)

type CacheProvider interface {
	SaveToken(c context.Context, userID uint64, token string, expirationInSeconds uint64) error
	GetToken(c context.Context, userID uint64) (string, error)
	IsTokenValid(c context.Context, userID uint64, token string) (bool, error)
	DeleteToken(c context.Context, userID uint64) error
	Increment(c context.Context, key string) error
}
