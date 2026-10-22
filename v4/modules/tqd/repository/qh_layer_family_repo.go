package repository

import (
	"context"

	qh "mdn4/modules/tqd/entity/qh"
)

// QHLayerFamilyRepository truy cập bảng qh_layer_families.
type QHLayerFamilyRepository interface {
	Create(ctx context.Context, row *qh.QHLayerFamily) error
	Update(ctx context.Context, row *qh.QHLayerFamily) error
	Delete(ctx context.Context, id uint64) error
	GetByID(ctx context.Context, id uint64) (*qh.QHLayerFamily, error)
	List(ctx context.Context, offset, limit int, search, sort string) ([]qh.QHLayerFamily, int64, error)
	ListClient(ctx context.Context, offset, limit int, search string) ([]qh.QHLayerFamily, int64, error)
}
