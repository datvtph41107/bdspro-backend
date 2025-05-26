package repo

import (
	"bdspro/internal/domain"
	"context"
)

type AssetCostRepo interface {
	Create(c context.Context, entity *domain.AssetCost) error
	Update(c context.Context, id uint64, entity *domain.AssetCost) error
	Delete(c context.Context, id uint64) error
	GetData(c context.Context) ([]domain.AssetCost, error)
	GetByID(c context.Context, id uint64) (*domain.AssetCost, error)
}
