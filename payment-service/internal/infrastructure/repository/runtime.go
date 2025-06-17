package postgres

import (
	"log"

	"github.com/tuandq2112/payment-service/internal/domain/repository"
	"github.com/tuandq2112/payment-service/internal/infrastructure/repository/postgres"
	"gorm.io/gorm"
)

func NewPaymentRepository(db *gorm.DB) repository.PaymentMethodRepository {
	err := db.AutoMigrate(&postgres.PaymentMethodModel{})
	if err != nil {
		log.Fatalf("Failed to migrate payment method table: %v", err)
	}
	return postgres.NewPaymentMethodRepository(db)
}

func NewWalletRepository(db *gorm.DB) repository.WalletsRepository {
	err := db.AutoMigrate(&postgres.WalletModel{})
	if err != nil {
		log.Fatalf("Failed to migrate wallet table: %v", err)
	}
	return postgres.NewWalletsRepository(db)
}

func NewWalletTransactionRepository(db *gorm.DB) repository.WalletTransactionsRepository {
	err := db.AutoMigrate(&postgres.WalletTransactionModel{})
	if err != nil {
		log.Fatalf("Failed to migrate wallet transaction table: %v", err)
	}
	return postgres.NewWalletTransactionRepository(db)
}

func NewWalletAuditLogRepository(db *gorm.DB) repository.WalletAuditLogRepository {
	err := db.AutoMigrate(&postgres.WalletAuditLogModel{})
	if err != nil {
		log.Fatalf("Failed to migrate wallet audit log table: %v", err)
	}
	return postgres.NewWalletAuditLogRepository(db)
}

func NewWithdrawalRequestRepository(db *gorm.DB) repository.WithdrawalRequestsRepository {
	err := db.AutoMigrate(&postgres.WithdrawalRequestModel{})
	if err != nil {
		log.Fatalf("Failed to migrate withdrawal request table: %v", err)
	}
	return postgres.NewWithdrawalRequestsRepository(db)
}