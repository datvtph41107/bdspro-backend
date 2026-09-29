package domain

import (
	_models "common/models"
)

// Bảng phường/xã
type Region struct {
	_models.BaseEntity
	Name     string  `gorm:"size:100;not null" json:"name"`
	Code     uint    `json:"code"`
	CodeName string  `json:"codeName"`
	Level    uint    `json:"level"`
	ParentID *uint64 `json:"parentId"`
	Unit     string  `json:"unit"`
	// Parent   *Region `gorm:"foreignKey:ParentID"`
}

func (Region) TableName() string {
	return "region"
}
