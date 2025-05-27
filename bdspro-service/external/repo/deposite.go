package postgre_repo

import (
	"bdspro/internal/domain"
	"context"
	"time"

	"gorm.io/gorm"
)

// @bind: bdspro/internal/repo.DepositeRepo
type GormDepositeRepo struct {
	DB *gorm.DB
}

func NewGormDepositeRepo(DB *gorm.DB) *GormDepositeRepo {
	return &GormDepositeRepo{
		DB: DB,
	}
}

func (r *GormDepositeRepo) GetDepositeByID(ctx context.Context, id uint64) (*domain.Deposite, error) {
	var deposite domain.Deposite
	err := r.DB.WithContext(ctx).Where("id = ?", id).First(&deposite).Error
	if err != nil {
		return nil, err
	}
	return &deposite, nil
}

func (r *GormDepositeRepo) CancelDeposite(ctx context.Context, id uint64) error {
	return r.DB.WithContext(ctx).
		Model(&domain.Deposite{}).
		Where("id = ? and deleted_at is null", id).
		Update("deleted_at", time.Now()).Error
}
