package domain

import (
	_models "common/models"
	"crm/internal/enums"
)

type HistoryEntity struct {
	_models.BaseEntity
	TargetId   uint64         `json:"targetId"`
	TargetType int8           `json:"targetType"`
	ActionType enums.EHistory `json:"actionType"`
	Note       string         `json:"note,omitempty"`
	PreStage   string         `json:"preStage,omitempty"`
	AfterStage string         `json:"afterStage,omitempty"`
}

func (HistoryEntity) TableName() string {
	return "crm_histories"
}
