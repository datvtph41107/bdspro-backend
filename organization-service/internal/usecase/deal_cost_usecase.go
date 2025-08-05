package usecase

import (
	"context"
	"organization/internal/domain/entity"
	"organization/internal/domain/repository"
	"organization/internal/dto"
	"organization/internal/enums"
)

type DealCostUsecase interface {
	Create(ctx context.Context, dealCost *entity.DealCost) error
	Update(ctx context.Context, dealCost *entity.DealCost) error
	Delete(ctx context.Context, id uint64) error
	GetByID(ctx context.Context, id uint64) (*entity.DealCost, error)
	ListByDealID(ctx context.Context, search *dto.DealCostSearchDTO) ([]*entity.DealCost, error)
}

type dealCostUsecase struct {
	repo         repository.DealCostRepository
	costTypeRepo repository.CostTypeRepository
}

func NewDealCostUsecase(repo repository.DealCostRepository, costTypeRepo repository.CostTypeRepository) DealCostUsecase {
	return &dealCostUsecase{repo: repo}
}

func (u *dealCostUsecase) Create(ctx context.Context, dealCost *entity.DealCost) error {
	costType, err := u.costTypeRepo.GetByID(ctx, dealCost.CostTypeID)
	if err != nil {
		return err
	}
	dealCost.Status = enums.ApprovedStatusApproved
	if costType.CostType == enums.ECostTypeExpense {
		dealCost.Status = enums.ApprovedStatusApproved
	} else {
		dealCost.Status = enums.ApprovedStatusPending
	}
	return u.repo.Create(ctx, dealCost)
}

func (u *dealCostUsecase) Update(ctx context.Context, dealCost *entity.DealCost) error {
	return u.repo.Update(ctx, dealCost)
}

func (u *dealCostUsecase) Delete(ctx context.Context, id uint64) error {
	return u.repo.Delete(ctx, id)
}

func (u *dealCostUsecase) GetByID(ctx context.Context, id uint64) (*entity.DealCost, error) {
	return u.repo.GetByID(ctx, id)
}

func (u *dealCostUsecase) ListByDealID(ctx context.Context, search *dto.DealCostSearchDTO) ([]*entity.DealCost, error) {
	return u.repo.ListByDealID(ctx, search)
}
