package usecases

import (
	_utils "common/utils"
	"context"
	"errors"
	"time"

	"bdspro/internal/domain/transaction"

	"bdspro/internal/repo"
)

var (
	ErrNotFound     = errors.New("transaction not found")
	ErrInvalidState = errors.New("invalid transaction state")
)

type TransactionUsecase interface {
	CreateTransaction(ctx context.Context, tx *transaction.Transaction) (*transaction.Transaction, error)
	GetTransaction(ctx context.Context, id uint32) (*transaction.Transaction, error)
	UpdateTransaction(ctx context.Context, tx *transaction.Transaction) (*transaction.Transaction, error)
	DeleteTransaction(ctx context.Context, id uint32) error
	ListTransactions(ctx context.Context, filters *repo.TransactionFilters) ([]*transaction.Transaction, int64, error)
	ApproveTransaction(ctx context.Context, id uint32, approvedBy uint32) error
	RejectTransaction(ctx context.Context, id uint32) error
	GetTransactionTypes(ctx context.Context) ([]string, error)
}

type transactionUsecase struct {
	transactionRepo repo.TransactionRepository
}

func NewTransactionUsecase(transactionRepo repo.TransactionRepository) TransactionUsecase {
	return &transactionUsecase{
		transactionRepo: transactionRepo,
	}
}

func (uc *transactionUsecase) CreateTransaction(ctx context.Context, tx *transaction.Transaction) (*transaction.Transaction, error) {
	// Set default values
	tx.TransactionStatus = transaction.TransactionStatusDraft
	tx.ApprovalStatus = transaction.ApprovalStatusPending
	tx.IsDeleted = false
	tx.CreatedAt = time.Now()
	tx.UpdatedAt = time.Now()
	tx.OwnerId = uint32(_utils.GetProfileIdWithContext(ctx))
	entity, err := uc.transactionRepo.Create(ctx, tx)
	if err != nil {
		return nil, err
	}
	return entity, nil
}

func (uc *transactionUsecase) GetTransaction(ctx context.Context, id uint32) (*transaction.Transaction, error) {
	tx, err := uc.transactionRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if tx == nil {
		return nil, ErrNotFound
	}
	return tx, nil
}

func (uc *transactionUsecase) UpdateTransaction(ctx context.Context, tx *transaction.Transaction) (*transaction.Transaction, error) {
	// Get existing transaction
	existingTx, err := uc.transactionRepo.GetByID(ctx, tx.Id)
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
	existingTx.UpdatedAt = time.Now()

	entity, err := uc.transactionRepo.Update(ctx, existingTx)
	if err != nil {
		return nil, err
	}
	return entity, nil
}

func (uc *transactionUsecase) DeleteTransaction(ctx context.Context, id uint32) error {
	// Check if transaction exists
	existingTx, err := uc.transactionRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if existingTx == nil {
		return ErrNotFound
	}

	return uc.transactionRepo.Delete(ctx, id)
}

func (uc *transactionUsecase) ListTransactions(ctx context.Context, filters *repo.TransactionFilters) ([]*transaction.Transaction, int64, error) {
	return uc.transactionRepo.List(ctx, filters)
}

func (uc *transactionUsecase) ApproveTransaction(ctx context.Context, id uint32, approvedBy uint32) error {
	// Get existing transaction
	tx, err := uc.transactionRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if tx == nil {
		return ErrNotFound
	}

	// Check if transaction can be approved
	if tx.ApprovalStatus != transaction.ApprovalStatusPending {
		return ErrInvalidState
	}

	err = uc.transactionRepo.Approve(ctx, id, approvedBy)
	if err != nil {
		return err
	}
	return nil
}

func (uc *transactionUsecase) RejectTransaction(ctx context.Context, id uint32) error {
	// Get existing transaction
	tx, err := uc.transactionRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if tx == nil {
		return ErrNotFound
	}

	// Check if transaction can be rejected
	if tx.ApprovalStatus != transaction.ApprovalStatusPending {
		return ErrInvalidState
	}

	err = uc.transactionRepo.Reject(ctx, id)
	if err != nil {
		return err
	}
	return nil
}

func (uc *transactionUsecase) GetTransactionTypes(ctx context.Context) ([]string, error) {
	return uc.transactionRepo.GetTransactionTypes(ctx)
}
