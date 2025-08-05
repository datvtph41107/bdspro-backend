package repository

import (
	"context"
	"organization/internal/domain/entity"
	"organization/internal/dto"
)

type CostTypeRepository interface {
	Create(ctx context.Context, costType *entity.CostType) error
	Update(ctx context.Context, costType *entity.CostType) error
	Delete(ctx context.Context, id uint64) error
	List(ctx context.Context, search *dto.CostTypeSearchDTO) ([]*entity.CostType, error)
	GetByID(ctx context.Context, id uint64) (*entity.CostType, error)
}
