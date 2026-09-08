package domain

import (
	_models "common/models"
	"time"
)

// ApiKeyEntity đại diện cho bảng api_keys trong database
// Lưu trữ thông tin API key phục vụ cho việc xác thực dịch vụ nội bộ.
type ApiKeyEntity struct {
	_models.BaseEntity
	Name        string     `gorm:"size:100;not null" json:"name"`
	AppName     string     `gorm:"size:100;not null" json:"appName"`
	Description string     `gorm:"size:255" json:"description"`
	ApiKey      string     `gorm:"size:255;not null" json:"apiKey"`
	ExpiredAt   *time.Time `json:"expiredAt"`
}

// TableName đặt tên bảng cho GORM
func (ApiKeyEntity) TableName() string {
	return "api_keys"
}
