package models

import (
	_models "common/models"
)

// MainAreaEntity đại diện cho nhãn khu vực hoạt động được gắn cho người dùng.
type MainAreaEntity struct {
	_models.BaseEntity
	Name          string `gorm:"column:name;type:varchar(255);not null" json:"name"`
	Description   string `gorm:"column:description;type:varchar(512)" json:"description"`
	NumberProfile uint64 `gorm:"column:number_profile;not null;default:0" json:"numberProfile"`
}

func (MainAreaEntity) TableName() string {
	return "main_area"
}
