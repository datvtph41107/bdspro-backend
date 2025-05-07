package a_product_access

import (
	"bdspro/enums"
	"bdspro/v1/models"
	"common/case/crud2"

	_routes "common/routes"

	"github.com/gin-gonic/gin"
)

type ProductAccessRouter struct {
	crud2.BaseRouter[models.ProductAccess, *ProductAccessService]
}

func NewProductAccessRouter(service *ProductAccessService) *ProductAccessRouter {
	return &ProductAccessRouter{
		BaseRouter: crud2.BaseRouter[models.ProductAccess, *ProductAccessService]{
			Service: service,
		},
	}
}

func (r *ProductAccessRouter) RegisterRoutes(router *gin.RouterGroup) {
	group := router.Group("/product-access")
	{
		group.POST("", r.BaseRouter.Create)
		group.PUT(":id", r.BaseRouter.Update)
		group.DELETE(":id", r.BaseRouter.Delete)
		group.GET(":id", r.BaseRouter.GetByID)
		group.GET("", r.BaseRouter.GetAll)
		group.GET("/private-fields", r.GetFields)
	}
}

// @Summary Tạo Dự án mới
// @Description Thêm một Dự án mới vào hệ thống
// @Tags Admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param ProductAccess body models.ProductAccess true "Thông tin Dự án"
// @Success 201 {object} models.ProductAccess
// @Failure 400 {object} map[string]string "Dữ liệu không hợp lệ"
// @Failure 500 {object} map[string]string "Không thể tạo dữ liệu"
// @Router /v1/bdspro/admin/product-access [post]
func (r *ProductAccessRouter) Create(c *gin.Context) {}

// @Summary Cập nhật thông tin Dự án
// @Description Cập nhật thông tin Dự án theo ID
// @Tags Admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID của Dự án"
// @Param ProductAccess body models.ProductAccess true "Thông tin Dự án mới"
// @Success 200 {object} models.ProductAccess
// @Failure 400 {object} map[string]string "ID không hợp lệ"
// @Failure 500 {object} map[string]string "Không thể cập nhật dữ liệu"
// @Router /v1/bdspro/admin/product-access/{id} [put]
func (r *ProductAccessRouter) Update(c *gin.Context) {}

// @Summary Xóa Dự án
// @Description Xóa Dự án theo ID (soft delete)
// @Tags Admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID của Dự án"
// @Success 200 {object} map[string]bool "Trạng thái xóa"
// @Failure 400 {object} map[string]string "ID không hợp lệ"
// @Failure 500 {object} map[string]string "Không thể xóa dữ liệu"
// @Router /v1/bdspro/admin/product-access/{id} [delete]
func (r *ProductAccessRouter) Delete(c *gin.Context) {}

// @Summary Lấy thông tin Dự án theo ID
// @Description Trả về thông tin chi tiết của một Dự án
// @Tags Admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID của Dự án"
// @Success 200 {object} models.ProductAccess
// @Failure 400 {object} map[string]string "ID không hợp lệ"
// @Failure 404 {object} map[string]string "Không tìm thấy dữ liệu"
// @Router /v1/bdspro/admin/product-access/{id} [get]
func (r *ProductAccessRouter) GetByID(c *gin.Context) {}

// @Summary Lấy danh sách Dự án
// @Description Trả về danh sách tất cả các Dự án
// @Tags Admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {array} models.ProductAccess
// @Failure 500 {object} map[string]string "Không thể lấy danh sách"
// @Router /v1/bdspro/admin/product-access [get]
func (r *ProductAccessRouter) GetAll(c *gin.Context) {}

// @Summary Lấy danh sách Dự án
// @Description Trả về danh sách tất cả các Dự án
// @Tags Admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {array} models.ProductAccess
// @Failure 500 {object} map[string]string "Không thể lấy danh sách"
// @Router /v1/bdspro/admin/product-access/private-fields [get]
func (r *ProductAccessRouter) GetFields(c *gin.Context) {
	// fields := []struct {
	// 	Value string `json:"value"`
	// 	Label string `json:"label"`
	// }{
	// 	{Value: "importPrice", Label: "Giá nhập"},
	// 	{Value: "operatingCost", Label: "Chi phí vận hành"},
	// }

	_routes.RouteResult(c, enums.PrivateFields, nil)
}
