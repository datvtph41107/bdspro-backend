package u_asset_income_type

import (
	"bdspro/v1/models"
	_category "common/case/category"

	"gorm.io/gorm"
)

type AssetIncomeTypeRepo struct {
	_category.CategoryRepo[models.AssetIncomeType]
}

func NewAssetIncomeTypeRepo(DB *gorm.DB) *AssetIncomeTypeRepo {
	return &AssetIncomeTypeRepo{
		CategoryRepo: _category.CategoryRepo[models.AssetIncomeType]{
			DB: DB,
		},
	}
}
