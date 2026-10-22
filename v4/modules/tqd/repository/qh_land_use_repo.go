package repository

import (
	_dto "common/domain/dto"
	"context"

	qh "mdn4/modules/tqd/entity/qh"
)

type QHLandUseRepository interface {
	Create(ctx context.Context, row *qh.QHLandUse) error
	Update(ctx context.Context, row *qh.QHLandUse) error
	Delete(ctx context.Context, id uint64) error
	GetByID(ctx context.Context, id uint64) (*qh.QHLandUse, error)
	GetByLayerAndLandUse(ctx context.Context, layerID, landUseID uint64) (*qh.QHLandUse, error)
	LinkLayerLandUse(ctx context.Context, layerID, landUseID uint64) error
	List(ctx context.Context, pagable *_dto.Pagable, layerID *uint64) ([]qh.QHLandUse, int64, error)
	ListVisibleByLayer(ctx context.Context, layerID uint64, offset, limit int) ([]qh.QHLandUse, int64, error)
	SoftDeleteByLayerID(ctx context.Context, layerID uint64) error
	HardDeleteByLayerID(ctx context.Context, layerID uint64) error
}
