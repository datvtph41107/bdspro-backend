package providers

import "context"

type CacheProvider interface {
	LRange(c context.Context, key string, start, stop int64) ([]string, error)
	LTrim(c context.Context, key string, start, stop int64) error
	RPush(c context.Context, key string, id uint64) error
	Get(c context.Context, key string) (string, error)
	Set(c context.Context, key string, value string) error
}
