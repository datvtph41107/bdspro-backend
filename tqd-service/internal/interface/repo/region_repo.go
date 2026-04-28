package repo

import (
	_dto "common/domain/dto"
	"context"
	qh_domain "tqd/internal/domain/qh"
)

// RegionRepository - Interface cho region repository
type RegionRepository interface {
	// Basic CRUD
	Create(ctx context.Context, region *qh_domain.QHRegion) error
	CreateBatch(ctx context.Context, regions []*qh_domain.QHRegion) error
	Update(ctx context.Context, region *qh_domain.QHRegion) error
	Delete(ctx context.Context, id uint64) error
	GetByID(ctx context.Context, id uint64) (*qh_domain.QHRegion, error)
	GetByIDs(ctx context.Context, ids []uint64) ([]qh_domain.QHRegion, error)

	// List with filters
	List(ctx context.Context, filter *RegionFilter) ([]qh_domain.QHRegion, int64, error)
	ListClient(ctx context.Context, filter *ClientRegionFilter) ([]qh_domain.QHRegion, int64, error)

	// Spatial queries
	FindByPoint(ctx context.Context, lat, lng float64) (*qh_domain.QHRegion, error)
	FindByBBox(ctx context.Context, minLng, minLat, maxLng, maxLat float64, filter *RegionFilter) ([]qh_domain.QHRegion, int64, error)

	// Statistics
	GetStatsByLayer(ctx context.Context, layerID uint64) (*RegionStats, error)
}

// RegionFilter - Filter cho admin list
type RegionFilter struct {
	LayerID     *uint64
	LandUseCode *string
	Status      *int
	Search      *string
	MinAreaHa   *float64
	MaxAreaHa   *float64
	OrderBy     string
	OrderDir    string
	Pagable     *_dto.Pagable
}

// ClientRegionFilter - Filter cho client list
type ClientRegionFilter struct {
	LayerID     *uint64
	LandUseCode *string
	BBox        *BBox
	Pagable     *_dto.Pagable
}

// BBox - Bounding box
type BBox struct {
	MinLng float64
	MinLat float64
	MaxLng float64
	MaxLat float64
}

// RegionStats - Thống kê regions
type RegionStats struct {
	LayerID      uint64        `json:"layerId"`
	TotalRegions int64         `json:"totalRegions"`
	TotalAreaHa  float64       `json:"totalAreaHa"`
	AvgAreaHa    float64       `json:"avgAreaHa"`
	MinAreaHa    float64       `json:"minAreaHa"`
	MaxAreaHa    float64       `json:"maxAreaHa"`
	LandUseStats []LandUseStat `json:"landUseStats"`
}

// LandUseStat - Thống kê theo loại đất
type LandUseStat struct {
	LandUseCode string  `json:"landUseCode"`
	LandUseName string  `json:"landUseName"`
	Count       int64   `json:"count"`
	TotalAreaHa float64 `json:"totalAreaHa"`
}
