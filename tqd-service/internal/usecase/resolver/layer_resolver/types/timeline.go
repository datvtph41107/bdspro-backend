package types

import "time"

// TimelineEventType — loại sự kiện trên timeline
type TimelineEventType string

const (
	EventLayerCreated     TimelineEventType = "layer_created"
	EventLayerEffective   TimelineEventType = "layer_effective"
	EventLayerModified    TimelineEventType = "layer_modified"
	EventLayerExpired     TimelineEventType = "layer_expired"
	EventLayerReplaced    TimelineEventType = "layer_replaced"
	EventLayerRevoked     TimelineEventType = "layer_revoked"
	EventConflictDetected TimelineEventType = "conflict_detected"
	EventRiskChanged      TimelineEventType = "risk_changed"
)

// TimelineEvent — 1 sự kiện trên timeline
type TimelineEvent struct {
	ID          string
	Date        time.Time
	EventType   TimelineEventType
	LayerID     uint64
	LayerName   string
	Description string
	Impact      *TimelineImpact
	Metadata    map[string]any
}

// TimelineImpact — đánh giá tác động của sự kiện
type TimelineImpact struct {
	RiskDelta         int
	AreaAffectedPct   float64
	NewConflicts      int
	ResolvedConflicts int
	Severity          string // "positive", "negative", "neutral"
}

// ParcelTimeline — toàn bộ timeline của 1 thửa đất
type ParcelTimeline struct {
	ParcelID  uint64
	Events    []*TimelineEvent
	RiskTrend []*RiskSnapshot
	Summary   *TimelineSummary
}

// RiskSnapshot — risk score tại 1 thời điểm
type RiskSnapshot struct {
	Date         time.Time
	RiskScore    uint32
	RiskLevel    string
	NumLayers    int
	NumConflicts int
}

// TimelineSummary — tóm tắt timeline
type TimelineSummary struct {
	TotalEvents      int
	FirstEventAt     time.Time
	LastEventAt      time.Time
	CurrentRiskLevel string
	RiskTrend        string // "improving", "stable", "worsening"
	MajorChanges     []string
}

// TimelineConfig — cấu hình cho Timeline Engine
type TimelineConfig struct {
	CacheTTLSeconds int
	MaxEvents       int
}

// DefaultTimelineConfig — cấu hình timeline mặc định
func DefaultTimelineConfig() *TimelineConfig {
	return &TimelineConfig{
		CacheTTLSeconds: 3600,
		MaxEvents:       100,
	}
}
