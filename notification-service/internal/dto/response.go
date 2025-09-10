package dto

import (
	shared_enum "pb/enums"
	"time"

	"notification/internal/domain"
	"notification/internal/enums"
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
func NewSendNotiResponse(e domain.HistoryFcmEntity) SendNotiResponse {
	return SendNotiResponse{
		Response: e.Response,
		Error:    e.Error,
	}
}

// AdminHistoryDTO dùng để trả về thông tin lịch sử admin
type AdminHistoryDTO struct {
	ID         uint64                     `json:"id"`
	TargetId   uint64                     `json:"targetId"`
	TargetType shared_enum.ETargetHistory `json:"targetType"`
	ActionType shared_enum.EHistory       `json:"actionType"`
	Title      string                     `json:"title,omitempty"`
	Note       []string                   `json:"note,omitempty"`
	PreStage   string                     `json:"preStage,omitempty"`
	AfterStage string                     `json:"afterStage,omitempty"`
	AdminID    uint64                     `json:"adminId"`
	OwnerID    *uint64                    `json:"ownerId,omitempty"`
	OwnerType  shared_enum.EOwnerType     `json:"ownerType"`
	IsInternal bool                       `json:"isInternal"`
	AdminRole  string                     `json:"adminRole,omitempty"`
	IPAddress  string                     `json:"ipAddress,omitempty"`
	UserAgent  string                     `json:"userAgent,omitempty"`
	CreatedAt  string                     `json:"createdAt"`
	UpdatedAt  string                     `json:"updatedAt"`
	// Thêm thông tin user
	AdminName   *string `json:"adminName,omitempty"`
	AdminAvatar *string `json:"adminAvatar,omitempty"`
	// Thêm trạng thái thành công/thất bại
	Status *string `json:"status,omitempty"`
}
