package a_project

import (
	"bdspro/v1/models"
	"common/case/crud2"

	"github.com/gin-gonic/gin"
)

type ProjectRouter struct {
	crud2.BaseRouter[models.Project, *ProjectService]
}

func NewProjectRouter(service *ProjectService) *ProjectRouter {
	return &ProjectRouter{
		BaseRouter: crud2.BaseRouter[models.Project, *ProjectService]{
			Service: service,
		},
	}
}

func (r *ProjectRouter) RegisterRoutes(router *gin.RouterGroup) {
	ProjectGroup := router.Group("/project")
	{
		ProjectGroup.POST("", r.BaseRouter.Create)
		ProjectGroup.PUT(":id", r.BaseRouter.Update)
		ProjectGroup.DELETE(":id", r.BaseRouter.Delete)
		ProjectGroup.GET(":id", r.BaseRouter.GetByID)
		ProjectGroup.GET("", r.GetAll)
	}
}

// @Summary Tạo Dự án mới
// @Description Thêm một Dự án mới vào hệ thống
// @Tags Admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param Project body models.Project true "Thông tin Dự án"
// @Success 201 {object} models.Project
// @Failure 400 {object} map[string]string "Dữ liệu không hợp lệ"
// @Failure 500 {object} map[string]string "Không thể tạo dữ liệu"
// @Router /admin/project [post]
// func (r *ProjectRouter) Create(c *gin.Context) {}

// @Summary Cập nhật thông tin Dự án
// @Description Cập nhật thông tin Dự án theo ID
// @Tags Admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID của Dự án"
// @Param Project body models.Project true "Thông tin Dự án mới"
// @Success 200 {object} models.Project
// @Failure 400 {object} map[string]string "ID không hợp lệ"
// @Failure 500 {object} map[string]string "Không thể cập nhật dữ liệu"
// @Router /admin/project/{id} [put]
// func (r *ProjectRouter) Update(c *gin.Context) {}

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
// @Router /admin/project/{id} [delete]
// func (r *ProjectRouter) Delete(c *gin.Context) {}

// @Summary Lấy thông tin Dự án theo ID
// @Description Trả về thông tin chi tiết của một Dự án
// @Tags Admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID của Dự án"
// @Success 200 {object} models.Project
// @Failure 400 {object} map[string]string "ID không hợp lệ"
// @Failure 404 {object} map[string]string "Không tìm thấy dữ liệu"
// @Router /admin/project/{id} [get]
// func (r *ProjectRouter) GetByID(c *gin.Context) {}

// @Summary Lấy danh sách Dự án
// @Description Trả về danh sách tất cả các Dự án
// @Tags Admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {array} models.Project
// @Failure 500 {object} map[string]string "Không thể lấy danh sách"
// @Router /admin/project [get]
// func (r *ProjectRouter) GetAll(c *gin.Context) {}
