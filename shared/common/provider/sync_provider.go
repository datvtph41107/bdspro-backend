package _provider

import (
	_redis "common/redis"
	_utils "common/utils"
	"context"
	"fmt"
	"strconv"
	"time"
)

// Prefix key Redis cho sync timestamp (format: prefix:profileId:resourceId)
const (
	SyncKeyProductMe           = "prd:me"
	SyncKeyProductDetail       = "prd:id"
	SyncKeyProductPriceHist    = "prd:phis"
	SyncKeyProductDistrs       = "prd:dist"
	SyncKeyProductPosts        = "prd:post"
	SyncKeyProductDeals        = "prd:deal"
	SyncKeyProductNotes        = "prd:note"
	SyncKeyProductAppointments = "prd:appt"
	SyncKeyProductContacts     = "prd:ct"
	SyncKeyPropertyHistories   = "ppt:his"
	SyncKeyPropertyDetail      = "ppt:id"
	SyncKeyPropertyMe          = "ppt:me"
	SyncKeyContactMe           = "ct:me"
	SyncKeyContactDetail       = "ct:id"
	// SyncKeyUserGuideByKey: 1 key chung cho toàn bộ user guides có key (timestamp = max updatedAt).
	SyncKeyUserGuideByKey = "hub:ug:bykey"
	// SyncKeyNotificationMe: ds thông báo + count chưa đọc của user.
	SyncKeyNotificationMe = "nt:me"
	// Conversation sync keys - client so sánh timestamp để biết API nào cần gọi cập nhật
	SyncKeyConversationDetail   = "cnv:id"
	SyncKeyConversationFiles    = "cnv:fs"
	SyncKeyConversationLinks    = "cnv:ls"
	SyncKeyConversationImages   = "cnv:imgs"
	SyncKeyConversationMembers  = "cnv:mems"
	SyncKeyConversationSetting  = "cnv:set"
	SyncKeyConversationPinned   = "cnv:pi"
	SyncKeyConversationMessages = "cnv:msgs"
)

type SyncProvider struct {
	redis *_redis.RedisService
}

func NewSyncProvider(rd *_redis.RedisService) *SyncProvider {
	return &SyncProvider{
		redis: rd,
	}
}

// GetKey trả về key Redis theo format {prefix}:{profileId}:{resourceId}
// Ví dụ: GetKey(ctx, "product", 123) -> "product:1:123" (1 là profileId từ context)
func (s *SyncProvider) GetKey(ctx context.Context, prefix string, resourceId uint64) string {
	profileId := _utils.GetProfileIdWithContext(ctx)
	return fmt.Sprintf("%s:%d:%d", prefix, profileId, resourceId)
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
	s.redis.SetWithTime(key, timestamp, 48*time.Hour)
	return nil
}

func (s *SyncProvider) PutTimeRequest(ctx context.Context, key string, timestamp int64) {
	ts := time.Now().UnixMilli() - 5000
	if timestamp > 0 {
		ts = timestamp
	}
	s.PutTimestamp(ctx, key, ts)
}

// Del xóa key trong Redis (để client sync lại ds/chi tiết).
func (s *SyncProvider) Del(ctx context.Context, key string) error {
	return s.redis.Delete(key)
}

// GetTimestamp lấy timestamp từ Redis, trả về 0 nếu key không tồn tại
func (s *SyncProvider) GetTimestamp(ctx context.Context, key string) int64 {
	str, err := s.redis.Get(key)
	if err != nil {
		return 0
	}
	ts, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		return 0
	}
	return ts
}

func (s *SyncProvider) MGet(ctx context.Context, keys []string) ([]interface{}, error) {
	vals, err := s.redis.Client.MGet(ctx, keys...).Result()
	if err != nil {
		return nil, err
	}
	return vals, nil
}

// func (s *SyncProvider) DelTimestamp(ctx context.Context, key string, timestamp int64) error {
// 	s.redis.Set(key, strconv.FormatInt(timestamp, 10))
// 	return nil
// }
