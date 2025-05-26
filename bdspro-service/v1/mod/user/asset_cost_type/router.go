package u_asset_cost_type

import (
	"bdspro/v1/models"
	_category "common/case/category"
	_routes "common/routes"
	_utils "common/utils"

	"github.com/gin-gonic/gin"
)

type AssetCostTypeRouter struct {
	_category.CategoryRouter[models.AssetCostType, SearchDTO]
	Service *AssetCostTypeService
}

func NewAssetCostTypeRouter(Service *AssetCostTypeService) *AssetCostTypeRouter {
	return &AssetCostTypeRouter{
		CategoryRouter: _category.CategoryRouter[models.AssetCostType, SearchDTO]{
			Service: Service,
		},
		Service: Service,
	}
}

func (r *AssetCostTypeRouter) RegisterRoutes(router *gin.RouterGroup, path string) {
	assetGroup := router.Group(path)
	{
		assetGroup.GET("", r.Search2)
		assetGroup.POST("", r.Create)
		assetGroup.PUT(":id", r.Update)
		assetGroup.DELETE(":id", r.Delete)
		assetGroup.GET(":id", r.GetByID)
	}
}

func (r *AssetCostTypeRouter) Search2(c *gin.Context) {
	var dto SearchDTO
	if err := _utils.ParseQuery2(c, &dto); err != nil {
		_routes.RouteResult(c, nil, &_routes.Except{
			Code:    400,
			Message: "Dữ liệu không hợp lệ",
		})

		return
	}

	result, err := r.Service.Search2(c, dto)
	_routes.RouteResult(c, result, err)
}
