package qh_domain

import (
	"time"
	"tqd/internal/enums"
)

type QHLayer struct {
	ID            uint              `gorm:"primaryKey;autoIncrement" json:"id"`
	Name          string            `gorm:"type:varchar(255);not null" json:"name"`
	DisplayName   string            `gorm:"type:varchar(255);not null" json:"display_name"`
	Type          enums.LayerType   `gorm:"type:int;not null;index" json:"type"`
	Status        enums.LayerStatus `gorm:"type:int;default:10;index" json:"status"`
	DisplayOrder  int               `gorm:"column:display_order;default:0;index" json:"display_order"`
	MinZoom       int               `gorm:"default:0" json:"min_zoom"`
	MaxZoom       int               `gorm:"default:22" json:"max_zoom"`
	Visible       int32             `gorm:"default:0" json:"visible"`
	LegalDoc      string            `gorm:"type:text" json:"legal_doc"`
	EffectiveDate *time.Time        `gorm:"index" json:"effective_date,omitempty"`
	CreatedAt     time.Time         `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt     time.Time         `gorm:"autoUpdateTime" json:"updated_at"`
	Avatar        string            `gorm:"type:text" json:"avatar"`
	LayerUrl      string            `gorm:"type:text" json:"layerUrl"`
}

func (QHLayer) TableName() string {
	return "qh_layers"
}
