package qh_domain

import "time"

type LandUseDict struct {
	// Primary identifier
	Code string `gorm:"primaryKey;type:varchar(20)" json:"code"`

	Name      string `gorm:"type:varchar(255);not null" json:"name"`
	GroupName string `gorm:"type:varchar(100)" json:"groupName"`
	Category  string `gorm:"type:varchar(50)" json:"category"`

	Color       string  `gorm:"type:varchar(7);default:'#CCCCCC'" json:"color"`
	FillOpacity float64 `gorm:"default:0.6" json:"fillOpacity"`
	StrokeColor string  `gorm:"type:varchar(7);default:'#000000'" json:"strokeColor"`
	StrokeWidth int     `gorm:"default:1" json:"strokeWidth"`

	DisplayPriority int  `gorm:"default:0" json:"displayPriority"`
	IsActive        bool `gorm:"default:true" json:"isActive"`

	CreatedAt time.Time `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updatedAt"`
}

func (LandUseDict) TableName() string {
	return "land_use_dict"
}
