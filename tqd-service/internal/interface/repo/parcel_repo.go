package repo

import (
	_dto "common/domain/dto"
	"context"
	qh_domain "tqd/internal/domain/qh"
)

type IParcelRepo interface {
	// Preview & Search
	GetParcelPreview(ctx context.Context, parcelID uint64) (*qh_domain.QHParcelInfo, error)
	SearchByAddress(ctx context.Context, keyword string, pagable _dto.Pagable) ([]qh_domain.QHParcelInfo, int64, error)

	// Basic CRUD
	GetByID(ctx context.Context, id uint64) (*qh_domain.Parcel, error)
	GetByIDs(ctx context.Context, ids []uint64) ([]qh_domain.Parcel, error)
	CreateBatch(ctx context.Context, parcels []qh_domain.Parcel) error
	FindByLocation(ctx context.Context, lat, lng float64) (*qh_domain.Parcel, error)
	FindIntersectPolygon(ctx context.Context, points []Point, limit, offset int) ([]qh_domain.Parcel, int64, error)
	FindWithinPolygon(ctx context.Context, points []Point, limit, offset int) ([]qh_domain.Parcel, int64, error)

	// Detail API
	GetParcelDetail(ctx context.Context, parcelID uint64) (*ParcelDetailResult, error)
	GetSpatialImpacts(ctx context.Context, parcelID uint64) ([]*SpatialImpactRow, error)
	GetLandUseByCode(ctx context.Context, code string) (*qh_domain.LandUseDict, error)
}

type Point struct {
	Lat float64
	Lng float64
}

// ParcelDetailResult - Kết quả join từ database
type ParcelDetailResult struct {
	ID       uint64  `gorm:"column:id"`
	Lat      float64 `gorm:"column:lat"`
	Lon      float64 `gorm:"column:lon"`
	Geometry []byte  `gorm:"column:geometry"`

	ParcelID          uint64  `gorm:"column:parcel_id"`
	PropertyCode      string  `gorm:"column:property_code"`
	PropertyUUID      string  `gorm:"column:property_uuid"`
	MapNumber         string  `gorm:"column:map_number"`
	LandNumber        string  `gorm:"column:land_number"`
	TotalAreaSqm      float64 `gorm:"column:total_area_sqm"`
	ShapeType         string  `gorm:"column:shape_type"`
	Direction         string  `gorm:"column:direction"`
	ActualLandUseCode string  `gorm:"column:actual_land_use_code"`
	ProvinceCode      string  `gorm:"column:province_code"`
	DistrictCode      string  `gorm:"column:district_code"`
	WardCode          string  `gorm:"column:ward_code"`
	AddressText       string  `gorm:"column:address_text"`

	ProvinceName string `gorm:"column:province_name"`
	WardName     string `gorm:"column:ward_name"`
}

// SpatialImpactRow - Spatial impact từ regions
type SpatialImpactRow struct {
	LayerID          uint64 `gorm:"column:layer_id"`
	LayerName        string `gorm:"column:layer_name"`
	LayerDisplayName string `gorm:"column:layer_display_name"`
	LayerType        uint32 `gorm:"column:layer_type"`
	EffectiveDate    string `gorm:"column:effective_date"`
	LegalDoc         string `gorm:"column:legal_doc"`

	LabelID          uint64 `gorm:"column:label_id"`
	LabelName        string `gorm:"column:label_name"`
	LabelDisplayName string `gorm:"column:label_display_name"`
	LabelColor       string `gorm:"column:label_color"`
	LabelFontSize    uint32 `gorm:"column:label_font_size"`
	LabelFontWeight  string `gorm:"column:label_font_weight"`
	LabelMinZoom     uint32 `gorm:"column:label_min_zoom"`
	LabelMaxZoom     uint32 `gorm:"column:label_max_zoom"`

	RegionID          uint64  `gorm:"column:region_id"`
	RegionName        string  `gorm:"column:region_name"`
	RegionDisplayName string  `gorm:"column:region_display_name"`
	LandUseCode       string  `gorm:"column:land_use_code"`
	LandUseName       string  `gorm:"column:land_use_name"`
	LandUseGroup      string  `gorm:"column:land_use_group"`
	ColorRed          uint8   `gorm:"column:color_red"`
	ColorGreen        uint8   `gorm:"column:color_green"`
	ColorBlue         uint8   `gorm:"column:color_blue"`
	OverlapArea       float64 `gorm:"column:overlap_area"`
	OverlapPct        float64 `gorm:"column:overlap_pct"`
	Geometry          []byte  `gorm:"column:geometry"`
}
