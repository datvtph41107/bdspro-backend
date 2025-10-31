package domain

import (
	_models "common/domain/entity"
	"hub/internal/enums"
)

// SystemConfigEntity chứa cấu hình hệ thống
type SystemConfigEntity struct {
	_models.BaseEntity
	Name        string                   `gorm:"size:255;not null" json:"name"`
	Key         string                   `gorm:"size:100;unique;not null" json:"key"`
	Value       string                   `gorm:"type:text" json:"value"`
	GroupConfig enums.ESystemConfigGroup `gorm:"column:group_key;default:10" json:"group"`
}

// TableName đặt tên bảng trong DB
func (SystemConfigEntity) TableName() string {
	return "system_config"
}
