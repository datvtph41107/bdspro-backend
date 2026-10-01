package provider

import (
	_dto "common/domain/dto"
	_enum "common/domain/enum"
	"context"
)

type NotificationProvider interface {
	CreateHistory(ctx context.Context, history *_dto.HistoryDTO) error

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

	RemoveNotification(
		ctx context.Context,
		ownerOf _enum.EOwnerOf,
		ownerId uint64,
		targetId uint64,
		notificationType _enum.ENotificationType,
	) error
}

// type INotificationClient interface {
// 	SendReminder(c context.Context, customerID uint64) error
// 	SendNoti(c context.Context, noti *delivery_dto.NotificationDTO) error
// 	CreateHistory(c context.Context, entity *delivery_dto.NotificationDTO) error
// 	ExistedFromDate(c context.Context, customerID uint64, fromDate time.Time) (bool, error)
// }

// type INotificationClient interface {
// 	_provider.NotificationProvider
// }
