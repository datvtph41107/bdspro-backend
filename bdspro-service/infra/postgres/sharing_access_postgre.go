package postgres

import (
	"bdspro/internal/domain"
	"bdspro/internal/dto"
	"bdspro/internal/enums"
	"common/case/crud3"
	"context"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// @bind: bdspro/internal/repo.SharingAccessRepo
type SharingAccessPostgre struct {
	// Db gorm
	crud3.BaseRepo[domain.SharingAccess]
	// PrivateFieldRepo *product_private.ProductPrivateRepo
}

func NewSharingAccessPostgre(db *gorm.DB) *SharingAccessPostgre {
	return &SharingAccessPostgre{
		BaseRepo: crud3.BaseRepo[domain.SharingAccess]{DB: db},
		// PrivateFieldRepo: PrivateFieldRepo,
	}
}

// ExistsByUserAndProduct kiểm tra xem UserId và ProductId có tồn tại không
func (repo *SharingAccessPostgre) ExistsByUserAndProduct(userId, productId uint64) (bool, error) {
	var count int64
	err := repo.DB.Model(&domain.SharingAccess{}).
		Where("to_type = 10 AND to_id = ? AND product_id = ? and deleted_at is null", userId, productId).
		Count(&count).Error

	if err != nil {
		return false, err
	}

	return count > 0, nil
}

// ExistsByUserAndProduct kiểm tra xem UserId và ProductId có tồn tại không
func (repo *SharingAccessPostgre) GetByUserAndProduct(userId, productId uint64) (*domain.SharingAccess, error) {
	var entity *domain.SharingAccess
	err := repo.DB.
		Where("to_type = 10 AND to_id = ? AND product_id = ? and deleted_at is null", userId, productId).
		First(&entity).Error

	if err != nil {
		return nil, err
	}

	return entity, nil
}

// ExistsByUserAndProduct kiểm tra xem UserId và ProductId có tồn tại không
func (repo *SharingAccessPostgre) GetByTargetAndProduct(userId, target, productId uint64) (*domain.SharingAccess, error) {
	var entity *domain.SharingAccess
	err := repo.DB.
		Where("to_type = ? AND to_id = ? AND product_id = ? and deleted_at is null", target, userId, productId).
		First(&entity).Error

	if err != nil {
		return nil, err
	}

	return entity, nil
}

func (r *SharingAccessPostgre) GetAll() ([]domain.SharingAccess, error) {
	var entities []domain.SharingAccess
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

func (r *SharingAccessPostgre) OwnerSearchProduct(c context.Context, productId uint64, fromType enums.EOwnerOf, dto *dto.SharingAccessSearch) (*[]domain.SharingAccess, error) {
	var entities *[]domain.SharingAccess
	query := r.QueryList(r.DB, productId, fromType, dto)
	if err := query.
		Offset(dto.GetOffset()).
		Limit(dto.GetLimit()).
		Find(&entities).Error; err != nil {
		return nil, err
	}
	// err := r.DB.
	// 	Table("sharing_access").
	// 	Select("sharing_access.*, pf.full_name as target_name, pf.avatar as avatar").
	// 	Joins("left join profile_transfer pf on sharing_access.to_id = pf.profile_id and sharing_access.to_type = 10").
	// 	Where(`deleted_at IS NULL
	// 	and sharing_access.domain_id = ?
	// 	and sharing_access.domain = ?
	// 	and sharing_access.from_type = ?
	// 	and sharing_access.to_type = ?`,
	// 		productId, enums.EDomainProduct, fromType, dto.ToType).
	// 	Offset(dto.GetOffset()).
	// 	Limit(dto.GetLimit()).
	// 	Scan(&entities).Error
	return entities, nil
}

func (r *SharingAccessPostgre) QueryList(txt *gorm.DB,
	assetId uint64,
	fromType enums.EOwnerOf,
	dto *dto.SharingAccessSearch,
) *gorm.DB {
	query := txt.
		Where("deleted_at IS NULL")
	if dto.DomainID != 0 {
		query = query.Where("domain_id = ?", dto.DomainID)
	}
	if dto.Domain != 0 {
		query = query.Where("domain = ?", dto.Domain)
	}
	if dto.FromType != 0 {
		query = query.Where("from_type = ?", dto.FromType)
	}
	if dto.ToType != 0 {
		query = query.Where("to_type = ?", dto.ToType)
	}
	return query
}

func (r *SharingAccessPostgre) OwnerSearchAsset(c context.Context, assetId uint64, fromType enums.EOwnerOf, dto *dto.SharingAccessSearch) (*[]domain.SharingAccess, error) {
	var entities *[]domain.SharingAccess
	query := r.DB.
		Model(&domain.SharingAccess{}).
		Where(`deleted_at IS NULL`)

	query = r.QueryList(query, assetId, fromType, dto)

	query = query.Offset(dto.GetOffset()).
		Limit(dto.GetLimit())

	err := query.Scan(&entities).Error
	return entities, err
}

func (r *SharingAccessPostgre) UpdateCommission(c context.Context,
	profileId uint64,
	dto dto.CommissionUpdate,
) error {
	err := r.DB.WithContext(c).
		Debug().
		Model(&domain.SharingAccess{}).
		Where("deleted_at is null and to_id = ? and to_type = ? and domain_id = ? and domain = ?",
			dto.ToID, dto.ToType, dto.ProductID, enums.EDomainProduct).
		Update("commission", dto.Commission).Error

	return err
}

// func (r *GormProductAccessRepo) GetByID(c context.Context, id uint64) (*domain.ProductAccess, error){

// }
// func (r *GormProductAccessRepo)	Delete(c context.Context, id uint64) error{

// }

func (r *SharingAccessPostgre) GetProductAccessByOrgID(ctx context.Context, orgID int64, dto dto.SharingAccessSearch) ([]*domain.SharingAccess, error) {
	var entities []*domain.SharingAccess
	err := r.DB.WithContext(ctx).
		Where("deleted_at is null and to_type = ? and to_id = ?", enums.EOwnerOfOrganization, orgID).
		Offset(dto.GetOffset()).
		Limit(dto.GetLimit()).
		Find(&entities).Error
	return entities, err
}

func (r *SharingAccessPostgre) BulkSave(c context.Context, id uint64, dto dto.SharingAccessBulk) (*[]domain.SharingAccess, error) {
	var entities []domain.SharingAccess

	assignments := append(
		clause.AssignmentColumns([]string{"fields", "commission"}),
		clause.Assignment{
			Column: clause.Column{Name: "deleted_at"},
			Value:  gorm.Expr("NULL"),
		},
	)

	err := r.DB.WithContext(c).
		Clauses(clause.OnConflict{
			Columns: []clause.Column{
				{Name: "product_id"},
				{Name: "to_id"},
				{Name: "to_type"},
			},
			DoUpdates: clause.Set(assignments),
		}).
		Create(&entities).Error

	return &entities, err
}

func (r *SharingAccessPostgre) Delete(c context.Context, id uint64) error {
	return r.DB.WithContext(c).
		Model(&domain.SharingAccess{}).
		Where("id = ?", id).
		Update("deleted_at", time.Now()).Error
}

func (r *SharingAccessPostgre) GetByID(c context.Context, id uint64) (*domain.SharingAccess, error) {
	var entity *domain.SharingAccess
	err := r.DB.WithContext(c).
		Where("id = ?", id).
		First(&entity).Error
	return entity, err
}

// CountSharingAccessByProductID đếm số sharing access theo productId
func (r *SharingAccessPostgre) CountSharingAccessByProductID(ctx context.Context, productID uint64) (uint32, error) {
	var count int64

	err := r.DB.WithContext(ctx).Model(&domain.SharingAccess{}).
		Where("domain_id = ? AND domain = ? AND deleted_at IS NULL", productID, enums.EDomainProduct).
		Count(&count).Error
	if err != nil {
		return 0, err
	}

	return uint32(count), nil
}
