package user_router

import (
	"bdspro/internal/usecases"

	"github.com/gin-gonic/gin"
)

type UserProjectBuildRouter struct {
	Usecase usecases.ProjectBuildUsecase
}

func NewUserProjectBuildRouter(uc usecases.ProjectBuildUsecase) *UserProjectBuildRouter {
	return &UserProjectBuildRouter{
		Usecase: uc,
	}
}

func (r *UserProjectBuildRouter) RegisterRoutes(router *gin.RouterGroup) {
	group := router.Group("/project-build")
	{
		group.POST("", r.Create)
		group.PUT(":id", r.Update)
		group.DELETE(":id", r.Delete)
		group.GET(":id", r.GetByID)
		group.GET("", r.GetAll)
	}
}

// todo

// @Summary Tạo Tòa nhà trong dự án mới
// @Description Thêm một Tòa nhà trong dự án mới vào hệ thống
// @Tags Admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param UserProjectBuild body domain.ProjectBuild true "Thông tin Tòa nhà trong dự án"
// @Success 201 {object} domain.ProjectBuild
// @Failure 400 {object} map[string]string "Dữ liệu không hợp lệ"
// @Failure 500 {object} map[string]string "Không thể tạo dữ liệu"
// @Router /v1/bdspro/admin/project-build [post]
func (r *UserProjectBuildRouter) Create(c *gin.Context) {}

// @Summary Cập nhật thông tin Tòa nhà trong dự án
// @Description Cập nhật thông tin Tòa nhà trong dự án theo ID
// @Tags Admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID của Tòa nhà trong dự án"
// @Param Project body domain.ProjectBuild true "Thông tin Tòa nhà trong dự án mới"
// @Success 200 {object} domain.ProjectBuild
// @Failure 400 {object} map[string]string "ID không hợp lệ"
// @Failure 500 {object} map[string]string "Không thể cập nhật dữ liệu"
// @Router /v1/bdspro/admin/project-build/{id} [put]
func (r *UserProjectBuildRouter) Update(c *gin.Context) {}

// @Summary Xóa Tòa nhà trong dự án
// @Description Xóa Tòa nhà trong dự án theo ID (soft delete)
// @Tags Admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID của Tòa nhà trong dự án"
// @Success 200 {object} map[string]bool "Trạng thái xóa"
// @Failure 400 {object} map[string]string "ID không hợp lệ"
// @Failure 500 {object} map[string]string "Không thể xóa dữ liệu"
// @Router /v1/bdspro/admin/project-build/{id} [delete]
func (r *UserProjectBuildRouter) Delete(c *gin.Context) {}

// @Summary Lấy thông tin Tòa nhà trong dự án theo ID
// @Description Trả về thông tin chi tiết của một Tòa nhà trong dự án
// @Tags Admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID của Tòa nhà trong dự án"
// @Success 200 {object} domain.ProjectBuild
// @Failure 400 {object} map[string]string "ID không hợp lệ"
// @Failure 404 {object} map[string]string "Không tìm thấy dữ liệu"
// @Router /v1/bdspro/admin/project-build/{id} [get]
func (r *UserProjectBuildRouter) GetByID(c *gin.Context) {}

// @Summary Lấy danh sách Tòa nhà trong dự án
// @Description Trả về danh sách tất cả các Tòa nhà trong dự án
// @Tags Admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {array} domain.ProjectBuild
// @Failure 500 {object} map[string]string "Không thể lấy danh sách"
// @Router /v1/bdspro/admin/project-build [get]
func (r *UserProjectBuildRouter) GetAll(c *gin.Context) {}
