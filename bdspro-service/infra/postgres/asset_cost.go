package postgres

import (
	"bdspro/infra/shared"
	"bdspro/internal/domain"
	"bdspro/internal/dto"
	"bdspro/internal/enums"
	"context"

	"gorm.io/gorm"
)

// @bind: bdspro/internal/repo.AssetCostRepo
type GormAssetCostRepo struct {
	shared.PostgreCrud[domain.AssetCost, *dto.AssetCostSearchDTO]
}

func NewGormAssetCostRepo(DB *gorm.DB) *GormAssetCostRepo {
	return &GormAssetCostRepo{
		PostgreCrud: shared.PostgreCrud[domain.AssetCost, *dto.AssetCostSearchDTO]{
			DB: DB,
		},
	}
}

func (r *GormAssetCostRepo) QuerySearch(c context.Context, profileId uint64, ownerType enums.EOwnerOf, dto *dto.AssetCostSearchDTO) (*gorm.DB, error) {
	query := r.DB.Table("asset_costs").
		Select("asset_costs.*, asct.type_name as type_name").
		Joins("left join asset_cost_types asct on asct.id = asset_costs.cost_type_id").
		Where("asset_costs.deleted_at IS NULL and asset_costs.owner_id = ? and asset_costs.owner_type = ?", profileId, ownerType)
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
	return query, nil
}

func (r *GormAssetCostRepo) Find(c context.Context, profileId uint64, ownerType enums.EOwnerOf, dto *dto.AssetCostSearchDTO) ([]domain.AssetCost, int64, error) {
	var entities []domain.AssetCost
	var total int64

	query, err := r.QuerySearch(c, profileId, ownerType, dto)
	if err != nil {
		return nil, 0, err
	}

	if err := query.
		Debug().
		Preload("CostType").
		Offset(dto.GetOffset()).
		Limit(dto.GetLimit()).
		Find(&entities).Error; err != nil {
		return nil, 0, nil
	}
	query, _ = r.QuerySearch(c, profileId, ownerType, dto)
	err = query.Count(&total).Error
	return entities, total, err
}
