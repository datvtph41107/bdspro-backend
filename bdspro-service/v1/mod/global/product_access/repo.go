package product_accecss

import (
	"bdspro/v1/models"
	"common/case/crud2"

	"gorm.io/gorm"
)

type ProductAccessRepo struct {
	crud2.BaseRepo[models.ProductAccess]
	// PrivateFieldRepo *product_private.ProductPrivateRepo
}

func NewProductAccessRepo(db *gorm.DB) *ProductAccessRepo {
	return &ProductAccessRepo{
		BaseRepo: crud2.BaseRepo[models.ProductAccess]{DB: db},
		// PrivateFieldRepo: PrivateFieldRepo,
	}
}

// ExistsByUserAndProduct kiểm tra xem UserId và ProductId có tồn tại không
func (repo *ProductAccessRepo) ExistsByUserAndProduct(userId, productId uint64) (bool, error) {
	var count int64
	err := repo.DB.Model(&models.ProductAccess{}).
		Where("target_type = 1 AND target_id = ? AND product_id = ? and deleted_at is null", userId, productId).
		Count(&count).Error

	if err != nil {
		return false, err
	}

	return count > 0, nil
}

// ExistsByUserAndProduct kiểm tra xem UserId và ProductId có tồn tại không
func (repo *ProductAccessRepo) GetByUserAndProduct(userId, productId uint64) (*models.ProductAccess, error) {
	var entity *models.ProductAccess
	err := repo.DB.
		Where("target_type = 1 AND target_id = ? AND product_id = ? and deleted_at is null", userId, productId).
		First(&entity).Error

	if err != nil {
		return nil, err
	}

	return entity, nil
}

func (r *ProductAccessRepo) GetAll() ([]models.ProductAccess, error) {
	var entities []models.ProductAccess
	err := r.DB.
		Preload("Product").
		Where("deleted_at IS NULL").
		Find(&entities).Error
	return entities, err
}

// func (repo *ProductAccessRepo) GetPrivateFields(access *models.ProductAccess) ( error) {
// 	repo.PrivateFieldRepo.GetPrivateFields(access)
// 	return &
// }
