package postgres

import (
	"context"

	"userv4/modules/payment/entity"
	"userv4/modules/payment/repository"

	"gorm.io/gorm"
)

type BankPostgres struct {
	db *gorm.DB
}

func NewBankPostgres(db *gorm.DB) *BankPostgres {
	return &BankPostgres{db: db}
}

func (r *BankPostgres) ListActive(ctx context.Context) ([]*entity.Bank, error) {
	var banks []*entity.Bank
	err := r.db.WithContext(ctx).
		Where("active = ?", true).
		Order("name asc").
		Find(&banks).Error
	if err != nil {
		return nil, err
	}
	return banks, nil
}

var _ repository.BankRepository = (*BankPostgres)(nil)
