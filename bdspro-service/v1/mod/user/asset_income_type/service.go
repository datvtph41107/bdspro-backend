package u_asset_income_type

import (
	"bdspro/v1/models"
	_category "common/case/category"
)

type AssetIncomeTypeService struct {
	_category.CategoryService[models.AssetIncomeType, SearchDTO]
}

func NewAssetIncomeTypeService(Repo *AssetIncomeTypeRepo) *AssetIncomeTypeService {
	return &AssetIncomeTypeService{
		CategoryService: _category.CategoryService[models.AssetIncomeType, SearchDTO]{
			CategoryRepo: Repo,
		},
	}
}
