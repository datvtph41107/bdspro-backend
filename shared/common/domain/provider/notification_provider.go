package _provider

import (
	_enum "common/domain/enum"
	"context"
)

type NotificationProvider interface {
	CreateNotification(
		ctx context.Context,
		avatar string,
		title string,
		message []string,
		notificationType _enum.ENotificationType,
		targetId *uint64,
		ownerId uint64,
		ownerOf _enum.EOwnerOf,
		attachData []string,
	) error
}
