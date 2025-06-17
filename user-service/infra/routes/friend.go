package routes

import (
	_routes "common/routes"
	"user/enums"
	user_service "user/internal/services/user"

	"github.com/gin-gonic/gin"
)

type FriendRouter struct {
	service *user_service.FriendService
}

func NewFriendRouter(service *user_service.FriendService) *FriendRouter {
	return &FriendRouter{
		service: service,
	}
}

func FriendRoutes(r *gin.Engine, router *FriendRouter) {
	api := r.Group("/user/friend")

	api.GET("/list", router.FriendList)
	api.POST("/request/:receiverId", router.FriendRequest)
	api.PUT("/request/:requestId/accept", router.FriendAccept)
	api.PUT("/group/:requestId/:groupId", router.FriendUpdateGroup)
	api.DELETE("/request/:requestId/reject", router.FriendReject)
	api.GET("/request/sent", router.FriendSent)
	api.GET("/request/received", router.FriendReceived)
	api.GET("/info", router.FriendInfo)

}

// @Summary Lấy danh sách bạn bè
// @Description Lấy danh sách tất cả bạn bè của người dùng
// @Tags Bạn bè
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {array} dto.FriendDTO "Danh sách bạn bè"
// @Failure 400 {string} string "Yêu cầu không hợp lệ"
// @Failure 404 {string} string "Không tìm thấy"
// @Router /v1/user/friend/list [get]
func (r *FriendRouter) FriendList(c *gin.Context) {
	result, err := r.service.FriendList(c)

	_routes.RouteResult(c, result, err)
}

// @Summary Gửi lời mời kết bạn
// @Description Gửi lời mời kết bạn đến người dùng khác
// @Tags Bạn bè
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param receiverId path int true "ID của người nhận lời mời"
// @Success 200 {string} string "Lời mời đã được gửi"
// @Failure 400 {string} string "Yêu cầu không hợp lệ"
// @Failure 404 {string} string "Không tìm thấy người nhận"
// @Router /v1/user/friend/request/{receiverId} [post]
func (r *FriendRouter) FriendRequest(c *gin.Context) {
	// var dto dto.FriendDTO
	// if err := c.ShouldBindJSON(&dto); err != nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	// 	return
	// }

	result, err := r.service.Request(c)
	_routes.RouteResult(c, result, err)

}

// @Summary Đồng ý kết bạn
// @Description Đồng ý lời mời kết bạn đang chờ
// @Tags Bạn bè
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param requestId path int true "ID của lời mời kết bạn"
// @Success 200 {string} string "Lời mời đã được chấp nhận"
// @Failure 400 {string} string "Yêu cầu không hợp lệ"
// @Failure 404 {string} string "Không tìm thấy lời mời"
// @Router /v1/user/friend/request/{requestId}/accept [put]
func (r *FriendRouter) FriendAccept(c *gin.Context) {
	result, err := r.service.ChangeStatus(c, enums.EFriendStatusAccepted)
	_routes.RouteResult(c, result, err)
}

// @Summary Gán nhóm bạn bè
// @Description Gán nhóm cho bạn bè
// @Tags Bạn bè
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param requestId path int true "ID của lời mời kết bạn"
// @Param groupId path int true "ID của nhóm"
// @Success 200 {string} string "Lời mời đã được chấp nhận"
// @Failure 400 {string} string "Yêu cầu không hợp lệ"
// @Failure 404 {string} string "Không tìm thấy lời mời"
// @Router /v1/user/friend/group/{requestId}/{groupId} [put]
func (r *FriendRouter) FriendUpdateGroup(c *gin.Context) {
	result, err := r.service.ChangeGroup(c)
	_routes.RouteResult(c, result, err)
}

// @Summary Hủy kết bạn
// @Description Hủy lời mời kết bạn đang chờ
// @Tags Bạn bè
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param requestId path int true "ID của lời mời kết bạn"
// @Success 200 {string} string "Lời mời đã bị từ chối"
// @Failure 400 {string} string "Yêu cầu không hợp lệ"
// @Failure 404 {string} string "Không tìm thấy lời mời"
// @Router /v1/user/friend/request/{requestId}/reject [delete]
func (r *FriendRouter) FriendReject(c *gin.Context) {
	result, err := r.service.ChangeStatus(c, enums.EFriendStatusReject)

	_routes.RouteResult(c, result, err)
}

// @Summary Lấy danh sách các lời mời kết bạn đã gửi
// @Description Lấy danh sách các lời mời kết bạn mà người dùng đã gửi
// @Tags Bạn bè
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {array} dto.FriendDTO "Danh sách lời mời đã gửi"
// @Failure 400 {string} string "Yêu cầu không hợp lệ"
// @Failure 404 {string} string "Không tìm thấy"
// @Router /v1/user/friend/request/sent [get]
func (r *FriendRouter) FriendSent(c *gin.Context) {
	result, err := r.service.RequestSent(c)

	_routes.RouteResult(c, result, err)
}

// @Summary Lấy danh sách các lời mời kết bạn đã nhận
// @Description Lấy danh sách các lời mời kết bạn mà người dùng đã nhận
// @Tags Bạn bè
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {array} dto.FriendDTO "Danh sách lời mời đã nhận"
// @Failure 400 {string} string "Yêu cầu không hợp lệ"
// @Failure 404 {string} string "Không tìm thấy"
// @Router /v1/user/friend/request/received [get]
func (r *FriendRouter) FriendReceived(c *gin.Context) {
	result, err := r.service.RequestReceived(c)

	_routes.RouteResult(c, result, err)
}

// @Summary Thông tin số người theo dõi, đang theo dõi và số bạn bè
// @Description Thông tin số người theo dõi, đang theo dõi và số bạn bè
// @Tags Bạn bè
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} dto.FollowInfoRes "Số lượng theo dõi và bạn bè"
// @Failure 400 {string} string "Yêu cầu không hợp lệ"
// @Failure 404 {string} string "Không tìm thấy"
// @Router /v1/user/friend/info [get]
func (r *FriendRouter) FriendInfo(c *gin.Context) {
	result, err := r.service.FriendInfo(c)

	_routes.RouteResult(c, result, err)
}
