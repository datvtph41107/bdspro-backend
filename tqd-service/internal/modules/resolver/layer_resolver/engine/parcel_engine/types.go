package parcel_engine

import (
	"tqd/internal/dto"
	"tqd/internal/modules/resolver/layer_resolver/types"
)

// UnifiedRow kết hợp tất cả field từ ParcelLayerRow và ParcelLayerRowV2
type UnifiedRow struct {
	ParcelID      uint64
	ParcelAreaSqm float64

	// Layer info
	LayerID            uint64
	LayerName          string
	LayerDisplayName   string
	LayerType          uint32
	LayerStatus        uint32
	LayerLegalStatus   uint32
	LayerTrustValue    float32
	LayerEffectiveDate string
	LayerExpiryDate    string
	LayerLegalDoc      string
	LayerAuthority     string

	// Authority issuing
	AuthorityIssuringID   uint64
	AuthorityIssuringName string
	AuthorityIssuringCode string
	AuthorityIssuringDesc string

	// Label info
	LabelID             uint64
	LabelName           string
	LabelDisplayName    string
	LabelColor          string
	LabelGroupCode      uint32
	LabelCanBuild       bool
	LabelPriority       int
	LabelBuildCondition string

	// Group info
	GroupCode     string
	GroupName     string
	GroupColor    string
	GroupCanBuild bool
	GroupPriority int

	// Region info
	RegionID           uint64
	RegionName         string
	RegionDisplayName  string
	RegionLandUseCode  string
	RegionLandUseName  string
	RegionLandUseGroup string
	RegionColorRed     uint8
	RegionColorGreen   uint8
	RegionColorBlue    uint8
	RegionLandUseColor string

	// Overlap
	OverlapAreaSqm float64
	OverlapPct     float64

	// Spatial
	CenterLat float64
	CenterLng float64
	Geometry  []byte
}

// UnifiedResult - kết quả đầu ra duy nhất cho cả 2 view
type UnifiedResult struct {
	ParcelID      uint64
	ParcelAreaSqm float64

	// Thông tin định danh parcel
	ParcelInfo *dto.ParcelInfoResponse

	// Quy hoạch chính (primary sau resolve)
	PrimaryPlanning *PrimaryPlanningInfo

	// Quy hoạch phụ (cho compare mode, build từ ranked[1:])
	SecondaryPlans []*SecondaryPlanInfo

	// Nhóm đất tổng hợp
	LandUseGroups []*LandUseGroupInfo

	// Toàn bộ layers (dữ liệu thô cho V1)
	Layers []*LayerDetail

	// Xung đột quy hoạch
	Conflicts []*ConflictInfo

	// Đánh giá rủi ro
	Risk *RiskInfo

	// Cảnh báo
	Warnings []string

	// Metadata
	Metadata *MetadataInfo
}

type SecondaryPlanInfo struct {
	LayerID        uint64
	LayerName      string
	LandUseCode    string
	LandUseName    string
	LandUseColor   string
	GroupCode      string
	OverlapPercent float64
	OverlapAreaSqm float64
}

type PrimaryPlanningInfo struct {
	Code        string
	Name        string
	Color       string
	GroupCode   string
	GroupName   string
	Ratio       float64
	AreaSqm     float64
	LayerID     uint64
	LayerName   string
	LegalStatus uint32
	CanBuild    bool
}

type LandUseGroupInfo struct {
	Code     string
	Name     string
	Color    string
	CanBuild bool
	Ratio    float64
	AreaSqm  float64
	Details  []*LandUseDetail
}

type LandUseDetail struct {
	LandUseCode string
	LandUseName string
	Ratio       float64
	AreaSqm     float64
}

type LayerDetail struct {
	ID            uint64
	Name          string
	DisplayName   string
	Type          uint32
	Status        uint32
	LegalStatus   uint32
	EffectiveDate string
	ExpiryDate    string
	Authority     *AuthorityInfo
	Indicators    []*IndicatorInfo
	Documents     []*DocumentInfo
	Labels        []*LabelDetail
	TotalArea     float64
	TotalPct      float64
	Geometry      []byte
}

type AuthorityInfo struct {
	ID          uint64
	Name        string
	Code        string
	Description string
}

type IndicatorInfo struct {
	Code  string
	Name  string
	Value string
	Unit  string
}

type DocumentInfo struct {
	ID       uint64
	Name     string
	FileUrl  string
	FileType string
}

type LabelDetail struct {
	ID          uint64
	Code        string
	Name        string
	Color       string
	Ratio       float64
	AreaSqm     float64
	Regions     []*RegionDetail
}

