package usecase

import (
	"context"
	"organization/internal/domain/entity"
	"organization/internal/domain/repository"
	"organization/internal/dto"
)

type CostTypeUsecase interface {
	Create(ctx context.Context, costType *entity.CostType) error
	Update(ctx context.Context, costType *entity.CostType) error
	Delete(ctx context.Context, id uint64) error
	List(ctx context.Context, search *dto.CostTypeSearchDTO) ([]*entity.CostType, error)
}

type costTypeUsecase struct {
	repo repository.CostTypeRepository
}

func NewCostTypeUsecase(repo repository.CostTypeRepository) CostTypeUsecase {
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
