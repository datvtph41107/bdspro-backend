package qh_domain

import "time"

type QHLandType struct {
	ID          uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	Code        string    `gorm:"type:varchar(20);uniqueIndex;not null" json:"code"` // ODT, TMD, LNP, DCC...
	Name        string    `gorm:"type:varchar(255);not null" json:"name"`            // Đất ở đô thị
	GroupName   string    `gorm:"type:varchar(100)" json:"groupName"`                // Nhóm đất
	Category    string    `gorm:"type:varchar(50)" json:"category"`                  // Residential, Commercial...
	Color       string    `gorm:"type:varchar(7);default:'#CCCCCC'" json:"color"`    // Mã màu hex
	FillOpacity float64   `gorm:"default:0.6" json:"fillOpacity"`
	StrokeColor string    `gorm:"type:varchar(7);default:'#000000'" json:"strokeColor"`
	StrokeWidth int       `gorm:"default:1" json:"strokeWidth"`
	Priority    int       `gorm:"default:0" json:"priority"`
	IsActive    bool      `gorm:"default:true" json:"isActive"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updatedAt"`
}

func (QHLandType) TableName() string {
	return "qh_land_type"
}
