package usecases

import (
	_utils "common/utils"
	"context"
	"errors"

	"bdspro/internal/domain"
	"bdspro/internal/enums"
	"bdspro/internal/repo"
	iusecase "bdspro/internal/usecases/interface"
)

var (
	ErrNotFound     = errors.New("transaction not found")
	ErrInvalidState = errors.New("invalid transaction state")
)

type TransactionUsecase interface {
	CreateTransaction(ctx context.Context, tx *domain.Transaction) (*domain.Transaction, error)
	GetTransaction(ctx context.Context, id uint64) (*domain.Transaction, error)
	UpdateTransaction(ctx context.Context, tx *domain.Transaction) (*domain.Transaction, error)
	DeleteTransaction(ctx context.Context, id uint64) error
	ListTransactions(ctx context.Context, filters *repo.TransactionFilters) ([]*domain.Transaction, int64, error)
	ApproveTransaction(ctx context.Context, id uint64, approvedBy uint32) error
	RejectTransaction(ctx context.Context, id uint64) error
	GetTransactionTypes(ctx context.Context) ([]string, error)
}

type transactionUsecase struct {
	transactionRepo repo.TransactionRepository
	Transaction     iusecase.ITransaction
}

func NewTransactionUsecase(
	transactionRepo repo.TransactionRepository,
	transaction iusecase.ITransaction,
) TransactionUsecase {
	return &transactionUsecase{
		transactionRepo: transactionRepo,
		Transaction:     transaction,
	}
}

func (uc *transactionUsecase) CreateTransaction(ctx context.Context, tx *domain.Transaction) (*domain.Transaction, error) {
	// Set default values
	tx.TransactionStatus = enums.TransactionStatusDraft
	tx.ApprovalStatus = enums.ApprovalStatusPending
	tx.OwnerId = _utils.GetProfileIdWithContext(ctx)
	entity, err := uc.transactionRepo.Create(ctx, tx)
	if err != nil {
		return nil, err
	}
	return entity, nil
}

func (uc *transactionUsecase) GetTransaction(ctx context.Context, id uint64) (*domain.Transaction, error) {
	tx, err := uc.transactionRepo.GetByID(ctx, &id)
	if err != nil {
		return nil, err
	}
	if tx == nil {
		return nil, ErrNotFound
	}
	return tx, nil
}

func (uc *transactionUsecase) UpdateTransaction(ctx context.Context, tx *domain.Transaction) (*domain.Transaction, error) {
	// Get existing transaction
	existingTx, err := uc.transactionRepo.GetByID(ctx, &tx.ID)
	if err != nil {
		return nil, err
	}
	if existingTx == nil {
		return nil, ErrNotFound
	}

	// Update only allowed fields
	existingTx.Amount = tx.Amount
	existingTx.Currency = tx.Currency
	existingTx.TransactionName = tx.TransactionName
	existingTx.Description = tx.Description
	existingTx.CategoryId = tx.CategoryId
	existingTx.PaymentMethodId = tx.PaymentMethodId
	existingTx.RelatedDealId = tx.RelatedDealId
	existingTx.RelatedTransactionId = tx.RelatedTransactionId
	existingTx.TransactionDate = tx.TransactionDate

	entity, err := uc.transactionRepo.Update(ctx, existingTx)
	if err != nil {
		return nil, err
	}
	return entity, nil
}

func (uc *transactionUsecase) DeleteTransaction(ctx context.Context, id uint64) error {
	// Check if transaction exists
	existingTx, err := uc.transactionRepo.GetByID(ctx, &id)
	if err != nil {
		return err
	}
	if existingTx == nil {
		return ErrNotFound
	}

	return uc.transactionRepo.Delete(ctx, uint32(id))
}

func (uc *transactionUsecase) ListTransactions(ctx context.Context, filters *repo.TransactionFilters) ([]*domain.Transaction, int64, error) {
	return uc.transactionRepo.List(ctx, filters)
}

func (uc *transactionUsecase) ApproveTransaction(ctx context.Context, id uint64, approvedBy uint32) error {
	// Get existing transaction
	tx, err := uc.transactionRepo.GetByID(ctx, &id)
	if err != nil {
		return err
	}
	if tx == nil {
		return ErrNotFound
	}

	// Check if transaction can be approved
	if tx.ApprovalStatus != enums.ApprovalStatusPending {
		return ErrInvalidState
	}

	err = uc.transactionRepo.Approve(ctx, uint32(id), approvedBy)
	if err != nil {
		return err
	}
	return nil
}

func (uc *transactionUsecase) RejectTransaction(ctx context.Context, id uint64) error {
	// Get existing transaction
	tx, err := uc.transactionRepo.GetByID(ctx, &id)
	if err != nil {
		return err
	}
	if tx == nil {
		return ErrNotFound
	}

	// Check if transaction can be rejected
	if tx.ApprovalStatus != enums.ApprovalStatusPending {
		return ErrInvalidState
	}

	err = uc.transactionRepo.Reject(ctx, uint32(id))
	if err != nil {
		return err
	}
	return nil
}

func (uc *transactionUsecase) GetTransactionTypes(ctx context.Context) ([]string, error) {
	return uc.transactionRepo.GetTransactionTypes(ctx)
}
