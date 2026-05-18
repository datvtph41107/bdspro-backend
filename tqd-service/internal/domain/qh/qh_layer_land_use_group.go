package qh_domain

import _models "common/domain/entity"

type QHLayerLandUseGroup struct {
	_models.BaseEntity
	LayerID        uint64 `gorm:"column:layer_id;not null;uniqueIndex:idx_layer_land_use_group" json:"layerId"`
	LandUseGroupID uint64 `gorm:"column:land_use_group_id;not null;uniqueIndex:idx_layer_land_use_group" json:"landUseGroupId"`

	Note string `gorm:"type:text" json:"note,omitempty"`

	DisplayOrder int `gorm:"default:0" json:"displayOrder"`

	IsVisible bool `gorm:"default:true" json:"isVisible"`

	Color string `gorm:"type:varchar(10)" json:"color,omitempty"`

	Layer        *QHLayer      `gorm:"foreignKey:LayerID" json:"layer,omitempty"`
	LandUseGroup *LandUseGroup `gorm:"foreignKey:LandUseGroupID" json:"landUseGroup,omitempty"`
}

func (QHLayerLandUseGroup) TableName() string {
	return "qh_layer_land_use_groups"
}

func (e *QHLayerLandUseGroup) GetDisplayColor() string {
	if e.Color != "" {
		return e.Color
	}
	if e.LandUseGroup != nil && e.LandUseGroup.Color != "" {
		return e.LandUseGroup.Color
	}
	return "#CCCCCC"
}
