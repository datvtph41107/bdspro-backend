package walletpostgres

import (
	walletdomain "payment/internal/domain/wallet"

	"gorm.io/gorm"
)

// AutoMigrateModels là registry duy nhất của phần Wallet/Bank compatibility.
func AutoMigrateModels() []any {
	return []any{
		&PaymentMethodModel{},
		&WalletModel{},
		&WalletTransactionModel{},
		&TransactionTypeModel{},
		&WalletAuditLogModel{},
		&WithdrawalRequestModel{},
		&walletdomain.DashboardMetric{},
	}
}

func AutoMigrate(database *gorm.DB) error {
	return database.AutoMigrate(AutoMigrateModels()...)
}
