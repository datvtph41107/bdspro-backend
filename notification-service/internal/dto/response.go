package dto

import (
	"noti/enums"
	"noti/internal/domain"
	"time"
)

// NotificationResponse dùng để trả về thông tin thông báo
type NotificationResponse struct {
	ID        uint64         `json:"id"`
	UserID    uint64         `json:"user_id"`
	Title     string         `json:"title"`
	Message   string         `json:"message"`
	Link      string         `json:"link,omitempty"`
	Type      enums.TypeEnum `json:"type"`
	IsRead    bool           `json:"isRead,omitempty"`
	Timestamp time.Time      `json:"timestamp,omitempty"`
}

// SendNotiResponse dùng để trả về kết quả khi gửi thông báo
type SendNotiResponse struct {
	Response string `json:"response,omitempty"`
	Error    string `json:"error,omitempty"`
}

// NewSendNotiResponse tạo một `SendNotiResponse` từ `HistoryEntity`
func NewSendNotiResponse(e domain.HistoryEntity) SendNotiResponse {
	return SendNotiResponse{
		Response: e.Response,
		Error:    e.Error,
	}
}
