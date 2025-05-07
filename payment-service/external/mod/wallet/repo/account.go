package a_repo

import (
	"payment/internal/wallet/domain"

	"gorm.io/gorm"
)

// @bind: payment/internal/wallet/repo.AccountRepo
type AccountGormRepo struct {
	DB *gorm.DB
}

func NewAccountGormRepo(db *gorm.DB) *AccountGormRepo {
	return &AccountGormRepo{
		DB: db,
	}
}

func (r AccountGormRepo) Create(account *domain.Account) (*domain.Account, error) {
	err := r.DB.Model(&domain.Account{}).Create(account).Error
	return account, err
}
func (r AccountGormRepo) Update(account *domain.Account) (*domain.Account, error) {
	return nil, nil
}
func (r AccountGormRepo) Lock(accountId uint64) error {
	return nil
}
func (r AccountGormRepo) ExistByID(userId uint64) (bool, error) {
	var count int64
	err := r.DB.Model(&domain.Account{}).
		Where("user_id = ? and deleted_at is null", userId).
		Count(&count).
		Error
	return count > 0, err
}
