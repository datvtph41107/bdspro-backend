package repo

import (
	"context"

	qh_domain "tqd/internal/domain/qh"
)

type QHLayerLandUseGroupRepository interface {
	Create(ctx context.Context, row *qh_domain.QHLayerLandUseGroup) error
	Update(ctx context.Context, row *qh_domain.QHLayerLandUseGroup) error
	Delete(ctx context.Context, id uint64) error
	GetByID(ctx context.Context, id uint64) (*qh_domain.QHLayerLandUseGroup, error)
	GetByLayerAndGroup(ctx context.Context, layerID, groupID uint64) (*qh_domain.QHLayerLandUseGroup, error)
	List(ctx context.Context, offset, limit int, layerID, groupID *uint64) ([]qh_domain.QHLayerLandUseGroup, int64, error)
	ListVisibleByLayer(ctx context.Context, layerID uint64, offset, limit int) ([]qh_domain.QHLayerLandUseGroup, int64, error)
	SoftDeleteByLayerID(ctx context.Context, layerID uint64) error
}
