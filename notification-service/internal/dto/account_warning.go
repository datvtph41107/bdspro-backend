package dto

import (
	"time"
)

// SendAccountWarningRequest đại diện cho request gửi cảnh báo tài khoản
type SendAccountWarningRequest struct {
	TargetID    uint64     `json:"targetId" binding:"required" validate:"required"`
	TargetType  string     `json:"targetType" binding:"required" validate:"required,oneof=user organization group"`
	WarningType string     `json:"warningType" binding:"required" validate:"required,oneof=violation reminder guide security system"`
	Title       string     `json:"title" binding:"required" validate:"required,min=5,max=255"`
	Content     string     `json:"content" binding:"required" validate:"required,min=20"`
	Severity    string     `json:"severity" validate:"omitempty,oneof=low medium high critical"`
	SendEmail   bool       `json:"sendEmail"`
	RelatedID   *uint64    `json:"relatedId"`
	RelatedType string     `json:"relatedType"`
	ExpiresAt   *time.Time `json:"expiresAt"`
	TemplateID  *uint64    `json:"templateId"` // ID của mẫu cảnh báo nếu sử dụng
}

// SendAccountWarningResponse đại diện cho response gửi cảnh báo
type SendAccountWarningResponse struct {
	ID          uint64    `json:"id"`
	TargetID    uint64    `json:"targetId"`
	TargetType  string    `json:"targetType"`
	WarningType string    `json:"warningType"`
	Title       string    `json:"title"`
	Content     string    `json:"content"`
	Severity    string    `json:"severity"`
	Status      string    `json:"status"`
	SentBy      uint64    `json:"sentBy"`
	SentAt      time.Time `json:"sentAt"`
	EmailSent   bool      `json:"emailSent"`
	Message     string    `json:"message"`
}

// GetAccountWarningListRequest đại diện cho request lấy danh sách cảnh báo
type GetAccountWarningListRequest struct {
	Page        int                    `json:"page" validate:"omitempty,min=1"`
	Size        int                    `json:"size" validate:"omitempty,min=1,max=100"`
	TargetID    *uint64                `json:"targetId"`
	TargetType  string                 `json:"targetType" validate:"omitempty,oneof=user organization group"`
	WarningType string                 `json:"warningType" validate:"omitempty,oneof=violation reminder guide security system"`
	Severity    string                 `json:"severity" validate:"omitempty,oneof=low medium high critical"`
	Status      string                 `json:"status" validate:"omitempty,oneof=sent read acknowledged expired"`
	SentBy      *uint64                `json:"sentBy"`
	IsRead      *bool                  `json:"isRead"`
	Filters     map[string]interface{} `json:"filters"`
}

// GetAccountWarningListResponse đại diện cho response danh sách cảnh báo
type GetAccountWarningListResponse struct {
	Data  []*AccountWarningItem `json:"data"`
	Total int64                 `json:"total"`
	Page  int                   `json:"page"`
	Size  int                   `json:"size"`
}

// AccountWarningItem đại diện cho item cảnh báo
type AccountWarningItem struct {
	ID          uint64     `json:"id"`
	TargetID    uint64     `json:"targetId"`
	TargetType  string     `json:"targetType"`
	WarningType string     `json:"warningType"`
	Title       string     `json:"title"`
	Content     string     `json:"content"`
	Severity    string     `json:"severity"`
	Status      string     `json:"status"`
	IsRead      bool       `json:"isRead"`
	ReadAt      *time.Time `json:"readAt"`
	SentBy      uint64     `json:"sentBy"`
	SentAt      time.Time  `json:"sentAt"`
	EmailSent   bool       `json:"emailSent"`
	EmailSentAt *time.Time `json:"emailSentAt"`
	RelatedID   *uint64    `json:"relatedId"`
	RelatedType string     `json:"relatedType"`
	ExpiresAt   *time.Time `json:"expiresAt"`
	CreatedAt   time.Time  `json:"createdAt"`
	UpdatedAt   time.Time  `json:"updatedAt"`
}

// MarkWarningAsReadRequest đại diện cho request đánh dấu đã đọc
type MarkWarningAsReadRequest struct {
	WarningID uint64 `json:"warningId" binding:"required" validate:"required"`
	TargetID  uint64 `json:"targetId" binding:"required" validate:"required"`
}

