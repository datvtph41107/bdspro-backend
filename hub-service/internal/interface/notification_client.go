package iprovider

import (
	_enum "common/domain/enum"
	"context"
	sharepb "pb/types/shared"
)

// @bind: hub/infra/client.NotificationClient
type INotificationClient interface {
	CreateToOwner(
		ctx context.Context,
		avatar string,
		title string,
		message []string,
		notificationType _enum.ENotificationType,
		targetId *uint64,
		ownerId uint64,
		ownerOf _enum.EOwnerOf,
		attachData []string,
		isMerge bool,
	) error
	SendBatch(
		ctx context.Context,
		requests []*sharepb.NotificationRequest,
	) error
}
