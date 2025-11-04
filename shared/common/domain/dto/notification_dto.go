package _dto

import (
	_enum "common/domain/enum"
)

type NotificationDTO struct {
	Avatar           string                  `json:"avatar"`
	Title            string                  `json:"title"`
	Message          []string                `json:"message"`
	OwnerID          uint64                  `json:"ownerId"`
	OwnerOf          _enum.EOwnerOf          `json:"ownerOf"`
	NotificationType _enum.ENotificationType `json:"notificationType"`
	TargetID         uint64                  `json:"targetId"`
	IsMerge          bool                    `json:"isMerge"`
	AttachData       []string                `json:"attachData"`
	SendToDevice     bool                    `json:"sendToDevice"`
}

type NotificationBatchDTO struct {
	Datas []NotificationDTO `json:"datas"`
}
