package routes

import (
	"bdspro/internal/dto"
	"bdspro/internal/usecases"
	_routes "common/routes"
	_utils "common/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

type UProductRouter struct {
	UserProductService *usecases.UProductUsecase
}

func NewUProductRoute(product *usecases.UProductUsecase) *UProductRouter {
	return &UProductRouter{
		UserProductService: product,
	}
}

func (route *UProductRouter) RegisterRoutes(r *gin.RouterGroup) {
	api := r.Group("/product")
	api.GET("/market", route.SearchProductMarket)
	api.GET("/total", route.SearchSmartProduct)

	api.GET("/me", route.SearchProductWithOwner)
	api.GET("/link/:id", route.SearchProductWithOwner)
	api.GET("/public", route.SearchProductWithOwner)
	api.POST("/new", route.NewProduct)
	api.POST("/suggest", route.NewProductWithDetect)
	api.GET("/:id", route.Detail)
	api.GET("/:id/posts", route.Posts)
	api.GET("/:id/members", route.Members)
	api.GET("/:id/history", route.History)
	api.PUT("/edit/:id", route.Edit)
	api.PUT("/publish", route.PublishProduct)
	api.PUT("/revert", route.RevertProduct)

	api.PUT("/archived", route.Archived)
	api.DELETE("/:id", route.Delete)
	api.PUT("/sale-status", route.SaleStatusUpdate)
	api.PUT("/rent-status", route.RentStatusUpdate)
	api.PUT("/sale-visibility", route.SaleVisibilityUpdate)
	api.PUT("/rent-visibility", route.RentVisibilityUpdate)
	api.POST("/deposite", route.Deposite)
	api.DELETE("/deposite/{productId}", route.CancelDeposite)

	// api.DELETE("/:id", route.Delete)
}

// @Summary Tạo sản phẩm mới
// @Description API để tạo một sản phẩm mới
// @Tags User
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param product body data.NewProductRequest true "Thông tin sản phẩm"
// @Success 201 {object} data.NewProductRequest "Sản phẩm được tạo thành công"
// @Router /v1/bdspro/user/product/new [post]
func (route *UProductRouter) NewProduct(c *gin.Context) {
	body, err := _utils.ParseBody[dto.ProductSaveRequest](c)
	if err != nil {
		_routes.RouteResult(c, nil, err)
		return
	}
	result, err := route.UserProductService.CreateProduct(c, body)
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
// @Router /v1/bdspro/user/product/{id} [post]
func (route *UProductRouter) Detail(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		_routes.RouteResult(c, nil, &_routes.Except{
			Code:    400,
			Message: "ID không hợp lệ",
		})
		return
	}

	result, err := route.UserProductService.Detail(c, id)
	_routes.RouteResult(c, result, err)
}

func (route *UProductRouter) Posts(c *gin.Context) {
	var dto dto.PostLinkSearch
	if err := _utils.ParseQuery2(c, &dto); err != nil {
		_routes.RouteResult(c, nil, err)
		return
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		_routes.RouteResult(c, nil, &_routes.Except{
			Code:    400,
			Message: "ID không hợp lệ",
		})
		return
	}

	result, err := route.UserProductService.Posts(c, id, dto)
	_routes.RouteResult(c, result, err)
}

func (route *UProductRouter) Members(c *gin.Context) {
	var dto dto.ProductAccessSearch
	if err := _utils.ParseQuery2(c, &dto); err != nil {
		_routes.RouteResult(c, nil, err)
		return
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		_routes.RouteResult(c, nil, &_routes.Except{
			Code:    400,
			Message: "ID không hợp lệ",
		})
		return
	}

	result, err := route.UserProductService.Members(c, id, dto)
	_routes.RouteResult(c, result, err)
}

func (route *UProductRouter) History(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		_routes.RouteResult(c, nil, &_routes.Except{
			Code:    400,
			Message: "ID không hợp lệ",
		})
		return
	}

	dto := dto.ProductHistorySearch{}
	if err := _utils.ParseQuery2(c, &dto); err != nil {
		_routes.RouteResult(c, nil, &_routes.Except{
			Code:    400,
			Message: "ID không hợp lệ",
		})
		return
	}

	result, total, err := route.UserProductService.History(c, id, dto)
	_routes.RouteResult(c, _routes.ResponseDTO{
		Code:          0,
		Data:          result,
		TotalElements: &total,
	}, err)
}

// @Summary Đăng tin sản phẩm
// @Description API để đăng tin sản phẩm lên hệ thống
// @Tags User
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param publish body data.StatusRequest true "Thông tin đăng tin"
// @Success 200 {object} data.StatusResponse "Tin được đăng thành công"
// @Router /v1/bdspro/user/product/publish [put]
func (route *UProductRouter) PublishProduct(e *gin.Context) {

}

