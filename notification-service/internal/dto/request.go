package dto

import (
	"errors"
	"notification/internal/enums"
	"time"

	_dto "common/domain/dto"
	_enum "common/domain/enum"
	shared_enum "pb/enums"
)

type NotiNewRequest struct {
	Avatar  string `json:"avatar,omitempty"`
	OwnerID uint64 `json:"ownerId" binding:"required"`
	// @deprecated: không dùng nữa thay bằng OwnerOf
	OwnerType    enums.OwnerType `json:"ownerType" binding:"required"`
	OwnerOf      _enum.EOwnerOf  `json:"ownerOf" binding:"required"`
	Title        string          `json:"title" binding:"required"`
	Message      []string        `json:"message" binding:"required"`
	AttachData   []string        `json:"attachData,omitempty"`
	SendToDevice bool            `json:"sendToDevice,omitempty"`
	// id của post, comment, deal, etc.
	// kết hợp với notificationType để không tạo quá nhiều thông báo.
	// ví dụ 10 user cùng like 1 bài thì chỉ only 1 thông báo thôi
	TargetID         uint64                  `json:"targetId" binding:"required"`
	NotificationType _enum.ENotificationType `json:"notificationType" binding:"required"`
	IsMerge          bool                    `json:"isMerge,omitempty"`

	PrefixMessage string         `json:"prefixMessage,omitempty"`
	Link          string         `json:"link,omitempty"`
	Type          enums.TypeEnum `json:"type" binding:"required"`
	VisibleAt     *time.Time     `json:"visibleAt,omitempty"`
}

// SearchNotiRequest đại diện cho request tìm kiếm thông báo
type SearchNotiRequest struct {
	_dto.Pagable
	ID      *uint64 `json:"id,omitempty"`
	UserID  *uint64 `json:"user_id,omitempty"`
	Message *string `json:"message,omitempty"`
	Type    *string `json:"type,omitempty"` // 1: request, 2: follow
	IsRead  *bool   `json:"is_read,omitempty"`
}

// SendNotiRequest đại diện cho request gửi thông báo
type SendNotiRequest struct {
	Target string `json:"target"` // Token của thiết bị nhận
	Title  string `json:"title"`
	Body   string `json:"body"`
	Image  string `json:"image"`
}

// SendNewNotiRequest khởi tạo một đối tượng SendNotiRequest
func SendNewNotiRequest(target string, title string, body string, image string) *SendNotiRequest {
	return &SendNotiRequest{
		Target: target,
		Title:  title,
		Body:   body,
		Image:  image,
	}
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

type PushTokenRequest struct {
	Token string `json:"token" binding:"required"`
	Title string `json:"title" binding:"required"`
	Body  string `json:"body" binding:"required"`
}

type PushTopicRequest struct {
	Topic string `json:"topic" binding:"required"`
	Title string `json:"title" binding:"required"`
	Body  string `json:"body" binding:"required"`
}

type PushByUserIdRequest struct {
	UserID uint64            `json:"userId" binding:"required"`
	Title  string            `json:"title" binding:"required"`
	Body   string            `json:"body" binding:"required"`
	Data   map[string]string `json:"data,omitempty"`
}

type HistorySearchDTO struct {
	_dto.Pagable
	UserID     uint64                      `form:"userId"`
	FromDate   *time.Time                  `form:"fromDate"`
	ToDate     *time.Time                  `form:"toDate"`
	ActionType *shared_enum.EHistory       `form:"actionType"`
	TargetType *shared_enum.ETargetHistory `form:"targetType"`
	TargetId   *uint64                     `form:"targetId"`
}
