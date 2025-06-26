package dto

import "organization/internal/enums"

type NotificationDTO struct {
	Type       enums.NotificationType
	UserId     *uint64
	AttachData []string
	TargetID   uint64
	IsMerge    bool
}
