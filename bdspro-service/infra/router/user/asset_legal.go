package user_router

import (
	"bdspro/internal/dto"
	"bdspro/internal/usecases"
	_routes "common/routes"
	_utils "common/utils"

	"github.com/gin-gonic/gin"
)

type UserAssetLegalRouter struct {
	UC *usecases.AssetLegalUsecase
}

func NewUAssetLegalRoute(uc *usecases.AssetLegalUsecase) *UserAssetLegalRouter {
	return &UserAssetLegalRouter{
		UC: uc,
	}
}

func (r *UserAssetLegalRouter) RegisterRoutes(router *gin.RouterGroup, path string) {
	group := router.Group(path)
	{
		group.GET("", r.Search)
	}
}

// @Summary Lấy danh sách
// @Tags User/asset-legal
// @Produce json
// @Param query query dto.AssetLegalGetDTO true "Size"
// @Security BearerAuth
// @Router /v1/bdspro/v2/user/asset/legal [get]
func (r *UserAssetLegalRouter) Search(c *gin.Context) {
	dto := dto.AssetLegalGetDTO{}
	if err := _utils.ParseQuery2(c, &dto); err != nil {
		_routes.RouteResult(c, nil, err)
		return
	}

	result, total, err := r.UC.GetData(c, dto)

	_routes.RouteResult(c, _routes.ResponseDTO{
		Code:          0,
		Data:          result,
		TotalElements: &total,
	}, err)
}

// // @Summary Lấy danh sách
// // @Tags Thu - chi
// // @Produce json
// // @Param query query dto.AssetLegalGetDTO true "Size"
// // @Security BearerAuth
// // @Router /v1/bdspro/v1/user/asset-legal [get]
// func GetAsset() {}

// // @Summary Xem chi tiết
// // @Tags Thu - chi
// // @Produce json
// // @Param id path int true "ID"
// // @Security BearerAuth
// // @Router /v1/bdspro/v1/user/asset-legal/{id} [get]
// func DetailAsset() {}

// // @Summary Tạo mới
// // @Tags Thu - chi
// // @Produce json
// // @Param body body domain.AssetLegal true "Body"
// // @Security BearerAuth
// // @Router /v1/bdspro/v1/user/asset-legal [post]
// func CreateAsset() {}

// // @Summary Lấy theo ID
// // @Tags Thu - chi
// // @Produce json
// // @Param id path int true "ID"
// // @Param body body domain.AssetLegal true "Body"
// // @Security BearerAuth
// // @Router /v1/bdspro/v1/user/asset-legal/{id} [put]
// func UpdateAsset() {}

// // @Summary Xóa
// // @Tags Thu - chi
// // @Produce json
// // @Param id path int true "ID"
// // @Security BearerAuth
// // @Router /v1/bdspro/v1/user/asset-legal/{id} [delete]
// func DeleteAsset() {}
