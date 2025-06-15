package group_router

import (
	shared_usecase "bdspro/internal/usecases/shared"

	"github.com/gin-gonic/gin"
)

type GroupProductRouter struct {
	UC *shared_usecase.ProductUsecase
}

func NewGroupProductRouter(groupProductUsecase *shared_usecase.ProductUsecase) *GroupProductRouter {
	return &GroupProductRouter{
		UC: groupProductUsecase,
	}
}

func (r *GroupProductRouter) RegisterRoutes(group *gin.RouterGroup) {
	// api := group.Group("/product")
	// {
	// 	api.GET("/:groupId", r.GetProduct)
	// 	api.POST("/share", r.ShareProduct)
	// }
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
func (r *GroupProductRouter) GetProduct(c *gin.Context) {
	// groupId, err := strconv.ParseUint(c.Param("groupId"), 10, 64)
	// if err != nil {
	// 	_routes.RouteResult(c, nil, err)
	// 	return
	// }

	// dto := dto.ProductSearchRequest{}
	// if err := _utils.ParseQuery2(c, &dto); err != nil {
	// 	_routes.RouteResult(c, nil, err)
	// 	return
	// }

	// products, total, err := r.GroupProductUsecase.GetProduct(c, groupId, dto)
	// _routes.RouteResult(c, _routes.ResponseDTO{
	// 	Code:          0,
	// 	Data:          products,
	// 	TotalElements: &total,
	// }, err)
}

// @Summary Chia sẻ sản phẩm
// @Description Chia sẻ sản phẩm
// @Tags Group
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param productShareRequest body org_dto.ProductShareRequest true "Product Share Request"
// @Success 200 {object} domain.AssetShare "Thành công"
// @Failure 400 {object} map[string]string "Dữ liệu không hợp lệ"
// @Failure 500 {object} map[string]string "Lỗi server"
// @Router /v1/bdspro/v2/group/product/share [post]
func (r *GroupProductRouter) ShareProduct(c *gin.Context) {
	// dto := org_dto.ProductShareRequest{}
	// if err := _utils.ParseBodyWithValidator(c, &dto); err != nil {
	// 	_routes.RouteResult(c, nil, err)
	// 	return
	// }

	// productAccess, err := r.GroupProductUsecase.ShareProduct(c, dto)
	// _routes.RouteResult(c, productAccess, err)
}
