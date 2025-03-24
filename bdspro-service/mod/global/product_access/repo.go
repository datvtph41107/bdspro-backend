package product_accecss

import (
	"bdspro/models"

	"gorm.io/gorm"
)

type ProductAccessRepo struct {
	DB *gorm.DB
}

func NewProductAccessRepo(db *gorm.DB) *ProductAccessRepo {
	return &ProductAccessRepo{
		DB: db,
	}
}

// ExistsByUserAndProduct kiểm tra xem UserId và ProductId có tồn tại không
func (repo *ProductAccessRepo) ExistsByUserAndProduct(userId, productId uint64) (bool, error) {
	var count int64
	err := repo.DB.Model(&models.ProductAccess{}).
		Where("user_id = ? AND product_id = ? and deleted_at is null", userId, productId).
		Count(&count).Error

	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (repo *ProductAccessRepo) Create(access *models.ProductAccess) error {
	return repo.DB.Create(&access).Error
}
