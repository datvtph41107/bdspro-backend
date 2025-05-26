package routes

import (
	_jwt "common/jwt"
	_routes "common/routes"
	user_service "user/services/user"

	"github.com/gin-gonic/gin"
)

type FollowRouter struct {
	service *user_service.FollowService
}

func NewFollowRouter(service *user_service.FollowService) *FollowRouter {
	return &FollowRouter{
		service: service,
	}
}

func FollowRoutes(r *gin.Engine, router *FollowRouter) {
	api := r.Group("/user/follow")

	api.Use(_jwt.JWTAuthMiddleware(nil))
	api.GET("/followers", router.FollowerList)
	api.GET("/following", router.FollowingList)
	api.POST("/:profileId", router.Add)
	api.DELETE("/:profileId", router.Remove)
	// api.PUT("/:id", router.Update)
}

// @Summary Lấy danh sách người theo dõi của người dùng
// @Description Lấy danh sách những người theo dõi người dùng cụ thể
// @Tags Theo dõi
// @Accept json
// @Produce json
// @Security BearerAuth
// @Failure 400 {string} string "Yêu cầu không hợp lệ"
// @Failure 404 {string} string "Không tìm thấy người dùng"
// @Router /follow/followers [get]
func (r *FollowRouter) FollowerList(c *gin.Context) {
	result, err := r.service.FollowerUser(c)
	_routes.RouteResult(c, result, err)
}

// @Summary Lấy danh sách người mà người dùng đang theo dõi
// @Description Lấy danh sách những người mà người dùng đang theo dõi
// @Tags Theo dõi
// @Accept json
// @Produce json
// @Security BearerAuth
// @Failure 400 {string} string "Yêu cầu không hợp lệ"
// @Failure 404 {string} string "Không tìm thấy người dùng"
// @Router /follow/following [get]
func (r *FollowRouter) FollowingList(c *gin.Context) {
	result, err := r.service.FollowingUser(c)
	_routes.RouteResult(c, result, err)
}

// @Summary Thêm người vào danh sách theo dõi
// @Description Thêm một người vào danh sách theo dõi của người dùng
// @Tags Theo dõi
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param profileId path int true "ID của người dùng"
// @Success 200 {string} string "Đã thêm người vào danh sách theo dõi"
// @Failure 400 {string} string "Yêu cầu không hợp lệ"
// @Failure 404 {string} string "Không tìm thấy người dùng hoặc người dùng"
// @Router /follow/{profileId} [post]
func (r *FollowRouter) Add(c *gin.Context) {
	result, err := r.service.FollowUser(c)

	_routes.RouteResult(c, result, err)
}

// @Summary Xóa người khỏi danh sách theo dõi
// @Description Xóa một người khỏi danh sách theo dõi của người dùng
// @Tags Theo dõi
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param profileId path int true "ID của người dùng"
// @Success 200 {string} string "Đã xóa người khỏi danh sách theo dõi"
// @Failure 400 {string} string "Yêu cầu không hợp lệ"
// @Failure 404 {string} string "Không tìm thấy người dùng hoặc người dùng"
// @Router /follow/{profileId} [delete]
func (r *FollowRouter) Remove(c *gin.Context) {
	result, err := r.service.UnfollowUser(c)

	_routes.RouteResult(c, result, err)
}
