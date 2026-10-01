package repository

import (
	"context"

	qh "userv4/modules/tqd/entity/qh"
)

type QHLayerLegendRepository interface {
	Create(ctx context.Context, row *qh.QHLayerLegend) error
	Update(ctx context.Context, row *qh.QHLayerLegend) error
	Delete(ctx context.Context, id uint64) error
	GetByID(ctx context.Context, id uint64) (*qh.QHLayerLegend, error)
	GetByLayerAndLabel(ctx context.Context, layerID, labelID uint64) (*qh.QHLayerLegend, error)
	List(ctx context.Context, offset, limit int, layerID, labelID *uint64, legendType *string, landUseGroupID *uint64, isVisible *bool) ([]qh.QHLayerLegend, int64, error)
	ListVisibleByLayer(ctx context.Context, layerID uint64, offset, limit int) ([]qh.QHLayerLegend, int64, error)
	SoftDeleteByLayerID(ctx context.Context, layerID uint64) error
	HardDeleteByLayerID(ctx context.Context, layerID uint64) error
	SoftDeleteByLabelID(ctx context.Context, labelID uint64) error
}
