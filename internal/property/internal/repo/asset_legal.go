package repo

import (
	"bdspro/internal/domain"
	"bdspro/internal/dto"
	"context"
)

type AssetLegalRepo interface {
	Create(c context.Context, entity *domain.AssetLegal) error
	Update(c context.Context, id uint64, entity *domain.AssetLegal) error
	Delete(c context.Context, id uint64) error
	GetData(c context.Context, dto dto.AssetLegalGetDTO) ([]domain.AssetLegal, int64, error)
	GetByID(c context.Context, id uint64) (*domain.AssetLegal, error)
	UpdateLegalItems(c context.Context, assetID uint64, items []domain.AssetLegal) error
}
