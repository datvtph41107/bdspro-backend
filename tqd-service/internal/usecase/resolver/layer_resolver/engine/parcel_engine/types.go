package parcel_engine

import (
	"tqd/internal/dto"
	"tqd/internal/usecase/resolver/layer_resolver/types"
)

// ============================================================
// UnifiedRow — nguồn đầu vào duy nhất cho engine.Process()
//
// Naming convention thống nhất:
//   - Authority*  : thông tin cơ quan ban hành
//   - LandUse*   : từ qh_land_use (PRIMARY, align QuickLayers)
//   - Group*     : từ land_use_groups (secondary, display only)
//   - Region*    : thông tin region địa lý
//   - Overlap*   : kết quả tính toán diện tích giao
//
// ============================================================
type UnifiedRow struct {
	// ── Parcel ──────────────────────────────────────────────
	ParcelID      uint64
	ParcelAreaSqm float64

	// ── Layer ───────────────────────────────────────────────
	LayerID            uint64
	LayerName          string
	LayerDisplayName   string
	LayerType          uint32
	LayerStatus        uint32
	LayerLegalStatus   uint32
	LayerTrustValue    float32
	LayerEffectiveDate string
	LayerExpiryDate    string
	LayerAvatar        string // ✅ THÊM – đồng bộ với proto Layer.layerAvatar

	// ── Authority ───────────────────────────────────────────
	// Đặt tên AuthorityID/Name/Code/Desc (không phải AuthorityIssuringXxx)
	// để nhất quán với tất cả consumer (engine, builder, converter)
	AuthorityID   uint64
	AuthorityName string
	AuthorityCode string
	AuthorityDesc string

	// ── Label ───────────────────────────────────────────────
	LabelID          uint64
	LabelName        string
	LabelDisplayName string
	LabelColor       string

	// ── Land Use (PRIMARY: qh_land_use via r.land_use_id) ──
	// Đây là nguồn sự thật duy nhất, align với GetParcelQuickInfo.
	// KHÔNG dùng luc.name / lg.name làm primary.
	LandUseCode  string
	LandUseName  string
	LandUseColor string // lu.color — hex string (e.g. "#FF5733")
	CanBuild     bool   // lu.can_build (COALESCE với luc, lg ở SQL)
	WarnLevel    int32  // lu.warn_level
	Priority     int    // COALESCE(lu.priority, luc.priority, lg.priority, 50)

	// ── Land Use Group (secondary, for display grouping) ────
	GroupCode  string
	GroupName  string
	GroupColor string

	// Build assessment
	BuildCondition string

	// ── Region ──────────────────────────────────────────────
	RegionID          uint64
	RegionName        string
	RegionDisplayName string

	// ── Overlap ─────────────────────────────────────────────
	OverlapAreaSqm float64
	OverlapPct     float64

	// ── Spatial ─────────────────────────────────────────────
	CenterLat float64
	CenterLng float64
	Geometry  []byte
}

// ============================================================
// UnifiedResult — kết quả đầu ra của engine.Process()
// ============================================================
type UnifiedResult struct {
	ParcelID      uint64
	ParcelAreaSqm float64

	// Thông tin định danh parcel
	ParcelInfo *dto.ParcelInfoResponse

	// Quy hoạch chính (primary sau resolve)
	PrimaryPlanning *PrimaryPlanningInfo

	// Quy hoạch phụ (cho compare/detail mode)
	SecondaryPlans []*SecondaryPlanInfo

	// Nhóm đất tổng hợp (cho overview)
	LandUseGroups []*LandUseGroupInfo

	// Toàn bộ layers (dữ liệu đầy đủ cho detail mode)
	Layers []*LayerDetail

	// Xung đột quy hoạch
	Conflicts []*ConflictInfo

	// Đánh giá rủi ro
	Risk *RiskInfo

	// Cảnh báo
	Warnings []string

	// Compare result (mode=compare)
	CompareResult *CompareResultInfo

	// Historical risk trend (mode=detail)
	HistoricalRisk *HistoricalRiskInfo

	// Resolution status
	ResolutionStatus string

	// Metadata
	Metadata *MetadataInfo
}

