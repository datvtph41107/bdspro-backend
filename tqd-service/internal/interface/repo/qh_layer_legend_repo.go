package repo

import (
	"context"

	qh_domain "tqd/internal/domain/qh"
)

// QHLayerLegendRepository truy cập bảng qh_layer_legends (chú giải).
type QHLayerLegendRepository interface {
	Create(ctx context.Context, row *qh_domain.QHLayerLegend) error
	Update(ctx context.Context, row *qh_domain.QHLayerLegend) error
	Delete(ctx context.Context, id uint64) error
	GetByID(ctx context.Context, id uint64) (*qh_domain.QHLayerLegend, error)
	GetByLayerAndLabel(ctx context.Context, layerID, labelID uint64) (*qh_domain.QHLayerLegend, error)
	List(ctx context.Context, offset, limit int, layerID, labelID *uint64) ([]qh_domain.QHLayerLegend, int64, error)
}
