package province

import (
	"bdspro/v1/models"
	"common/case/crud"

	"github.com/gin-gonic/gin"
)

type ProvinceRouter struct {
	crud.BaseRouter[models.Province]
}

func NewProvinceRouter(service *ProvinceService) *ProvinceRouter {
	return &ProvinceRouter{
		BaseRouter: crud.BaseRouter[models.Province]{
			Service: service.BaseService,
		},
	}
}

func (r *ProvinceRouter) RegisterRoutes(router *gin.RouterGroup) {
	provinceGroup := router.Group("/provinces")
	{
		provinceGroup.POST("", r.BaseRouter.Create)
		provinceGroup.PUT(":id", r.BaseRouter.Update)
		provinceGroup.DELETE(":id", r.BaseRouter.Delete)
		provinceGroup.GET(":id", r.BaseRouter.GetByID)
		provinceGroup.GET("", r.BaseRouter.GetAll)
	}
}

// @BasePath /api

// @Summary Tạo tỉnh mới
// @Description Tạo một bản ghi tỉnh mới
// @Tags Admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param province body models.Province true "Thông tin tỉnh"
// @Success 201 {object} models.Province
// @Failure 400 {object} map[string]string "Dữ liệu không hợp lệ"
// @Failure 500 {object} map[string]string "Không thể tạo dữ liệu"
// @Router /admin/provinces [post]
func (r *ProvinceRouter) Create(c *gin.Context) {

}

// @Summary Cập nhật thông tin tỉnh
// @Description Cập nhật thông tin tỉnh theo ID
// @Tags Admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID của tỉnh"
// @Param province body models.Province true "Thông tin tỉnh mới"
// @Success 200 {object} models.Province
// @Failure 400 {object} map[string]string "ID không hợp lệ"
// @Failure 500 {object} map[string]string "Không thể cập nhật dữ liệu"
// @Router /admin/provinces/{id} [put]
func (r *ProvinceRouter) Update(c *gin.Context) {

}

// @Summary Xóa tỉnh
// @Description Xóa tỉnh theo ID (soft delete)
// @Tags Admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID của tỉnh"
// @Success 200 {object} map[string]bool "Trạng thái xóa"
// @Failure 400 {object} map[string]string "ID không hợp lệ"
// @Failure 500 {object} map[string]string "Không thể xóa dữ liệu"
// @Router /admin/provinces/{id} [delete]
func (r *ProvinceRouter) Delete(c *gin.Context) {

}

// @Summary Lấy thông tin tỉnh theo ID
// @Description Lấy thông tin chi tiết của một tỉnh
// @Tags Admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID của tỉnh"
// @Success 200 {object} models.Province
// @Failure 400 {object} map[string]string "ID không hợp lệ"
// @Failure 404 {object} map[string]string "Không tìm thấy dữ liệu"
// @Router /admin/provinces/{id} [get]
func (r *ProvinceRouter) GetByID(c *gin.Context) {

}

// @Summary Lấy danh sách các tỉnh
// @Description Trả về danh sách tất cả các tỉnh
// @Tags Admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {array} models.Province
// @Failure 500 {object} map[string]string "Không thể lấy danh sách"
// @Router /admin/provinces [get]
func (r *ProvinceRouter) GetAll(c *gin.Context) {

}
