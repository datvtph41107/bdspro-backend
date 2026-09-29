package repo

import (
	"context"

	qh_domain "tqd/internal/domain/qh"
)

// QHLayerFamilyRepository truy cập bảng qh_layer_families.
type QHLayerFamilyRepository interface {
	Create(ctx context.Context, row *qh_domain.QHLayerFamily) error
	Update(ctx context.Context, row *qh_domain.QHLayerFamily) error
	Delete(ctx context.Context, id uint64) error
	GetByID(ctx context.Context, id uint64) (*qh_domain.QHLayerFamily, error)
	List(ctx context.Context, offset, limit int, search, sort string) ([]qh_domain.QHLayerFamily, int64, error)
	ListClient(ctx context.Context, offset, limit int, search string) ([]qh_domain.QHLayerFamily, int64, error)
}
