package base_dto

import base_enum "base/enum"

type HistoryDTO struct {
	Title      string                   `json:"title"`
	Note       []string                 `json:"note"`
	TargetId   uint64                   `json:"targetId"`
	TargetType base_enum.ETargetHistory `json:"targetType"`
	ActionType base_enum.EHistory       `json:"actionType"`
	OwnerID    *uint64                  `json:"ownerId"`
	OwnerOf    base_enum.EOwnerOf       `json:"ownerOf"`
}
