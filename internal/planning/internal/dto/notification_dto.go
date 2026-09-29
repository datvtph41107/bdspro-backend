package dto

import "time"

// NotificationDTO represents notification data
type NotificationDTO struct {
	ID        uint64                 `json:"id"`
	UserID    uint64                 `json:"user_id"`
	Title     string                 `json:"title"`
	Message   string                 `json:"message"`
	Type      string                 `json:"type"`
	Priority  string                 `json:"priority"`
	Data      map[string]interface{} `json:"data"`
	IsRead    bool                   `json:"is_read"`
	CreatedAt time.Time              `json:"created_at"`
	UpdatedAt time.Time              `json:"updated_at"`
}

// NotificationBatchDTO represents batch notification data
type NotificationBatchDTO struct {
	UserIDs  []uint64               `json:"user_ids"`
	Title    string                 `json:"title"`
	Message  string                 `json:"message"`
	Type     string                 `json:"type"`
	Priority string                 `json:"priority"`
	Data     map[string]interface{} `json:"data"`
}

// HistoryDTO represents notification history data
type HistoryDTO struct {
	ID             uint64    `json:"id"`
	UserID         uint64    `json:"user_id"`
	NotificationID uint64    `json:"notification_id"`
	Action         string    `json:"action"`
	CreatedAt      time.Time `json:"created_at"`
}

type SendNotificationRequest struct {
	UserID     uint64   `json:"userId"`
	Title      string   `json:"title"`
	Message    []string `json:"message"`
	Type       uint32   `json:"type"`
	TargetID   uint64   `json:"targetId,omitempty"`
	AttachData []string `json:"attachData,omitempty"`
}
