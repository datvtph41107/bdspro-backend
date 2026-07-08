package domain

import (
	_models "common/models"
	"crm/internal/enums"
)

type StageEntity struct {
	_models.BaseEntity
	StageName   string      `json:"stageName" binding:"required"`
	OrderNumber int         `json:"orderNumber" binding:"required"`
	PipelineID  uint64      `json:"pipelineId" binding:"required"`
	RuleID      *uint64     `json:"ruleId"`
	Step        enums.EStep `json:"step"`
	Active      bool        `json:"active"`
	Color       uint32      `gorm:"default:10" json:"color"`
	ColorRGB    string      `gorm:"default:10" json:"colorRGB"`

	Pipeline *PipelineEntity `gorm:"references:ID" json:"pipeline"`
}

func (StageEntity) TableName() string {
	return "stage"
}