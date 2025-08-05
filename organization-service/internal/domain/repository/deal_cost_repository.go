package repository

import (
	"context"
	"organization/internal/domain/entity"
	"organization/internal/dto"
)

type DealCostRepository interface {
	Create(ctx context.Context, dealCost *entity.DealCost) error
	Update(ctx context.Context, dealCost *entity.DealCost) error
	Delete(ctx context.Context, id uint64) error
	GetByID(ctx context.Context, id uint64) (*entity.DealCost, error)
	ListByDealID(ctx context.Context, search *dto.DealCostSearchDTO) ([]*entity.DealCost, error)
}
