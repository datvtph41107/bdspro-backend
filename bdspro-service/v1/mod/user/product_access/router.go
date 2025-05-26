package u_product_access

import (
	"bdspro/internal/enums"
	user_dto "bdspro/v1/mod/user/dto/request"
	"bdspro/v1/models"
	"common/case/crud2"
	_routes "common/routes"
	_utils "common/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

type UserProductAccessRoute struct {
	crud2.BaseRouter[models.ProductAccess, *ProductAccessService]
}

func NewUserProductAccessRoute(service *ProductAccessService) *UserProductAccessRoute {
	return &UserProductAccessRoute{
		BaseRouter: crud2.BaseRouter[models.ProductAccess, *ProductAccessService]{
			Service: service,
		},
	}
}

func (r *UserProductAccessRoute) RegisterRoutes(router *gin.RouterGroup) {
	group := router.Group("/product-access")
	{
		group.POST("/:productId/bulk", r.BulkSave)
		group.POST("/:productId/batch", r.CreateBatch)
		// group.PUT(":id", r.BaseRouter.Update)
		group.DELETE(":id", r.BaseRouter.Delete)
		// group.GET(":id", r.BaseRouter.GetByID)
		group.GET("", r.Search)
		group.GET("/private-fields", r.GetFields)
	}
}

// @Summary Tạo Dự án mới
// @Description Thêm một Dự án mới vào hệ thống
// @Tags User/Product-share
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param productId path int true "Thông tin Dự án"
// @Param ProductAccess body []AccessRequest true "Thông tin Dự án"
// @Success 201 {object} models.ProductAccess
// @Failure 400 {object} map[string]string "Dữ liệu không hợp lệ"
// @Failure 500 {object} map[string]string "Không thể tạo dữ liệu"
// @Router /v1/bdspro/user/product-access/{productId}/bulk [post]
func (r *UserProductAccessRoute) BulkSave(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("productId"), 10, 64)
	if err != nil {
		_routes.RouteResult(c, nil, err)
		return
	}

	body := BulkRequest{}

	if err := _utils.ParseBodyWithValidator(c, &body); err != nil {
		_routes.RouteResult(c, nil, err)
		return
	}

	result, err := r.Service.BulkSave(c, id, body)
	_routes.RouteResult(c, result, err)
}

// @Summary Tạo Dự án mới
// @Description Thêm một Dự án mới vào hệ thống
// @Tags User/Product-share
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param productId path int true "Thông tin Dự án"
// @Param ProductAccess body []AccessRequest true "Thông tin Dự án"
// @Success 201 {object} models.ProductAccess
// @Failure 400 {object} map[string]string "Dữ liệu không hợp lệ"
// @Failure 500 {object} map[string]string "Không thể tạo dữ liệu"
// @Router /v1/bdspro/user/product-access/{productId}/batch [post]
func (r *UserProductAccessRoute) CreateBatch(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("productId"), 10, 64)
	if err != nil {
		_routes.RouteResult(c, nil, err)
		return
	}

	body := []AccessRequest{}
	if err := _utils.ParseBodyWithValidator(c, &body); err != nil {
		_routes.RouteResult(c, nil, err)
		return
	}

	result, err := r.Service.CreateBatch(c, id, body)
	_routes.RouteResult(c, result, err)
}

// // @Summary Cập nhật thông tin Dự án
// // @Description Cập nhật thông tin Dự án theo ID
// // @Tags User/Product-share
// // @Accept json
// // @Produce json
// // @Security BearerAuth
// // @Param id path int true "ID của Dự án"
// // @Param ProductAccess body models.ProductAccess true "Thông tin Dự án mới"
// // @Success 200 {object} models.ProductAccess
// // @Failure 400 {object} map[string]string "ID không hợp lệ"
// // @Failure 500 {object} map[string]string "Không thể cập nhật dữ liệu"
// // @Router /v1/bdspro/user/product-access/{id} [put]
// func (r *UserProductAccessRoute) Update(c *gin.Context) {}

// @Summary Xóa Dự án
// @Description Xóa Dự án theo ID (soft delete)
// @Tags User/Product-share
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID của Dự án"
// @Success 200 {object} map[string]bool "Trạng thái xóa"
// @Failure 400 {object} map[string]string "ID không hợp lệ"
// @Failure 500 {object} map[string]string "Không thể xóa dữ liệu"
// @Router /v1/bdspro/user/product-access/{id} [delete]
func (r *UserProductAccessRoute) Delete(c *gin.Context) {}

// // @Summary Lấy thông tin Dự án theo ID
// // @Description Trả về thông tin chi tiết của một Dự án
// // @Tags User/Product-share
// // @Accept json
// // @Produce json
// // @Security BearerAuth
// // @Param id path int true "ID của Dự án"
// // @Success 200 {object} models.ProductAccess
// // @Failure 400 {object} map[string]string "ID không hợp lệ"
// // @Failure 404 {object} map[string]string "Không tìm thấy dữ liệu"
// // @Router /v1/bdspro/user/product-access/{id} [get]
// func (r *UserProductAccessRoute) GetByID(c *gin.Context) {}

// @Summary Lấy danh sách Dự án
// @Description Trả về danh sách tất cả các Dự án
// @Tags User/Product-share
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {array} models.ProductAccess
// @Failure 500 {object} map[string]string "Không thể lấy danh sách"
// @Router /v1/bdspro/user/product-access [get]
func (r *UserProductAccessRoute) Search(c *gin.Context) {
	var dto user_dto.ProductAccessSearch
	if err := _utils.ParseQuery2(c, &dto); err != nil {
		_routes.RouteResult(c, nil, err)
		return
	}

	results, err := r.Service.Search(c, dto)
	_routes.RouteResult(c, results, err)
}

// @Summary Lấy danh sách Dự án
// @Description Trả về danh sách tất cả các Dự án
// @Tags User/Product-share
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {array} models.ProductAccess
// @Failure 500 {object} map[string]string "Không thể lấy danh sách"
// @Router /v1/bdspro/user/product-access/private-fields [get]
func (r *UserProductAccessRoute) GetFields(c *gin.Context) {
	// fields := []struct {
	// 	Value string `json:"value"`
	// 	Label string `json:"label"`
	// }{
	// 	{Value: "importPrice", Label: "Giá nhập"},
	// 	{Value: "operatingCost", Label: "Chi phí vận hành"},
	// }

	_routes.RouteResult(c, enums.PrivateFields, nil)
}
