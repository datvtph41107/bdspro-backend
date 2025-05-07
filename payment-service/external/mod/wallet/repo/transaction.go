package a_repo

import (
	"database/sql"
	"payment/internal/wallet/domain"

	"gorm.io/gorm"
)

// @bind: payment/internal/wallet/repo.TransactionRepo
type TransactionGormRepo struct {
	DB *gorm.DB
}

func NewTransactionRepo(db *gorm.DB) *TransactionGormRepo {
	return &TransactionGormRepo{
		DB: db,
	}
}

func (r *TransactionGormRepo) GetByID(id uint64) (*domain.Transaction, error) {
	return nil, nil
}
func (r *TransactionGormRepo) ExistByID(id uint64) (bool, error) {
	var count int64
	err := r.DB.Model(&domain.Transaction{}).
		Where("transaction_id = ? and deleted_at is null", &id).
		Count(&count).
		Error

	return count > 0, err
}

func (r *TransactionGormRepo) Create(transaction *domain.Transaction) (*domain.Transaction, error) {
	err := r.DB.Model(&domain.Transaction{}).Create(transaction).Error
	return transaction, err
}
func (r *TransactionGormRepo) Balance(accountId *uint64) (float64, error) {
	var sum sql.NullFloat64
	err := r.DB.Model(&domain.Transaction{}).
		Select("SUM(amount)").
		Where("account_id = ? and deleted_at is null", &accountId).
		Scan(&sum).
		Error
	if sum.Valid {
		return sum.Float64, err
	}
	return 0, nil
}
