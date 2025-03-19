package routes

import (
	"net/http"
	"strconv"

	_db "common/db"
	_jwt "common/jwt"
	_routes "common/routes"
	"noti/data"
	_ "noti/docs"
	"noti/services"

	"github.com/gin-gonic/gin"
)

// NotificationRouter chứa các API cho thông báo
type NotificationRouter struct {
	Service  *services.NotificationService
	Firebase *services.FirebaseService
}

func NewNotificationRouter(service *services.NotificationService, firebase *services.FirebaseService) *NotificationRouter {
	return &NotificationRouter{
		Service:  service,
		Firebase: firebase,
	}
}

func (r *NotificationRouter) RegisterRoutes(router *gin.Engine) {
	notificationGroup := router.Group("/notification")
	{
		notificationGroup.Use(_jwt.JWTAuthMiddleware(nil))
		notificationGroup.GET("/count", r.CountUnread)
		notificationGroup.POST("/new", r.SendNotification)
		notificationGroup.PUT("/read/:id", r.MarkAsRead)
		notificationGroup.GET("", r.Search)
		notificationGroup.POST("/app/push/token", r.SendToToken)
		notificationGroup.POST("/app/push/topic", r.SendToTopic)
	}
}

// CountUnread godoc
// @Summary Số thông báo mới chưa đọc
// @Description Hiển thị số lượng thông báo chưa đọc của người dùng
// @Tags Notifications
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]int "count"
// @Failure 500 {object} map[string]string "error"
// @Router /notification/count [get]
func (r *NotificationRouter) CountUnread(c *gin.Context) {
	count, err := r.Service.Count(c)
	_routes.RouteResult(c, count, err)
}

// SendNotification godoc
// @Summary Tạo thông báo mới
// @Description API này không bắn thông báo tới thiết bị của người dùng, chỉ tạo thông báo
// @Tags Notifications
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param notification body models.NotificationEntity true "Chi tiết thông báo"
// @Success 200 {object} map[string]string "message"
// @Failure 500 {object} map[string]string "error"
// @Router /notification/new [post]
func (r *NotificationRouter) SendNotification(c *gin.Context) {
	var notification data.NewNotiRequest
	if err := c.ShouldBindJSON(&notification); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	result, err := r.Service.CreateNotification(c, notification)
	_routes.RouteResult(c, result, err)
}

// MarkAsRead godoc
// @Summary Đánh dấu thông báo đã đọc
// @Description Đánh dấu thông báo với ID cụ thể là đã đọc
// @Tags Notifications
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID của thông báo"
// @Success 200 {object} map[string]string "message"
// @Failure 500 {object} map[string]string "error"
// @Router /notification/read/{id} [put]
func (r *NotificationRouter) MarkAsRead(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	result, err := r.Service.Read(c, id)
	_routes.RouteResult(c, result, err)
}

// Search godoc
// @Summary Lấy danh sách thông báo
// @Description Trả về danh sách thông báo của người dùng
// @Tags Notifications
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "Trang cần lấy (default = 0)"
// @Param size query int false "Số lượng trên mỗi trang (default = 200)"
// @Success 200 {object} []models.NotificationEntity "Danh sách thông báo"
// @Failure 500 {object} map[string]string "error"
// @Router /notification [get]
func (r *NotificationRouter) Search(c *gin.Context) {
	userID, _ := strconv.ParseUint(c.Query("user_id"), 10, 32)
	notifications, err := r.Service.Search(c, data.SearchNotiRequest{UserID: &userID}, *_db.NewPage(0, 10))
	if err != nil {
		_routes.RouteResult(c, nil, err)
		return
	}
	_routes.RouteResult(c, notifications, err)
}

// SendToToken godoc
// @Summary Gửi thông báo đến một thiết bị qua token
// @Description API chỉ lưu lịch sử, không lưu thông báo
// @Tags Notifications
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body data.SendNotiRequest true "Thông tin thông báo"
// @Success 200 {object} map[string]string "message"
// @Failure 500 {object} map[string]string "error"
// @Router /notification/app/push/token [post]
func (r *NotificationRouter) SendToToken(c *gin.Context) {
	var request data.SendNotiRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		_routes.RouteResult(c, nil, err)
		return
	}
	result, err := r.Firebase.SendToToken(c, &request)
	_routes.RouteResult(c, result, err)
}

// SendToTopic godoc
// @Summary Gửi thông báo đến topic
// @Description API chỉ lưu lịch sử, không lưu thông báo
// @Tags Notifications
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body data.SendNotiRequest true "Thông tin thông báo"
// @Success 200 {object} map[string]string "message"
// @Failure 500 {object} map[string]string "error"
// @Router /notification/app/push/topic [post]
func (r *NotificationRouter) SendToTopic(c *gin.Context) {
	var request data.SendNotiRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		_routes.RouteResult(c, nil, err)
		return
	}
	result, err := r.Firebase.SendToTopic(c, &request)
	_routes.RouteResult(c, result, err)
}
