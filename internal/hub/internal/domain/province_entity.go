package domain

import (
	_models "common/domain/entity"
)

// Province entity represents a province/city in Vietnam
type Province struct {
	_models.BaseEntity
	Name     string `gorm:"column:name;type:varchar(255);not null" json:"name"`
	Type     int    `gorm:"column:type;type:int;not null" json:"type"`
	TypeText string `gorm:"column:type_text;type:varchar(50);not null" json:"typeText"`
	Slug     string `gorm:"column:slug;type:varchar(100);not null;uniqueIndex" json:"slug"`
}

// TableName specifies the table name for Province entity
func (Province) TableName() string {
	return "provinces"
}
