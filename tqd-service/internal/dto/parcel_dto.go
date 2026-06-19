package dto

import (
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

	Analysis *PolygonAnalysisSummary
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
	ParcelID      uint64
	ParcelAreaSqm float64

	LayerID            uint64
	LayerName          string
	LayerDisplayName   string
	LayerType          uint32
	LayerStatus        uint32
	LayerLegalStatus   uint32
	LayerTrustValue    float32
	LayerEffectiveDate string
	LayerExpiryDate    string
	LayerAuthority     string

	AuthorityIssuringID   uint64 `gorm:"column:authority_issuring_id"`
	AuthorityIssuringName string `gorm:"column:authority_issuring_name"`
	AuthorityIssuringCode string `gorm:"column:authority_issuring_code"`
	AuthorityIssuringDesc string `gorm:"column:authority_issuring_desc"`

	LabelID             uint64
	LabelName           string
	LabelDisplayName    string
	LabelColor          string
	LabelGroupCode      string
	LabelCanBuild       bool
	LabelPriority       int
	LabelBuildCondition string

	GroupCode     string
	GroupName     string
	GroupColor    string
	GroupCanBuild bool
	GroupPriority int

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
	Geometry  []byte
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
	ColorRGB       *RGB
	OverlapAreaSqm float64
	OverlapPct     float64
	CenterLat      float64
	CenterLng      float64
	Geometry       []byte
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
