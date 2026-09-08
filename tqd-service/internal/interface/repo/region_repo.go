package repo

import (
	"context"
	"database/sql"

	_dto "common/domain/dto"
	qh_domain "tqd/internal/domain/qh"
)

type RegionRepository interface {
	// Basic CRUD
	Create(ctx context.Context, region *qh_domain.QHRegion) error
	CreateBatch(ctx context.Context, regions []*qh_domain.QHRegion) error
	Update(ctx context.Context, region *qh_domain.QHRegion) error
	Delete(ctx context.Context, id uint64) error
	// SoftDeleteAllByLayerID xóa mềm mọi region thuộc layer (layer_id).
	SoftDeleteAllByLayerID(ctx context.Context, layerID uint64) error
	// HardDeleteAllByLayerID xóa thật mọi region thuộc layer (layer_id).
	HardDeleteAllByLayerID(ctx context.Context, layerID uint64) error
	GetByID(ctx context.Context, id uint64) (*qh_domain.QHRegion, error)
	GetByIDs(ctx context.Context, ids []uint64) ([]qh_domain.QHRegion, error)

	// List with filters
	List(ctx context.Context, filter *RegionFilter) ([]qh_domain.QHRegion, int64, error)
	ListByImportBatchID(ctx context.Context, batchID string) ([]qh_domain.QHRegion, error)
	// ListImportErrorLogsByLayerID — log batch insert lỗi (regions_payload JSON), mới nhất trước
	ListImportErrorLogsByLayerID(ctx context.Context, layerID uint64, offset, limit int) ([]qh_domain.QHRegionImportErrorLog, int64, error)
	GetImportErrorLogByID(ctx context.Context, id uint64) (*qh_domain.QHRegionImportErrorLog, error)
	UpdateImportErrorLog(ctx context.Context, row *qh_domain.QHRegionImportErrorLog) error
	// TryBeginImportErrorRetry — CAS pending → processing; true nếu chiếm được lock
	TryBeginImportErrorRetry(ctx context.Context, id uint64) (bool, error)
	// ResetImportErrorRetryIfProcessing — processing → pending (rollback lock / lỗi giữa chừng)
	ResetImportErrorRetryIfProcessing(ctx context.Context, id uint64) error
	HardDeleteImportErrorLogsByLayerID(ctx context.Context, layerID uint64) error

	ListClient(ctx context.Context, filter *ClientRegionFilter) ([]qh_domain.QHRegion, int64, error)

	// Spatial queries
	FindByPoint(ctx context.Context, lat, lng float64, layerID *uint64) (*qh_domain.QHRegion, error)
	FindByBBox(ctx context.Context, minLng, minLat, maxLng, maxLat float64, layerID *uint64, labelID *uint64, limit, offset int) ([]qh_domain.QHRegion, int64, error)

	//Update regions label_id
	UpdateLabelForAllRegions(ctx context.Context, oldLabelID uint64, newLabelID uint64) error
	// UpdateLabelForRegionsWhereLabelIn set label_id cho mọi region đang trỏ tới một trong labelIDs (một query).
	UpdateLabelForRegionsWhereLabelIn(ctx context.Context, labelIDs []uint64, newLabelID uint64) error
	UpdateLabelForRegion(ctx context.Context, regionID uint64, labelID uint64) error
	GetSyncReferenceByLayerLabel(ctx context.Context, layerID, labelID uint64) (*RegionSyncReference, error)
	SyncLandUseAndLegendByLayerLabel(ctx context.Context, layerID, labelID, landUseID, legendID uint64) (int64, error)

	// Statistics
	GetStatsByLayer(ctx context.Context, layerID uint64) (*RegionStats, error)
	GetByBBox(ctx context.Context, minX, minY, maxX, maxY float64, regionType *uint32, layerID *uint64) ([]qh_domain.QHRegion, error)
	GetIntersections(ctx context.Context, id uint64, bufferMeters float64) ([]qh_domain.QHRegion, error)
	GetBufferGeometry(ctx context.Context, sourceID uint64, distanceMeters float64) (string, error)
	GetLayerBBox(ctx context.Context, layerID uint64) (minLon, minLat, maxLon, maxLat float64, err error)
	GetMVTTile(ctx context.Context, layerID uint64, z uint8, x, y uint32) ([]byte, error)
	StreamGeoJSON(ctx context.Context, layerID uint64) (*sql.Rows, error)
	StreamGeoJSONByFamilyID(ctx context.Context, familyID uint64) (*sql.Rows, error)
}

type RegionFilter struct {
	LayerID          *uint64
	LabelID          *uint64
	LandUseCode      *string
	Status           *uint32
	ProcessingStatus *uint32
	Search           *string
	MinAreaHa        *float64
	MaxAreaHa        *float64
	OrderBy          string
	OrderDir         string
	Pagable          *_dto.Pagable
}

type ClientRegionFilter struct {
	LayerID     *uint64
	LabelID     *uint64
	LandUseCode *string
	BBox        *BBox
	Pagable     *_dto.Pagable
}

type BBox struct {
	MinLng float64
	MinLat float64
	MaxLng float64
	MaxLat float64
}

type RegionStats struct {
	LayerID      uint64        `json:"layerId"`
	TotalRegions int64         `json:"totalRegions"`
	TotalAreaSqm float64       `json:"totalAreaHa"`
	AvgAreaSqm   float64       `json:"avgAreaHa"`
	MinAreaSqm   float64       `json:"minAreaHa"`
	MaxAreaSqm   float64       `json:"maxAreaHa"`
	LandUseStats []LandUseStat `json:"landUseStats"`
}

type RegionSyncReference struct {
	LegendID  uint64
	LandUseID uint64
}

type LandUseStat struct {
	LandUseCode  string  `json:"landUseCode"`
	LandUseName  string  `json:"landUseName"`
	Count        int64   `json:"count"`
	TotalAreaSqm float64 `json:"totalAreaSqm"`
}
