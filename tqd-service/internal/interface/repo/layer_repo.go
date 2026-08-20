package repo

import (
	_dto "common/domain/dto"
	"context"
	qh_domain "tqd/internal/domain/qh"
	"tqd/internal/dto"
	"tqd/internal/enums"
)

type LayerRepository interface {
	// Create: nếu replaceLayerIDs không rỗng, chạy trong transaction — tạo layer rồi đánh dấu các layer cũ bị thay thế.
	Create(ctx context.Context, layer *qh_domain.QHLayer, replaceLayerIDs []uint64, updatedBy uint64) error
	Update(ctx context.Context, layer *qh_domain.QHLayer) error
	Delete(ctx context.Context, id uint64) error
	HardDelete(ctx context.Context, id uint64) error
	GetByID(ctx context.Context, id uint64) (*qh_domain.QHLayer, error)
	GetByName(ctx context.Context, name string) (*qh_domain.QHLayer, error)

	List(ctx context.Context, filter *LayerFilter) ([]qh_domain.QHLayer, int64, error)
	ListClient(ctx context.Context, filter *ClientLayerFilter) ([]qh_domain.QHLayer, int64, error)

	UpdateStatus(ctx context.Context, id uint64, status enums.LayerStatus) error
	UpdateVisibility(ctx context.Context, id uint64, visible int32) error

	UpdateImportStatus(ctx context.Context, id uint64, st enums.LayerImportStatus, importBatchID *string) error
	GetByImportBatchID(ctx context.Context, batchID string) (*qh_domain.QHLayer, error)
	ResetLayerImportTracking(ctx context.Context, id uint64) error
	GetZoomRangeByFamilyID(ctx context.Context, familyID uint64) (minZoom, maxZoom uint32, err error)

	GetLayerSeoSource(ctx context.Context, layerID uint64) (*dto.LayerSeoSource, error)
	GetLayerSeoSourcesByIDs(ctx context.Context, layerIDs []uint64) ([]dto.LayerSeoSource, error)
	UpdateLayerSeoID(ctx context.Context, layerID uint64, seoID uint64) error
}

// LayerFilter - Filter cho admin list
type LayerFilter struct {
	_dto.Pagable
	Status            *enums.LayerStatus
	Type              *enums.LayerType
	Visible           *int32
	Search            *string
	OrderBy           string
	OrderDir          string
	PlanningProjectID *uint64
}

// ClientLayerFilter - Filter cho client list
type ClientLayerFilter struct {
	_dto.Pagable
	Type   *enums.LayerType
	Search *string
	DocID  uint64
	// Page   int
	// Limit  int
}
