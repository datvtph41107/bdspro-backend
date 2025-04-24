package a_doc_type

import (
	"bdspro/v1/models"
	"common/case/crud"

	"github.com/gin-gonic/gin"
)

type DocTypeRouter struct {
	crud.BaseRouter[models.DocType]
}

func NewDocTypeRouter(service *DocTypeService) *DocTypeRouter {
	return &DocTypeRouter{
		BaseRouter: crud.BaseRouter[models.DocType]{
			Service: service.BaseService,
		},
	}
}

func (r *DocTypeRouter) RegisterRoutes(router *gin.RouterGroup) {
	DocTypeGroup := router.Group("/doc-type")
	{
		DocTypeGroup.POST("", r.BaseRouter.Create)
		DocTypeGroup.PUT(":id", r.BaseRouter.Update)
		DocTypeGroup.DELETE(":id", r.BaseRouter.Delete)
		DocTypeGroup.GET(":id", r.BaseRouter.GetByID)
		DocTypeGroup.GET("", r.BaseRouter.GetAll)
	}
}

// @Summary Tạo Loại giấy tờ pháp lý mới
// @Description Thêm một Loại giấy tờ pháp lý mới vào hệ thống
// @Tags Admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param DocType body models.DocType true "Thông tin Loại giấy tờ pháp lý"
// @Success 201 {object} models.DocType
// @Failure 400 {object} map[string]string "Dữ liệu không hợp lệ"
// @Failure 500 {object} map[string]string "Không thể tạo dữ liệu"
// @Router /admin/doc-type [post]
func (r *DocTypeRouter) Create(c *gin.Context) {}

// @Summary Cập nhật thông tin Loại giấy tờ pháp lý
// @Description Cập nhật thông tin Loại giấy tờ pháp lý theo ID
// @Tags Admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID của Loại giấy tờ pháp lý"
// @Param DocType body models.DocType true "Thông tin Loại giấy tờ pháp lý mới"
// @Success 200 {object} models.DocType
// @Failure 400 {object} map[string]string "ID không hợp lệ"
// @Failure 500 {object} map[string]string "Không thể cập nhật dữ liệu"
// @Router /admin/doc-type/{id} [put]
func (r *DocTypeRouter) Update(c *gin.Context) {}

// @Summary Xóa Loại giấy tờ pháp lý
// @Description Xóa Loại giấy tờ pháp lý theo ID (soft delete)
// @Tags Admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID của Loại giấy tờ pháp lý"
// @Success 200 {object} map[string]bool "Trạng thái xóa"
// @Failure 400 {object} map[string]string "ID không hợp lệ"
// @Failure 500 {object} map[string]string "Không thể xóa dữ liệu"
// @Router /admin/doc-type/{id} [delete]
func (r *DocTypeRouter) Delete(c *gin.Context) {}

// @Summary Lấy thông tin Loại giấy tờ pháp lý theo ID
// @Description Trả về thông tin chi tiết của một Loại giấy tờ pháp lý
// @Tags Admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID của Loại giấy tờ pháp lý"
// @Success 200 {object} models.DocType
// @Failure 400 {object} map[string]string "ID không hợp lệ"
// @Failure 404 {object} map[string]string "Không tìm thấy dữ liệu"
// @Router /admin/doc-type/{id} [get]
func (r *DocTypeRouter) GetByID(c *gin.Context) {}

// @Summary Lấy danh sách Loại giấy tờ pháp lý
// @Description Trả về danh sách tất cả các Loại giấy tờ pháp lý
// @Tags Admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {array} models.DocType
// @Failure 500 {object} map[string]string "Không thể lấy danh sách"
// @Router /admin/doc-type [get]
func (r *DocTypeRouter) GetAll(c *gin.Context) {}
