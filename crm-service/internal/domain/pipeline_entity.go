package domain

import (
	_models "common/models"
	"crm/internal/enums"
)

type PipelineEntity struct {
	_models.BaseEntity
	ID             uint64           `gorm:"primaryKey" json:"id"`
	PipelineName   string           `json:"pipelineName" binding:"required"`
	Color          uint32           `gorm:"default:10" json:"color"`
	Active         bool             `json:"active"`
	OwnerID        uint64           `json:"ownerId" binding:"required"`
	OwnerType      enums.EOwnerType `json:"ownerType" binding:"required"`
	StageDefaultID *uint64          `json:"stageDefaultId"`
	IsDefault      bool             `gorm:"not null;default:false" json:"isDefault"`
	IsCustom       bool             `gorm:"not null;default:false" json:"isCustom"`
	DefaultStage   *StageEntity     `gorm:"foreignKey:StageDefaultID" json:"defaultStage"`
	Stages         []StageEntity    `gorm:"foreignKey:PipelineID" json:"stages"`
}

func (PipelineEntity) TableName() string {
	return "pipeline"
}
