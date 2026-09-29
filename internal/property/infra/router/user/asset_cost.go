package user_router

import (
	"bdspro/internal/domain"
	"bdspro/internal/dto"
	"bdspro/internal/enums"
	shared_usecase "bdspro/internal/usecases/shared"
	_routes "common/routes"
	_utils "common/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

type UserAssetCostRouter struct {
	UC *shared_usecase.AssetCostUsecase
}

func NewUserAssetCostRouter(uc *shared_usecase.AssetCostUsecase) *UserAssetCostRouter {
	return &UserAssetCostRouter{
		UC: uc,
	}
}

func (r *UserAssetCostRouter) RegisterRoutes(router *gin.RouterGroup, path string) {
	assetGroup := router.Group(path)
	{
		assetGroup.GET("", r.Search)
		assetGroup.POST("", r.Create)
		assetGroup.PUT(":id", r.Update)
		assetGroup.DELETE(":id", r.Delete)
		assetGroup.GET(":id", r.GetByID)
		// assetGroup.GET(":id/detail", r.Detail)
	}
}

// @Summary Lấy danh sách loại chi phí
// @Tags Hợp đồng khai thác
// @Produce json
// @Param query query dto.AssetCostSearchDTO true "Size"
// @Security BearerAuth
// @Router /v1/bdspro/v2/user/assets/cost [get]
func (r *UserAssetCostRouter) Search(c *gin.Context) {
	var dto dto.AssetCostSearchDTO
	if err := _utils.ParseQuery2(c, &dto); err != nil {
		_routes.RouteResult(c, nil, err)
		return
	}

	result, total, err := r.UC.Search(c, enums.EOwnerOfMember, &dto)
	_routes.RouteResult(c, _routes.ResponseDTO{
		Code:          0,
		Data:          result,
		TotalElements: &total,
	}, err)
}

// @Summary Xem chi tiết loại chi phí
// @Tags Hợp đồng khai thác
// @Produce json
// @Param id path int true "ID"
// @Security BearerAuth
// @Router /v1/bdspro/v2/user/assets/cost/{id} [get]
func (r *UserAssetCostRouter) GetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		_routes.RouteResult(c, nil, err)
		return
	}

	assetCost, err := r.UC.Detail(c, id)
	_routes.RouteResult(c, assetCost, err)
}

// @Summary Tạo mới loại chi phí
// @Tags Hợp đồng khai thác
// @Produce json
// @Param body body domain.AssetCost true "Body"
// @Security BearerAuth
// @Router /v1/bdspro/v2/user/assets/cost [post]
func (r *UserAssetCostRouter) Create(c *gin.Context) {
	var dto domain.AssetCost
	if err := _utils.ParseBodyWithValidator(c, &dto); err != nil {
		_routes.RouteResult(c, nil, err)
		return
	}

	err := r.UC.Create(c, &dto)
	_routes.RouteResult(c, nil, err)
}

// @Summary Lấy theo ID loại chi phí
// @Tags Hợp đồng khai thác
// @Produce json
// @Param id path int true "ID"
// @Param body body domain.AssetCost true "Body"
// @Security BearerAuth
// @Router /v1/bdspro/v2/user/assets/cost/{id} [put]
func (r *UserAssetCostRouter) Update(c *gin.Context) {
	var dto domain.AssetCost
	if err := _utils.ParseBodyWithValidator(c, &dto); err != nil {
		_routes.RouteResult(c, nil, err)
		return
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		_routes.RouteResult(c, nil, err)
		return
	}

	err = r.UC.Update(c, id, &dto)
	_routes.RouteResult(c, nil, err)
}

// @Summary Xóa loại chi phí
// @Tags Hợp đồng khai thác
// @Produce json
// @Param id path int true "ID"
// @Security BearerAuth
// @Router /v1/bdspro/v2/user/assets/cost/{id} [delete]
func (r *UserAssetCostRouter) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		_routes.RouteResult(c, nil, err)
		return
	}

	err = r.UC.Delete(c, id)
	_routes.RouteResult(c, nil, err)
}
