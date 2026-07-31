package domain

import (
	_models "common/models"
)

type Module struct {
	_models.BaseEntity
	Name        string `gorm:"size:100;not null" json:"name"`
	Key         string `gorm:"size:50;uniqueIndex;not null" json:"key"`
	Code        string `gorm:"size:20;uniqueIndex;not null" json:"code"`
	Description string `gorm:"size:500" json:"description"`
	IsActive    bool   `gorm:"default:true" json:"isActive"`
	Priority    uint32 `gorm:"type:int;default:10" json:"priority"`
}

func (Module) TableName() string {
	return "modules"
}
