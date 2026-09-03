package rule

import (
	"context"
	"fmt"
	"sort"
	"time"

	"tqd/internal/usecase/resolver/layer_resolver/types"
)

// TimelineRepository — interface cho Timeline Engine
type TimelineRepository interface {
	GetParcelLayerRowsAt(ctx context.Context, parcelID uint64, asOfDate string) ([]*TimelineRow, error)
}

// TimelineRow — dữ liệu layer tại 1 thời điểm
type TimelineRow struct {
	LayerID          uint64
	LayerName        string
	LayerDisplayName string
	OverlapPct       float64
	OverlapAreaSqm   float64
}

// TimelineEngine — Engine xây dựng timeline biến động layer trên thửa đất
type TimelineEngine struct {
	config       *types.TimelineConfig
	auditRepo    AuditRepository
	timelineRepo TimelineRepository
}

// NewTimelineEngine tạo TimelineEngine
func NewTimelineEngine(config *types.TimelineConfig, auditRepo AuditRepository, timelineRepo TimelineRepository) *TimelineEngine {
	if config == nil {
		config = types.DefaultTimelineConfig()
	}
	return &TimelineEngine{config: config, auditRepo: auditRepo, timelineRepo: timelineRepo}
}

// Name trả về tên engine
func (e *TimelineEngine) Name() string { return "TimelineEngine" }

// BuildTimeline xây dựng timeline cho 1 thửa đất
func (e *TimelineEngine) BuildTimeline(ctx context.Context, parcelID uint64) (*types.ParcelTimeline, error) {
	// B1: Lấy toàn bộ audit entries liên quan đến parcel
	auditEntries, err := e.auditRepo.QueryByParcel(ctx, parcelID)
	if err != nil {
		return nil, fmt.Errorf("query audit by parcel: %w", err)
	}

	// B2: Chuyển audit entries → timeline events
	events := e.auditToTimelineEvents(auditEntries)

	// B3: Sắp xếp theo thời gian
	sort.Slice(events, func(i, j int) bool {
		return events[i].Date.Before(events[j].Date)
	})

	// B4: Giới hạn số lượng events
	if len(events) > e.config.MaxEvents {
		events = events[len(events)-e.config.MaxEvents:]
	}

	// B5: Build summary
	summary := e.buildSummary(events)

	return &types.ParcelTimeline{
		ParcelID: parcelID,
		Events:   events,
		Summary:  summary,
	}, nil
}

// auditToTimelineEvents chuyển audit entries thành timeline events
func (e *TimelineEngine) auditToTimelineEvents(entries []*types.AuditEntry) []*types.TimelineEvent {
	events := make([]*types.TimelineEvent, 0, len(entries))
	for _, entry := range entries {
		events = append(events, &types.TimelineEvent{
			ID:          fmt.Sprintf("evt_%d", entry.ID),
			Date:        entry.Timestamp,
			EventType:   e.mapActionToEventType(entry.Action),
			LayerID:     entry.EntityID,
			Description: fmt.Sprintf("%s %s bởi %s", entry.Action, entry.EntityType, entry.ActorName),
			Metadata: map[string]any{
				"changes": entry.Changes,
				"reason":  entry.Reason,
			},
		})
	}
	return events
}

// mapActionToEventType ánh xạ AuditAction → TimelineEventType
func (e *TimelineEngine) mapActionToEventType(action types.AuditAction) types.TimelineEventType {
	switch action {
	case types.ActionCreated:
		return types.EventLayerCreated
	case types.ActionUpdated:
		return types.EventLayerModified
	case types.ActionDeleted:
		return types.EventLayerRevoked
	case types.ActionReplaced:
		return types.EventLayerReplaced
	case types.ActionStateChanged:
		return types.EventRiskChanged
	default:
		return types.EventLayerModified
	}
}

// buildSummary xây dựng tóm tắt timeline
func (e *TimelineEngine) buildSummary(events []*types.TimelineEvent) *types.TimelineSummary {
	if len(events) == 0 {
		return &types.TimelineSummary{}
	}

	majorChanges := make([]string, 0)
	for _, evt := range events {
		switch evt.EventType {
		case types.EventLayerCreated, types.EventLayerReplaced, types.EventLayerRevoked:
			majorChanges = append(majorChanges, evt.Description)
		}
	}

	return &types.TimelineSummary{
		TotalEvents:  len(events),
		FirstEventAt: events[0].Date,
		LastEventAt:  events[len(events)-1].Date,
		MajorChanges: majorChanges,
	}
}

// Validate kiểm tra config
func (e *TimelineEngine) Validate() error {
	if e.auditRepo == nil {
		return fmt.Errorf("TimelineEngine: auditRepo is nil")
	}
	return nil
}

// Helper: abs
func abs(n int) int {
	if n < 0 {
		return -n
	}
	return n
}

// Ensure time import is used
var _ = time.Now
