package repo

import (
	"bdspro/internal/domain"
	"bdspro/internal/dto"
	"context"
)

type TxCostTypeRepository interface {
	Create(ctx context.Context, costType *domain.TxCostType) error
	Update(ctx context.Context, costType *domain.TxCostType) error
	Delete(ctx context.Context, id uint64) error
	List(ctx context.Context, search *dto.TxCostTypeSearchDTO) ([]*domain.TxCostType, error)
	GetByID(ctx context.Context, id uint64) (*domain.TxCostType, error)
}
