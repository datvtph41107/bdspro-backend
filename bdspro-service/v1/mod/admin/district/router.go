package district

import (
	"bdspro/v1/models"
	"common/case/crud"

	"github.com/gin-gonic/gin"
)

type DistrictRouter struct {
	crud.BaseRouter[models.District]
}

func NewDistrictRouter(service *DistrictService) *DistrictRouter {
	return &DistrictRouter{
		BaseRouter: crud.BaseRouter[models.District]{
			Service: service.BaseService,
		},
	}
}

func (r *DistrictRouter) RegisterRoutes(router *gin.RouterGroup) {
	districtGroup := router.Group("/districts")
	{
		districtGroup.POST("", r.BaseRouter.Create)
		districtGroup.PUT(":id", r.BaseRouter.Update)
		districtGroup.DELETE(":id", r.BaseRouter.Delete)
		districtGroup.GET(":id", r.BaseRouter.GetByID)
		districtGroup.GET("", r.BaseRouter.GetAll)
	}
}

// @BasePath /api

// @Summary Tạo quận/huyện mới
// @Description Thêm một quận/huyện mới vào hệ thống
// @Tags Admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param district body models.District true "Thông tin quận/huyện"
// @Success 201 {object} models.District
// @Failure 400 {object} map[string]string "Dữ liệu không hợp lệ"
// @Failure 500 {object} map[string]string "Không thể tạo dữ liệu"
// @Router /v1/bdspro/admin/districts [post]
func (r *DistrictRouter) Create(c *gin.Context) {}

// @Summary Cập nhật thông tin quận/huyện
// @Description Cập nhật thông tin quận/huyện theo ID
// @Tags Admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID của quận/huyện"
// @Param district body models.District true "Thông tin quận/huyện mới"
// @Success 200 {object} models.District
// @Failure 400 {object} map[string]string "ID không hợp lệ"
// @Failure 500 {object} map[string]string "Không thể cập nhật dữ liệu"
// @Router /v1/bdspro/admin/districts/{id} [put]
func (r *DistrictRouter) Update(c *gin.Context) {}

// @Summary Xóa quận/huyện
// @Description Xóa quận/huyện theo ID (soft delete)
// @Tags Admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID của quận/huyện"
// @Success 200 {object} map[string]bool "Trạng thái xóa"
// @Failure 400 {object} map[string]string "ID không hợp lệ"
// @Failure 500 {object} map[string]string "Không thể xóa dữ liệu"
// @Router /v1/bdspro/admin/districts/{id} [delete]
func (r *DistrictRouter) Delete(c *gin.Context) {}

// @Summary Lấy thông tin quận/huyện theo ID
// @Description Trả về thông tin chi tiết của một quận/huyện
// @Tags Admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID của quận/huyện"
// @Success 200 {object} models.District
// @Failure 400 {object} map[string]string "ID không hợp lệ"
// @Failure 404 {object} map[string]string "Không tìm thấy dữ liệu"
// @Router /v1/bdspro/admin/districts/{id} [get]
func (r *DistrictRouter) GetByID(c *gin.Context) {}

// @Summary Lấy danh sách quận/huyện
// @Description Trả về danh sách tất cả các quận/huyện
// @Tags Admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {array} models.District
// @Failure 500 {object} map[string]string "Không thể lấy danh sách"
// @Router /v1/bdspro/admin/districts [get]
func (r *DistrictRouter) GetAll(c *gin.Context) {}
