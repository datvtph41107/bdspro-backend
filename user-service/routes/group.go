package routes

import (
	_jwt "common/jwt"
	_routes "common/routes"
	user_service "user/services/user"

	"github.com/gin-gonic/gin"
)

type GroupRouter struct {
	service *user_service.GroupService
}

func NewGroupRouter(service *user_service.GroupService) *GroupRouter {
	return &GroupRouter{
		service: service,
	}
}

// func GroupRoutes(r *gin.Engine, router *GroupRouter) {
// 	api := r.Group("/crm/group")

// 	api.Use(_jwt.JWTAuthMiddleware(nil))
// 	api.GET("", router.Get)
// 	api.POST("", router.Create)
// 	api.PUT("/:id", router.Update)
// 	api.DELETE("/:id", router.Delete)
// }

// Khai báo API
func GroupRoutes(r *gin.Engine, router *GroupRouter) {
	api := r.Group("/user/groups")

	api.Use(_jwt.JWTAuthMiddleware(nil))
	api.GET("", router.ListGroups)
	api.POST("", router.CreateGroup)
	api.PUT("/:id", router.UpdateGroup)
	api.DELETE("/:id", router.DeleteGroup)
}

// @Summary Lấy danh sách nhóm
// @Description Lấy danh sách
// @Tags Nhóm
// @Accept json
// @Produce json
// @Security BearerAuth
// @Failure 400 {string} string "Yêu cầu không hợp lệ"
// @Failure 404 {string} string "Không tìm thấy nhóm"
// @Router /groups [get]
// Lấy danh sách nhóm
func (r *GroupRouter) ListGroups(c *gin.Context) {
	result, err := r.service.ListGroups(c)
	_routes.RouteResult(c, result, err)
}

// @Summary Tạo nhóm mới
// @Description Tạo một nhóm mới với thông tin được cung cấp
// @Tags Nhóm
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body user_model.GroupEntity true "Thông tin cập nhật"
// @Failure 400 {string} string "Yêu cầu không hợp lệ"
// @Failure 409 {string} string "Nhóm đã tồn tại"
// @Router /groups [post]
// Thêm nhóm
func (r *GroupRouter) CreateGroup(c *gin.Context) {
	result, err := r.service.CreateGroup(c)
	_routes.RouteResult(c, result, err)
}

// @Summary Cập nhật thông tin nhóm
// @Description Cập nhật thông tin nhóm hiện tại
// @Tags Nhóm
// @Accept json
// @Produce json
// @Param request body user_model.GroupEntity true "Thông tin cập nhật"
// @Param id path int true "ID của lời mời kết bạn"
// @Security BearerAuth
// @Failure 400 {string} string "Yêu cầu không hợp lệ"
// @Failure 404 {string} string "Không tìm thấy nhóm"
// @Router /groups/{id} [put]
// Cập nhật nhóm
func (r *GroupRouter) UpdateGroup(c *gin.Context) {
	result, err := r.service.UpdateGroup(c)
	_routes.RouteResult(c, result, err)
}

// @Summary Xóa nhóm
// @Description Xóa một nhóm khỏi hệ thống
// @Tags Nhóm
// @Accept json
// @Produce json
// @Param id path int true "ID của nhóm"
// @Security BearerAuth
// @Param id path int true "ID của lời mời kết bạn"
// @Success 200 {string} string "Nhóm đã bị xóa"
// @Failure 400 {string} string "Yêu cầu không hợp lệ"
// @Failure 404 {string} string "Không tìm thấy nhóm"
// @Router /groups/{id} [delete]
// Xóa nhóm
func (r *GroupRouter) DeleteGroup(c *gin.Context) {
	err := r.service.DeleteGroup(c)
	_routes.RouteResult(c, nil, err)
}
