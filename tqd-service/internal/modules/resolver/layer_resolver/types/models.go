package types

import "time"

type LayerCandidate struct {
	LayerID          uint64
	LayerName        string
	LayerDisplayName string
	LayerType        uint32
	LayerStatus      uint32

	LegalStatus   uint32
	EffectiveDate *time.Time
	ExpiryDate    *time.Time
	TrustValue    float32
	Authority     string
	LegalDocRef   string

	LabelID          uint64
	LabelName        string
	LabelDisplayName string
	LabelColor       string

	LandUseCode    string
	LandUseName    string
	LandUseColor   string
	GroupCode      string
	GroupName      string
	GroupColor     string
	CanBuild       bool
	Priority       int
	BuildCondition string

	OverlapPercent float64
	OverlapAreaSqm float64
	ParcelAreaSqm  float64

	CenterLat float64
	CenterLng float64
	Geometry  []byte

	SourceType string
}

type EvaluatedLayer struct {
	*LayerCandidate
	LegalStatusCanonical string
	LegalWeight          float64
}

type RankedLayer struct {
	*EvaluatedLayer
	Rank  int
	Score float64
}

type ConflictDetail struct {
	Type            string
	Severity        ConflictSeverity
	Between         []uint64
	Description     string
	Recommendation  string
	ActionRequired  string
	AffectedPercent float64
}

type LandUseGroupDetail struct {
	Code           string
	Name           string
	Color          string
	CanBuild       bool
	Priority       int
	LandUseCode    string
	LandUseName    string
	LandUseColor   string
	OverlapPercent float64
	OverlapAreaSqm float64
}

type ResolveResult struct {
	Primary            *RankedLayer
	Secondary          []*RankedLayer
	HasConflict        bool
	CriticalConflicts  []*ConflictDetail
	InfoConflicts      []*ConflictDetail
	CanBuild           bool
	BuildStatus        string
	Reason             string
	Mode               ResolveMode
	ProcessingTimeMs   int64
	LayerLandUseGroups map[uint64][]*LandUseGroupDetail

	ParcelAreaSqm  float64
	RiskScore      uint32
	RiskLevel      string
	RiskReasons    []string
	Recommendation string

	ResolutionStatus string

	Layers []*LayerDetail

	CompareResult  *CompareResultInfo
	HistoricalRisk *HistoricalRiskInfo
}

type CompareResultInfo struct {
	TotalLayers int
	Changes     []*CompareChangeInfo
}

type CompareChangeInfo struct {
	LayerID    uint64
	LayerName  string
	ChangeType string
	OldValue   string
	NewValue   string
}

type HistoricalRiskInfo struct {
	Trend          string
	TrendDetail    string
	ActiveLayers   int
	ExpiringLayers int
	AnalyzedAt     string
}

// ─── Layer/Region detail types (dùng cho assessment mode=detail, map → proto Layer) ───

type LayerDetail struct {
	ID            uint64
	Name          string
	DisplayName   string
	Type          uint32
	EffectiveDate string
	Labels        []*LabelDetail
	LegalDocs     []*LegalDocDetail
	Authority     *AuthorityDetail
	TotalArea     float64
	TotalPct      float64
}

type LabelDetail struct {
	ID          uint64
	Name        string
	DisplayName string
	Color       string
	Regions     []*RegionDetail
	TotalArea   float64
	TotalPct    float64
}

type RegionDetail struct {
	ID             uint64
	Name           string
	DisplayName    string
	LandUseCode    string
	LandUseName    string
	ColorRed       uint8
	ColorGreen     uint8
	ColorBlue      uint8
	OverlapAreaSqm float64
	OverlapPct     float64
	Geometry       []byte
}

type LegalDocDetail struct {
	ID       uint64
	Name     string
	FileUrl  string
	FileType string
}

type AuthorityDetail struct {
	ID          uint64
	Name        string
	Code        string
	Description string
}
