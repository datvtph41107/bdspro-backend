package _provider

import (
	"context"
)

type CacheTimeProvider interface {
	GetCacheTime(ctx context.Context, key string) int64
}
