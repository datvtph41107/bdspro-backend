package dto

import "time"

// NotificationWithDetails represents a notification with additional details
type NotificationWithDetails struct {
	ID        uint      `json:"id"`
	UserID    uint      `json:"user_id"`
	Title     string    `json:"title"`
	Message   string    `json:"message"`
	Type      string    `json:"type"`
	Data      string    `json:"data"`
	IsRead    bool      `json:"is_read"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// NotificationRequest represents a request to create notification
type NotificationRequest struct {
	UserID  uint   `json:"user_id" binding:"required"`
	Title   string `json:"title" binding:"required,max=255"`
	Message string `json:"message" binding:"required"`
	Type    string `json:"type" binding:"required,max=50"`
	Data    string `json:"data"`
}

// NotificationFilter represents filter criteria for notifications
type NotificationFilter struct {
	UserID   *uint   `json:"user_id,omitempty"`
	Type     *string `json:"type,omitempty"`
	IsRead   *bool   `json:"is_read,omitempty"`
	IsActive *bool   `json:"is_active,omitempty"`
	DateFrom *string `json:"date_from,omitempty"`
	DateTo   *string `json:"date_to,omitempty"`
}

// NotificationTemplateRequest represents a request to create/update notification template
type NotificationTemplateRequest struct {
	Name      string `json:"name" binding:"required,max=100"`
	Title     string `json:"title" binding:"required,max=255"`
	Message   string `json:"message" binding:"required"`
	Type      string `json:"type" binding:"required,max=50"`
	Variables string `json:"variables"`
	IsActive  bool   `json:"is_active"`
}

// NotificationPreferenceRequest represents a request to update notification preferences
type NotificationPreferenceRequest struct {
	EmailEnabled          bool `json:"email_enabled"`
	PushEnabled           bool `json:"push_enabled"`
	SMSEnabled            bool `json:"sms_enabled"`
	EventNotifications    bool `json:"event_notifications"`
	ReviewNotifications   bool `json:"review_notifications"`
	LocationNotifications bool `json:"location_notifications"`
	SystemNotifications   bool `json:"system_notifications"`
}

// NotificationSummary represents a summary of notifications for a user
type NotificationSummary struct {
	UserID             uint `json:"user_id"`
	TotalNotifications int  `json:"total_notifications"`
	UnreadCount        int  `json:"unread_count"`
	ReadCount          int  `json:"read_count"`
	ArchivedCount      int  `json:"archived_count"`
}

// NotificationListResponse represents a list of notifications with pagination
type NotificationListResponse struct {
	Data  []NotificationWithDetails `json:"data"`
	Total int64                     `json:"total"`
	Page  int                       `json:"page"`
	Size  int                       `json:"size"`
}

// NotificationResponse represents a notification response
type NotificationResponse struct {
	ID        uint      `json:"id"`
	UserID    uint      `json:"user_id"`
	Title     string    `json:"title"`
	Message   string    `json:"message"`
	Type      string    `json:"type"`
	Data      string    `json:"data"`
	IsRead    bool      `json:"is_read"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// NotificationTemplateResponse represents a notification template response
type NotificationTemplateResponse struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	Title     string    `json:"title"`
	Message   string    `json:"message"`
	Type      string    `json:"type"`
	Variables string    `json:"variables"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// NotificationPreferenceResponse represents notification preferences response
type NotificationPreferenceResponse struct {
	UserID                uint      `json:"user_id"`
	EmailEnabled          bool      `json:"email_enabled"`
	PushEnabled           bool      `json:"push_enabled"`
	SMSEnabled            bool      `json:"sms_enabled"`
	EventNotifications    bool      `json:"event_notifications"`
	ReviewNotifications   bool      `json:"review_notifications"`
	LocationNotifications bool      `json:"location_notifications"`
	SystemNotifications   bool      `json:"system_notifications"`
	CreatedAt             time.Time `json:"created_at"`
	UpdatedAt             time.Time `json:"updated_at"`
}
