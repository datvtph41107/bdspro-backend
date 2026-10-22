package services

import (
	_dto "common/domain/dto"
	"fmt"
	"sort"
	"strings"
	"time"

	"mdn4/modules/crm/dto"
)

type appointmentListMemCacheKey struct {
	profileID uint64
	filterKey string
	page      uint32
	size      uint32
	sort      string
}

type appointmentListMemCacheEntry struct {
	list      []*dto.AppointmentDTO
	total     int64
	expiresAt time.Time
}

const appointmentListMemCacheTTL = 30 * time.Minute

func appointmentListFilterKey(filter map[string]any) string {
	if len(filter) == 0 {
		return ""
	}
	keys := make([]string, 0, len(filter))
	for k := range filter {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	parts := make([]string, 0, len(keys))
	for _, k := range keys {
		parts = append(parts, fmt.Sprintf("%s=%v", k, filter[k]))
	}
	return strings.Join(parts, "&")
}

func (uc *appointmentUsecase) appointmentListCacheKey(profileID uint64, filter map[string]any, pagable _dto.Pagable) appointmentListMemCacheKey {
	return appointmentListMemCacheKey{
		profileID: profileID,
		filterKey: appointmentListFilterKey(filter),
		page:      pagable.Page,
		size:      pagable.Size,
		sort:      pagable.Sort,
	}
}

func (uc *appointmentUsecase) getCachedAppointmentList(profileID uint64, filter map[string]any, pagable _dto.Pagable) ([]*dto.AppointmentDTO, int64, bool) {
	key := uc.appointmentListCacheKey(profileID, filter, pagable)
	now := time.Now()

	uc.listMemCacheMu.RLock()
	entry, ok := uc.listByKey[key]
	uc.listMemCacheMu.RUnlock()
	if !ok || now.After(entry.expiresAt) {
		return nil, 0, false
	}

	out := make([]*dto.AppointmentDTO, len(entry.list))
	copy(out, entry.list)
	return out, entry.total, true
}

func (uc *appointmentUsecase) setCachedAppointmentList(profileID uint64, filter map[string]any, pagable _dto.Pagable, list []*dto.AppointmentDTO, total int64) {
	key := uc.appointmentListCacheKey(profileID, filter, pagable)
	stored := make([]*dto.AppointmentDTO, len(list))
	copy(stored, list)

	uc.listMemCacheMu.Lock()
	uc.listByKey[key] = appointmentListMemCacheEntry{
		list:      stored,
		total:     total,
		expiresAt: time.Now().Add(appointmentListMemCacheTTL),
	}
	uc.listMemCacheMu.Unlock()
}

func (uc *appointmentUsecase) invalidateAppointmentListMemCache(profileID uint64) {
	uc.listMemCacheMu.Lock()
	defer uc.listMemCacheMu.Unlock()

	for k := range uc.listByKey {
		if k.profileID == profileID {
			delete(uc.listByKey, k)
		}
	}
}
