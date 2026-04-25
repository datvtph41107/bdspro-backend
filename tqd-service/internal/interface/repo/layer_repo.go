// repo/layer_repo.go
package repo

import (
	"context"
	qh_domain "tqd/internal/domain/qh"
	"tqd/internal/enums"
)

type LayerRepository interface {
	Create(ctx context.Context, layer *qh_domain.QHLayer) error
	Update(ctx context.Context, layer *qh_domain.QHLayer) error
	Delete(ctx context.Context, id uint) error
	GetByID(ctx context.Context, id uint) (*qh_domain.QHLayer, error)
	GetByName(ctx context.Context, name string) (*qh_domain.QHLayer, error)
	List(ctx context.Context, filter *LayerFilter) ([]qh_domain.QHLayer, int64, error)
	UpdateStatus(ctx context.Context, id uint, status enums.LayerStatus) error
}

type LayerFilter struct {
	Status  *enums.LayerStatus // Dùng enum thay vì string
	Search  *string
	Type    *enums.LayerType // Dùng enum thay vì string
	Visible *bool
	Page    int
	Limit   int
}