// ── Planning info ──────────────────────────────────────────

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
	// Normalized legal
	LegalCanonical   string
	LegalConfidence  string
	LegalActiveState string
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

// ── Land Use Group ─────────────────────────────────────────

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
	WarnLevel   int // lu.warn_level — align với QuickLayers
}

// ── Layer hierarchy ────────────────────────────────────────

type LayerDetail struct {
	ID              uint64
	Name            string
	DisplayName     string
	Type            uint32
	Status          uint32
	LegalStatus     uint32
	LegalStatusName string // ✅ THÊM – tên hiển thị của legal status (ví dụ: "Đã ban hành", "Hết hiệu lực")
	LayerAvatar     string // ✅ THÊM – avatar của layer (đồng bộ với proto Layer.layerAvatar)
	EffectiveDate   string
	ExpiryDate      string
	Authority       *AuthorityInfo
	Documents       []*DocumentInfo
	Labels          []*LabelDetail
	TotalArea       float64
	TotalPct        float64
}

type AuthorityInfo struct {
	ID          uint64
	Name        string
	Code        string
	Description string
}

type DocumentInfo struct {
	ID       uint64
	Name     string
	FileUrl  string
	FileType string
}

type LabelDetail struct {
	ID      uint64
	Code    string
	Name    string
	Color   string
	Ratio   float64
	AreaSqm float64
	Regions []*RegionDetail
}

type RegionDetail struct {
	ID             uint64
	Name           string
	DisplayName    string
	LandUseCode    string
	LandUseName    string
	LandUseGroup   string
	LandUseColor   string
	OverlapAreaSqm float64
	OverlapPct     float64
	CenterLat      float64
	CenterLng      float64
	Geometry       []byte
	WarnLevel      int32
	CanBuild       bool
}

// ── Risk & Conflict ────────────────────────────────────────

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

// ── Compare & Historical ───────────────────────────────────

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

// ── Metadata ───────────────────────────────────────────────

type MetadataInfo struct {
	DataSource     string
	QueriedAt      string
	ResponseTimeMs int64
	Cached         bool
}

// ── Engine config types ────────────────────────────────────

type RiskThreshold struct {
	MinScore       float64
	Level          string
	CanBuild       bool
	Recommendation string
}

type LandUseRiskWeight struct {
	Multiplier float64
	Label      string
}

type EngineConfig struct {
	LegalWeight    float64
	PriorityWeight float64
	AreaWeight     float64

	LegalScores          map[uint32]float64
	ConflictMatrix       map[string]map[string]bool
	DefaultGroupPriority map[string]int

	LandUseRiskWeights map[string]*LandUseRiskWeight
	DefaultRiskWeight  *LandUseRiskWeight
	MultiLayerPenalty  float64
	RiskThresholds     []*RiskThreshold
	DefaultPriority    int

	// Sub-engine configs
	LegalNormalization     *types.LegalNormalizationConfig
	ConflictClassification *types.ConflictClassificationConfig
	SemanticCanonicalMap   map[string]string
	SemanticLabels         map[string]string
	LifecycleConfig        *types.LifecycleConfig
	AuditConfig            *types.AuditConfig
	TimelineConfig         *types.TimelineConfig
	CompareConfig          *types.CompareConfig
	HistoricalRiskConfig   *types.HistoricalRiskConfig
	PipelineConfig         *types.PipelineConfig

	// Warning/status messages
	WarningOverlapExceeded string
	WarningMultiLayer      string
	WarningConflict        string
	NoPlanningImpactMsg    string
	NoSignificantImpactMsg string

	// Conflict default labels
	ConflictType           string
	ConflictSeverity       string
	ConflictRecommendation string
	ConflictActionRequired string

	// Build status labels
	StatusProhibited      string
	StatusConditional     string
	StatusAllowed         string
	ConflictStatusSuffix  string
	ConflictSummarySuffix string

	// Land type display names
	LandTypeNames map[uint64]string
}

