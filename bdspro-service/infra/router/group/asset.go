package group_router

import (
	"bdspro/internal/dto"
	org_dto "bdspro/internal/dto/org"
	group_usecase "bdspro/internal/usecases/group"
	_routes "common/routes"
	_utils "common/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

type GroupAssetRouter struct {
	GroupAssetUsecase *group_usecase.GroupAssetUsecase
}

func NewGroupAssetRouter(groupAssetUsecase *group_usecase.GroupAssetUsecase) *GroupAssetRouter {
	return &GroupAssetRouter{
		GroupAssetUsecase: groupAssetUsecase,
	}
}

func (r *GroupAssetRouter) RegisterRoutes(group *gin.RouterGroup) {
	api := group.Group("/asset")
	{
		api.GET("/:groupId", r.GetAsset)
		api.POST("/share", r.ShareAsset)
	}
}

// @Summary Danh sách sản phẩm
// @Description Danh sách sản phẩm của nhóm
// @Tags Group
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page"
// @Param limit query int false "Limit"
// @Param productId query int false "Product ID"
// @Param groupId path int true "Group ID"
// @Success 200 {object} []domain.Product
// @Failure 400 {object} map[string]string "Dữ liệu không hợp lệ"
// @Failure 500 {object} map[string]string "Lỗi server"
// @Router /v1/bdspro/v2/group/product/{groupId} [get]
func (r *GroupAssetRouter) GetAsset(c *gin.Context) {
	groupId, err := strconv.ParseUint(c.Param("groupId"), 10, 64)
	if err != nil {
		_routes.RouteResult(c, nil, err)
		return
	}

	dto := dto.ProductSearchRequest{}
	if err := _utils.ParseQuery2(c, &dto); err != nil {
		_routes.RouteResult(c, nil, err)
		return
	}

	products, total, err := r.GroupAssetUsecase.GetAsset(c, groupId, dto)
	_routes.RouteResult(c, _routes.ResponseDTO{
		Code:          0,
		Data:          products,
		TotalElements: &total,
	}, err)
}

// @Summary Chia sẻ sản phẩm
// @Description Chia sẻ sản phẩm
// @Tags Group
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param productShareRequest body org_dto.ProductShareRequest true "Product Share Request"
// @Success 200 {object} domain.ProductAccess "Thành công"
// @Failure 400 {object} map[string]string "Dữ liệu không hợp lệ"
// @Failure 500 {object} map[string]string "Lỗi server"
// @Router /v1/bdspro/v2/group/asset/share [post]
func (r *GroupAssetRouter) ShareAsset(c *gin.Context) {
	dto := org_dto.ProductShareRequest{}
	if err := _utils.ParseBodyWithValidator(c, &dto); err != nil {
		_routes.RouteResult(c, nil, err)
		return
	}

	productAccess, err := r.GroupAssetUsecase.ShareProduct(c, dto)
	_routes.RouteResult(c, productAccess, err)
}
