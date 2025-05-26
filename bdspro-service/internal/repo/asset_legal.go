package repo

import (
	"bdspro/internal/domain"
	"context"
)

type AssetLegalRepo interface {
	Create(c context.Context, entity *domain.AssetLegal) error
	Update(c context.Context, id uint64, entity *domain.AssetLegal) error
	Delete(c context.Context, id uint64) error
	GetData(c context.Context) ([]domain.AssetLegal, error)
	GetByID(c context.Context, id uint64) (*domain.AssetLegal, error)
}
