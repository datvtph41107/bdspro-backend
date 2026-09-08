package qh_domain

import (
	_models "common/domain/entity"
	_enum "common/domain/enum"
)

type QHLandUse struct {
	_models.BaseEntity

	// LandUseGroupID uint64 `gorm:"column:land_use_group_id;not null;uniqueIndex:idx_layer_land_use_group" json:"landUseGroupId"`

	Note string `gorm:"type:text" json:"note,omitempty"`
	Name string `gorm:"type:varchar(100)" json:"name"`
	Code string `gorm:"type:varchar(10);" json:"code"`

	DisplayOrder int    `gorm:"default:0" json:"displayOrder"`
	Description  string `gorm:"type:text" json:"description"`

	IsVisible      bool   `gorm:"default:true" json:"isVisible"`
	Priority       int    `gorm:"default:50" json:"priority"`
	BuildCondition string `gorm:"type:text" json:"buildCondition"`

	Color    string `gorm:"type:varchar(10)" json:"color,omitempty"`
	CanBuild bool   `gorm:"default:false" json:"canBuild"`

	IsActive bool `gorm:"default:true" json:"isActive"`

	WarnLevel _enum.EWarningLevel `gorm:"default:0" json:"warnLevel"`

	Layers []*QHLayer `gorm:"many2many:qh_layer_land_use;joinForeignKey:land_use_id;joinReferences:layer_id" json:"layers,omitempty"`
	// LandUseGroup *LandUseGroup `gorm:"foreignKey:LandUseGroupID" json:"landUseGroup,omitempty"`
}

func (QHLandUse) TableName() string {
	return "qh_land_use"
}

func (e *QHLandUse) GetDisplayColor() string {
	if e.Color != "" {
		return e.Color
	}
	// if e.LandUseGroup != nil && e.LandUseGroup.Color != "" {
	// 	return e.LandUseGroup.Color
	// }
	return "#CCCCCC"
}
