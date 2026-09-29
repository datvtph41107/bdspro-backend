package dto

import (
	"bdspro/internal/enums"
	"time"
)

// DealHistoryDTO định nghĩa cấu trúc dữ liệu cho lịch sử thương vụ
type DealHistoryDTO struct {
	DealID      uint64                     `json:"deal_id"`            // ID của thương vụ
	ActorID     uint64                     `json:"actor_id"`           // ID người thực hiện hành động
	ActorName   string                     `json:"actor_name"`         // Tên người thực hiện
	ActorAvatar string                     `json:"actor_avatar"`       // Avatar người thực hiện
	ActionType  enums.DealHistoryEventType `json:"action_type"`        // Loại hành động (CREATE_DEAL, UPDATE_DEAL, etc.)
	ActionName  string                     `json:"action_name"`        // Tên hành động hiển thị
	Content     string                     `json:"content"`            // Nội dung chi tiết
	CreatedAt   time.Time                  `json:"created_at"`         // Thời gian tạo
	Metadata    map[string]interface{}     `json:"metadata,omitempty"` // Dữ liệu bổ sung
}

// DealHistoryRequest định nghĩa request cho việc tạo lịch sử
type DealHistoryRequest struct {
	DealID      uint64                 `json:"deal_id"`
	ActorID     uint64                 `json:"actor_id"`
	ActorName   string                 `json:"actor_name"`
	ActorAvatar string                 `json:"actor_avatar"`
	ActionType  string                 `json:"action_type"`
	ActionName  string                 `json:"action_name"`
	Content     string                 `json:"content"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
}

// DealHistoryResponse định nghĩa response từ notification service
type DealHistoryResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	ID      uint64 `json:"id,omitempty"`
}
