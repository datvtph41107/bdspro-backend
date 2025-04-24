package a_product

import (
	"bdspro/v1/models"
	"common/case/crud2"

	"github.com/gin-gonic/gin"
)

type AdminProductRouter struct {
	crud2.BaseRouter[models.Product, *AdminProductService]
}

func NewAdminProductRouter(service *AdminProductService) *AdminProductRouter {
	return &AdminProductRouter{
		BaseRouter: crud2.BaseRouter[models.Product, *AdminProductService]{
			Service: service,
		},
	}
}

func (r *AdminProductRouter) RegisterRoutes(router *gin.RouterGroup) {
	AdminProductGroup := router.Group("/product")
	{
		AdminProductGroup.POST("", r.Create)
		AdminProductGroup.PUT(":id", r.Update)
		AdminProductGroup.DELETE(":id", r.Delete)
		AdminProductGroup.GET(":id", r.GetByID)
		AdminProductGroup.GET("", r.GetAll)
	}
}

// @Summary Tạo Dự án mới
// @Description Thêm một Dự án mới vào hệ thống
// @Tags Admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param AdminProduct body models.AdminProduct true "Thông tin Dự án"
// @Success 201 {object} models.AdminProduct
// @Failure 400 {object} map[string]string "Dữ liệu không hợp lệ"
// @Failure 500 {object} map[string]string "Không thể tạo dữ liệu"
// @Router /admin/AdminProduct [post]
// func (r *AdminProductRouter) Create(c *gin.Context) {}

// @Summary Cập nhật thông tin Dự án
// @Description Cập nhật thông tin Dự án theo ID
// @Tags Admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID của Dự án"
// @Param AdminProduct body models.AdminProduct true "Thông tin Dự án mới"
// @Success 200 {object} models.AdminProduct
// @Failure 400 {object} map[string]string "ID không hợp lệ"
// @Failure 500 {object} map[string]string "Không thể cập nhật dữ liệu"
// @Router /admin/AdminProduct/{id} [put]
// func (r *AdminProductRouter) Update(c *gin.Context) {}

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
// @Router /admin/AdminProduct/{id} [delete]
// func (r *AdminProductRouter) Delete(c *gin.Context) {}

// @Summary Lấy thông tin Dự án theo ID
// @Description Trả về thông tin chi tiết của một Dự án
// @Tags Admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID của Dự án"
// @Success 200 {object} models.AdminProduct
// @Failure 400 {object} map[string]string "ID không hợp lệ"
// @Failure 404 {object} map[string]string "Không tìm thấy dữ liệu"
// @Router /admin/AdminProduct/{id} [get]
// func (r *AdminProductRouter) GetByID(c *gin.Context) {}

// @Summary Lấy danh sách Dự án
// @Description Trả về danh sách tất cả các Dự án
// @Tags Admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {array} models.AdminProduct
// @Failure 500 {object} map[string]string "Không thể lấy danh sách"
// @Router /admin/AdminProduct [get]
// func (r *AdminProductRouter) GetAll(c *gin.Context) {}
