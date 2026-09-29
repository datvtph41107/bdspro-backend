package dto

import (
	"time"

	_enum "common/domain/enum"
	qh_domain "tqd/internal/domain/qh"
)

//	type ParcelDTO struct {
//		ID       uint64  `json:"id"`
//		Name     string  `json:"name"`
//		Address  string  `json:"address"`
//		Category uint64  `json:"categoryId"`
//		Latitude  float64 `json:"latitude"`
//		Longitude float64 `json:"longitude"`
//	}
type ParcelList struct {
	Data []qh_domain.Parcel `json:"data"`
}

type ParcelResponse struct {
	ID         uint64  `json:"id"`
	Name       string  `json:"name"`
	Address    string  `json:"address"`
	CategoryID uint64  `json:"categoryId"`
	Latitude   float64 `json:"latitude"`
	Longitude  float64 `json:"longitude"`
}

// parcel_dto.go
type ParcelLayerInfo struct {
	ID       uint64 `gorm:"column:id"`
	ParcelID uint64 `gorm:"column:parcel_id"`

	LayerID          uint64    `gorm:"column:layer_id"`
	LayerName        string    `gorm:"column:layer_name"`
	LayerDisplayName string    `gorm:"column:layer_display_name"`
	LayerAvatar      string    `gorm:"column:layer_avatar"`
	LayerLegalStatus int       `gorm:"column:layer_legal_status"`
	LayerType        uint32    `gorm:"column:layer_type"`
	LayerUpdatedAt   time.Time `gorm:"column:layer_updated_at"`

	PlanningProjectID *uint64 `gorm:"column:planning_project_id"`

	CanBuild bool `gorm:"column:can_build"`

	LayerOrder       int32
	LandUseID        uint64
	LandUseName      string
	LandUseColor     string
	RelationType     string
	WarnLevel        _enum.EWarningLevel `gorm:"column:warn_level"`
	IntersectAreaSqm float64             `gorm:"column:intersect_area_sqm"`
	Percent          float64             `gorm:"column:percent"`
}

type PolygonMapTargetsResponse struct {
	ResultType uint32
	Parcels    []qh_domain.Parcel
	Regions    []RegionInfoResponse
	Total      int64
	Z          uint32

	Page     uint32
	PageSize uint32

	AdjustedGeoJSON   string
	AdjustmentReason  string
	OriginalAreaKm2   float64
	AdjustedAreaKm2   float64
	MaxAllowedAreaKm2 float64
	Warning           string

	Analysis       *PolygonAnalysisSummary
	QuickLayerRows []ParcelLayerInfo
}

type PolygonTargetFilter struct {
	LayerIDs     []uint64
	LabelIDs     []uint64
	LandUseIDs   []uint64
	LandUseCodes []string
	CanBuild     *bool
	Keyword      string
	Zoom         *uint32
}

func (f *PolygonTargetFilter) HasPlanningFilters() bool {
	if f == nil {
		return false
	}

	return len(f.LayerIDs) > 0 ||
		len(f.LabelIDs) > 0 ||
		len(f.LandUseIDs) > 0 ||
		len(f.LandUseCodes) > 0 ||
		f.CanBuild != nil ||
		f.Keyword != ""
}

type PolygonLandUseStat struct {
	LandUseID    uint64  `gorm:"column:land_use_id"`
	LandUseCode  string  `gorm:"column:land_use_code"`
	LandUseName  string  `gorm:"column:land_use_name"`
	LandUseColor string  `gorm:"column:land_use_color"`
	CanBuild     bool    `gorm:"column:can_build"`
	RegionCount  int64   `gorm:"column:region_count"`
	AreaSqm      float64 `gorm:"column:area_sqm"`
	AreaPct      float64 `gorm:"-"`
}

type PolygonLayerStat struct {
	LayerID          uint64  `gorm:"column:layer_id"`
	LayerName        string  `gorm:"column:layer_name"`
	LayerDisplayName string  `gorm:"column:layer_display_name"`
	RegionCount      int64   `gorm:"column:region_count"`
	AreaSqm          float64 `gorm:"column:area_sqm"`
	AreaPct          float64 `gorm:"-"`
}

type PolygonBuildabilityStat struct {
	CanBuild    bool    `gorm:"column:can_build"`
	RegionCount int64   `gorm:"column:region_count"`
	AreaSqm     float64 `gorm:"column:area_sqm"`
	AreaPct     float64 `gorm:"-"`
}

type PolygonAnalysisWarning struct {
	Type     string
	Count    int64
	Severity string
	Message  string
}

