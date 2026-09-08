package domain

import (
	_models "common/domain/entity"
)

// Ward entity represents a ward/commune in Vietnam
type Ward struct {
	_models.BaseEntity
	Name       string `gorm:"column:name;type:varchar(255);not null" json:"name"`
	DistrictID string `gorm:"column:district_id;type:varchar(10);not null;index" json:"districtId"`
	Type       int    `gorm:"column:type;type:int;not null" json:"type"`
	TypeText   string `gorm:"column:type_text;type:varchar(50);not null" json:"typeText"`
}

// TableName specifies the table name for Ward entity
func (Ward) TableName() string {
	return "wards"
}
