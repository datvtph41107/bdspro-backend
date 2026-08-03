package _provider

import (
	_redis "common/redis"

	"github.com/google/wire"
)

var SyncProviderSet = wire.NewSet(
	NewRedisService,
)

func NewRedisService() *_redis.RedisService {
	return _redis.NewRedisService()
}
