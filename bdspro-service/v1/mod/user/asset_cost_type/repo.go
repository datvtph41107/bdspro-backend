package u_asset_cost_type

import (
	"bdspro/v1/models"
	_category "common/case/category"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AssetCostTypeRepo struct {
	_category.CategoryRepo[models.AssetCostType]
}

func NewAssetCostTypeRepo(DB *gorm.DB) *AssetCostTypeRepo {
	return &AssetCostTypeRepo{
		CategoryRepo: _category.CategoryRepo[models.AssetCostType]{
			DB: DB,
		},
	}
}

func (r *AssetCostTypeRepo) SearchQuery(g *gorm.DB, profileId uint64, dto SearchDTO) *gorm.DB {
	query := g.Model(&models.AssetCostType{}).
		Where("deleted_at IS NULL and (owner_id = ? or owner_id is null)", profileId)
	if dto.Text != "" {
		query = query.Where("type_name like ?", "%"+dto.Text+"%")
	}
	if dto.Types != nil {
		query = query.Where("type in (?)", dto.Types)
	}
	return query

}

func (r *AssetCostTypeRepo) Search2(c *gin.Context, profileId uint64, dto SearchDTO) ([]models.AssetCostType, int64, error) {
	var entities []models.AssetCostType
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
