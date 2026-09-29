package qh_domain

import (
	"time"
	"tqd/internal/enums"
)

type QHLabel struct {
	ID          uint64 `gorm:"primaryKey" json:"id"`
	LayerID     uint64 `gorm:"not null;index" json:"layerId"`
	Name        string `gorm:"type:varchar(100);not null" json:"name"`
	DisplayName string `gorm:"type:varchar(255)" json:"displayName"`
	Description string `gorm:"type:text" json:"description"`

	// Style
	Color           string  `gorm:"type:varchar(7);default:'#CCCCCC'" json:"color"`
	FillOpacity     float64 `gorm:"default:0.6" json:"fillOpacity"`
	StrokeColor     string  `gorm:"type:varchar(7);default:'#000000'" json:"strokeColor"`
	StrokeWidth     int     `gorm:"default:1" json:"strokeWidth"`
	StrokeDashArray string  `gorm:"type:varchar(50)" json:"strokeDashArray,omitempty"`

	// Display
	DisplayOrder int  `gorm:"default:0;index" json:"displayOrder"`
	IsVisible    bool `gorm:"default:true" json:"isVisible"`
	MinZoom      int  `gorm:"default:0" json:"minZoom"`
	MaxZoom      int  `gorm:"default:22" json:"maxZoom"`

	RegionCount int64             `gorm:"default:0" json:"regionCount"`
	Status      enums.LabelStatus `gorm:"default:10" json:"status"`

	LandCodeID *uint64 `gorm:"index" json:"standardId,omitempty"`

	// LandCode   *LandUseCode `gorm:"-" json:"standard,omitempty"`

	SourceType  string `gorm:"type:varchar(50);default:'unknown'" json:"sourceType"`
	SourceValue string `gorm:"type:varchar(255)" json:"sourceValue"`

	StandardAt      *time.Time `gorm:"column:standard_at" json:"standardAt,omitempty"`
	ClearStandardAt bool       `gorm:"-" json:"-"`

	LandUseID *uint64 `gorm:"-:migration;->;column:land_use_id" json:"landUseId,omitempty"`
	LegendID  *uint64 `gorm:"-:migration;->;column:legend_id" json:"legendId,omitempty"`

	CreatedAt time.Time  `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt time.Time  `gorm:"autoUpdateTime" json:"updatedAt"`
	DeletedAt *time.Time `gorm:"index" json:"deletedAt,omitempty"`
}

func (QHLabel) TableName() string {
	return "qh_labels"
}

// Đã tham chiếu làm chuẩn hóa với hệ thống chưa
func (l *QHLabel) HasStandard() bool {
	return l.LandCodeID != nil && *l.LandCodeID > 0
}

func (l *QHLabel) GetColorHex() string {
	if l.Color == "" {
		return "#CCCCCC"
	}
	return l.Color
}

func (l *QHLabel) IsActive() bool {
	return l.Status == enums.LabelStatusActive && l.DeletedAt == nil
}

func (l *QHLabel) IsVisibleAtZoom(zoom int) bool {
	if !l.IsVisible {
		return false
	}
	if zoom < l.MinZoom || zoom > l.MaxZoom {
		return false
	}
	return true
}
