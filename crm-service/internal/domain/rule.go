package domain

import (
	base_enum "base/enum"
	_models "common/models"
	"crm/internal/enums"
)

type RuleEntity struct {
	_models.BaseEntity
	OwnerID        uint64             `json:"ownerId"`
	OwnerType      base_enum.EOwnerOf `json:"ownerType"`
	Active         bool               `json:"active"`
	RuleName       string             `json:"ruleName" binding:"required"`
	Note           string             `json:"note"`
	Condition      enums.ECondition   `json:"condition" binding:"required"`
	ConditionValue string             `json:"conditionValue" binding:"required"`
	Trigger        enums.ERuleTrigger `json:"trigger" binding:"required"`
	TriggerValue   string             `json:"triggerValue" binding:"required"`
	// CycleSeconds int64            `json:"cycleSeconds"`
	// Deadline *time.Time `json:"deadline"`
}

func (e *RuleEntity) TableName() string {
	return "rules"
}