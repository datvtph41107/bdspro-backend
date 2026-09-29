package _dto

import _enum "common/domain/enum"

type HistoryDTO struct {
	Title      string               `json:"title"`
	Note       []string             `json:"note"`
	TargetId   uint64               `json:"targetId"`
	TargetType _enum.ETargetHistory `json:"targetType"`
	ActionType _enum.EHistory       `json:"actionType"`
	OwnerID    *uint64              `json:"ownerId"`
	OwnerOf    _enum.EOwnerOf       `json:"ownerOf"`
}
