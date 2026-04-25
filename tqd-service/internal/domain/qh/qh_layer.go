package qh_domain

import (
	"encoding/json"
	"time"
	"tqd/internal/domain/jsonb"
	"tqd/internal/enums"
)

type QHLayer struct {
	ID          uint64 `gorm:"primaryKey;autoIncrement" json:"id"`
	Name        string `gorm:"type:varchar(255);not null;uniqueIndex:idx_name" json:"name"`
	DisplayName string `gorm:"type:varchar(255);not null" json:"displayName"`
	Description string `gorm:"type:text" json:"description,omitempty"`

	Type   enums.LayerType   `gorm:"type:int;not null;index" json:"type"`
	Status enums.LayerStatus `gorm:"type:int;default:10;index" json:"status"`

	DisplayOrder int   `gorm:"column:display_order;default:0;index" json:"displayOrder"`
	Visible      int32 `gorm:"default:0;index" json:"visible"` // 0: ẩn, 1: hiển thị
	MinZoom      int   `gorm:"default:0" json:"minZoom"`
	MaxZoom      int   `gorm:"default:22" json:"maxZoom"`

	LayerURL     string `gorm:"type:text;not null" json:"layerUrl"`
	Avatar       string `gorm:"type:text" json:"avatar,omitempty"`
	ImageURL     string `gorm:"type:text" json:"imageUrl,omitempty"`
	ThumbnailURL string `gorm:"type:text" json:"thumbnailUrl,omitempty"`

	// Styling
	StyleConfig jsonb.JSONB `gorm:"type:jsonb;default:'{}'" json:"styleConfig,omitempty"`

	LegalDoc         string `gorm:"type:text" json:"legalDoc,omitempty"`
	LegalDocURL      string `gorm:"type:text" json:"legalDocUrl,omitempty"`
	DocumentNumber   string `gorm:"type:varchar(100)" json:"documentNumber,omitempty"`
	IssuingAuthority string `gorm:"type:varchar(255)" json:"issuingAuthority,omitempty"`

	EffectiveDate *time.Time `gorm:"type:date;index" json:"effectiveDate,omitempty"`
	ExpiryDate    *time.Time `gorm:"type:date;index" json:"expiryDate,omitempty"`

	CreatedBy uint64     `gorm:"default:0" json:"createdBy"`
	UpdatedBy uint64     `gorm:"default:0" json:"updatedBy"`
	CreatedAt time.Time  `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt time.Time  `gorm:"autoUpdateTime" json:"updatedAt"`
	DeletedAt *time.Time `gorm:"index" json:"deletedAt,omitempty"`
}

func (QHLayer) TableName() string {
	return "qh_layers"
}

type StyleConfig struct {
	FillColor     string  `json:"fillColor,omitempty"`
	FillOpacity   float64 `json:"fillOpacity,omitempty"`
	StrokeColor   string  `json:"strokeColor,omitempty"`
	StrokeWidth   int     `json:"strokeWidth,omitempty"`
	LineColor     string  `json:"lineColor,omitempty"`
	LineWidth     int     `json:"lineWidth,omitempty"`
	LineOpacity   float64 `json:"lineOpacity,omitempty"`
	CircleColor   string  `json:"circleColor,omitempty"`
	CircleRadius  int     `json:"circleRadius,omitempty"`
	CircleOpacity float64 `json:"circleOpacity,omitempty"`
}

func (s StyleConfig) ToJSONB() jsonb.JSONB {
	data, _ := json.Marshal(s)
	var result jsonb.JSONB
	json.Unmarshal(data, &result)
	return result
}

// Business logic methods
func (l *QHLayer) IsActive() bool {
	if l.DeletedAt != nil {
		return false
	}
	if l.Status != enums.LayerStatusActive {
		return false
	}
	if l.Visible != 1 {
		return false
	}
	now := time.Now()
	if l.EffectiveDate != nil && l.EffectiveDate.After(now) {
		return false
	}
	if l.ExpiryDate != nil && l.ExpiryDate.Before(now) {
		return false
	}
	return true
}

func (l *QHLayer) IsClientVisible() bool {
	return l.IsActive() && l.Visible == 1
}

func (l *QHLayer) GetStyleConfig() StyleConfig {
	if l.StyleConfig == nil {
		// Default style theo type
		defaults := map[enums.LayerType]StyleConfig{
			enums.LayerTypeLandUse: {
				FillColor:   "#CCCCCC",
				FillOpacity: 0.6,
				StrokeColor: "#000000",
				StrokeWidth: 1,
			},
			enums.LayerTypeInfrastructure: {
				LineColor:   "#0000FF",
				LineWidth:   2,
				LineOpacity: 0.8,
			},
		}
		if def, ok := defaults[l.Type]; ok {
			return def
		}
		return StyleConfig{}
	}

	var config StyleConfig
	data, _ := json.Marshal(l.StyleConfig)
	json.Unmarshal(data, &config)
	return config
}