type PolygonAnalysisSummary struct {
	ResultType          uint32
	Z                   uint32
	QueryAreaSqm        float64
	TotalMatchedAreaSqm float64
	RegionCount         int64
	ParcelCount         int64
	LandUseStats        []PolygonLandUseStat
	LayerStats          []PolygonLayerStat
	BuildabilityStats   []PolygonBuildabilityStat
	LegalKnownPct       float64
	LegalUnknownPct     float64
	Warnings            []PolygonAnalysisWarning
}

// API 1: thông tin nhẹ, dùng để render header / summary base.
type ParcelSeoSource struct {
	ParcelID  uint64  `json:"parcelId" gorm:"column:parcel_id"`
	AdrSearch string  `json:"adrSearch" gorm:"column:adr_search"`
	SeoID     *uint64 `json:"seoId" gorm:"column:seo_id"`
}

type ParcelInfoResponse struct {
	ParcelID uint64

	Lat     float64
	Lon     float64
	AreaSqm float64

	AddressText string
	Geometry    []byte

	MapNumber    string
	LandNumber   string
	PropertyCode string
	PropertyUUID string
	ShapeType    string
	Direction    string

	LandUseCode  string
	LandUseName  string
	LandUseColor string

	Province     string
	ProvinceCode string
	WardCode     string
	SeoID        uint64
	IsSeo        bool
}

// API 2: layers/impact - data chính để vẽ map + đánh giá rủi ro.
type ParcelLayersResponse struct {
	ParcelID      uint64
	ParcelAreaSqm float64

	ParcelInfo    *ParcelInfoResponse
	Summary       *ParcelSummary
	PrimaryPlan   *ParcelPlanInfo
	LandUseGroups []*PlanLandUseGroupDetail
	Layers        []*LayerImpact
	ConcludeLabel *RegionImpact
	Statistics    *ImpactStats
	Risk          *RiskAssessment
	Metadata      *ResponseMeta
}

type ParcelSummary struct {
	DominantLandUse          string
	DominantPercentage       float64
	TotalAffectedPercentage  float64
	TopLayerName             string
	Buildable                bool
	DominantLandUseCode      string
	DominantLandUseColor     string
	DominantLandUseGroupCode string
	DominantLandUseGroupName string
	DominantAreaSqm          float64
	DominantCanBuild         bool
}

type ParcelPlanInfo struct {
	LayerID          uint64
	LayerName        string
	DisplayName      string
	LandUseCode      string
	LandUseName      string
	LandUseColor     string
	GroupCode        string
	GroupName        string
	CanBuild         bool
	ImpactPercent    float64
	ImpactAreaSqm    float64
	LegalCanonical   string
	LegalConfidence  string
	LegalActiveState string
}

// Raw row từ SQL.
// Chỉ là “nguyên liệu” flat để usecase group lại, không phải response cuối.
type ParcelLayerRow struct {
	ParcelID      uint64
	ParcelAreaSqm float64

	LayerID            uint64
	LayerName          string
	LayerDisplayName   string
	LayerType          uint32
	LayerEffectiveDate string
	LayerLegalDocs     string `gorm:"column:layer_legal_docs"`
	LayerLegalDoc      string

	AuthorityIssuringID   uint64 `gorm:"column:authority_issuring_id"`
	AuthorityIssuringName string `gorm:"column:authority_issuring_name"`
	AuthorityIssuringCode string `gorm:"column:authority_issuring_code"`
	AuthorityIssuringDesc string `gorm:"column:authority_issuring_desc"`

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

	OverlapAreaSqm float64
	OverlapPct     float64

	CenterLat float64
	CenterLng float64

	Geometry []byte
}

