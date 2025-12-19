package usecases

import (
	"bdspro/internal/domain"
	"bdspro/internal/dto"
	"bdspro/internal/enums"
	"bdspro/internal/repo"
	"context"
)

type DealCostUsecase interface {
	Create(ctx context.Context, dealCost *domain.TxDealCost) error
	Update(ctx context.Context, dealCost *domain.TxDealCost) error
	Delete(ctx context.Context, id uint64) error
	GetByID(ctx context.Context, id uint64) (*domain.TxDealCost, error)
	ListByDealID(ctx context.Context, search *dto.TxDealCostSearchDTO) ([]*domain.TxDealCost, error)
}

type dealCostUsecase struct {
	repo         repo.TxDealCostRepository
	costTypeRepo repo.TxCostTypeRepository
}

func NewDealCostUsecase(repo repo.TxDealCostRepository, costTypeRepo repo.TxCostTypeRepository) DealCostUsecase {
	return &dealCostUsecase{repo: repo, costTypeRepo: costTypeRepo}
}

func (u *dealCostUsecase) Create(ctx context.Context, dealCost *domain.TxDealCost) error {
	costType, err := u.costTypeRepo.GetByID(ctx, dealCost.CostTypeID)
	if err != nil {
		return err
	}
	if costType != nil && costType.CostType == enums.TxCostTypeExpense {
		dealCost.Status = enums.TxApprovedStatusApproved
	} else {
		dealCost.Status = enums.TxApprovedStatusPending
	}

	return u.repo.Create(ctx, dealCost)
}

func (u *dealCostUsecase) Update(ctx context.Context, dealCost *domain.TxDealCost) error {
	return u.repo.Update(ctx, dealCost)
}

func (u *dealCostUsecase) Delete(ctx context.Context, id uint64) error {
	return u.repo.Delete(ctx, id)
}

func (u *dealCostUsecase) GetByID(ctx context.Context, id uint64) (*domain.TxDealCost, error) {
	return u.repo.GetByID(ctx, id)
}

func (u *dealCostUsecase) ListByDealID(ctx context.Context, search *dto.TxDealCostSearchDTO) ([]*domain.TxDealCost, error) {
	return u.repo.ListByDealID(ctx, search)
}
