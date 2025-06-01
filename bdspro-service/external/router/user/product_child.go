package user_router

import (
	"bdspro/internal/dto"
	"bdspro/internal/usecases"
	_routes "common/routes"
	_utils "common/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

type UProductChildRoute struct {
	productService *usecases.UProductChildUsecase
}

func NewUProductChildRoute(product *usecases.UProductChildUsecase) *UProductChildRoute {
	return &UProductChildRoute{
		productService: product,
	}
}

func (route *UProductChildRoute) RegisterRoute(r *gin.RouterGroup) {
	api := r.Group("/product")
	// api.GET("", route.CreateChild)
	api.POST("/split", route.CreateChild)
	api.POST("/merge", route.MergeChild)
	api.POST("/merge-all/:productId", route.MergeAll)
	api.POST("/devide", route.DevideChild)
	api.GET("/history-child/:id", route.History)
	api.GET("/area-info/:id", route.InfoArea)

}

// func (route *UserProductChildRoute) ListChild(c *gin.Context) {
// 	_body := data.ProductSearchRequest{}

// 	if err := _utils.ParseQuery2(c, &_body); err != nil {
// 		_routes.RouteResult(c, nil, err)
// 		return
// 	}
// 	result, err := route.productService.Search(c, &_body)
// 	_routes.RouteResult(c, result, err)
// }

// @Summary Tạo sản phẩm con
// @Description API
// @Tags User/Product-child
// @Accept json
// @Produce json
// @Param productId path int false "Product ID"
// @Security BearerAuth
// @Param body body dto.SaveProductChildRequest true "Thông tin sản phẩm con"
// @Router /v1/bdspro/user/product/split [post]
func (route *UProductChildRoute) CreateChild(c *gin.Context) {
	_body := dto.SaveProductChildRequest{}

	if err := _utils.ParseBodyWithValidator(c, &_body); err != nil {
		_routes.RouteResult(c, nil, err)
		return
	}
	result, err := route.productService.Create(c, &_body)
	_routes.RouteResult(c, result, err)
}

// @Summary Tạo sản phẩm con
// @Description API
// @Tags User/Product-child
// @Accept json
// @Produce json
// @Param productId path int false "Product ID"
// @Security BearerAuth
// @Router /v1/bdspro/user/product/split [post]
func (route *UProductChildRoute) InfoArea(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		_routes.RouteResult(c, nil, &_routes.Except{
			Code:    400,
			Message: "ID không hợp lệ",
		})
		return
	}

	result, err := route.productService.InfoArea(c, &id)
	_routes.RouteResult(c, result, err)
}

// @Summary Gộp sản phẩm con
// @Description API
// @Tags User/Product-child
// @Accept json
// @Produce json
// @Param productId path int false "Product ID"
// @Security BearerAuth
// @Param body body dto.MergeProductChildRequest true "Thông tin sản phẩm con"
// @Router /v1/bdspro/user/product/merge [post]
func (route *UProductChildRoute) MergeChild(c *gin.Context) {
	_body := dto.MergeProductChildRequest{}
	if err := _utils.ParseBodyWithValidator(c, &_body); err != nil {
		_routes.RouteResult(c, nil, err)
		return
	}

	result, err := route.productService.MergeProductChild(c, &_body)
	_routes.RouteResult(c, result, err)
}

// @Summary Gộp toàn bộ sản phẩm con
// @Description API
// @Tags User/Product-child
// @Accept json
// @Produce json
// @Param productId path int true "Product ID"
// @Security BearerAuth
// @Router /v1/bdspro/user/product/merge-all/{productId} [post]
func (route *UProductChildRoute) MergeAll(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("productId"), 10, 64)
	if err != nil {
		_routes.RouteResult(c, nil, err)
		return
	}
	result, err := route.productService.MergeAllProductChild(c, &id)
	_routes.RouteResult(c, result, err)
}

// @Summary Gộp toàn bộ sản phẩm con
// @Description API
// @Tags User/Product-child
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param body body dto.DevideChildRequest true "Thông tin sản phẩm con"
// @Router /v1/bdspro/user/product/devide [post]
func (route *UProductChildRoute) DevideChild(c *gin.Context) {
	_body := dto.DevideChildRequest{}
	if err := _utils.ParseBodyWithValidator(c, &_body); err != nil {
		_routes.RouteResult(c, nil, err)
		return
	}
	result, err := route.productService.DevideProductChild(c, &_body)
	_routes.RouteResult(c, result, err)
}

// @Summary Gộp toàn bộ sản phẩm con
// @Description API
// @Tags User/Product-child
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param body query dto.ProductHistorySearch true "Thông tin sản phẩm con"
// @Router /v1/bdspro/user/product/history-child/:id [get]
func (route *UProductChildRoute) History(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		_routes.RouteResult(c, nil, &_routes.Except{
			Code:    400,
			Message: "ID không hợp lệ",
		})
		return
	}

	_body := dto.ProductHistorySearch{}
	if err := _utils.ParseQuery2(c, &_body); err != nil {
		_routes.RouteResult(c, nil, err)
		return
	}
	result, err := route.productService.History(c, id, &_body)
	_routes.RouteResult(c, result, err)
}