type ParcelLayerRowV2 struct {
	// ── Parcel ──────────────────────────────────────────────
	ParcelID      uint64  `gorm:"column:parcel_id"`
	ParcelAreaSqm float64 `gorm:"column:parcel_area_sqm"`

	// ── Layer ───────────────────────────────────────────────
	LayerID            uint64  `gorm:"column:layer_id"`
	LayerName          string  `gorm:"column:layer_name"`
	LayerDisplayName   string  `gorm:"column:layer_display_name"`
	LayerType          uint32  `gorm:"column:layer_type"`
	LayerStatus        uint32  `gorm:"column:layer_status"`
	LayerLegalStatus   uint32  `gorm:"column:layer_legal_status"`
	LayerTrustValue    float32 `gorm:"column:layer_trust_value"`
	LayerEffectiveDate string  `gorm:"column:layer_effective_date"`
	LayerExpiryDate    string  `gorm:"column:layer_expiry_date"`
	LayerAvatar        string  `gorm:"column:layer_avatar"`
	// ── Authority ───────────────────────────────────────────
	AuthorityIssuringID   uint64 `gorm:"column:authority_issuring_id"`
	AuthorityIssuringName string `gorm:"column:authority_issuring_name"`
	AuthorityIssuringCode string `gorm:"column:authority_issuring_code"`
	AuthorityIssuringDesc string `gorm:"column:authority_issuring_desc"`

	// ── Label ───────────────────────────────────────────────
	LabelID          uint64 `gorm:"column:label_id"`
	LabelName        string `gorm:"column:label_name"`
	LabelDisplayName string `gorm:"column:label_display_name"`
	LabelColor       string `gorm:"column:label_color"`

	// Business fields (derived from qh_land_use PRIMARY)
	LabelGroupCode      string `gorm:"column:label_group_code"`
	LabelCanBuild       bool   `gorm:"column:label_can_build"`
	LabelPriority       int    `gorm:"column:label_priority"`
	LabelBuildCondition string `gorm:"column:label_build_condition"`

	// ── Land Use Group ──────────────────────────────────────
	GroupCode     string `gorm:"column:group_code"`
	GroupName     string `gorm:"column:group_name"`
	GroupColor    string `gorm:"column:group_color"`
	GroupCanBuild bool   `gorm:"column:group_can_build"`
	GroupPriority int    `gorm:"column:group_priority"`

	// ── Land Use (PRIMARY: qh_land_use via r.land_use_id) ──
	// Đây là nguồn sự thật duy nhất, align với GetParcelQuickInfo
	RegionLandUseCode  string `gorm:"column:region_land_use_code"`
	RegionLandUseName  string `gorm:"column:region_land_use_name"`
	RegionLandUseGroup string `gorm:"column:region_land_use_group"`
	RegionLandUseColor string `gorm:"column:region_land_use_color"`

	WarnLevel int32 `gorm:"column:warn_level"`

	// ── Region ──────────────────────────────────────────────
	RegionID          uint64 `gorm:"column:region_id"`
	RegionName        string `gorm:"column:region_name"`
	RegionDisplayName string `gorm:"column:region_display_name"`

	// ── Overlap ─────────────────────────────────────────────
	OverlapAreaSqm float64 `gorm:"column:overlap_area_sqm"`
	OverlapPct     float64 `gorm:"column:overlap_pct"`

	// ── Spatial ─────────────────────────────────────────────
	CenterLat float64 `gorm:"column:center_lat"`
	CenterLng float64 `gorm:"column:center_lng"`
	Geometry  []byte  `gorm:"column:geometry"`
}

type PlanLandUseGroupDetail struct {
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

type AuthorityIssuingDTO struct {
	ID          uint64 `json:"id"`
	Name        string `json:"name"`
	Code        string `json:"code"`
	Description string `json:"description"`
}

type LayerImpact struct {
	ID               uint64
	Name             string
	DisplayName      string
	Type             uint32
	Severity         string
	EffectiveDate    string
	LegalDoc         string
	LegalDocs        []*LayerLegalDTO
	AuthorityIssuing *AuthorityIssuingDTO
	Labels           []*LabelImpact
	TotalArea        float64
	TotalPct         float64
	LegalStatus      uint32
	LegalStatusName  string
	LayerAvatar      string
}

type LabelImpact struct {
	ID          uint64
	Name        string
	DisplayName string
	Color       string
	FontSize    uint32
	FontWeight  string
	MinZoom     uint32
	MaxZoom     uint32
	Regions     []*RegionImpact
	TotalArea   float64
	TotalPct    float64
}

type RegionImpact struct {
	ID             uint64
	Name           string
	DisplayName    string
	LandUseCode    string
	LandUseName    string
	LandUseGroup   string
	ColorRGB       string
	OverlapAreaSqm float64
	OverlapPct     float64
	CenterLat      float64
	CenterLng      float64
	Geometry       []byte
	WarnLevel      int32
	CanBuild       bool
}

type RGB struct {
	Red   uint8
	Green uint8
	Blue  uint8
}

type ImpactStats struct {
	TotalOverlapArea     float64
	TotalOverlapPct      float64
	LayersAffected       uint32
	LabelsAffected       uint32
	RegionsAffected      uint32
	HasConflict          bool
	SpatialDataAvailable bool
	Warning              string
}

type RiskAssessment struct {
	Level          string
	Score          uint32
	CanBuild       bool
	Reasons        []string
	Recommendation string
}

type ResponseMeta struct {
	DataSource     string
	QueriedAt      string
	ResponseTimeMs int64
	Cached         bool
}