func DefaultEngineConfig() *EngineConfig {
	riskWeights := make(map[string]*LandUseRiskWeight, len(types.DefaultLandUseRiskWeights))
	for k, v := range types.DefaultLandUseRiskWeights {
		riskWeights[k] = &LandUseRiskWeight{Multiplier: v.Multiplier, Label: v.Label}
	}
	thresholds := make([]*RiskThreshold, len(types.DefaultRiskThresholds))
	for i, t := range types.DefaultRiskThresholds {
		thresholds[i] = &RiskThreshold{
			MinScore:       t.MinScore,
			Level:          t.Level,
			CanBuild:       t.CanBuild,
			Recommendation: t.Recommendation,
		}
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
		DefaultRiskWeight: &LandUseRiskWeight{
			Multiplier: types.DefaultRiskWeightValue.Multiplier,
			Label:      types.DefaultRiskWeightValue.Label,
		},
		MultiLayerPenalty: types.DefaultMultiLayerPenaltyValue,
		RiskThresholds:    thresholds,
		DefaultPriority:   types.DefaultPriorityValue,
		LegalNormalization: &types.LegalNormalizationConfig{
			CanonicalMap:     types.DefaultCanonicalMap(),
			ConfidenceWeight: types.DefaultConfidenceWeight(),
			ActiveRules:      types.DefaultActiveRules(),
		},
		ConflictClassification: &types.ConflictClassificationConfig{
			SeverityThresholds: types.DefaultSeverityThresholds(),
			TypeWeights:        types.DefaultConflictTypeWeights(),
			GuidanceTemplates:  types.DefaultGuidanceTemplates(),
		},
		SemanticCanonicalMap: make(map[string]string),
		SemanticLabels:       make(map[string]string),
		LifecycleConfig: &types.LifecycleConfig{
			AllowedTransitions:   types.DefaultAllowedTransitions(),
			AutoIncrementVersion: true,
		},
		AuditConfig:            types.DefaultAuditConfig(),
		TimelineConfig:         types.DefaultTimelineConfig(),
		CompareConfig:          types.DefaultCompareConfig(),
		HistoricalRiskConfig:   types.DefaultHistoricalRiskConfig(),
		PipelineConfig:         types.DefaultPipelineConfig(),
		ConflictType:           types.DefaultConflictType,
		ConflictSeverity:       types.DefaultConflictSeverity,
		ConflictRecommendation: types.DefaultConflictRecommendation,
		ConflictActionRequired: types.DefaultConflictActionRequired,
		WarningOverlapExceeded: types.DefaultWarningOverlapExceeded,
		WarningMultiLayer:      types.DefaultWarningMultiLayer,
		WarningConflict:        types.DefaultWarningConflict,
		NoPlanningImpactMsg:    types.DefaultNoPlanningImpactMsg,
		NoSignificantImpactMsg: types.DefaultNoSignificantImpactMsg,
		LandTypeNames:          landTypes,
		StatusProhibited:       types.DefaultStatusProhibited,
		StatusConditional:      types.DefaultStatusConditional,
		StatusAllowed:          types.DefaultStatusAllowed,
		ConflictStatusSuffix:   types.DefaultConflictStatusSuffix,
		ConflictSummarySuffix:  types.DefaultConflictSummarySuffix,
	}
}

// ── Internal accumulator (private) ────────────────────────

type layerGroupAccum struct {
	code     string
	name     string
	color    string
	canBuild bool
	ratio    float64
	areaSqm  float64
	details  map[string]*LandUseDetail
}

// ── Row converters ─────────────────────────────────────────
//
// ToUnifiedRowsFromV2 chuyển ParcelLayerRowV2 → []UnifiedRow.
// Field mapping 1-1, không có fallback logic — đã được resolve đúng ở SQL.

