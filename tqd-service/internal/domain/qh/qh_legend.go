package qh_domain

import _models "common/domain/entity"

const (
	LegendTypeColorFill  = "color_fill"
	LegendTypeLineStroke = "line_stroke"
	LegendTypeIcon       = "icon"
	LegendTypePattern    = "pattern"
	LegendTypeImage      = "image"
)

var validLegendTypes = map[string]bool{
	LegendTypeColorFill:  true,
	LegendTypeLineStroke: true,
	LegendTypeIcon:       true,
	LegendTypePattern:    true,
	LegendTypeImage:      true,
}

func IsValidLegendType(t string) bool {
	return validLegendTypes[t]
}

const (
	GeometryTypePolygon = "polygon"
	GeometryTypeLine    = "line"
	GeometryTypePoint   = "point"
	GeometryTypeMulti   = "multi"
)

var validGeometryTypes = map[string]bool{
	GeometryTypePolygon: true,
	GeometryTypeLine:    true,
	GeometryTypePoint:   true,
	GeometryTypeMulti:   true,
}

func IsValidGeometryType(t string) bool {
	return validGeometryTypes[t]
}

type QHLayerLegend struct {
	_models.BaseEntity

	LayerID     uint64     `gorm:"column:layer_id;not null;" json:"layerId"`
	LandUseID   *uint64    `gorm:"column:land_use_id;index" json:"landUseId,omitempty"`
	LandUse     *QHLandUse `gorm:"foreignKey:LandUseID" json:"landUse,omitempty"`
	LabelID     uint64     `gorm:"column:label_id;" json:"labelId"`
	Description string     `gorm:"type:text" json:"description,omitempty"`
	Note        string     `gorm:"type:text" json:"note,omitempty"`
	Color       string     `gorm:"type:varchar(10)" json:"color,omitempty"`

	DisplayOrder int  `gorm:"default:0" json:"displayOrder"`
	IsVisible    bool `gorm:"default:true" json:"isVisible"`

	LegendType   string `gorm:"type:varchar(20);default:'color_fill'" json:"legendType"`
	GeometryType string `gorm:"type:varchar(20);default:'polygon'" json:"geometryType"`

	StrokeDashArray string `gorm:"type:varchar(50)" json:"strokeDashArray,omitempty"`
	IconURL         string `gorm:"type:text" json:"iconUrl,omitempty"`
	ImageURL        string `gorm:"type:text" json:"imageUrl,omitempty"`

	Label *QHLabel `gorm:"foreignKey:LabelID" json:"label,omitempty"`

	// LandUseGroupID *uint64       `gorm:"column:land_use_group_id;index" json:"landUseGroupId,omitempty"`
	// LandUseGroup   *LandUseGroup `gorm:"foreignKey:LandUseGroupID" json:"landUseGroup,omitempty"`
}

func (QHLayerLegend) TableName() string {
	return "qh_legends"
}

func (e *QHLayerLegend) GetLegendType() string {
	if e.LegendType == "" {
		return LegendTypeColorFill
	}
	return e.LegendType
}

func (e *QHLayerLegend) GetGeometryType() string {
	if e.GeometryType == "" {
		return GeometryTypePolygon
	}
	return e.GeometryType
}

func (e *QHLayerLegend) GetDisplayColor() string {
	if e.Color != "" {
		return e.Color
	}
	if e.LandUse != nil && e.LandUse.Color != "" {
		return e.LandUse.Color
	}
	return "#CCCCCC"
}

func (e *QHLayerLegend) GetStrokeColor() string {
	// if e.Label != nil && e.Label.StrokeColor != "" {
	// 	return e.Label.StrokeColor
	// }
	return "#000000"
}

func (e *QHLayerLegend) GetEffectiveStrokeDash() string {
	if e.StrokeDashArray != "" {
		return e.StrokeDashArray
	}
	// if e.Label != nil {
	// 	return e.Label.StrokeDashArray
	// }
	return ""
}
