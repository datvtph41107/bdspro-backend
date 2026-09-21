package user_router

import (
	"bdspro/internal/domain"
	"bdspro/internal/dto"
	"bdspro/internal/enums"
	shared_usecase "bdspro/internal/usecases/shared"
	_errors "common/errors"
	_routes "common/routes"
	_utils "common/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

type UserCostTypeRouter struct {
	// Router shared.ICrudOwnerRouter[domain.AssetCostType, *dto.AssetCostTypeSearchDTO]
	UC *shared_usecase.AssetCostTypeUsecase
}

func NewUserCostTypeRouter(uc *shared_usecase.AssetCostTypeUsecase) *UserCostTypeRouter {
	return &UserCostTypeRouter{
		UC: uc,
		// Router: &shared.CrudOwnerRouter[domain.AssetCostType, *dto.AssetCostTypeSearchDTO]{
		// 	UC: uc,
		// },
	}
}

func (r *UserCostTypeRouter) RegisterRoutes(router *gin.RouterGroup, path string) {
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
// @Param query query dto.AssetCostTypeSearchDTO true "Size"
// @Security BearerAuth
// @Router /v1/bdspro/v2/user/assets/cost-type [get]
func (r *UserCostTypeRouter) Search(c *gin.Context) {
	var dto dto.AssetCostTypeSearchDTO
	if err := _utils.ParseQuery2(c, &dto); err != nil {
		_routes.RouteResult(c, nil, _errors.ReturnError(_errors.RequestValidationFailed))

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
// @Router /v1/bdspro/v2/user/assets/cost-type/{id} [get]
func (r *UserCostTypeRouter) GetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		_routes.RouteResult(c, nil, _errors.ReturnError(_errors.RequestValidationFailed))
		return
	}

	assetCostType, err := r.UC.Detail(c, id)
	_routes.RouteResult(c, assetCostType, err)
}

// @Summary Tạo mới loại chi phí
// @Tags Hợp đồng khai thác
// @Produce json
// @Param body body domain.AssetCostType true "Body"
// @Security BearerAuth
// @Router /v1/bdspro/v2/user/assets/cost-type [post]
func (r *UserCostTypeRouter) Create(c *gin.Context) {
	var dto domain.AssetCostType
	if err := _utils.ParseBodyWithValidator(c, &dto); err != nil {
		_routes.RouteResult(c, nil, _errors.ReturnError(_errors.RequestValidationFailed))
		return
	}

	err := r.UC.Create(c, &dto)
	_routes.RouteResult(c, nil, err)
}

// @Summary Cập nhật loại chi phí
// @Tags Hợp đồng khai thác
// @Produce json
// @Param id path int true "ID"
// @Param body body domain.AssetCostType true "Body"
// @Security BearerAuth
// @Router /v1/bdspro/v2/user/assets/cost-type/{id} [put]
func (r *UserCostTypeRouter) Update(c *gin.Context) {
	var dto domain.AssetCostType
	if err := _utils.ParseBodyWithValidator(c, &dto); err != nil {
		_routes.RouteResult(c, nil, _errors.ReturnError(_errors.RequestValidationFailed))
		return
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		_routes.RouteResult(c, nil, _errors.ReturnError(_errors.RequestValidationFailed))
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
// @Router /v1/bdspro/v2/user/assets/cost-type/{id} [delete]
func (r *UserCostTypeRouter) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		_routes.RouteResult(c, nil, _errors.ReturnError(_errors.RequestValidationFailed))
		return
	}

	err = r.UC.Delete(c, id)
	_routes.RouteResult(c, nil, err)
}
