package qh_domain

import "time"

type QHLabel struct {
	ID          uint64     `gorm:"primaryKey;autoIncrement" json:"id"`
	LayerID     uint64     `gorm:"not null;index;comment:Thuộc layer" json:"layerId"`
	Name        string     `gorm:"type:varchar(255);not null" json:"name"`
	DisplayName string     `gorm:"type:varchar(255)" json:"displayName"`
	Color       string     `gorm:"type:varchar(7);default:'#000000'" json:"color"`
	FontSize    int        `gorm:"default:12" json:"fontSize"`
	FontWeight  string     `gorm:"type:varchar(20);default:'normal'" json:"fontWeight"`
	MinZoom     int        `gorm:"default:10" json:"minZoom"`
	MaxZoom     int        `gorm:"default:22" json:"maxZoom"`
	IsVisible   bool       `gorm:"default:true" json:"isVisible"`
	Status      int        `gorm:"default:10" json:"status"`
	CreatedAt   time.Time  `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt   time.Time  `gorm:"autoUpdateTime" json:"updatedAt"`
	DeletedAt   *time.Time `gorm:"index" json:"deletedAt,omitempty"`
}

func (QHLabel) TableName() string {
	return "qh_labels"
}
