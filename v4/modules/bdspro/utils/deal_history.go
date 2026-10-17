package utils

import (
	"context"
	"fmt"
	"log"
	"time"

	_utils "common/utils"
	"userv4/internal/interface/providers"
	"userv4/modules/bdspro/enums"
	notifentity "userv4/modules/notification/entity"
	notifrepo "userv4/modules/notification/repository"

	"gorm.io/datatypes"
)

// DealHistoryLogger writes deal_history rows and bumps deal:his sync (async, non-blocking).
type DealHistoryLogger struct {
	history  notifrepo.DealHistoryRepository
	sync     *_utils.SyncUtil
	profiles providers.ProfileProvider
}

func NewDealHistoryLogger(
	history notifrepo.DealHistoryRepository,
	sync *_utils.SyncUtil,
	profiles providers.ProfileProvider,
) *DealHistoryLogger {
	return &DealHistoryLogger{history: history, sync: sync, profiles: profiles}
}

// ListByDeal returns paginated deal_history rows for a deal (newest first).
func (l *DealHistoryLogger) ListByDeal(ctx context.Context, dealID uint64, page, size uint32) ([]*notifentity.DealHistory, int64, error) {
	if l == nil || l.history == nil || dealID == 0 {
		return nil, 0, nil
	}
	return l.history.ListByDealID(ctx, dealID, page, size)
}

// LogDealHistory records a deal timeline event (parity with bdspro EventHistoryUsecase.LogDealHistory).
func (l *DealHistoryLogger) LogDealHistory(
	ctx context.Context,
	dealID uint64,
	eventType enums.DealHistoryEventType,
	content string,
	metadata map[string]interface{},
) {
	if l == nil || l.history == nil || dealID == 0 {
		return
	}

	actorID := _utils.GetProfileIdWithContext(ctx)
	actionName := eventType.GetEventName()

	go func() {
		cloneCtx := _utils.CloneContext(ctx)
		row := &notifentity.DealHistory{
			DealID:     dealID,
			ActorID:    actorID,
			ActionType: actionName,
			ActionName: actionName,
			Content:    content,
			Metadata:   toJSONMap(metadata),
		}
		if l.profiles != nil && actorID > 0 {
			if p, err := l.profiles.GetByProfileID(cloneCtx, actorID); err == nil && p != nil {
				row.ActorName = p.FullName
				row.ActorAvatar = p.Avatar
			}
		}
		if err := l.history.Create(cloneCtx, row); err != nil {
			log.Printf("deal history: create failed deal=%d event=%s: %v", dealID, actionName, err)
			return
		}
		if l.sync != nil {
			now := time.Now().UnixMilli()
			key := l.sync.GetKey(cloneCtx, _utils.SyncKeyDealHistory, dealID)
			_ = l.sync.PutTimestamp(cloneCtx, key, now)
		}
	}()
}

func toJSONMap(m map[string]interface{}) datatypes.JSONMap {
	if len(m) == 0 {
		return datatypes.JSONMap{}
	}
	out := make(datatypes.JSONMap, len(m))
	for k, v := range m {
		out[k] = v
	}
	return out
}

// FormatDealHistoryContent is a small helper for consistent log messages.
func FormatDealHistoryContent(format string, args ...any) string {
	return fmt.Sprintf(format, args...)
}

// EventTypeFromActionName maps stored action_type / action_name back to DealHistoryEventType.
func EventTypeFromActionName(name string) uint32 {
	for k, v := range enums.DealHistoryEventNameMap {
		if v == name {
			return uint32(k)
		}
	}
	return 0
}
