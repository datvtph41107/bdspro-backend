package postgre

import (
	"bdspro/internal/domain"
	"bdspro/internal/dto"
	"context"
	"time"

	"gorm.io/gorm"
)

// bind: bdspro/internal/repo.AssetCostRepo
type GormAssetCostRepo struct {
	DB *gorm.DB
}

func NewGormAssetCostRepo(DB *gorm.DB) *GormAssetCostRepo {
	return &GormAssetCostRepo{
		DB: DB,
	}
}

func (r *GormAssetCostRepo) Create(c context.Context, entity *domain.AssetCost) error {
	return r.DB.WithContext(c).Create(entity).Error
}

func (r *GormAssetCostRepo) Update(c context.Context, id uint64, entity *domain.AssetCost) error {
	return r.DB.WithContext(c).
		Model(entity).
		Where("id = ? and deleted_at is null", id).
		Updates(entity).Error
}

func (r *GormAssetCostRepo) Delete(c context.Context, id uint64) error {
	return r.DB.WithContext(c).
		Model(&domain.AssetCost{}).
		Where("id = ?", id).
		Update("deleted_at", time.Now()).Error
}

func (r *GormAssetCostRepo) GetData(c context.Context) ([]domain.AssetCost, error) {
	var entities []domain.AssetCost
	err := r.DB.Where("deleted_at IS NULL").Find(&entities).Error
	return entities, err
}

func (r *GormAssetCostRepo) GetByID(c context.Context, id uint64) (*domain.AssetCost, error) {
	var entity domain.AssetCost
	err := r.DB.Where("id = ? AND deleted_at IS NULL", id).First(&entity).Error
	if err != nil {
		return nil, err
	}
	return &entity, nil
}

func (r *GormAssetCostRepo) SearchQuery(g *gorm.DB, profileId uint64, dto dto.AssetCostSearchDTO) *gorm.DB {
	query := g.Table("asset_costs").
		Select("asset_costs.*, asct.type_name as type_name").
		Joins("left join asset_cost_types asct on asct.id = asset_costs.cost_type_id").
		Where("asset_costs.deleted_at IS NULL and asset_costs.owner_id = ?", profileId)
	if dto.AssetID != nil {
		query = query.Where("asset_id = ?", dto.AssetID)
	}
	if dto.Types != nil {
		query = query.Where("type in (?)", dto.Types)
	}
	if dto.TypeIDs != nil {
		query = query.Where("type_id in (?)", dto.TypeIDs)
	}
	// if dto.Text != nil {
	// 	query = query.Where("type_id in (?)", dto.TypeIDs)
	// }
	return query

}

func (r *GormAssetCostRepo) Search2(c context.Context, profileId uint64, dto dto.AssetCostSearchDTO) ([]domain.AssetCostQuery, int64, error) {
	var entities []domain.AssetCostQuery
	var total int64

	query := r.SearchQuery(r.DB, profileId, dto)

	if err := query.
		Offset(dto.GetOffset()).
		Limit(dto.GetLimit()).
		Scan(&entities).Error; err != nil {
		return nil, 0, nil
	}
	err := r.SearchQuery(r.DB, profileId, dto).Count(&total).Error
	return entities, total, err
}
