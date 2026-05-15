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

// API 1: thông tin nhẹ, dùng để render header / summary base.
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
}

// API 2: layers/impact - data chính để vẽ map + đánh giá rủi ro.
type ParcelLayersResponse struct {
	ParcelID      uint64
	ParcelAreaSqm float64

	Summary       *ParcelSummary
	Layers        []*LayerImpact
	ConcludeLabel *RegionImpact
	Statistics    *ImpactStats
	Risk          *RiskAssessment
	Metadata      *ResponseMeta
}

type ParcelSummary struct {
	DominantLandUse         string
	DominantPercentage      float64
	TotalAffectedPercentage float64
	TopLayerName            string
	Buildable               bool
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

	// Layer fields
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

	LabelID             uint64
	LabelName           string
	LabelDisplayName    string
	LabelColor          string
	LabelGroupCode      uint32
	LabelCanBuild       bool
	LabelPriority       int
	LabelBuildCondition string

	// Region fields
	RegionID          uint64
	RegionName        string
	RegionDisplayName string
	RegionLandUseCode string
	RegionLandUseName string
	RegionColorRed    uint8
	RegionColorGreen  uint8
	RegionColorBlue   uint8

	// Impact
	OverlapAreaSqm float64
	OverlapPct     float64

	// Spatial
	CenterLat float64
	CenterLng float64
	Geometry  []byte
}

type LayerImpact struct {
	ID            uint64
	Name          string
	DisplayName   string
	Type          uint32
	Severity      string
	EffectiveDate string
	LegalDoc      string
	LegalDocs     []*LayerLegalDTO
	Labels        []*LabelImpact
	TotalArea     float64
	TotalPct      float64
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
