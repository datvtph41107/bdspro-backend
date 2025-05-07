package u_product_child

import (
	_routes "common/routes"
	_utils "common/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

type UserProductChildRoute struct {
	productService *ProductChildService
}

func NewUserProductChildRoute(product *ProductChildService) *UserProductChildRoute {
	return &UserProductChildRoute{
		productService: product,
	}
}

func (route *UserProductChildRoute) RegisterRoute(r *gin.RouterGroup) {
	api := r.Group("/product")
	// api.GET("", route.CreateChild)
	api.POST("/split", route.CreateChild)
	api.POST("/merge", route.MergeChild)
	api.POST("/merge-all/:productId", route.MergeAll)
	api.POST("/devide", route.DevideChild)
	api.GET("/history", route.History)
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
// @Param body body SaveProductChildRequest true "Thông tin sản phẩm con"
// @Router /v1/bdspro/user/product/split [post]
func (route *UserProductChildRoute) CreateChild(c *gin.Context) {
	_body := SaveProductChildRequest{}

	if err := _utils.ParseBodyWithValidator(c, &_body); err != nil {
		_routes.RouteResult(c, nil, err)
		return
	}
	result, err := route.productService.Create(c, &_body)
	_routes.RouteResult(c, result, err)
}

// @Summary Gộp sản phẩm con
// @Description API
// @Tags User/Product-child
// @Accept json
// @Produce json
// @Param productId path int false "Product ID"
// @Security BearerAuth
// @Param body body MergeProductChildRequest true "Thông tin sản phẩm con"
// @Router /v1/bdspro/user/product/merge [post]
func (route *UserProductChildRoute) MergeChild(c *gin.Context) {
	_body := MergeProductChildRequest{}
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
func (route *UserProductChildRoute) MergeAll(c *gin.Context) {
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
// @Param body body DevideChildRequest true "Thông tin sản phẩm con"
// @Router /v1/bdspro/user/product/devide [post]
func (route *UserProductChildRoute) DevideChild(c *gin.Context) {
	_body := DevideChildRequest{}
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
// @Param body query ProductHistorySearch true "Thông tin sản phẩm con"
// @Router /v1/bdspro/user/product/history [get]
func (route *UserProductChildRoute) History(c *gin.Context) {
	_body := ProductHistorySearch{}
	if err := _utils.ParseQuery2(c, &_body); err != nil {
		_routes.RouteResult(c, nil, err)
		return
	}
	result, err := route.productService.History(c, &_body)
	_routes.RouteResult(c, result, err)
}