// @Summary Hoàn tác trạng thái sản phẩm
// @Description API để hoàn tác trạng thái sản phẩm về trạng thái trước đó
// @Tags User
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param revert body data.StatusRequest true "Thông tin hoàn tác"
// @Success 200 {object} data.StatusResponse "Trạng thái sản phẩm đã được hoàn tác"
// @Router /v1/bdspro/user/product/revert/:id [put]
func (route *UProductRouter) RevertProduct(e *gin.Context) {

}

// SearchProductWithOwner godoc
// @Summary Lấy danh sách sản phẩm của user hiện tại
// @Description Returns a list of products based on search filters
// @Tags User/Product-search
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param name query string false "Product name"
// @Param code query string false "Product code"
// @Param area query number false "Product area"
// @Param price query integer false "Product price"
// @Param parentId query integer false "null=>all; 0=>parent; >0=>childrent"
// @Param description query string false "Product description"
// @Param note query string false "Additional notes"
// @Param limit query int false "Limit results" default(10)
// @Param skip query int false "Skip results" default(0)
// @Param sort query string false "Sort order" default("price_asc")
// @Router /v1/bdspro/user/product/me [get]
func (route *UProductRouter) SearchProductWithOwner(c *gin.Context) {
	body := dto.ProductSearchRequest{}
	err := _utils.ParseQuery2(c, &body)
	if err != nil {
		_routes.RouteResult(c, nil, err)
		return
	}
	result, err := route.UserProductService.SearchProductWithOwner(c, &body)
	_routes.RouteResult(c, result, err)
}

// @Summary Lấy danh sách sản phẩm trang chủ
// @Description Returns a list of products based on search filters
// @Tags User/Product-search
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
// @Router /v1/bdspro/user/product/market [get]
func (route *UProductRouter) SearchProductMarket(c *gin.Context) {
	// _routes.RouteResult(c, nil, nil)
	body, _ := _utils.ParseQueryToStruct[dto.ProductSearchRequest](c.Request.URL.Query())
	result, err := route.UserProductService.SearchProductMarket(c, body)
	_routes.RouteResult(c, result, err)
}

// @Summary Tìm kiếm tổng hợp
// @Description Returns a list of products based on search filters
// @Tags User/Product-search
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param text query string true "Text"
// @Param size query int false "Sizet results" default(10)
// @Param page query int false "Page" default(0)
// @Param sort query string false "Sort order" default("price_asc")
// @Router /v1/bdspro/user/product/total [get]
func (route *UProductRouter) SearchSmartProduct(c *gin.Context) {
	body, _ := _utils.ParseQueryToStruct[dto.TextSearchRequest](c.Request.URL.Query())
	result, err := route.UserProductService.SmartSearch(c, body.Text, body.GetPage(), body.GetSize())
	_routes.RouteResult(c, result, err)
}

// @Summary Gợi ý các trường thông tin
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
// @Router /v1/bdspro/user/product/suggest [post]
func (route *UProductRouter) NewProductWithDetect(c *gin.Context) {
	_body, _ := _utils.ParseBody[dto.ProductSaveDetect](c)
	result, err := route.UserProductService.CreateProductWithDetect(c, _body)
	_routes.RouteResult(c, result, err)
}

// @Summary Cập nhật sản phẩm
// @Description API này cho phép cập nhật thông tin sản phẩm dựa trên dữ liệu được gửi lên.
// @Tags User
// @Accept json
// @Produce json
// @Param productId path int true "Product ID"
// @Security BearerAuth
// @Param body body dto.ProductSaveRequest true "Thông tin sản phẩm cần cập nhật"
// @Router /v1/bdspro/user/product/{productId} [put]
func (route *UProductRouter) Edit(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	body, err := _utils.ParseBody[dto.ProductSaveRequest](c)
	if err != nil {
		_routes.RouteResult(c, nil, err)
		return
	}
	result, err := route.UserProductService.Update(c, id, body)
	_routes.RouteResult(c, result, err)
}

func (route *UProductRouter) Rent(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	body, err := _utils.ParseBody[dto.ProductSaveRequest](c)
	if err != nil {
		_routes.RouteResult(c, nil, err)
		return
	}
	result, err := route.UserProductService.Rent(c, id, body)
	_routes.RouteResult(c, result, err)
}

// @Summary Cập nhật sản phẩm
// @Description API này cho phép cập nhật thông tin sản phẩm dựa trên dữ liệu được gửi lên.
// @Tags User
// @Accept json
// @Produce json
// @Param productId path int true "Product ID"
// @Security BearerAuth
// @Param body body dto.ArchivedRequest true "Thông tin sản phẩm cần cập nhật"
// @Router /v1/bdspro/user/product/archived [put]
func (route *UProductRouter) Archived(c *gin.Context) {
	body := dto.ArchivedRequest{}

	if err := _utils.ParseBodyWithValidator(c, &body); err != nil {
		_routes.RouteResult(c, nil, err)
		return
	}
	err := route.UserProductService.Archived(c, body)
	_routes.RouteResult(c, nil, err)
}

