package u_asset_cost

import (
	"bdspro/v1/models"
	_owner "common/case/owner"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AssetCostRepo struct {
	_owner.OwnerRepo[models.AssetCost]
}

func NewAssetCostRepo(DB *gorm.DB) *AssetCostRepo {
	return &AssetCostRepo{
		OwnerRepo: _owner.OwnerRepo[models.AssetCost]{
			DB: DB,
		},
	}
}

func (r *AssetCostRepo) SearchQuery(g *gorm.DB, profileId uint64, dto SearchDTO) *gorm.DB {
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

func (r *AssetCostRepo) Search2(c *gin.Context, profileId uint64, dto SearchDTO) ([]models.AssetCostQuery, int64, error) {
	var entities []models.AssetCostQuery
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
