package data

import (
	"errors"
	"noti/enums"
)

// SendNotiRequest đại diện cho request gửi thông báo
type SendNotiRequest struct {
	Target string `json:"target"` // Token của thiết bị nhận
	Title  string `json:"title"`
	Body   string `json:"body"`
	Image  string `json:"image"`
}

// NewSendNotiRequest khởi tạo một đối tượng SendNotiRequest
func NewSendNotiRequest(target, title, body, image string) *SendNotiRequest {
	return &SendNotiRequest{
		Target: target,
		Title:  title,
		Body:   body,
		Image:  image,
	}
}

// SearchNotiRequest đại diện cho request tìm kiếm thông báo
type SearchNotiRequest struct {
	ID      *uint64 `json:"id,omitempty"`
	UserID  *uint64 `json:"user_id,omitempty"`
	Message *string `json:"message,omitempty"`
	Type    *string `json:"type,omitempty"` // 1: request, 2: follow
	IsRead  *int    `json:"is_read,omitempty"`
}

// NewNotiRequest đại diện cho request tạo thông báo
type NewNotiRequest struct {
	UserID  uint64         `json:"userId" binding:"required"`
	Title   string         `json:"title" binding:"required"`
	Message string         `json:"message" binding:"required"`
	Link    string         `json:"link,omitempty"`
	Type    enums.TypeEnum `json:"type" binding:"required"`
}

// Validate kiểm tra các trường dữ liệu của NewNotiRequest
func (n *NewNotiRequest) Validate() error {
	if n.UserID == 0 {
		return errors.New("userId không được để trống")
	}
	if n.Title == "" {
		return errors.New("title không được để trống")
	}
	if n.Message == "" {
		return errors.New("message không được để trống")
	}
	if n.Type == 0 {
		return errors.New("type không được để trống")
	}
	return nil
}
