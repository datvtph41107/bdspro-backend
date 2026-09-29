package cmd

// CountUnread godoc
// @Summary Số thông báo mới chưa đọc
// @Description Hiển thị số lượng thông báo chưa đọc của người dùng
// @Tags Notifications
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]int "count"
// @Failure 500 {object} map[string]string "error"
// @Router /v2/notification/count [get]
func CountUnread() {
}

// SendNotification godoc
// @Summary Tạo thông báo mới
// @Description API này không bắn thông báo tới thiết bị của người dùng, chỉ tạo thông báo
// @Tags Notifications
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]string "message"
// @Failure 500 {object} map[string]string "error"
// @Router /v2/notification/new [post]
func SendNotification() {
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
// @Router /v2/notification/read/{id} [put]
func MarkAsRead() {
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
// @Failure 500 {object} map[string]string "error"
// @Router /v2/notification [get]
func Search() {
}

// @Summary Gửi thông báo đến một thiết bị qua token
// @Description API chỉ lưu lịch sử, không lưu thông báo
// @Tags Notifications
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]string "message"
// @Failure 500 {object} map[string]string "error"
// @Router /v2/notification/app/push/token [post]
func SendToToken() {
}

// SendToTopic godoc
// @Summary Gửi thông báo đến topic
// @Description API chỉ lưu lịch sử, không lưu thông báo
// @Tags Notifications
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]string "message"
// @Failure 500 {object} map[string]string "error"
// @Router /v2/notification/app/push/topic [post]
func SendToTopic() {

}

// @Summary Xóa thông báo
// @Description Xóa thông báo theo ID
// @Tags Notifications
// @Param id path int true "Notification ID"
// @Security BearerAuth
// @Success 200 {object} gin.H
// @Failure 400 {object} gin.H
// @Failure 404 {object} gin.H
// @Router /v2/notification/remove/{id} [delete]
func DeleteNotification() {

}
