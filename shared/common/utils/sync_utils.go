package _utils

import (
	_redis "common/redis"
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"
)

var errSyncRedisUnavailable = errors.New("redis unavailable")

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
	// SyncKeyAppointmentMe: ds lịch hẹn của user (ListAppointments).
	SyncKeyAppointmentMe = "appt:me"
	// SyncKeyDealMe: ds thương vụ của user (GetDealsUser).
	SyncKeyDealMe = "deal:me"
	// Deal-scoped sync keys (resourceId = dealId) — /v2.1/org/deal/{id}/ts
	SyncKeyDealDetail     = "deal:id"
	SyncKeyDealMembers    = "deal:mems"
	SyncKeyDealProducts   = "deal:prds"
	SyncKeyDealInvestment = "deal:inv"
	SyncKeyDealBank       = "deal:bank"
	SyncKeyDealProcess    = "deal:proc"
	SyncKeyDealCommission = "deal:comm"
	SyncKeyDealHistory    = "deal:his"
	SyncKeyDealSummary    = "deal:sum"
	// SyncKeyDealPermission: /v2/auth/permission/deal/{id}
	SyncKeyDealPermission = "deal:perm"
	// SyncKeyDealContracts: /v2.1/transaction/contract/deal/list/{dealId} — FE key trong /ts: "transaction"
	SyncKeyDealContracts = "deal:ctr"
	// SyncKeyDealInvitation: /v2.1/org/deal-invitation/list/{dealId}
	SyncKeyDealInvitation = "deal:invt"
	// SyncKeyDealStatistics: /v2.1/transaction/contract/deal/statistics/{dealId}
	SyncKeyDealStatistics = "deal:stat"
	// SyncKeyDealCost: /v2.1/transaction/deal-cost/list/{dealId}
	SyncKeyDealCost = "deal:cost"
	// SyncKeyPermissionMe: quyền của user (/v2/auth/permission/me).
	SyncKeyPermissionMe = "perm:me"
	// SyncKeyProfileInfoMe: profile info me (/v2/user/profile/info/me).
	SyncKeyProfileInfoMe = "usr:info:me"
	// SyncKeyBdsproEnums: enum bdspro public (global, không theo profile).
	SyncKeyBdsproEnums = "bds:enums"
	// SyncKeyPaymentBankAll: danh mục ngân hàng GET /v2/payment/bank/all (global).
	SyncKeyPaymentBankAll = "pay:bank:all"
	// SyncKeyTxCostTypes: danh mục loại chi phí GET /v4/transaction/cost-type (global, mọi costType).
	SyncKeyTxCostTypes = "tx:cost:types"
	// SyncKeyPoiCategoryIcons: icon POI categories (global).
	SyncKeyPoiCategoryIcons = "tqd:poi:icons"
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
	SyncKeyChatUnreadMe         = "cnv:ur"
	SyncKeyConversationsMe      = "cnv:me"
	SyncKeyConversationCommon   = "cnv:c"
	// TQD sync keys
	SyncKeyTQDLayerDetail      = "tqd:layer:id"
	SyncKeyTQDLayerList        = "tqd:layer:me"
	SyncKeyTQDLayerFamilyList  = "tqd:layer-family:me"
	SyncKeyTQDRegionDetail     = "tqd:region:id"
	SyncKeyTQDRegionList       = "tqd:region:me"
	SyncKeyTQDParcelDetail     = "tqd:parcel:id"
	SyncKeyTQDParcelLayers     = "tqd:parcel:layers"
	SyncKeyTQDLabelDetail      = "tqd:label:id"
	SyncKeyTQDLabelList        = "tqd:label:me"
	SyncKeyTQDSubscriptionList = "tqd:sub:me"
	SyncKeyTQDNotificationList = "tqd:notif:me"
	SyncKeyTQDPOIDetail        = "tqd:poi:id"
	SyncKeyTQDPOIList          = "tqd:poi:me"
	SyncKeyTQDLocationSearch   = "tqd:loc:me"
	// SyncKeyTQDProvincesWards: 1 key chung cho danh mục province + ward công khai.
	SyncKeyTQDProvincesWards = "tqd:loc:provinces-wards"
	SyncKeyTQDReportList     = "tqd:report:me"
	SyncKeyTQDFeatureList    = "tqd:feature:me"
)

type SyncUtil struct {
	redis *_redis.RedisService
}

