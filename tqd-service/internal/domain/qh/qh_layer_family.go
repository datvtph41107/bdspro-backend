package qh_domain

import _models "common/domain/entity"

type QHLayerFamily struct {
	_models.BaseEntity
	Name   string    `gorm:"type:varchar(100);" json:"name"`
	Layers []QHLayer `gorm:"foreignKey:FamilyID;references:ID" json:"layers"`
}

func (QHLayerFamily) TableName() string {
	return "qh_layer_families"
}
