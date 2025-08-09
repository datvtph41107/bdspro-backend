package postgres

import (
	"context"

	"payment/internal/domain/entity"

	"gorm.io/gorm"
)

// @bind: payment/internal/domain/repository.BankRepository
type BankPostgresRepository struct {
	db *gorm.DB
}

func NewBankPostgresRepository(db *gorm.DB) *BankPostgresRepository {
	return &BankPostgresRepository{
		db: db,
	}
}

func (r *BankPostgresRepository) GetAll(ctx context.Context) ([]*entity.Bank, error) {
	var banks []*entity.Bank
	if err := r.db.WithContext(ctx).Find(&banks).Error; err != nil {
		return nil, err
	}
	return banks, nil
}

func (r *BankPostgresRepository) GetByID(ctx context.Context, id uint32) (*entity.Bank, error) {
	var bank entity.Bank
	if err := r.db.WithContext(ctx).First(&bank, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &bank, nil
}

func (r *BankPostgresRepository) Create(ctx context.Context, bank *entity.Bank) (*entity.Bank, error) {
	if err := r.db.WithContext(ctx).Create(bank).Error; err != nil {
		return nil, err
	}
	return bank, nil
}

func (r *BankPostgresRepository) Update(ctx context.Context, bank *entity.Bank) (*entity.Bank, error) {
	if err := r.db.WithContext(ctx).Save(bank).Error; err != nil {
		return nil, err
	}
	return bank, nil
}

func (r *BankPostgresRepository) Delete(ctx context.Context, id uint32) error {
	return r.db.WithContext(ctx).Delete(&entity.Bank{}, id).Error
}
