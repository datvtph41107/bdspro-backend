package postgres

import (
	"bdspro/internal/domain"
	"context"

	"gorm.io/gorm"
)

// @bind: bdspro/internal/repo.HouseInfoRepo
type GormHouseInfoRepo struct {
	DB *gorm.DB
}

func NewGormHouseInfoRepo(db *gorm.DB) *GormHouseInfoRepo {
	return &GormHouseInfoRepo{
		DB: db,
	}
}

// ExistsByUserAndProduct kiểm tra xem UserId và ProductId có tồn tại không
func (repo *GormHouseInfoRepo) ExistsByUserAndProduct(userId, productId uint64) (bool, error) {
	var count int64
	err := repo.DB.Model(&domain.HouseInfo{}).
		Where("created_by = ? AND product_id = ? and deleted_at is null", userId, productId).
		Count(&count).Error

	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (repo *GormHouseInfoRepo) Create(access *domain.HouseInfo) error {
	return repo.DB.Create(&access).Error
}

func (r *GormHouseInfoRepo) UpdateHouseInfo(ctx context.Context, houseInfo *domain.HouseInfo) error {
	// Bắt đầu transaction để đảm bảo toàn vẹn dữ liệu
	// tx := r.DB.Begin()

	// houseInfo.ProductID = productID
	return GetDB(ctx, r.DB).Save(&houseInfo).Error
}
