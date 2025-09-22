package postgres

import (
	"gorm.io/gorm"
)

type GormDepositeRepo struct {
	DB *gorm.DB
}

// func NewGormDepositeRepo(DB *gorm.DB) repo.DepositeRepo {
// 	return &GormDepositeRepo{
// 		DB: DB,
// 	}
// }

// func (r *GormDepositeRepo) GetDepositeByID(ctx context.Context, id uint64) (*domain.Deposite, error) {
// 	var deposite domain.Deposite
// 	err := r.DB.WithContext(ctx).Where("id = ?", id).First(&deposite).Error
// 	if err != nil {
// 		return nil, err
// 	}
// 	return &deposite, nil
// }