func NewSyncUtil(rd *_redis.RedisService) *SyncUtil {
	return &SyncUtil{
		redis: rd,
	}
}

// GetKey trả về key Redis theo format {profileId}:{resourceId}:{prefix}
// Ví dụ: GetKey(ctx, "product", 123) -> "1:123:product" (1 là profileId từ context)
func (s *SyncUtil) GetKey(ctx context.Context, prefix string, resourceId uint64) string {
	profileId := GetProfileIdWithContext(ctx)
	return s.GetKeyForProfile(profileId, resourceId, prefix)
}

// GetKey trả về key Redis theo format {profileId}:{resourceId}:{prefix}
// Ví dụ: GetKey(ctx, "product", 123) -> "1:123:product" (1 là profileId từ context)
func (s *SyncUtil) GetKeyWithoutMe(ctx context.Context, prefix string, resourceId uint64) string {
	return fmt.Sprintf("%d:%s", resourceId, prefix)
}

// GetKeyForProfile trả về key Redis theo profileId cụ thể (dùng khi cần invalidate key của user khác)
func (s *SyncUtil) GetKeyForProfile(profileId, resourceId uint64, prefix string) string {
	return fmt.Sprintf("%d:%d:%s", profileId, resourceId, prefix)
}

func (s *SyncUtil) HasUpdated(ctx context.Context, key string, timestamp int64) bool {
	str, err := s.redis.Get(key)
	if err != nil {
		fmt.Println("err", err)
		return true
	}
	fmt.Println("str", str)
	ts, err := strconv.ParseInt(str, 10, 64)
	if ts > 0 {
		return timestamp < ts
	}

	return true
}

// return: isChanged, cache exists
func (s *SyncUtil) IsCacheChanged(ctx context.Context, key string, timestamp int64) (bool, bool) {
	str, err := s.redis.Get(key)
	if err != nil {
		return true, false
	}
	ts, err := strconv.ParseInt(str, 10, 64)
	if ts > 0 {
		return timestamp < ts, true
	}

	return true, false
}

func (s *SyncUtil) PutTimestamp(ctx context.Context, key string, timestamp int64) error {
	if s.redis == nil {
		return errSyncRedisUnavailable
	}
	if timestamp <= 0 {
		timestamp = time.Now().UnixMilli()
	}
	return s.redis.SetWithTime(key, timestamp, 48*time.Hour)
}

func (s *SyncUtil) PutTimeRequest(ctx context.Context, key string, timestamp int64) {
	ts := time.Now().UnixMilli() - 5000
	if timestamp > 0 {
		ts = timestamp
	}
	s.PutTimestamp(ctx, key, ts)
}

// Del xóa key trong Redis (để client sync lại ds/chi tiết).
func (s *SyncUtil) Del(ctx context.Context, key string) error {
	return s.redis.Delete(key)
}

// GetTimestamp lấy timestamp từ Redis, trả về 0 nếu key không tồn tại
func (s *SyncUtil) GetTimestamp(ctx context.Context, key string) int64 {
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

// EnsureTimestamp trả timestamp hiện có trên Redis.
// Nếu key chưa có, hoặc tồn tại mà = 0 / không parse được → ghi now (ms) rồi trả về.
// Dùng cho catalog global (poi icons, provinces/wards, …) — không overwrite khi đã có ts > 0.
func (s *SyncUtil) EnsureTimestamp(ctx context.Context, key string) int64 {
	if ts := s.GetTimestamp(ctx, key); ts > 0 {
		return ts
	}
	now := time.Now().UnixMilli()
	if err := s.PutTimestamp(ctx, key, now); err != nil {
		return 0
	}
	if ts := s.GetTimestamp(ctx, key); ts > 0 {
		return ts
	}
	return 0
}

// MGet đọc nhiều timestamp. Redis down/unavailable → trả về slice nil (timestamp=0),
// không error — FE coi như chưa cache và gọi detail bình thường.
func (s *SyncUtil) MGet(ctx context.Context, keys []string) ([]interface{}, error) {
	empty := make([]interface{}, len(keys))
	if len(keys) == 0 || s.redis == nil {
		return empty, nil
	}
	vals, err := s.redis.MGet(ctx, keys...)
	if err != nil {
		return empty, nil
	}
	return vals, nil
}

// func (s *SyncProvider) DelTimestamp(ctx context.Context, key string, timestamp int64) error {
// 	s.redis.Set(key, strconv.FormatInt(timestamp, 10))
// 	return nil
// }