type RegionDetail struct {
	ID             uint64
	Name           string
	DisplayName    string
	LandUseCode    string
	LandUseName    string
	Color          *dto.RGB
	OverlapAreaSqm float64
	OverlapPct     float64
	CenterLat      float64
	CenterLng      float64
	Geometry       []byte
}

type ConflictInfo struct {
	Type            string
	Severity        string
	Between         []uint64
	Description     string
	Recommendation  string
	ActionRequired  string
	AffectedPercent float64
}

type RiskInfo struct {
	Score          uint32
	Level          string
	CanBuild       bool
	Reasons        []string
	Recommendation string
}

type MetadataInfo struct {
	DataSource     string
	QueriedAt      string
	ResponseTimeMs int64
	Cached         bool
}

// RiskThreshold - ngưỡng phân loại rủi ro (data-driven)
type RiskThreshold struct {
	MinScore       float64
	Level          string
	CanBuild       bool
	Recommendation string
}

// LandUseRiskWeight - trọng số rủi ro cho từng loại đất
type LandUseRiskWeight struct {
	Multiplier float64 // Hệ số nhân rủi ro (VD: DGT=2.0, ODT=0.35)
	Label      string  // Mô tả tiếng Việt (VD: "không được phép xây dựng")
}

// EngineConfig - cấu hình cho engine
type EngineConfig struct {
	LegalWeight          float64
	PriorityWeight       float64
	AreaWeight           float64
	LegalScores          map[uint32]float64
	ConflictMatrix       map[string]map[string]bool
	DefaultGroupPriority map[string]int

	LandUseRiskWeights map[string]*LandUseRiskWeight
	DefaultRiskWeight  *LandUseRiskWeight
	MultiLayerPenalty  float64
	RiskThresholds     []*RiskThreshold
	DefaultPriority    int

	// Conflict labels
	ConflictType            string
	ConflictSeverity        string
	ConflictRecommendation  string
	ConflictActionRequired  string

	// Warning messages
	WarningOverlapExceeded string
	WarningMultiLayer      string
	WarningConflict        string

	// Assess messages
	NoPlanningImpactMsg     string
	NoSignificantImpactMsg  string

	// Land type display names (key: landTypeId)
	LandTypeNames map[uint64]string

	// Build status labels (for old AssessEngine compatibility)
	StatusProhibited   string
	StatusConditional  string
	StatusAllowed      string
	ConflictStatusSuffix string

	// Conflict summary suffix
	ConflictSummarySuffix string
}

func DefaultEngineConfig() *EngineConfig {
	riskWeights := make(map[string]*LandUseRiskWeight, len(types.DefaultLandUseRiskWeights))
	for k, v := range types.DefaultLandUseRiskWeights {
		riskWeights[k] = &LandUseRiskWeight{Multiplier: v.Multiplier, Label: v.Label}
	}
	thresholds := make([]*RiskThreshold, len(types.DefaultRiskThresholds))
	for i, t := range types.DefaultRiskThresholds {
		thresholds[i] = &RiskThreshold{MinScore: t.MinScore, Level: t.Level, CanBuild: t.CanBuild, Recommendation: t.Recommendation}
	}
	landTypes := make(map[uint64]string, len(types.DefaultLandTypeNames))
	for k, v := range types.DefaultLandTypeNames {
		landTypes[k] = v
	}

	return &EngineConfig{
		LegalWeight:          types.DefaultLegalWeight,
		PriorityWeight:       types.DefaultPriorityWeight,
		AreaWeight:           types.DefaultAreaWeight,
		LegalScores:          types.DefaultLegalScores,
		ConflictMatrix:       make(map[string]map[string]bool),
		DefaultGroupPriority: make(map[string]int),
		LandUseRiskWeights:   riskWeights,
		DefaultRiskWeight:    &LandUseRiskWeight{Multiplier: types.DefaultRiskWeightValue.Multiplier, Label: types.DefaultRiskWeightValue.Label},
		MultiLayerPenalty:    types.DefaultMultiLayerPenaltyValue,
		RiskThresholds:       thresholds,
		DefaultPriority:      types.DefaultPriorityValue,
		ConflictType:            types.DefaultConflictType,
		ConflictSeverity:        types.DefaultConflictSeverity,
		ConflictRecommendation:  types.DefaultConflictRecommendation,
		ConflictActionRequired:  types.DefaultConflictActionRequired,
		WarningOverlapExceeded:  types.DefaultWarningOverlapExceeded,
		WarningMultiLayer:       types.DefaultWarningMultiLayer,
		WarningConflict:         types.DefaultWarningConflict,
		NoPlanningImpactMsg:     types.DefaultNoPlanningImpactMsg,
		NoSignificantImpactMsg:  types.DefaultNoSignificantImpactMsg,
		LandTypeNames:           landTypes,
		StatusProhibited:        types.DefaultStatusProhibited,
		StatusConditional:       types.DefaultStatusConditional,
		StatusAllowed:           types.DefaultStatusAllowed,
		ConflictStatusSuffix:    types.DefaultConflictStatusSuffix,
		ConflictSummarySuffix:   types.DefaultConflictSummarySuffix,
	}
}