// @Summary Xóa sản phẩm
// @Description API
// @Tags User
// @Accept json
// @Produce json
// @Param productId path int true "Product ID"
// @Security BearerAuth
// @Router /v1/bdspro/user/product/{productId} [delete]
func (route *UProductRouter) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		_routes.RouteResult(c, nil, err)
		return
	}
	err = route.UserProductService.Delete(c, id)
	_routes.RouteResult(c, nil, err)
}

// @Summary Xóa sản phẩm
// @Description API
// @Tags User
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param body body dto.StatusUpdate true "Thông tin sản phẩm cần cập nhật"
// @Router /v1/bdspro/user/product/rent-status [put]
func (route *UProductRouter) RentStatusUpdate(c *gin.Context) {
	status := dto.StatusUpdate{}

	if err := _utils.ParseBodyWithValidator(c, &status); err != nil {
		_routes.RouteResult(c, nil, err)
		return
	}
	err := route.UserProductService.RentStatusUpdate(c, status)
	_routes.RouteResult(c, nil, err)
}

// @Summary Xóa sản phẩm
// @Description API
// @Tags User
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param body body dto.StatusUpdate true "Thông tin sản phẩm cần cập nhật"
// @Router /v1/bdspro/user/product/sale-status [put]
func (route *UProductRouter) SaleStatusUpdate(c *gin.Context) {
	status := dto.StatusUpdate{}

	if err := _utils.ParseBodyWithValidator(c, &status); err != nil {
		_routes.RouteResult(c, nil, err)
		return
	}
	err := route.UserProductService.SaleStatusUpdate(c, status)
	_routes.RouteResult(c, nil, err)
}

// @Summary Trạng thái sản phẩm
// @Description API
// @Tags User
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param body body dto.VisibilityUpdate true "Thông tin sản phẩm cần cập nhật"
// @Router /v1/bdspro/user/product/rent-visibility [put]
func (route *UProductRouter) RentVisibilityUpdate(c *gin.Context) {
	status := dto.VisibilityUpdate{}

	if err := _utils.ParseBodyWithValidator(c, &status); err != nil {
		_routes.RouteResult(c, nil, err)
		return
	}
	err := route.UserProductService.RentVisibilityUpdate(c, status)
	_routes.RouteResult(c, nil, err)
}

// @Summary Trạng thái sản phẩm
// @Description API
// @Tags User
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param body body dto.VisibilityUpdate true "Thông tin sản phẩm cần cập nhật"
// @Router /v1/bdspro/user/product/sale-visibility [put]
func (route *UProductRouter) SaleVisibilityUpdate(c *gin.Context) {
	status := dto.VisibilityUpdate{}

	if err := _utils.ParseBodyWithValidator(c, &status); err != nil {
		_routes.RouteResult(c, nil, err)
		return
	}
	err := route.UserProductService.SaleVisibilityUpdate(c, status)
	_routes.RouteResult(c, nil, err)
}

// @Summary Cọc
// @Description API
// @Tags User
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param body body dto.DepositeRequest true "Thông tin sản phẩm cần cập nhật"
// @Router /v1/bdspro/user/product/deposite [post]
func (route *UProductRouter) Deposite(c *gin.Context) {
	dto := dto.DepositeRequest{}

	if err := _utils.ParseBodyWithValidator(c, &dto); err != nil {
		_routes.RouteResult(c, nil, err)
		return
	}
	err := route.UserProductService.ProductDeposite(c, dto)
	_routes.RouteResult(c, nil, err)
}

// @Summary Hủy cọc
// @Description API
// @Tags User
// @Accept json
// @Produce json
// @Security BearerAuth
// @Router /v1/bdspro/user/product/deposite/{productId} [delete]
func (route *UProductRouter) CancelDeposite(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("productId"), 10, 64)
	if err != nil {
		_routes.RouteResult(c, nil, err)
		return
	}

	err = route.UserProductService.CancelDeposite(c, id)
	_routes.RouteResult(c, nil, err)
}

// @Summary Chia sẻ co-broker
// @Description API
// @Tags User
// @Accept json
// @Produce json
// @Security BearerAuth
// @Router /v1/bdspro/user/product/access-share [post]
func (route *UProductRouter) Share(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("productId"), 10, 64)
	if err != nil {
		_routes.RouteResult(c, nil, err)
		return
	}

	err = route.UserProductService.CancelDeposite(c, id)
	_routes.RouteResult(c, nil, err)
}
