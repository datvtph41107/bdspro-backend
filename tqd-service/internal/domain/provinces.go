// internal/domain/province.go
package domain

import (
	"time"
)

type Province struct {
	ID        string    `json:"id" gorm:"column:id;primaryKey;size:50"`
	FullName  string    `json:"fullName" gorm:"column:full_name;size:255;not null;index:idx_provinces_name"`
	ShortName string    `json:"shortName" gorm:"column:short_name;size:100;index:idx_provinces_short_name"`
	Lat       float64   `json:"lat" gorm:"column:lat;type:decimal(10,8)"`
	Lng       float64   `json:"lng" gorm:"column:lng;type:decimal(11,8)"`
	Code      string    `json:"code" gorm:"column:code;size:20;uniqueIndex;index:idx_provinces_code"`
	WardCount int32     `json:"wardCount" gorm:"column:ward_count;default:0"`
	CreatedAt time.Time `json:"createdAt" gorm:"column:created_at;default:CURRENT_TIMESTAMP;index:idx_provinces_created_at"`
	UpdatedAt time.Time `json:"updatedAt" gorm:"column:updated_at;default:CURRENT_TIMESTAMP"`
}

func (Province) TableName() string {
	return "province_v2"
}
