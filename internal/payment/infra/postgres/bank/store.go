package bankpostgres

import (
	"context"

	"payment/internal/domain/bank"

	"gorm.io/gorm"
)

// @bind: payment/internal/domain/bankuc.Repository
type BankPostgresRepository struct {
	db *gorm.DB
}

func NewBankPostgresRepository(db *gorm.DB) *BankPostgresRepository {
	return &BankPostgresRepository{
		db: db,
	}
}

func (r *BankPostgresRepository) GetAll(ctx context.Context) ([]*bank.Bank, error) {
	var banks []*bank.Bank
	if err := r.db.WithContext(ctx).Find(&banks).Error; err != nil {
		return nil, err
	}
	return banks, nil
}

func (r *BankPostgresRepository) GetByID(ctx context.Context, id uint32) (*bank.Bank, error) {
	var bank bank.Bank
	if err := r.db.WithContext(ctx).First(&bank, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &bank, nil
}

func (r *BankPostgresRepository) Create(ctx context.Context, bank *bank.Bank) (*bank.Bank, error) {
	if err := r.db.WithContext(ctx).Create(bank).Error; err != nil {
		return nil, err
	}
	return bank, nil
}

func (r *BankPostgresRepository) Update(ctx context.Context, bank *bank.Bank) (*bank.Bank, error) {
	if err := r.db.WithContext(ctx).Save(bank).Error; err != nil {
		return nil, err
	}
	return bank, nil
}

func (r *BankPostgresRepository) Delete(ctx context.Context, id uint32) error {
	return r.db.WithContext(ctx).Delete(&bank.Bank{}, id).Error
}