// MarkWarningAsReadResponse đại diện cho response đánh dấu đã đọc
type MarkWarningAsReadResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// MarkWarningAsAcknowledgedRequest đại diện cho request đánh dấu đã xác nhận
type MarkWarningAsAcknowledgedRequest struct {
	WarningID uint64 `json:"warningId" binding:"required" validate:"required"`
	TargetID  uint64 `json:"targetId" binding:"required" validate:"required"`
}

// MarkWarningAsAcknowledgedResponse đại diện cho response đánh dấu đã xác nhận
type MarkWarningAsAcknowledgedResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// CreateWarningTemplateRequest đại diện cho request tạo mẫu cảnh báo
type CreateWarningTemplateRequest struct {
	Name        string `json:"name" binding:"required" validate:"required,min=3,max=255"`
	WarningType string `json:"warningType" binding:"required" validate:"required,oneof=violation reminder guide security system"`
	Title       string `json:"title" binding:"required" validate:"required,min=5,max=255"`
	Content     string `json:"content" binding:"required" validate:"required,min=20"`
	Severity    string `json:"severity" validate:"omitempty,oneof=low medium high critical"`
	IsActive    bool   `json:"isActive"`
}

// CreateWarningTemplateResponse đại diện cho response tạo mẫu cảnh báo
type CreateWarningTemplateResponse struct {
	ID          uint64    `json:"id"`
	Name        string    `json:"name"`
	WarningType string    `json:"warningType"`
	Title       string    `json:"title"`
	Content     string    `json:"content"`
	Severity    string    `json:"severity"`
	IsActive    bool      `json:"isActive"`
	CreatedBy   uint64    `json:"createdBy"`
	UpdatedBy   uint64    `json:"updatedBy"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

// GetWarningTemplateListRequest đại diện cho request lấy danh sách mẫu
type GetWarningTemplateListRequest struct {
	Page        int    `json:"page" validate:"omitempty,min=1"`
	Size        int    `json:"size" validate:"omitempty,min=1,max=100"`
	WarningType string `json:"warningType" validate:"omitempty,oneof=violation reminder guide security system"`
}

// GetWarningTemplateListResponse đại diện cho response danh sách mẫu
type GetWarningTemplateListResponse struct {
	Data  []*WarningTemplateItem `json:"data"`
	Total int64                  `json:"total"`
	Page  int                    `json:"page"`
	Size  int                    `json:"size"`
}

// WarningTemplateItem đại diện cho item mẫu cảnh báo
type WarningTemplateItem struct {
	ID          uint64    `json:"id"`
	Name        string    `json:"name"`
	WarningType string    `json:"warningType"`
	Title       string    `json:"title"`
	Content     string    `json:"content"`
	Severity    string    `json:"severity"`
	IsActive    bool      `json:"isActive"`
	CreatedBy   uint64    `json:"createdBy"`
	UpdatedBy   uint64    `json:"updatedBy"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

// GetWarningLogListRequest đại diện cho request lấy danh sách log
type GetWarningLogListRequest struct {
	Page      int     `json:"page" validate:"omitempty,min=1"`
	Size      int     `json:"size" validate:"omitempty,min=1,max=100"`
	WarningID *uint64 `json:"warningId"`
	TargetID  *uint64 `json:"targetId"`
}

// GetWarningLogListResponse đại diện cho response danh sách log
type GetWarningLogListResponse struct {
	Data  []*WarningLogItem `json:"data"`
	Total int64             `json:"total"`
	Page  int               `json:"page"`
	Size  int               `json:"size"`
}

// WarningLogItem đại diện cho item log cảnh báo
type WarningLogItem struct {
	ID          uint64    `json:"id"`
	WarningID   uint64    `json:"warningId"`
	TargetID    uint64    `json:"targetId"`
	Action      string    `json:"action"`
	Status      string    `json:"status"`
	Reason      string    `json:"reason"`
	PerformedBy uint64    `json:"performedBy"`
	PerformedAt time.Time `json:"performedAt"`
	IPAddress   string    `json:"ipAddress"`
	UserAgent   string    `json:"userAgent"`
}
