package postgre_repo

import (
	"bdspro/internal/domain"
	"bdspro/internal/dto"
	"common/case/crud3"
	"context"

	"gorm.io/gorm"
)

// @bind: bdspro/internal/repo.ProductAccessRepo
type GormProductAccessRepo struct {
	// Db gorm
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
		Where("target_type = 10 AND target_id = ? AND product_id = ? and deleted_at is null", userId, productId).
		First(&entity).Error

	if err != nil {
		return nil, err
	}

	return entity, nil
}

// ExistsByUserAndProduct kiểm tra xem UserId và ProductId có tồn tại không
func (repo *GormProductAccessRepo) GetByTargetAndProduct(userId, target, productId uint64) (*domain.ProductAccess, error) {
	var entity *domain.ProductAccess
	err := repo.DB.
		Where("target_type = ? AND target_id = ? AND product_id = ? and deleted_at is null", target, userId, productId).
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

func (r *GormProductAccessRepo) OwnerSearchProduct(c context.Context, id uint64, dto *dto.ProductAccessSearch) (*[]domain.ProductAccessQuery, error) {
	var entities *[]domain.ProductAccessQuery
	err := r.DB.
		Table("product_access").
		Select("product_access.*, pf.full_name as target_name, pf.avatar as avatar").
		Joins("left join profile_transfer pf on product_access.target_id = pf.profile_id and product_access.target_type = 10").
		Where("deleted_at IS NULL and product_access.product_id = ?", id).
		Offset(dto.GetOffset()).
		Limit(dto.GetLimit()).
		Scan(&entities).Error
	return entities, err
}

func (r *GormProductAccessRepo) UpdateCommission(c context.Context,
	profileId uint64,
	dto dto.ProductAccessUpdate,
) error {
	err := r.DB.WithContext(c).
		Debug().
		Model(&domain.ProductAccess{}).
		Where("deleted_at is null and target_id = ? and target_type = ? and product_id = ?",
			dto.TargetID, dto.TargetType, dto.ProductID).
		Update("commission", dto.Commission).Error

	return err
}

// func (r *GormProductAccessRepo) GetByID(c context.Context, id uint64) (*domain.ProductAccess, error){

// }
// func (r *GormProductAccessRepo)	Delete(c context.Context, id uint64) error{

// }
