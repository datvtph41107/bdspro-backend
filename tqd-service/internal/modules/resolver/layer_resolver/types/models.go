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
}
