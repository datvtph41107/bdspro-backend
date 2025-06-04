package user_postgre

import (
	"bdspro/infra/shared"
	"bdspro/internal/domain"
	"bdspro/internal/dto"
	"bdspro/internal/enums"
	"context"

	"gorm.io/gorm"
)

// @bind: bdspro/internal/repo/user.UserAssetCostRepo
type PostgreUserAssetCost struct {
	shared.PostgreCrud[domain.AssetCost, *dto.AssetCostSearchDTO]
}

func NewPostgreUserAssetCost(db *gorm.DB) *PostgreUserAssetCost {
	repo := &PostgreUserAssetCost{
		PostgreCrud: shared.PostgreCrud[domain.AssetCost, *dto.AssetCostSearchDTO]{
			DB: db,
		},
	}
	repo.Repo = repo
	return repo
}

func (r *PostgreUserAssetCost) QuerySearch(c context.Context, profileId uint64, dto *dto.AssetCostSearchDTO) (*gorm.DB, error) {
	query := r.DB.Table("asset_costs").
		Select("asset_costs.*, asct.type_name as type_name").
		Joins("left join asset_cost_types asct on asct.id = asset_costs.cost_type_id").
		Where("asset_costs.deleted_at IS NULL and (asset_costs.owner_id = ? and asset_costs.owner_type = ?)", profileId, enums.EOwnerTypeMember)
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
