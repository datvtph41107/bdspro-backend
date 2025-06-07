package domain

import (
	_models "common/models"
	"crm/internal/enums"
)

type HistoryEntity struct {
	_models.BaseEntity
	TargetId   uint64               `json:"targetId"`
	TargetType enums.ETargetHistory `json:"targetType"`
	ActionType enums.EHistory       `json:"actionType"`
	Note       string               `json:"note,omitempty"`
	PreStage   string               `json:"preStage,omitempty"`
	AfterStage string               `json:"afterStage,omitempty"`
	OwnerID    *uint64              `json:"ownerId"`
	OwnerType  enums.EOwnerType     `json:"ownerType"`
}

func (HistoryEntity) TableName() string {
	return "crm_histories"
}
