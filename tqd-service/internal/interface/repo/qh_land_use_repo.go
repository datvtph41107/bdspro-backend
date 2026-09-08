package repo

import (
	_dto "common/domain/dto"
	"context"

	qh_domain "tqd/internal/domain/qh"
)

type QHLandUseRepository interface {
	Create(ctx context.Context, row *qh_domain.QHLandUse) error
	Update(ctx context.Context, row *qh_domain.QHLandUse) error
	Delete(ctx context.Context, id uint64) error
	GetByID(ctx context.Context, id uint64) (*qh_domain.QHLandUse, error)
	GetByLayerAndLandUse(ctx context.Context, layerID, landUseID uint64) (*qh_domain.QHLandUse, error)
	LinkLayerLandUse(ctx context.Context, layerID, landUseID uint64) error
	List(ctx context.Context, pagable *_dto.Pagable, layerID *uint64) ([]qh_domain.QHLandUse, int64, error)
	ListVisibleByLayer(ctx context.Context, layerID uint64, offset, limit int) ([]qh_domain.QHLandUse, int64, error)
	SoftDeleteByLayerID(ctx context.Context, layerID uint64) error
	HardDeleteByLayerID(ctx context.Context, layerID uint64) error
}
