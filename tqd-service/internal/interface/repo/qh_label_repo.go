package repo

import (
	"context"
	qh_domain "tqd/internal/domain/qh"
)

type QHLabelRepository interface {
	Create(ctx context.Context, label *qh_domain.QHLabel) error
	Update(ctx context.Context, label *qh_domain.QHLabel) error
	Delete(ctx context.Context, id uint64) error
	GetByID(ctx context.Context, id uint64) (*qh_domain.QHLabel, error)
	ListByLayerID(ctx context.Context, layerID uint64) ([]qh_domain.QHLabel, error)
}
