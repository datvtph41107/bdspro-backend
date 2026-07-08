package provider

import (
	"context"
	notificationpb "pb/types/notification"
)

type NotificationProvider interface {
	CreateHistory(ctx context.Context, history *notificationpb.HistoryDTO) (*notificationpb.HistoryDTO, error)
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