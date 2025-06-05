package domain

import (
	_models "common/models"
	"crm/internal/enums"
)

type PipelineEntity struct {
	_models.BaseEntity
	PipelineName string           `json:"pipelineName" binding:"required"`
	Color        string           `gorm:"type:varchar(15)" json:"color"`
	Active       bool             `json:"active"`
	OwnerID      uint64           `json:"ownerId" binding:"required"`
	OwnerType    enums.EOwnerType `json:"ownerType" binding:"required"`
}

func (PipelineEntity) TableName() string {
	return "pipeline"
}
