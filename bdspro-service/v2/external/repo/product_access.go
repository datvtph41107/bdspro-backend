package postgre_repo

import (
	"bdspro/v2/internal/domain"
	"common/case/crud3"

	"gorm.io/gorm"
)

// @bind: bdspro/v2/repo.ProductAccessRepo
type GormProductAccessRepo struct {
	crud3.BaseRepo[domain.ProductAccess]
	// PrivateFieldRepo *product_private.ProductPrivateRepo
}

func NewGromProductAccessRepo(db *gorm.DB) *GormProductAccessRepo {
	return &GormProductAccessRepo{
		BaseRepo: crud3.BaseRepo[domain.ProductAccess]{DB: db},
		// PrivateFieldRepo: PrivateFieldRepo,
	}
}

// ExistsByUserAndProduct kiểm tra xem UserId và ProductId có tồn tại không
func (repo *GormProductAccessRepo) ExistsByUserAndProduct(userId, productId uint64) (bool, error) {
	var count int64
	err := repo.DB.Model(&domain.ProductAccess{}).
		Where("target_type = 1 AND target_id = ? AND product_id = ? and deleted_at is null", userId, productId).
		Count(&count).Error

	if err != nil {
		return false, err
	}

	return count > 0, nil
}

// ExistsByUserAndProduct kiểm tra xem UserId và ProductId có tồn tại không
func (repo *GormProductAccessRepo) GetByUserAndProduct(userId, productId uint64) (*domain.ProductAccess, error) {
	var entity *domain.ProductAccess
	err := repo.DB.
		Where("target_type = 1 AND target_id = ? AND product_id = ? and deleted_at is null", userId, productId).
		First(&entity).Error

	if err != nil {
		return nil, err
	}

	return entity, nil
}

func (r *GormProductAccessRepo) GetAll() ([]domain.ProductAccess, error) {
	var entities []domain.ProductAccess
	err := r.DB.
		Preload("Product").
		Where("deleted_at IS NULL").
		Find(&entities).Error
	return entities, err
}

// func (repo *ProductAccessRepo) GetPrivateFields(access *domain.ProductAccess) ( error) {
// 	repo.PrivateFieldRepo.GetPrivateFields(access)
// 	return &
// }
