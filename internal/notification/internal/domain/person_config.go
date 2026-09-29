package domain

import (
	_enum "common/domain/enum"
	_models "common/models"
)

// PersonConfigEntity đại diện cho cấu hình cá nhân được bật/tắt
type PersonConfigEntity struct {
	_models.BaseEntity
	ID        uint64 `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID    uint64 `gorm:"not null;index;uniqueIndex:idx_person_config_user_key" json:"userId"`
	Key       string `gorm:"size:255;not null;uniqueIndex:idx_person_config_user_key" json:"key"`
	Checked   bool   `gorm:"default:false" json:"checked"`
	IsDefault bool   `gorm:"default:false" json:"isDefault"`

	Channel _enum.EChannelNotification `gorm:"default:10;uniqueIndex:idx_person_config_user_key" json:"channel"`
}

func (PersonConfigEntity) TableName() string {
	return "person_configs"
}
