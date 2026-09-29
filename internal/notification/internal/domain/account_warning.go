package domain

import (
	_models "common/models"
	"notification/internal/enums"
	"time"
)

// AccountWarningEntity đại diện cho cảnh báo tài khoản
type AccountWarningEntity struct {
	_models.BaseEntity
	ID          uint64                  `gorm:"primaryKey;autoIncrement" json:"id"`
	TargetID    uint64                  `gorm:"index;not null" json:"targetId"`           // ID người nhận cảnh báo
	TargetType  enums.TargetTypeEnum    `gorm:"size:20;not null" json:"targetType"`       // user, organization, group
	WarningType enums.WarningTypeEnum   `gorm:"size:50;not null" json:"warningType"`      // violation, reminder, guide, security, system
	Title       string                  `gorm:"size:255;not null" json:"title"`           // Tiêu đề cảnh báo
	Content     string                  `gorm:"type:text;not null" json:"content"`        // Nội dung cảnh báo
	Severity    enums.SeverityEnum      `gorm:"size:20;default:'medium'" json:"severity"` // low, medium, high, critical
	Status      enums.WarningStatusEnum `gorm:"size:20;default:'sent'" json:"status"`     // sent, read, acknowledged
	IsRead      bool                    `gorm:"default:false" json:"isRead"`              // Đã đọc chưa
	ReadAt      *time.Time              `json:"readAt"`                                   // Thời gian đọc
	SentBy      uint64                  `gorm:"not null" json:"sentBy"`                   // ID người gửi
	SentAt      time.Time               `gorm:"autoCreateTime" json:"sentAt"`             // Thời gian gửi
	EmailSent   bool                    `gorm:"default:false" json:"emailSent"`           // Đã gửi email chưa
	EmailSentAt *time.Time              `json:"emailSentAt"`                              // Thời gian gửi email
	RelatedID   *uint64                 `json:"relatedId"`                                // ID liên quan (post, comment, etc.)
	RelatedType string                  `gorm:"size:50" json:"relatedType"`               // Loại liên quan
	ExpiresAt   *time.Time              `json:"expiresAt"`                                // Thời gian hết hạn
	TemplateID  *uint64                 `gorm:"index" json:"templateId"`                  // ID mẫu cảnh báo nếu sử dụng
	DeletedAt   *time.Time              `gorm:"index" json:"deletedAt"`                   // Soft delete
}

// TableName đặt tên bảng
func (AccountWarningEntity) TableName() string {
	return "account_warnings"
}

// AccountWarningTemplateEntity đại diện cho mẫu cảnh báo
type AccountWarningTemplateEntity struct {
	_models.BaseEntity
	ID          uint64                `gorm:"primaryKey;autoIncrement" json:"id"`
	Name        string                `gorm:"size:255;not null" json:"name"`            // Tên mẫu
	WarningType enums.WarningTypeEnum `gorm:"size:50;not null" json:"warningType"`      // Loại cảnh báo
	Title       string                `gorm:"size:255;not null" json:"title"`           // Tiêu đề mẫu
	Content     string                `gorm:"type:text;not null" json:"content"`        // Nội dung mẫu
	Severity    enums.SeverityEnum    `gorm:"size:20;default:'medium'" json:"severity"` // Mức độ nghiêm trọng
	IsActive    bool                  `gorm:"default:true" json:"isActive"`             // Có hoạt động không
	CreatedBy   uint64                `gorm:"not null" json:"createdBy"`                // Người tạo
	UpdatedBy   uint64                `gorm:"not null" json:"updatedBy"`                // Người cập nhật
	DeletedAt   *time.Time            `gorm:"index" json:"deletedAt"`                   // Soft delete
}

// TableName đặt tên bảng
func (AccountWarningTemplateEntity) TableName() string {
	return "account_warning_templates"
}

// AccountWarningLogEntity đại diện cho log cảnh báo
type AccountWarningLogEntity struct {
	_models.BaseEntity
	ID          uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	WarningID   uint64    `gorm:"index;not null" json:"warningId"`   // ID cảnh báo
	TargetID    uint64    `gorm:"index;not null" json:"targetId"`    // ID người nhận
	Action      string    `gorm:"size:50;not null" json:"action"`    // sent, read, acknowledged
	Status      string    `gorm:"size:20;not null" json:"status"`    // success, failed
	Reason      string    `gorm:"size:500" json:"reason"`            // Lý do
	PerformedBy uint64    `gorm:"not null" json:"performedBy"`       // Người thực hiện
	PerformedAt time.Time `gorm:"autoCreateTime" json:"performedAt"` // Thời gian thực hiện
	IPAddress   string    `gorm:"size:45" json:"ipAddress"`          // IP thực hiện
	UserAgent   string    `gorm:"size:500" json:"userAgent"`         // User agent
}

// TableName đặt tên bảng
func (AccountWarningLogEntity) TableName() string {
	return "account_warning_logs"
}
