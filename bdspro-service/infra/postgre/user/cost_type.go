package user_postgre

import (
	"bdspro/infra/shared"
	"bdspro/internal/domain"
	"bdspro/internal/dto"
	"bdspro/internal/enums"
	"context"

	"gorm.io/gorm"
)

// @bind: bdspro/internal/repo/user.UserCostTypeRepo
type PostgreUserCostType struct {
	shared.PostgreCrud[domain.AssetCostType, *dto.AssetCostTypeSearchDTO]
}

func NewPostgreUserCostType(DB *gorm.DB) *PostgreUserCostType {
	repo := &PostgreUserCostType{
		PostgreCrud: shared.PostgreCrud[domain.AssetCostType, *dto.AssetCostTypeSearchDTO]{
			DB: DB,
		},
	}
	repo.Repo = repo
	return repo
}

func (r *PostgreUserCostType) QuerySearch(c context.Context, profileId uint64, dto *dto.AssetCostTypeSearchDTO) (*gorm.DB, error) {
	query := r.DB.Model(&domain.AssetCostType{}).
		Where("deleted_at IS NULL and ((owner_id = ? and owner_type = ?) or owner_id is null)", profileId, enums.EOwnerTypeMember)
	if dto.Text != "" {
		query = query.Where("type_name like ?", "%"+dto.Text+"%")
	}
	if dto.Types != nil {
		query = query.Where("type in (?)", dto.Types)
	}
	return query, nil
}

// func (r *GormCostTypeRepo) Create(c context.Context, entity *domain.AssetCostType) error {
// 	return r.DB.WithContext(c).Create(entity).Error
// }

// func (r *GormCostTypeRepo) Update(c context.Context, id uint64, entity *domain.AssetCostType) error {
// 	return r.DB.WithContext(c).
// 		Model(entity).
// 		Where("id = ? and deleted_at is null", id).
// 		Updates(entity).Error
// }

// func (r *GormCostTypeRepo) Delete(c context.Context, id uint64) error {
// 	return r.DB.WithContext(c).
// 		Model(&domain.AssetCostType{}).
// 		Where("id = ?", id).
// 		Update("deleted_at", time.Now()).Error
// }

// func (r *GormAssetCostTypeRepo) GetData(c context.Context) ([]domain.AssetCostType, error) {
// 	var entities []domain.AssetCostType
// 	err := r.DB.Where("deleted_at IS NULL").Find(&entities).Error
// 	return entities, err
// }

// func (r *GormAssetCostTypeRepo) GetByID(c context.Context, id uint64) (*domain.AssetCostType, error) {
// 	var entity domain.AssetCostType
// 	err := r.DB.Where("id = ? AND deleted_at IS NULL", id).First(&entity).Error
// 	if err != nil {
// 		return nil, err
// 	}
// 	return &entity, nil
// }

// func (r *GormAssetCostTypeRepo) Search2(c *gin.Context, profileId uint64, dto dto.AssetCostTypeSearchDTO) ([]domain.AssetCostType, int64, error) {
// 	var entities []domain.AssetCostType
// 	var total int64

// 	query := r.SearchQuery(r.DB, profileId, dto)

// 	if err := query.
// 		Offset(dto.GetOffset()).
// 		Limit(dto.GetLimit()).
// 		Scan(&entities).Error; err != nil {
// 		return nil, 0, nil
// 	}
// 	err := r.SearchQuery(r.DB, profileId, dto).Count(&total).Error
// 	return entities, total, err
// }

// func (r *GormAssetCostTypeRepo) Search(c context.Context, profileId uint64, id uint64, dto dto.AssetCostTypeSearchDTO) ([]domain.AssetCostType, int64, error) {
// 	var entities []domain.AssetCostType
// 	err := r.DB.Where("deleted_at IS NULL").Find(&entities).Error
// 	return entities, err
// 	// return r.Search2(c, profileId, dto)
// }

// func (r *GormCostTypeRepo) GetByID(c context.Context, id uint64) (*domain.AssetCostType, error) {
// 	var entity domain.AssetCostType
// 	err := r.DB.Where("id = ? AND deleted_at IS NULL", id).First(&entity).Error
// 	if err != nil {
// 		return nil, err
// 	}
// 	return &entity, nil
// }

// func (r *GormAssetCostTypeRepo) Detail(c context.Context, id uint64) (*domain.AssetCostType, error) {
// 	return r.GetByID(c, id)
// }

// func (r *GormAssetCostTypeRepo) Create(c context.Context, entity *domain.AssetCostType) error {
// 	return r.Create(c, entity)
// }

// func (r *GormAssetCostTypeRepo) Update(c context.Context, id uint64, entity *domain.AssetCostType) error {
// 	return r.Update(c, id, entity)
// }

// func (r *GormAssetCostTypeRepo) Delete(c context.Context, id uint64) error {
// 	return r.Delete(c, id)
// }

// func (r *GormAssetCostTypeRepo) GetAll(c context.Context) ([]domain.AssetCostType, error) {
// 	// return r.GetData(c)
// }
