package repo

import (
	_dto "common/domain/dto"
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

// LayerFilter - Filter cho admin list
type LayerFilter struct {
	Status   *enums.LayerStatus
	Type     *enums.LayerType
	Visible  *int32
	Search   *string
	OrderBy  string
	OrderDir string
	Pagable  *_dto.Pagable
}

// ClientLayerFilter - Filter cho client list
type ClientLayerFilter struct {
	Type    *enums.LayerType
	Search  *string
	Pagable *_dto.Pagable
}
