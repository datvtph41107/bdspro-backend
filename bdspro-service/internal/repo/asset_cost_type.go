package repo

import (
	"bdspro/internal/domain"
	"context"
)

type AssetCostTypeRepo interface {
	Create(c context.Context, entity *domain.AssetCostType) error
	Update(c context.Context, id uint64, entity *domain.AssetCostType) error
	Delete(c context.Context, id uint64) error
	GetData(c context.Context) ([]domain.AssetCostType, error)
	GetByID(c context.Context, id uint64) (*domain.AssetCostType, error)
}
