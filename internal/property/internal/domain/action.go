package domain

import _models "common/models"

type ActionType int

const (
	ActionViewPost  ActionType = 1
	ActionClickPost ActionType = 2
)

type Action struct {
	_models.BaseEntity
	ActionType uint   `json:"action"`
	TargetID   uint64 `json:"targetId"`
}

func (Action) TableName() string {
	return "action"
}
