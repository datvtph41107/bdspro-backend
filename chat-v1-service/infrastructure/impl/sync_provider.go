package impl

import (
	_provider "common/provider"
	_redis "common/redis"
)

func NewSyncProvider(rd *_redis.RedisService) *_provider.SyncProvider {
	return _provider.NewSyncProvider(rd)
}
