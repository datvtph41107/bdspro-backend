package shared_dto

import shared_enum "pb/enums"

type HistoryDTO struct {
	Title      string                     `json:"title"`
	Note       []string                   `json:"note"`
	TargetId   uint64                     `json:"targetId"`
	TargetType shared_enum.ETargetHistory `json:"targetType"`
	ActionType shared_enum.EHistory       `json:"actionType"`
	OwnerID    *uint64                    `json:"ownerId"`
	OwnerOf    shared_enum.EOwnerType     `json:"ownerType"`
}
