package domain

import _models "common/models"

type Block struct {
	_models.BaseEntity
	Name      string     `gorm:"size:255;not null" json:"name"`
	ProjectID uint64     `gorm:"type:int8;not null;index" json:"projectId"`
	Project   *Project   `gorm:"foreignKey:ProjectID;references:ID" json:"project,omitempty"`
	Buildings []Building `gorm:"foreignKey:BlockID" json:"buildings,omitempty"`
}

func (Block) TableName() string {
	return "blocks"
}
