package services

import (
	"context"
	"strconv"

	_errors "common/errors"
	_utils "common/utils"

	"userv4/modules/bdspro/dto"
)

// SyncService owns the app-wide sync timestamp aggregator (GET /v2/bdspro/ts).
type SyncService interface {
	GetAppTimestamps(ctx context.Context) (*dto.AppTimestampsResponse, error)
}

// LocalSyncService ports ProductHandler.GetAppTimestamps onto SyncUtil.
type LocalSyncService struct {
	sync *_utils.SyncUtil
}

func NewLocalSyncService(sync *_utils.SyncUtil) *LocalSyncService {
	return &LocalSyncService{sync: sync}
}

func (s *LocalSyncService) GetAppTimestamps(ctx context.Context) (*dto.AppTimestampsResponse, error) {
	if s.sync == nil {
		return &dto.AppTimestampsResponse{Timestamps: map[string]int64{}}, nil
	}

	profileID := _utils.GetProfileIdWithContext(ctx)

	redisKeys := []string{
		_utils.SyncKeyBdsproEnums,
		_utils.SyncKeyPaymentBankAll,
		s.sync.GetKey(ctx, _utils.SyncKeyPermissionMe, 0),
		s.sync.GetKey(ctx, _utils.SyncKeyProfileInfoMe, 0),
		_utils.SyncKeyUserGuideByKey,
		s.sync.GetKey(ctx, _utils.SyncKeyProductMe, 0),
		_utils.SyncKeyPoiCategoryIcons,
		_utils.SyncKeyTQDProvincesWards,
		s.sync.GetKey(ctx, _utils.SyncKeyContactMe, 0),
		s.sync.GetKey(ctx, _utils.SyncKeyAppointmentMe, 0),
		s.sync.GetKey(ctx, _utils.SyncKeyDealMe, 0),
		s.sync.GetKey(ctx, _utils.SyncKeyNotificationMe, 0),
	}
	names := []string{
		"enum",
		"bank_all",
		"permission_me",
		"info_me",
		"user_guides_by_key",
		"product_me",
		"poi_category_icons",
		"provinces_wards",
		"contact_me",
		"appointment_me",
		"deal_me",
		"notification_me",
	}
	// Keys that are not profile-scoped — safe to EnsureTimestamp without a profile.
	globalKey := map[string]bool{
		_utils.SyncKeyBdsproEnums:       true,
		_utils.SyncKeyPaymentBankAll:    true,
		_utils.SyncKeyUserGuideByKey:    true,
		_utils.SyncKeyPoiCategoryIcons:  true,
		_utils.SyncKeyTQDProvincesWards: true,
	}

	vals, err := s.sync.MGet(ctx, redisKeys)
	if err != nil {
		return nil, _errors.InternalServerException(err.Error())
	}

	timestamps := make(map[string]int64, len(names))
	for i, v := range vals {
		ts := parseSyncTimestamp(v)
		if ts == 0 {
			// MGet can miss while GET still works (e.g. transient pool / nil slot).
			ts = s.sync.GetTimestamp(ctx, redisKeys[i])
		}
		if ts > 0 {
			timestamps[names[i]] = ts
			continue
		}
		// Missing/zero: seed only when key is global, or we have a real profile
		// so we write `{profileId}:0:nt:me` instead of polluting `0:0:nt:me`.
		if !globalKey[redisKeys[i]] && profileID == 0 {
			timestamps[names[i]] = 0
			continue
		}
		timestamps[names[i]] = s.sync.EnsureTimestamp(ctx, redisKeys[i])
	}

	return &dto.AppTimestampsResponse{Timestamps: timestamps}, nil
}

func parseSyncTimestamp(v any) int64 {
	if v == nil {
		return 0
	}
	switch x := v.(type) {
	case string:
		ts, _ := strconv.ParseInt(x, 10, 64)
		return ts
	case []byte:
		ts, _ := strconv.ParseInt(string(x), 10, 64)
		return ts
	case int64:
		return x
	case int:
		return int64(x)
	case uint64:
		return int64(x)
	case float64:
		return int64(x)
	default:
		return 0
	}
}
