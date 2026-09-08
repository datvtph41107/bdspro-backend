package domain

import (
	_models "common/domain/entity"
)

// District entity represents a district in Vietnam
type District struct {
	_models.BaseEntity
	Name       string `gorm:"column:name;type:varchar(255);not null" json:"name"`
	ProvinceID string `gorm:"column:province_id;type:varchar(10);not null;index" json:"provinceId"`
	Type       int    `gorm:"column:type;type:int;not null" json:"type"`
	TypeText   string `gorm:"column:type_text;type:varchar(50);not null" json:"typeText"`
}

// TableName specifies the table name for District entity
func (District) TableName() string {
	return "districts"
}

