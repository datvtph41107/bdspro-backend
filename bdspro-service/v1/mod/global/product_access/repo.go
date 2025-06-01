package product_accecss

import (
	"bdspro/internal/enums"
	user_dto "bdspro/v1/mod/user/dto/request"
	"bdspro/v1/models"
	"common/case/crud2"

	"github.com/gin-gonic/gin"
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

// func (r *ProductAccessRepo) OwnerGetByProduct(c *gin.Context, id uint64, targetType enums.TargetTypeEnum) (*[]models.ProductAccessQuery, error) {
// 	var entities *[]models.ProductAccessQuery
// 	err := r.DB.
// 		Table("product_access").
// 		Select("product_access.*, pf.full_name as target_name").
// 		Joins("left join profile_transfer pf on product_access.target_id = pf.profile_id and product_access.target_type = 10").
// 		// Preload("Product").
// 		Where("deleted_at IS NULL and product_id = ? and target_type = ?", id, targetType).
// 		Scan(&entities).Error
// 	return entities, err
// }

func (r *ProductAccessRepo) OwnerGetByProduct(c *gin.Context, id uint64, targetType enums.ETargetType) (*[]models.ProductAccessQuery, error) {
	var entities *[]models.ProductAccessQuery
	err := r.DB.
		Table("product_access").
		Select("product_access.*, pf.full_name as target_name, pf.avatar as avatar").
		Joins("left join profile_transfer pf on product_access.target_id = pf.profile_id and product_access.target_type = 10").
		Where("deleted_at IS NULL and product_access.product_id = ?", id).
		Scan(&entities).Error
	return entities, err
}

func (r *ProductAccessRepo) OwnerSearchProduct(c *gin.Context, id uint64, dto *user_dto.ProductAccessSearch) (*[]models.ProductAccessQuery, error) {
	var entities *[]models.ProductAccessQuery
	err := r.DB.
		Debug().
		Table("product_access").
		Select("product_access.*, pf.full_name as target_name, pf.avatar as avatar").
		Joins("left join profile_transfer pf on product_access.target_id = pf.profile_id and product_access.target_type = 10").
		Where("deleted_at IS NULL and product_access.product_id = ?", id).
		Offset(dto.GetOffset()).
		Limit(dto.GetLimit()).
		Scan(&entities).Error
	return entities, err
}

// func (repo *ProductAccessRepo) GetPrivateFields(access *models.ProductAccess) ( error) {
// 	repo.PrivateFieldRepo.GetPrivateFields(access)
// 	return &
// }
