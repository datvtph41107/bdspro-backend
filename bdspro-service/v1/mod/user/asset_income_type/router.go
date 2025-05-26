package u_asset_income_type

import (
	"bdspro/v1/models"
	_category "common/case/category"
)

type AssetIncomeTypeRouter struct {
	_category.CategoryRouter[models.AssetIncomeType, SearchDTO]
}

func NewAssetIncomeTypeRouter(Service *AssetIncomeTypeService) *AssetIncomeTypeRouter {
	return &AssetIncomeTypeRouter{
		CategoryRouter: _category.CategoryRouter[models.AssetIncomeType, SearchDTO]{
			Service: Service,
		},
	}
}
