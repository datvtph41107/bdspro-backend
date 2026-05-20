// internal/domain/ward.go
package domain

import (
	"time"
)

type Ward struct {
	ID           string    `json:"id" gorm:"column:id;primaryKey;size:50"`
	FullName     string    `json:"fullName" gorm:"column:full_name;size:255;not null;index:idx_wards_name"`
	ShortName    string    `json:"shortName" gorm:"column:short_name;size:100;index:idx_wards_short_name"`
	Lat          float64   `json:"lat" gorm:"column:lat;type:decimal(10,8)"`
	Lng          float64   `json:"lng" gorm:"column:lng;type:decimal(11,8)"`
	Code         string    `json:"code" gorm:"column:code;size:20;uniqueIndex;index:idx_wards_code"`
	ProvinceID   string    `json:"provinceId" gorm:"column:province_id;size:50;not null;index:idx_wards_province_id"`
	ProvinceCode string    `json:"provinceCode" gorm:"column:province_code;size:20"`
	CreatedAt    time.Time `json:"createdAt" gorm:"column:created_at;default:CURRENT_TIMESTAMP;index:idx_wards_created_at"`
	UpdatedAt    time.Time `json:"updatedAt" gorm:"column:updated_at;default:CURRENT_TIMESTAMP"`

	Province *Province `json:"province,omitempty" gorm:"foreignKey:ProvinceID;references:ID"`
}

func (Ward) TableName() string {
	return "ward_v2"
}
