package a_amenity

import (
	"bdspro/v1/models"
	"common/case/crud2"

	"github.com/gin-gonic/gin"
)

type AmenityRouter struct {
	crud2.BaseRouter[models.Amenity, *AmenityService]
}

func NewAmenityRouter(service *AmenityService) *AmenityRouter {
	return &AmenityRouter{
		BaseRouter: crud2.BaseRouter[models.Amenity, *AmenityService]{
			Service: service,
		},
	}
}

func (r *AmenityRouter) RegisterRoutes(router *gin.RouterGroup) {
	amenityGroup := router.Group("/amenity")
	{
		amenityGroup.POST("", r.BaseRouter.Create)
		amenityGroup.PUT(":id", r.BaseRouter.Update)
		amenityGroup.DELETE(":id", r.BaseRouter.Delete)
		amenityGroup.GET(":id", r.BaseRouter.GetByID)
		amenityGroup.GET("", r.BaseRouter.GetAll)
	}
}

// @Summary Tạo Tiện ích mới
// @Description Thêm một Tiện ích mới vào hệ thống
// @Tags Admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param Amenity body models.Amenity true "Thông tin Tiện ích"
// @Success 201 {object} models.Amenity
// @Failure 400 {object} map[string]string "Dữ liệu không hợp lệ"
// @Failure 500 {object} map[string]string "Không thể tạo dữ liệu"
// @Router /admin/amenity [post]
func (r *AmenityRouter) Create(c *gin.Context) {}

// @Summary Cập nhật thông tin Tiện ích
// @Description Cập nhật thông tin Tiện ích theo ID
// @Tags Admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID của Tiện ích"
// @Param Amenity body models.Amenity true "Thông tin Tiện ích mới"
// @Success 200 {object} models.Amenity
// @Failure 400 {object} map[string]string "ID không hợp lệ"
// @Failure 500 {object} map[string]string "Không thể cập nhật dữ liệu"
// @Router /admin/amenity/{id} [put]
func (r *AmenityRouter) Update(c *gin.Context) {}

// @Summary Xóa Tiện ích
// @Description Xóa Tiện ích theo ID (soft delete)
// @Tags Admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID của Tiện ích"
// @Success 200 {object} map[string]bool "Trạng thái xóa"
// @Failure 400 {object} map[string]string "ID không hợp lệ"
// @Failure 500 {object} map[string]string "Không thể xóa dữ liệu"
// @Router /admin/amenity/{id} [delete]
func (r *AmenityRouter) Delete(c *gin.Context) {}

// @Summary Lấy thông tin Tiện ích theo ID
// @Description Trả về thông tin chi tiết của một Tiện ích
// @Tags Admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID của Tiện ích"
// @Success 200 {object} models.Amenity
// @Failure 400 {object} map[string]string "ID không hợp lệ"
// @Failure 404 {object} map[string]string "Không tìm thấy dữ liệu"
// @Router /admin/amenity/{id} [get]
func (r *AmenityRouter) GetByID(c *gin.Context) {}

// @Summary Lấy danh sách Tiện ích
// @Description Trả về danh sách tất cả các Tiện ích
// @Tags Admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {array} models.Amenity
// @Failure 500 {object} map[string]string "Không thể lấy danh sách"
// @Router /admin/amenity [get]
func (r *AmenityRouter) GetAll(c *gin.Context) {}
