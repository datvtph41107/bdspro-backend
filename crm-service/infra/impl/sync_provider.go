package impl

import (
	_provider "common/provider"
	_redis "common/redis"
)

func NewCommonRedisProvider() *_redis.RedisService {
	return _redis.NewRedisService()
}

func NewSyncProvider(rd *_redis.RedisService) *_provider.SyncProvider {
	return _provider.NewSyncProvider(rd)
}