// Internal helpers
type layerGroupAccum struct {
	code     string
	name     string
	color    string
	canBuild bool
	ratio    float64
	areaSqm  float64
	details  map[string]*LandUseDetail
}

// Helpers to convert from old types
func ToUnifiedRowsFromV1(rows []dto.ParcelLayerRow) []UnifiedRow {
	result := make([]UnifiedRow, len(rows))
	for i, r := range rows {
		result[i] = UnifiedRow{
			ParcelID:              r.ParcelID,
			ParcelAreaSqm:         r.ParcelAreaSqm,
			LayerID:               r.LayerID,
			LayerName:             r.LayerName,
			LayerDisplayName:      r.LayerDisplayName,
			LayerType:             r.LayerType,
			LayerEffectiveDate:    r.LayerEffectiveDate,
			LayerLegalDoc:         r.LayerLegalDoc,
			AuthorityIssuringID:   r.AuthorityIssuringID,
			AuthorityIssuringName: r.AuthorityIssuringName,
			AuthorityIssuringCode: r.AuthorityIssuringCode,
			AuthorityIssuringDesc: r.AuthorityIssuringDesc,
			RegionID:              r.RegionID,
			RegionName:            r.RegionName,
			RegionDisplayName:     r.RegionDisplayName,
			RegionLandUseCode:     r.RegionLandUseCode,
			RegionLandUseName:     r.RegionLandUseName,
			RegionLandUseGroup:    r.RegionLandUseGroup,
			RegionColorRed:        r.RegionColorRed,
			RegionColorGreen:      r.RegionColorGreen,
			RegionColorBlue:       r.RegionColorBlue,
			RegionLandUseColor:    r.RegionLandUseColor,
			OverlapAreaSqm:        r.OverlapAreaSqm,
			OverlapPct:            r.OverlapPct,
			CenterLat:             r.CenterLat,
			CenterLng:             r.CenterLng,
			Geometry:              r.Geometry,
		}
	}
	return result
}

func ToUnifiedRowsFromV2(rows []dto.ParcelLayerRowV2) []UnifiedRow {
	result := make([]UnifiedRow, len(rows))
	for i, r := range rows {
		result[i] = UnifiedRow{
			ParcelID:              r.ParcelID,
			ParcelAreaSqm:         r.ParcelAreaSqm,
			LayerID:               r.LayerID,
			LayerName:             r.LayerName,
			LayerDisplayName:      r.LayerDisplayName,
			LayerType:             r.LayerType,
			LayerStatus:           r.LayerStatus,
			LayerLegalStatus:      r.LayerLegalStatus,
			LayerTrustValue:       r.LayerTrustValue,
			LayerEffectiveDate:    r.LayerEffectiveDate,
			LayerExpiryDate:       r.LayerExpiryDate,
			LayerAuthority:        r.LayerAuthority,
			LabelID:               r.LabelID,
			LabelName:             r.LabelName,
			LabelDisplayName:      r.LabelDisplayName,
			LabelColor:            r.LabelColor,
			LabelGroupCode:        r.LabelGroupCode,
			LabelCanBuild:         r.LabelCanBuild,
			LabelPriority:         r.LabelPriority,
			LabelBuildCondition:   r.LabelBuildCondition,
			GroupCode:             r.GroupCode,
			GroupName:             r.GroupName,
			GroupColor:            r.GroupColor,
			GroupCanBuild:         r.GroupCanBuild,
			GroupPriority:         r.GroupPriority,
			RegionID:              r.RegionID,
			RegionName:            r.RegionName,
			RegionDisplayName:     r.RegionDisplayName,
			RegionLandUseCode:     r.RegionLandUseCode,
			RegionLandUseName:     r.RegionLandUseName,
			RegionColorRed:        r.RegionColorRed,
			RegionColorGreen:      r.RegionColorGreen,
			RegionColorBlue:       r.RegionColorBlue,
			OverlapAreaSqm:        r.OverlapAreaSqm,
			OverlapPct:            r.OverlapPct,
			CenterLat:             r.CenterLat,
			CenterLng:             r.CenterLng,
			Geometry:              r.Geometry,
		}
	}
	return result
}
