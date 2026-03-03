package domain

import (
	_models "common/domain/entity"
)

type Region struct {
	_models.BaseEntity
	Name      string    `gorm:"column:name;size:255" json:"name"`
	Code      uint32    `gorm:"column:code;size:50" json:"code"`
	Active    bool      `gorm:"column:active;default:true;index" json:"active"`
}

func (Region) TableName() string {
	return "regions"
}
