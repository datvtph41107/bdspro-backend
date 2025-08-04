package repository

import (
	"context"
	"organization/internal/domain/entity"
	"organization/internal/dto"
	"organization/internal/enums"
)

type DealTransactionRepository interface {
	Create(ctx context.Context, transaction *entity.DealTransaction) error
	Update(ctx context.Context, transaction *entity.DealTransaction) error
	Delete(ctx context.Context, id uint64) error
	GetByID(ctx context.Context, id uint64) (*entity.DealTransaction, error)
	GetByIds(ctx context.Context, ids []uint64) ([]*entity.DealTransaction, error)
	GetByDealID(ctx context.Context, dealID uint64, transactionType *enums.TransactionType, page, limit int) ([]*entity.DealTransaction, int64, error)
	Approve(ctx context.Context, id uint64, approvedBy uint64, status enums.ApprovedStatus, rejectReason string) error
	
	// Statistics methods
	GetFinancialSummary(ctx context.Context, dealID uint64) (*dto.DealFinancialSummary, error)
	GetTotalAmountByType(ctx context.Context, dealID uint64, transactionType enums.TransactionType) (float64, error)
} 