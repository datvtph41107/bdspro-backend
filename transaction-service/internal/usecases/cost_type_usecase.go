package transaction_usecase

import (
	"context"
	"tx/internal/domain/entity"
	"tx/internal/dto"
	"tx/internal/repo"
)

type CostTypeUsecase interface {
	Create(ctx context.Context, costType *entity.CostType) error
	Update(ctx context.Context, costType *entity.CostType) error
	Delete(ctx context.Context, id uint64) error
	List(ctx context.Context, search *dto.CostTypeSearchDTO) ([]*entity.CostType, error)
}

type costTypeUsecase struct {
	repo repo.CostTypeRepository
}

func NewCostTypeUsecase(repo repo.CostTypeRepository) CostTypeUsecase {
	return &costTypeUsecase{repo: repo}
}

func (u *costTypeUsecase) Create(ctx context.Context, costType *entity.CostType) error {
	return u.repo.Create(ctx, costType)
}

func (u *costTypeUsecase) Update(ctx context.Context, costType *entity.CostType) error {
	return u.repo.Update(ctx, costType)
}

func (u *costTypeUsecase) Delete(ctx context.Context, id uint64) error {
	return u.repo.Delete(ctx, id)
}

func (u *costTypeUsecase) List(ctx context.Context, search *dto.CostTypeSearchDTO) ([]*entity.CostType, error) {
	return u.repo.List(ctx, search)
}
