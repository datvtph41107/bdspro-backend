package usecases

import (
	"bdspro/internal/domain"
	"bdspro/internal/dto"
	"bdspro/internal/repo"
	"context"
)

type TxCostTypeUsecase interface {
	Create(ctx context.Context, costType *domain.TxCostType) error
	Update(ctx context.Context, costType *domain.TxCostType) error
	Delete(ctx context.Context, id uint64) error
	List(ctx context.Context, search *dto.TxCostTypeSearchDTO) ([]*domain.TxCostType, error)
}

type costTypeUsecase struct {
	repo repo.TxCostTypeRepository
}

func NewCostTypeUsecase(repo repo.TxCostTypeRepository) TxCostTypeUsecase {
	return &costTypeUsecase{repo: repo}
}

func (u *costTypeUsecase) Create(ctx context.Context, costType *domain.TxCostType) error {
	return u.repo.Create(ctx, costType)
}

func (u *costTypeUsecase) Update(ctx context.Context, costType *domain.TxCostType) error {
	return u.repo.Update(ctx, costType)
}

func (u *costTypeUsecase) Delete(ctx context.Context, id uint64) error {
	return u.repo.Delete(ctx, id)
}

func (u *costTypeUsecase) List(ctx context.Context, search *dto.TxCostTypeSearchDTO) ([]*domain.TxCostType, error) {
	return u.repo.List(ctx, search)
}
