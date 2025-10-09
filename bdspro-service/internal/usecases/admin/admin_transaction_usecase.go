package admin_usecases

import (
	"bdspro/internal/domain"
	admin_repo "bdspro/internal/repo/admin"
	_utils "common/utils"
	"context"
	"fmt"
)

type AdminTransactionUsecase struct {
	AdminTransactionRepo admin_repo.IAdminTransactionRepo
}

func NewAdminTransactionUsecase(
	adminTransactionRepo admin_repo.IAdminTransactionRepo,
) *AdminTransactionUsecase {
	return &AdminTransactionUsecase{
		AdminTransactionRepo: adminTransactionRepo,
	}
}

func (uc *AdminTransactionUsecase) GetList(ctx context.Context, filter *admin_repo.TransactionFilter) ([]*domain.Transaction, int64, error) {
	// Validate pagination
	if filter.Page < 0 {
		filter.Page = 0
	}
	if filter.Size <= 0 {
		filter.Size = 20
	}
	if filter.Size > 100 {
		filter.Size = 100
	}

	return uc.AdminTransactionRepo.GetList(ctx, filter)
}

func (uc *AdminTransactionUsecase) GetDetail(ctx context.Context, id uint64) (*domain.Transaction, error) {
	transaction, err := uc.AdminTransactionRepo.GetDetail(ctx, id)
	if err != nil {
		return nil, err
	}

	if transaction == nil {
		return nil, fmt.Errorf("không tìm thấy thương vụ")
	}

	return transaction, nil
}

func (uc *AdminTransactionUsecase) ApproveTransaction(ctx context.Context, id uint64) error {
	// Get current transaction
	transaction, err := uc.AdminTransactionRepo.GetDetail(ctx, id)
	if err != nil {
		return err
	}

	if transaction == nil {
		return fmt.Errorf("không tìm thấy thương vụ")
	}

	// Get profile ID from context
	profileID := _utils.GetProfileIdWithContext(ctx)
	approvedBy := uint32(profileID)

	// Update approval status
	return uc.AdminTransactionRepo.UpdateApprovalStatus(ctx, id, "approved", &approvedBy)
}

func (uc *AdminTransactionUsecase) RejectTransaction(ctx context.Context, id uint64) error {
	// Get current transaction
	transaction, err := uc.AdminTransactionRepo.GetDetail(ctx, id)
	if err != nil {
		return err
	}

	if transaction == nil {
		return fmt.Errorf("không tìm thấy thương vụ")
	}

	// Get profile ID from context
	profileID := _utils.GetProfileIdWithContext(ctx)
	approvedBy := uint32(profileID)

	// Update approval status
	return uc.AdminTransactionRepo.UpdateApprovalStatus(ctx, id, "rejected", &approvedBy)
}

func (uc *AdminTransactionUsecase) DeleteTransaction(ctx context.Context, id uint64) error {
	// Get current transaction
	transaction, err := uc.AdminTransactionRepo.GetDetail(ctx, id)
	if err != nil {
		return err
	}

	if transaction == nil {
		return fmt.Errorf("không tìm thấy thương vụ")
	}

	return uc.AdminTransactionRepo.Delete(ctx, id)
}
