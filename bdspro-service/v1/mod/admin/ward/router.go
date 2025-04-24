package ward

import (
	"bdspro/v1/models"
	"common/case/crud"

	"github.com/gin-gonic/gin"
)

type WardRouter struct {
	crud.BaseRouter[models.Ward]
}

func NewWardRouter(service *WardService) *WardRouter {
	return &WardRouter{
		BaseRouter: crud.BaseRouter[models.Ward]{
			Service: service.BaseService,
		},
	}
}

func (r *WardRouter) RegisterRoutes(router *gin.RouterGroup) {
	wardGroup := router.Group("/ward")
	{
		wardGroup.POST("", r.BaseRouter.Create)
		wardGroup.PUT(":id", r.BaseRouter.Update)
		wardGroup.DELETE(":id", r.BaseRouter.Delete)
		wardGroup.GET(":id", r.BaseRouter.GetByID)
		wardGroup.GET("", r.BaseRouter.GetAll)
	}
}

// @BasePath /api

// @Summary Tạo tỉnh mới
// @Description Tạo một bản ghi tỉnh mới
// @Tags Admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param Ward body models.Ward true "Thông tin tỉnh"
// @Success 201 {object} models.Ward
// @Failure 400 {object} map[string]string "Dữ liệu không hợp lệ"
// @Failure 500 {object} map[string]string "Không thể tạo dữ liệu"
// @Router /admin/ward [post]
func (r *WardRouter) Create(c *gin.Context) {

}

// @Summary Cập nhật thông tin tỉnh
// @Description Cập nhật thông tin tỉnh theo ID
// @Tags Admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID của tỉnh"
// @Param Ward body models.Ward true "Thông tin tỉnh mới"
// @Success 200 {object} models.Ward
// @Failure 400 {object} map[string]string "ID không hợp lệ"
// @Failure 500 {object} map[string]string "Không thể cập nhật dữ liệu"
// @Router /admin/ward/{id} [put]
func (r *WardRouter) Update(c *gin.Context) {

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
// @Router /admin/ward/{id} [delete]
func (r *WardRouter) Delete(c *gin.Context) {

}

// @Summary Lấy thông tin tỉnh theo ID
// @Description Lấy thông tin chi tiết của một tỉnh
// @Tags Admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID của tỉnh"
// @Success 200 {object} models.Ward
// @Failure 400 {object} map[string]string "ID không hợp lệ"
// @Failure 404 {object} map[string]string "Không tìm thấy dữ liệu"
// @Router /admin/ward/{id} [get]
func (r *WardRouter) GetByID(c *gin.Context) {

}

// @Summary Lấy danh sách các tỉnh
// @Description Trả về danh sách tất cả các tỉnh
// @Tags Admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {array} models.Ward
// @Failure 500 {object} map[string]string "Không thể lấy danh sách"
// @Router /admin/ward [get]
func (r *WardRouter) GetAll(c *gin.Context) {

}
