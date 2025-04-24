package a_build_room

import (
	"bdspro/v1/models"
	"common/case/crud"

	"github.com/gin-gonic/gin"
)

type BuildRoomRouter struct {
	crud.BaseRouter[models.BuildRoom]
}

func NewBuildRoomRouter(service *BuildRoomService) *BuildRoomRouter {
	return &BuildRoomRouter{
		BaseRouter: crud.BaseRouter[models.BuildRoom]{
			Service: service.BaseService,
		},
	}
}

func (r *BuildRoomRouter) RegisterRoutes(router *gin.RouterGroup) {
	group := router.Group("/build-room")
	{
		group.POST("", r.BaseRouter.Create)
		group.PUT(":id", r.BaseRouter.Update)
		group.DELETE(":id", r.BaseRouter.Delete)
		group.GET(":id", r.BaseRouter.GetByID)
		group.GET("", r.BaseRouter.GetAll)
	}
}

// @Summary Tạo Phòng chung cho các tầng mới
// @Description Thêm một Phòng chung cho các tầng mới vào hệ thống
// @Tags Admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param BuildRoom body models.BuildRoom true "Thông tin Phòng chung cho các tầng"
// @Success 201 {object} models.BuildRoom
// @Failure 400 {object} map[string]string "Dữ liệu không hợp lệ"
// @Failure 500 {object} map[string]string "Không thể tạo dữ liệu"
// @Router /admin/build-room [post]
func (r *BuildRoomRouter) Create(c *gin.Context) {}

// @Summary Cập nhật thông tin Phòng chung cho các tầng
// @Description Cập nhật thông tin Phòng chung cho các tầng theo ID
// @Tags Admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID của Phòng chung cho các tầng"
// @Param Project body models.BuildRoom true "Thông tin Phòng chung cho các tầng mới"
// @Success 200 {object} models.BuildRoom
// @Failure 400 {object} map[string]string "ID không hợp lệ"
// @Failure 500 {object} map[string]string "Không thể cập nhật dữ liệu"
// @Router /admin/build-room/{id} [put]
func (r *BuildRoomRouter) Update(c *gin.Context) {}

// @Summary Xóa Phòng chung cho các tầng
// @Description Xóa Phòng chung cho các tầng theo ID (soft delete)
// @Tags Admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID của Phòng chung cho các tầng"
// @Success 200 {object} map[string]bool "Trạng thái xóa"
// @Failure 400 {object} map[string]string "ID không hợp lệ"
// @Failure 500 {object} map[string]string "Không thể xóa dữ liệu"
// @Router /admin/build-room/{id} [delete]
func (r *BuildRoomRouter) Delete(c *gin.Context) {}

// @Summary Lấy thông tin Phòng chung cho các tầng theo ID
// @Description Trả về thông tin chi tiết của một Phòng chung cho các tầng
// @Tags Admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID của Phòng chung cho các tầng"
// @Success 200 {object} models.BuildRoom
// @Failure 400 {object} map[string]string "ID không hợp lệ"
// @Failure 404 {object} map[string]string "Không tìm thấy dữ liệu"
// @Router /admin/build-room/{id} [get]
func (r *BuildRoomRouter) GetByID(c *gin.Context) {}

// @Summary Lấy danh sách Phòng chung cho các tầng
// @Description Trả về danh sách tất cả các Phòng chung cho các tầng
// @Tags Admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {array} models.BuildRoom
// @Failure 500 {object} map[string]string "Không thể lấy danh sách"
// @Router /admin/build-room [get]
func (r *BuildRoomRouter) GetAll(c *gin.Context) {}
