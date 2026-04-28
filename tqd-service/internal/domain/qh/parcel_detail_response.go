package qh_domain

// ParcelDetailResponse - Response cho API Detail VIP
type ParcelDetailResponse struct {
	// Basic parcel info
	ParcelID     uint64  `json:"parcelId"`
	Lat          float64 `json:"lat"`
	Lon          float64 `json:"lon"`
	TotalAreaSqm float64 `json:"totalAreaSqm"`
	AddressText  string  `json:"addressText"`
	Geometry     []byte  `json:"geometry,omitempty"`

	// Parcel detail info
	MapNumber    string `json:"mapNumber"`
	LandNumber   string `json:"landNumber"`
	PropertyCode string `json:"propertyCode"`
	PropertyUUID string `json:"propertyUuid"`
	ShapeType    string `json:"shapeType"`
	Direction    string `json:"direction"`

	// Current land use
	CurrentLandUse *LandUseBrief `json:"currentLandUse"`

	// Administrative
	Administrative *AdminBrief `json:"administrative"`

	// Planning impacts
	Layers []*LayerImpact `json:"layers"`

	// Statistics
	Statistics *ImpactStats `json:"statistics"`

	// Risk assessment
	Risk *RiskAssessment `json:"risk"`

	// Metadata
	Metadata *ResponseMeta `json:"metadata"`
}

// LandUseBrief - Loại đất hiện tại
type LandUseBrief struct {
	Code  string `json:"code"`
	Name  string `json:"name"`
	Group string `json:"group"`
	Color string `json:"color"`
}

// AdminBrief - Thông tin hành chính
type AdminBrief struct {
	Province     string `json:"province"`
	ProvinceCode string `json:"provinceCode"`
	WardCode     string `json:"wardCode"`
	Address      string `json:"address"`
}

// LayerImpact - Layer quy hoạch
type LayerImpact struct {
	ID            uint64         `json:"id"`
	Name          string         `json:"name"`
	DisplayName   string         `json:"displayName"`
	Type          uint32         `json:"type"`
	EffectiveDate string         `json:"effectiveDate,omitempty"`
	LegalDoc      string         `json:"legalDoc,omitempty"`
	Labels        []*LabelImpact `json:"labels"`
	TotalArea     float64        `json:"totalArea"`
	TotalPct      float64        `json:"totalPercentage"`
}

// LabelImpact - Nhãn (theo loại đất)
type LabelImpact struct {
	ID          uint64          `json:"id"`
	Name        string          `json:"name"`
	DisplayName string          `json:"displayName"`
	Color       string          `json:"color"`
	FontSize    uint32          `json:"fontSize"`
	FontWeight  string          `json:"fontWeight"`
	MinZoom     uint32          `json:"minZoom"`
	MaxZoom     uint32          `json:"maxZoom"`
	Regions     []*RegionImpact `json:"regions"`
	TotalArea   float64         `json:"totalArea"`
	TotalPct    float64         `json:"totalPercentage"`
}

// RegionImpact - Region cụ thể
type RegionImpact struct {
	ID             uint64  `json:"id"`
	Name           string  `json:"name"`
	DisplayName    string  `json:"displayName"`
	LandUseCode    string  `json:"landUseCode"`
	LandUseName    string  `json:"landUseName"`
	LandUseGroup   string  `json:"landUseGroup"`
	ColorRGB       *RGB    `json:"colorRgb"`
	OverlapAreaSqm float64 `json:"overlapAreaSqm"`
	OverlapPct     float64 `json:"overlapPercentage"`
	Geometry       []byte  `json:"geometry,omitempty"`
}

// RGB - Màu sắc
type RGB struct {
	Red   uint8 `json:"red"`
	Green uint8 `json:"green"`
	Blue  uint8 `json:"blue"`
}

// ImpactStats - Thống kê
type ImpactStats struct {
	TotalOverlapArea     float64 `json:"totalOverlapArea"`
	TotalOverlapPct      float64 `json:"totalOverlapPercentage"`
	LayersAffected       uint32  `json:"layersAffected"`
	LabelsAffected       uint32  `json:"labelsAffected"`
	RegionsAffected      uint32  `json:"regionsAffected"`
	HasConflict          bool    `json:"hasConflict"`
	SpatialDataAvailable bool    `json:"spatialDataAvailable"`
	Warning              string  `json:"warning,omitempty"`
}

// RiskAssessment - Đánh giá rủi ro
type RiskAssessment struct {
	Level          string   `json:"level"`
	Score          uint32   `json:"score"`
	CanBuild       bool     `json:"canBuild"`
	Reasons        []string `json:"reasons"`
	Recommendation string   `json:"recommendation"`
}

// ResponseMeta - Metadata
type ResponseMeta struct {
	DataSource     string `json:"dataSource"`
	QueriedAt      string `json:"queriedAt"`
	ResponseTimeMs int64  `json:"responseTimeMs"`
	Cached         bool   `json:"cached"`
}
