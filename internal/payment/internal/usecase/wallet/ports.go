package wallet

import (
	"context"
	"time"

	walletdomain "payment/internal/domain/wallet"
	paymentpb "pb/types/payment"
)

// Transaction is a caller-owned database transaction. Context binds every
// participating persistence adapter to the same database handle; this avoids
// the historical fake transaction where a tx was stored in context but
// repositories silently continued using their root DB handle.
type Transaction interface {
	Commit() error
	Rollback() error
	Context(context.Context) context.Context
}

type PaymentMethodRepository interface {
	GetPaymentMethods(ctx context.Context, organizationID uint32, page, size int) ([]*paymentpb.PaymentMethod, int64, error)
	CreatePaymentMethod(ctx context.Context, method *walletdomain.PaymentMethod) (*walletdomain.PaymentMethod, error)
	GetPaymentMethodById(ctx context.Context, id uint32) (*walletdomain.PaymentMethod, error)
	UpdatePaymentMethod(ctx context.Context, method *walletdomain.PaymentMethod) (*walletdomain.PaymentMethod, error)
	DeletePaymentMethod(ctx context.Context, id uint32) error
}

type WalletRepository interface {
	GetWalletById(ctx context.Context, id uint32) (*walletdomain.Wallet, error)
	GetWalletByIdForUpdate(ctx context.Context, id uint32) (*walletdomain.Wallet, error)
	UpdateWallet(ctx context.Context, wallet *walletdomain.Wallet) error
	BeginTx(ctx context.Context) Transaction
	CreateWallet(ctx context.Context, wallet *walletdomain.Wallet) (*walletdomain.Wallet, error)
	GetWalletByUserId(ctx context.Context, userID uint32) (*walletdomain.Wallet, error)
}

type WalletTransactionRepository interface {
	GetTransactions(ctx context.Context, walletID uint32, req *paymentpb.GetWalletTransactionsRequest) ([]*paymentpb.WalletTransaction, int64, error)
	GetTotalAmountByType(ctx context.Context, walletID uint32, transactionType string) (int64, error)
	GetTransactionCountByStatus(ctx context.Context, walletID uint32, status string) (int64, error)
	CreateTransaction(ctx context.Context, transaction *walletdomain.WalletTransaction) error
	UpdateTransaction(ctx context.Context, transaction *walletdomain.WalletTransaction) error
	GetTransactionByDealId(ctx context.Context, dealID uint32) (*walletdomain.WalletTransaction, error)
	GetTransactionByCode(ctx context.Context, code string) (*walletdomain.WalletTransaction, error)
	GetTransactionsReport(ctx context.Context, req *paymentpb.ReportRequest) ([]*paymentpb.WalletTransaction, int64, int64, error)
	ExportTransactions(ctx context.Context, req *paymentpb.ReportRequest) ([]byte, error)
	GetDealTransactions(ctx context.Context, req *paymentpb.GetDealTransactionsRequest) ([]*paymentpb.DealTransaction, int64, error)
	GetDealTransactionById(ctx context.Context, id uint32) (*paymentpb.DealTransaction, error)
	UpdateDealTransaction(ctx context.Context, transaction *paymentpb.DealTransaction) error
	GetTransactionTypes(ctx context.Context) ([]*paymentpb.TransactionType, error)
	CreateTransactionType(ctx context.Context, transactionType *paymentpb.TransactionType) error
	UpdateTransactionType(ctx context.Context, transactionType *paymentpb.TransactionType) error
	DeleteTransactionType(ctx context.Context, id uint32) error
	CountTransactionsByStatus(ctx context.Context, status string) (int64, error)
}

type WalletAuditRepository interface {
	CreateAuditLog(ctx context.Context, log *walletdomain.WalletAuditLog) error
}

type WithdrawalRepository interface {
	CreateWithdrawalRequest(ctx context.Context, request *walletdomain.WithdrawalRequest) error
	GetWithdrawalRequestById(ctx context.Context, id uint32) (*walletdomain.WithdrawalRequest, error)
	UpdateWithdrawalRequest(ctx context.Context, request *walletdomain.WithdrawalRequest) error
}

type DashboardMetricRepository interface {
	CreateOrUpdate(ctx context.Context, metric *walletdomain.DashboardMetric) (*walletdomain.DashboardMetric, error)
	GetDashboardMetrics(ctx context.Context, fromDate, toDate time.Time) ([]*walletdomain.DashboardMetric, error)
}
