package models

import (
	_models "common/domain/entity"
)

// PurposeUseEntity đại diện cho mục đích sử dụng được gắn với người dùng.
type PurposeUseEntity struct {
	_models.BaseEntity
	Name          string `gorm:"column:name;type:varchar(255);not null" json:"name"`
	Description   string `gorm:"column:description;type:text" json:"description,omitempty"`
	IsActive      bool   `gorm:"column:is_active;not null;default:true" json:"isActive"`
	NumberProfile uint64 `gorm:"column:number_profile;not null;default:0" json:"numberProfile"`
}

// TableName trả về tên bảng tương ứng trong cơ sở dữ liệu.
func (PurposeUseEntity) TableName() string {
	return "purpose_use"
}
