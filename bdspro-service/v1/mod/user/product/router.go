package u_product

import (
	"bdspro/v1/data"
	_routes "common/routes"
	_utils "common/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

type UserProductRoute struct {
	productService *ProductService
}

func NewUserProductRoute(product *ProductService) *UserProductRoute {
	return &UserProductRoute{
		productService: product,
	}
}

func (route *UserProductRoute) RegisterRouteUserProduct(r *gin.RouterGroup) {
	api := r.Group("/product")
	api.GET("/market", route.SearchProductMarket)
	api.GET("/total", route.SearchSmartProduct)

	api.GET("/me", route.SearchProductWithOwner)
	api.GET("/link/:id", route.SearchProductWithOwner)
	api.GET("/public", route.SearchProductWithOwner)
	api.POST("/new", route.NewProduct)
	api.POST("/suggest", route.NewProductWithDetect)
	api.GET("/:id", route.Detail)
	api.PUT("/edit/:id", route.Edit)
	api.PUT("/publish", route.PublishProduct)
	api.PUT("/revert", route.RevertProduct)
	api.PUT("/rent/:id", route.RevertProduct)
}

// @Summary Tạo sản phẩm mới
// @Description API để tạo một sản phẩm mới
// @Tags User
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param product body data.NewProductRequest true "Thông tin sản phẩm"
// @Success 201 {object} data.NewProductRequest "Sản phẩm được tạo thành công"
// @Router /user/product/new [post]
func (route *UserProductRoute) NewProduct(c *gin.Context) {
	body, _ := _utils.ParseBody[data.ProductSaveRequest](c)
	result, err := route.productService.CreateProduct(c, body)
	_routes.RouteResult(c, result, err)
}

// @Summary Tạo sản phẩm mới
// @Description API để tạo một sản phẩm mới
// @Tags User
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param product body data.NewProductRequest true "Thông tin sản phẩm"
// @Success 201 {object} data.NewProductRequest "Sản phẩm được tạo thành công"
// @Router /user/product/{id} [post]
func (route *UserProductRoute) Detail(c *gin.Context) {
	result, err := route.productService.Detail(c)
	_routes.RouteResult(c, result, err)
}

// @Summary Đăng tin sản phẩm
// @Description API để đăng tin sản phẩm lên hệ thống
// @Tags User
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param publish body data.StatusRequest true "Thông tin đăng tin"
// @Success 200 {object} data.StatusResponse "Tin được đăng thành công"
// @Router /user/product/publish [put]
func (route *UserProductRoute) PublishProduct(e *gin.Context) {

}

// @Summary Hoàn tác trạng thái sản phẩm
// @Description API để hoàn tác trạng thái sản phẩm về trạng thái trước đó
// @Tags User
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param revert body data.StatusRequest true "Thông tin hoàn tác"
// @Success 200 {object} data.StatusResponse "Trạng thái sản phẩm đã được hoàn tác"
// @Router /user/product/revert/:id [put]
func (route *UserProductRoute) RevertProduct(e *gin.Context) {

}

// SearchProductWithOwner godoc
// @Summary Lấy danh sách sản phẩm của user hiện tại
// @Description Returns a list of products based on search filters
// @Tags User
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param name query string false "Product name"
// @Param code query string false "Product code"
// @Param area query number false "Product area"
// @Param price query integer false "Product price"
// @Param description query string false "Product description"
// @Param note query string false "Additional notes"
// @Param limit query int false "Limit results" default(10)
// @Param skip query int false "Skip results" default(0)
// @Param sort query string false "Sort order" default("price_asc")
// @Router /user/product/me [get]
func (route *UserProductRoute) SearchProductWithOwner(c *gin.Context) {
	body, _ := _utils.ParseQuery[data.ProductSearchRequest](c)
	result, err := route.productService.SearchProductWithOwner(c, body)
	_routes.RouteResult(c, result, err)
}

func (route *UserProductRoute) SearchProductMarket(c *gin.Context) {
	body, _ := _utils.ParseQueryToStruct[data.ProductSearchRequest](c.Request.URL.Query())
	result, err := route.productService.SearchProductMarket(c, body)
	_routes.RouteResult(c, result, err)
}

func (route *UserProductRoute) SearchSmartProduct(c *gin.Context) {
	body, _ := _utils.ParseQueryToStruct[data.TextSearchRequest](c.Request.URL.Query())
	result, err := route.productService.SmartSearch(c, body.Text)
	_routes.RouteResult(c, result, err)
}

func (route *UserProductRoute) NewProductWithDetect(c *gin.Context) {
	_body, _ := _utils.ParseBody[data.ProductSaveDetect](c)
	result, err := route.productService.CreateProductWithDetect(c, _body)
	_routes.RouteResult(c, result, err)
}

// @Summary Cập nhật sản phẩm
// @Description API này cho phép cập nhật thông tin sản phẩm dựa trên dữ liệu được gửi lên.
// @Tags User
// @Accept json
// @Produce json
// @Param productId path int true "Product ID"
// @Security BearerAuth
// @Param body body data.ProductSaveRequest true "Thông tin sản phẩm cần cập nhật"
// @Router /user/product/{productId} [put]
func (route *UserProductRoute) Edit(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	body, err := _utils.ParseBody[data.ProductSaveRequest](c)
	if err != nil {
		_routes.RouteResult(c, nil, err)
	}
	result, err := route.productService.Update(c, id, body)
	_routes.RouteResult(c, result, err)
}

func (route *UserProductRoute) Rent(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	body, err := _utils.ParseBody[data.ProductSaveRequest](c)
	if err != nil {
		_routes.RouteResult(c, nil, err)
	}
	result, err := route.productService.Rent(c, id, body)
	_routes.RouteResult(c, result, err)
}
