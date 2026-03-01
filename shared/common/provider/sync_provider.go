package _provider

import (
	_redis "common/redis"
	"context"
	"strconv"
)

type SyncProvider struct {
	redis *_redis.RedisService
}

func NewSyncProvider(rd *_redis.RedisService) *SyncProvider {
	return &SyncProvider{
		redis: rd,
	}
}

func (s *SyncProvider) HasUpdated(ctx context.Context, key string, timestamp int64) bool {
	str, err := s.redis.Get(key)
	if err != nil {
		return true
	}
	ts, err := strconv.ParseInt(str, 10, 64)
	if ts > 0 {
		return timestamp < ts
	}

	return true
}

func (s *SyncProvider) PutTimestamp(ctx context.Context, key string, timestamp int64) error {
	s.redis.Set(key, strconv.FormatInt(timestamp, 10))
	return nil
}

// func (s *SyncProvider) DelTimestamp(ctx context.Context, key string, timestamp int64) error {
// 	s.redis.Set(key, strconv.FormatInt(timestamp, 10))
// 	return nil
// }
