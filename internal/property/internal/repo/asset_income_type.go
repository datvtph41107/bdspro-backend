package repo

import (
	"bdspro/internal/domain"
	"context"
)

type AssetIncomeTypeRepo interface {
	Create(c context.Context, entity *domain.AssetIncomeType) error
	Update(c context.Context, id uint64, entity *domain.AssetIncomeType) error
	Delete(c context.Context, id uint64) error
	GetData(c context.Context) ([]domain.AssetIncomeType, error)
	GetByID(c context.Context, id uint64) (*domain.AssetIncomeType, error)
}
