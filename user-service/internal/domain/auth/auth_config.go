package auth

import (
	_models "common/domain/entity"
)

// AuthConfig lưu cấu hình động của hệ thống auth
type AuthConfig struct {
	_models.BaseEntity
	ConfigKey   string `gorm:"column:config_key;type:varchar(20);uniqueIndex;not null" json:"configKey"`
	ConfigValue string `gorm:"column:config_value;type:text" json:"configValue"`
	IsActive    bool   `gorm:"column:is_active;default:true" json:"isActive"`
}

// TableName đặt tên bảng cho GORM
func (AuthConfig) TableName() string {
	return "auth_config"
}

// Các config key constants
const (
	// ZNS Config Keys
	ConfigKeyZNSToken       = "ZNS_TOKEN"
	ConfigKeyZNSRefresh     = "ZNS_REFRESH"
	ConfigKeyZNSTemplateID  = "ZNS_TEMPLATE_ID"
	ConfigKeyZNSOTPTemplate = "ZNS_OTP_TEMPLATE_ID"
)
