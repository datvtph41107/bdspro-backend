package postgres

import (
	"context"
	"organization/internal/domain/entity"

	"gorm.io/gorm"
)

// @bind: organization/internal/domain/repository.BankAccountRepository
type BankAccountPostgresRepository struct {
	db *gorm.DB
}

func NewBankAccountPostgresRepository(db *gorm.DB) *BankAccountPostgresRepository {
	return &BankAccountPostgresRepository{db: db}
}

func (r *BankAccountPostgresRepository) Create(ctx context.Context, bankAccount *entity.BankAccount) (*entity.BankAccount, error) {
	if err := r.db.WithContext(ctx).Create(bankAccount).Error; err != nil {
		return nil, err
	}
	return bankAccount, nil
}
