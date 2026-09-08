package qh_domain

import (
	"time"
)

type QHLayerLegal struct {
	ID        uint64     `gorm:"primaryKey;autoIncrement" json:"id"`
	LayerID   uint64     `gorm:"not null;index:idx_layer_legal" json:"layerId"`
	Name      string     `gorm:"type:varchar(255);not null" json:"name"`
	FileURL   string     `gorm:"type:text;not null" json:"fileUrl"`
	FileType  string     `gorm:"type:varchar(50)" json:"fileType"`
	CreatedAt time.Time  `gorm:"autoCreateTime" json:"createdAt"`
	DeletedAt *time.Time `gorm:"index" json:"deletedAt,omitempty"`
}

func (QHLayerLegal) TableName() string {
	return "qh_layer_legals"
}
