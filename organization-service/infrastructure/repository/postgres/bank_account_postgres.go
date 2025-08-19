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

func (r *BankAccountPostgresRepository) GetByID(ctx context.Context, id *uint64) (*entity.BankAccount, error) {
	var bankAccount entity.BankAccount
	if err := r.db.WithContext(ctx).
		Where("id = ? AND deleted_at IS NULL", id).
		First(&bankAccount).Error; err != nil {
		return nil, err
	}
	return &bankAccount, nil
}

func (r *BankAccountPostgresRepository) DealUpdateAccount(ctx context.Context, accountId *uint64, model *entity.BankAccount) (*entity.BankAccount, error) {
	tx := GetDB(ctx, r.db)
	if accountId == nil {
		// Nếu chưa có bankAccountId thì tạo mới bank account
		bankAccount := &entity.BankAccount{
			BankName:      model.BankName,
			AccountNumber: model.AccountNumber,
			AccountName:   model.AccountName,
		}
		if err := tx.Create(bankAccount).Error; err != nil {
			return nil, err
		}
		// Cập nhật lại deal với bankAccountId mới
		return bankAccount, nil
	}
	// Nếu đã có bankAccountId thì update thông tin bank account
	err := tx.Model(&entity.BankAccount{}).
		Where("id = ? AND deleted_at is null", accountId).
		Updates(map[string]interface{}{
			"bank_name":      model.BankName,
			"account_number": model.AccountNumber,
			"account_name":   model.AccountName,
			"bank_id":        model.BankId,
		}).Error
	if err != nil {
		return nil, err
	}

	bankAccount, err := r.GetByID(ctx, accountId)
	if err != nil {
		return nil, err
	}
	return bankAccount, nil
}
