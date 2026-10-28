package services

import (
	"context"

	_errors "common/errors"
	_utils "common/utils"

	"mdn4/modules/bdspro/dto"
)

// SyncService owns the app-wide sync timestamp aggregator (GET /v2/bdspro/ts).
type SyncService interface {
	GetAppTimestamps(ctx context.Context) (*dto.AppTimestampsResponse, error)
}

// LocalSyncService ports ProductHandler.GetAppTimestamps onto SyncUtil.
// @bind: mdn4/modules/bdspro/services.SyncService
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
		_utils.SyncKeyTxCostTypes,
		_utils.SyncKeyPropertyTypes,
		s.sync.GetKey(ctx, _utils.SyncKeyPermissionMe, 0),
		s.sync.GetKey(ctx, _utils.SyncKeyProfileInfoMe, 0),
		_utils.SyncKeyUserGuideByKey,
		s.sync.GetKey(ctx, _utils.SyncKeyProductMe, 0),
		s.sync.GetKey(ctx, _utils.SyncKeyListingMe, 0),
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
		"cost_types",
		"property_types",
		"permission_me",
		"info_me",
		"user_guides_by_key",
		"product_me",
		"listing_me",
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
		_utils.SyncKeyTxCostTypes:       true,
		_utils.SyncKeyPropertyTypes:     true,
		_utils.SyncKeyUserGuideByKey:    true,
		_utils.SyncKeyPoiCategoryIcons:  true,
		_utils.SyncKeyTQDProvincesWards: true,
	}

	values, err := s.readyTimestamps(ctx, redisKeys, func(key string) bool {
		// Seed profile-scoped keys only with a real profile, so we write
		// `{profileId}:0:nt:me` instead of polluting `0:0:nt:me`.
		return globalKey[key] || profileID != 0
	})
	if err != nil {
		return nil, err
	}

	timestamps := make(map[string]int64, len(names))
	for i, name := range names {
		timestamps[name] = values[i]
	}

	return &dto.AppTimestampsResponse{Timestamps: timestamps}, nil
}

// readyTimestamps reads every key in one MGET and seeds the missing seedable ones in one batch,
// returning timestamps in key order (0 for keys that are missing and not seedable).
func (s *LocalSyncService) readyTimestamps(ctx context.Context, keys []string, seedable func(key string) bool) ([]int64, error) {
	vals, err := s.sync.MGet(ctx, keys)
	if err != nil {
		return nil, _errors.InternalServerException(err.Error())
	}

	values := make([]int64, len(keys))
	var missing []int
	for i, v := range vals {
		if values[i] = parseSyncTimestamp(v); values[i] == 0 && seedable(keys[i]) {
			missing = append(missing, i)
		}
	}
	if len(missing) == 0 {
		return values, nil
	}

	missingKeys := make([]string, len(missing))
	for j, i := range missing {
		missingKeys[j] = keys[i]
	}
	for j, ts := range s.sync.SeedTimestamps(ctx, missingKeys) {
		values[missing[j]] = ts
	}
	return values, nil
}

func parseSyncTimestamp(v any) int64 {
	return _utils.ParseTimestamp(v)
}
