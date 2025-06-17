package repository

import (
	"context"

	paymentpb "pb/types/payment"

	"github.com/tuandq2112/payment-service/internal/domain/entity"
)

type WalletTransactionsRepository interface {
	GetTransactions(ctx context.Context, walletId uint32, filter *paymentpb.TransactionFilter, pagination *paymentpb.PaginationRequest) ([]*paymentpb.WalletTransaction, int64, error)
	GetTotalAmountByType(ctx context.Context, walletId uint32, transactionType string) (float64, error)
	GetTransactionCountByStatus(ctx context.Context, walletId uint32, status string) (int64, error)
	CreateTransaction(ctx context.Context, transaction *entity.WalletTransaction) error
	UpdateTransaction(ctx context.Context, transaction *entity.WalletTransaction) error
	GetTransactionByDealId(ctx context.Context, dealId uint32) (*entity.WalletTransaction, error)
	GetTransactionsReport(ctx context.Context, req *paymentpb.ReportRequest) ([]*paymentpb.WalletTransaction, int64, float64, error)
	ExportTransactions(ctx context.Context, req *paymentpb.ReportRequest) ([]byte, error)
	GetDealTransactions(ctx context.Context, filter *paymentpb.TransactionFilter, pagination *paymentpb.PaginationRequest) ([]*paymentpb.DealTransaction, int64, error)
	GetDealTransactionById(ctx context.Context, id uint32) (*paymentpb.DealTransaction, error)
	UpdateDealTransaction(ctx context.Context, transaction *paymentpb.DealTransaction) error
	GetTransactionTypes(ctx context.Context) ([]*paymentpb.TransactionType, error)
	CreateTransactionType(ctx context.Context, transactionType *paymentpb.TransactionType) error
	UpdateTransactionType(ctx context.Context, transactionType *paymentpb.TransactionType) error
	DeleteTransactionType(ctx context.Context, id uint32) error
}