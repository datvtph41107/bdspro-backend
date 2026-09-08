package repo

import (
	"bdspro/internal/domain"
	"bdspro/internal/dto"
	"context"
)

type TxDealCostRepository interface {
	Create(ctx context.Context, dealCost *domain.TxDealCost) error
	Update(ctx context.Context, dealCost *domain.TxDealCost) error
	Delete(ctx context.Context, id uint64) error
	GetByID(ctx context.Context, id uint64) (*domain.TxDealCost, error)
	ListByDealID(ctx context.Context, search *dto.TxDealCostSearchDTO) ([]*domain.TxDealCost, error)
	GetStatisticsByDealID(ctx context.Context, dealID uint64) (*dto.TxDealCostStatisticsDTO, error)
	GetStatisticsByDealIDs(ctx context.Context, dealIDs []uint64) (*dto.TxStatisticDealsDTO, error)
}
