package repo

import (
	"context"
	qh_domain "tqd/internal/domain/qh"
	"tqd/internal/enums"
)

type LayerRepository interface {
	Create(ctx context.Context, layer *qh_domain.QHLayer) error
	Update(ctx context.Context, layer *qh_domain.QHLayer) error
	Delete(ctx context.Context, id uint64) error
	GetByID(ctx context.Context, id uint64) (*qh_domain.QHLayer, error)
	GetByName(ctx context.Context, name string) (*qh_domain.QHLayer, error)

	List(ctx context.Context, filter *LayerFilter) ([]qh_domain.QHLayer, int64, error)
	ListClient(ctx context.Context, filter *ClientLayerFilter) ([]qh_domain.QHLayer, int64, error)

	UpdateStatus(ctx context.Context, id uint64, status enums.LayerStatus) error
	UpdateVisibility(ctx context.Context, id uint64, visible int32) error
}

// Filter admin list
type LayerFilter struct {
	Status   *enums.LayerStatus
	Type     *enums.LayerType
	Visible  *int32 // 0, 1, nil
	Search   *string
	OrderBy  string
	OrderDir string // "ASC" or "DESC"
	Page     int
	Limit    int
}

// Filter client list
type ClientLayerFilter struct {
	Type   *enums.LayerType
	Search *string
	Page   int
	Limit  int
}
