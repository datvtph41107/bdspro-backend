// repo/layer_repo.go
package repo

import (
	"context"
	"tqd/internal/domain"
	"tqd/internal/enums"
)

type LayerRepository interface {
	Create(ctx context.Context, layer *domain.Layer) error
	Update(ctx context.Context, layer *domain.Layer) error
	Delete(ctx context.Context, id uint) error
	GetByID(ctx context.Context, id uint) (*domain.Layer, error)
	GetByName(ctx context.Context, name string) (*domain.Layer, error)
	List(ctx context.Context, filter *LayerFilter) ([]domain.Layer, int64, error)
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
