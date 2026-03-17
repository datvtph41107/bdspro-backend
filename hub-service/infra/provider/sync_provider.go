package provider

import (
	_provider "common/provider"
	_redis "common/redis"
)

// NewSyncProvider tạo SyncProvider dùng Redis (cho sync timestamp user guide by key).
func NewSyncProvider() *_provider.SyncProvider {
	return _provider.NewSyncProvider(_redis.NewRedisService())
}