func ToUnifiedRowsFromV2(rows []dto.ParcelLayerRowV2) []UnifiedRow {
	result := make([]UnifiedRow, 0, len(rows))
	for _, r := range rows {
		result = append(result, UnifiedRow{
			ParcelID:      r.ParcelID,
			ParcelAreaSqm: r.ParcelAreaSqm,

			LayerID:            r.LayerID,
			LayerName:          r.LayerName,
			LayerDisplayName:   r.LayerDisplayName,
			LayerType:          r.LayerType,
			LayerStatus:        r.LayerStatus,
			LayerLegalStatus:   r.LayerLegalStatus,
			LayerTrustValue:    r.LayerTrustValue,
			LayerEffectiveDate: r.LayerEffectiveDate,
			LayerExpiryDate:    r.LayerExpiryDate,
			LayerAvatar:        r.LayerAvatar, // ✅ THÊM

			// Authority: ParcelLayerRowV2 dùng AuthorityIssuringXxx
			AuthorityID:   r.AuthorityIssuringID,
			AuthorityName: r.AuthorityIssuringName,
			AuthorityCode: r.AuthorityIssuringCode,
			AuthorityDesc: r.AuthorityIssuringDesc,

			LabelID:          r.LabelID,
			LabelName:        r.LabelName,
			LabelDisplayName: r.LabelDisplayName,
			LabelColor:       r.LabelColor,

			// PRIMARY from qh_land_use (via GetParcelLayerRowsV2Unified SQL)
			LandUseCode:  r.RegionLandUseCode,  // lu.code
			LandUseName:  r.RegionLandUseName,  // lu.name
			LandUseColor: r.RegionLandUseColor, // COALESCE(lu.color, label_color, group_color)
			CanBuild:     r.LabelCanBuild,      // COALESCE(lu.can_build, luc, lg) from SQL
			WarnLevel:    r.WarnLevel,          // lu.warn_level
			Priority:     r.LabelPriority,      // COALESCE(lu.priority, luc, lg) from SQL

			// Group (secondary)
			GroupCode:      r.GroupCode,
			GroupName:      r.GroupName,
			GroupColor:     r.GroupColor,
			BuildCondition: r.LabelBuildCondition,

			RegionID:          r.RegionID,
			RegionName:        r.RegionName,
			RegionDisplayName: r.RegionDisplayName,

			OverlapAreaSqm: r.OverlapAreaSqm,
			OverlapPct:     r.OverlapPct,

			CenterLat: r.CenterLat,
			CenterLng: r.CenterLng,
			Geometry:  r.Geometry,
		})
	}
	return result
}

// ToUnifiedRowsFromV1 — legacy compat, giữ nguyên cho các caller cũ
func ToUnifiedRowsFromV1(rows []dto.ParcelLayerRow) []UnifiedRow {
	result := make([]UnifiedRow, 0, len(rows))
	for _, r := range rows {
		result = append(result, UnifiedRow{
			ParcelID:      r.ParcelID,
			ParcelAreaSqm: r.ParcelAreaSqm,

			LayerID:            r.LayerID,
			LayerName:          r.LayerName,
			LayerDisplayName:   r.LayerDisplayName,
			LayerType:          r.LayerType,
			LayerEffectiveDate: r.LayerEffectiveDate,

			AuthorityID:   r.AuthorityIssuringID,
			AuthorityName: r.AuthorityIssuringName,
			AuthorityCode: r.AuthorityIssuringCode,
			AuthorityDesc: r.AuthorityIssuringDesc,

			// V1: land use info từ label (legacy, không có lu trực tiếp)
			LandUseCode:  r.RegionLandUseCode,
			LandUseName:  r.RegionLandUseName,
			LandUseColor: r.RegionLandUseColor,

			RegionID:          r.RegionID,
			RegionName:        r.RegionName,
			RegionDisplayName: r.RegionDisplayName,

			OverlapAreaSqm: r.OverlapAreaSqm,
			OverlapPct:     r.OverlapPct,
			CenterLat:      r.CenterLat,
			CenterLng:      r.CenterLng,
			Geometry:       r.Geometry,
		})
	}
	return result
}
