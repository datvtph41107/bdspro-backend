package user_router

import (
	"bdspro/internal/domain"
	"bdspro/internal/dto"
	"bdspro/internal/enums"
	"bdspro/internal/usecases"
	_routes "common/routes"
	_utils "common/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

type AssetShareRouter struct {
	uc *usecases.AssetShareUC
}

func NewAssetShareRouter(uc *usecases.AssetShareUC) *AssetShareRouter {
	return &AssetShareRouter{
		uc: uc,
	}
}

func (r *AssetShareRouter) RegisterRoutes(router *gin.RouterGroup, path string) {
	assetGroup := router.Group(path)
	{
		assetGroup.GET("", r.Search)
		assetGroup.POST("", r.Create)
		assetGroup.PUT(":id", r.Update)
		assetGroup.DELETE(":id", r.Delete)
	}
}

// @Summary Lấy danh sách chia sẻ tài sản
// @Tags User: Tài sản
// @Produce json
// @Param query query dto.AssetShareSearchDTO true "Size"
// @Security BearerAuth
// @Router /v1/bdspro/v1/user/asset/share [get]
func (r *AssetShareRouter) Search(c *gin.Context) {
	var dto dto.AssetShareSearchDTO
	if err := _utils.ParseQuery2(c, &dto); err != nil {
		_routes.RouteResult(c, nil, err)
		return
	}

	result, total, err := r.uc.Search(c, enums.EOwnerOfMember, &dto)
	_routes.RouteResult(c, _routes.ResponseDTO{
		Code:          0,
		Data:          result,
		TotalElements: &total,
	}, err)
}

// @Summary Xem chi tiết chia sẻ tài sản
// @Tags User: Tài sản
// @Produce json
// @Param id path int true "ID"
// @Security BearerAuth
// @Router /v1/bdspro/v1/user/asset/share/{id} [get]
func (r *AssetShareRouter) Detail(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	result, err := r.uc.Detail(c, id)
	_routes.RouteResult(c, result, err)
}

// @Summary Tạo mới chia sẻ tài sản
// @Tags User: Tài sản
// @Produce json
// @Param body body domain.AssetShare true "Body"
// @Security BearerAuth
// @Router /v1/bdspro/v1/user/asset/share [post]
func (r *AssetShareRouter) Create(c *gin.Context) {
	body := domain.AssetShare{}
	if err := _utils.ParseBodyWithValidator(c, &body); err != nil {
		_routes.RouteResult(c, nil, err)
		return
	}

	err := r.uc.Create(c, &body)
	_routes.RouteResult(c, nil, err)
}

// @Summary Cập nhật chia sẻ tài sản
// @Tags User: Tài sản
// @Produce json
// @Param id path int true "ID"
// @Param body body domain.AssetShare true "Body"
// @Security BearerAuth
// @Router /v1/bdspro/v1/user/asset/share/{id} [put]
func (r *AssetShareRouter) Update(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	body := domain.AssetShare{}
	if err := _utils.ParseBodyWithValidator(c, &body); err != nil {
		_routes.RouteResult(c, nil, err)
		return
	}

	err := r.uc.Update(c, id, &body)
	_routes.RouteResult(c, nil, err)
}

// @Summary Xóa chia sẻ tài sản
// @Tags User: Tài sản
// @Produce json
// @Param id path int true "ID"
// @Security BearerAuth
// @Router /v1/bdspro/v1/user/asset/share/{id} [delete]
func (r *AssetShareRouter) Delete(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	err := r.uc.Delete(c, id)
	_routes.RouteResult(c, nil, err)
}
