package types

import "time"

// LayerCandidate - Đầu vào resolver (đã được chuẩn hóa)
type LayerCandidate struct {
	// Layer info
	LayerID          uint64
	LayerName        string
	LayerDisplayName string
	LayerType        uint32
	LayerStatus      uint32

	// Legal info
	LegalStatus   uint32
	EffectiveDate *time.Time
	ExpiryDate    *time.Time
	TrustValue    float32
	Authority     string

	// Label info
	LabelID          uint64
	LabelName        string
	LabelDisplayName string
	LabelColor       string

	// Land use info
	LandUseCode    string
	LandUseName    string
	LandUseColor   string
	GroupCode      string
	CanBuild       bool
	Priority       int
	BuildCondition string

	// Impact data
	OverlapPercent float64
	OverlapAreaSqm float64
	ParcelAreaSqm  float64

	// Spatial
	CenterLat float64
	CenterLng float64
	Geometry  []byte

	SourceType string
}

// EvaluatedLayer - Sau bước đánh giá pháp lý
type EvaluatedLayer struct {
	*LayerCandidate
	LegalStatusCanonical string
	LegalWeight          float64
}

// RankedLayer - Sau bước xếp hạng
type RankedLayer struct {
	*EvaluatedLayer
	Rank  int
	Score float64
}

// ConflictDetail - Chi tiết xung đột
type ConflictDetail struct {
	Type            string
	Severity        ConflictSeverity
	Between         []uint64
	Description     string
	Recommendation  string
	ActionRequired  string
	AffectedPercent float64
}

// ResolveResult - Kết quả cuối cùng
type ResolveResult struct {
	Primary           *RankedLayer
	Secondary         []*RankedLayer
	HasConflict       bool
	CriticalConflicts []*ConflictDetail
	InfoConflicts     []*ConflictDetail
	CanBuild          bool
	BuildStatus       string
	Reason            string
	Mode              ResolveMode
	ProcessingTimeMs  int64
}
