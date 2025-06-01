package user_router

import (
	"bdspro/internal/dto"
	"bdspro/internal/usecases"
	_routes "common/routes"
	_utils "common/utils"

	"github.com/gin-gonic/gin"
)

type UAssetLegalRoute struct {
	UC *usecases.AssetLegalUsecase
}

func NewUAssetLegalRoute(uc *usecases.AssetLegalUsecase) *UAssetLegalRoute {
	return &UAssetLegalRoute{
		UC: uc,
	}
}

func (r *UAssetLegalRoute) RegisterRoutes(router *gin.RouterGroup) {
	group := router.Group("/asset-legal")
	{
		group.GET("", r.Search)
	}
}

func (s *UAssetLegalRoute) Search(c *gin.Context) {
	dto := dto.AssetLegalGetDTO{}
	if err := _utils.ParseQuery2(c, &dto); err != nil {
		_routes.RouteResult(c, nil, err)
		return
	}

	result, total, err := s.UC.GetData(c, dto)

	_routes.RouteResult(c, _routes.ResponseDTO{
		Code:          0,
		Data:          result,
		TotalElements: &total,
	}, err)
}
