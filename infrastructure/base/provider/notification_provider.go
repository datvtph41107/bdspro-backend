package base_provider

import (
	base_dto "base/dto"
	"context"
	"time"
)

type NotificationProvider interface {
	SendNotification(ctx context.Context, noti *base_dto.NotificationDTO) error
	SendBatchNotification(ctx context.Context, noti *base_dto.NotificationBatchDTO) error
	CreateCRMHistory(ctx context.Context, history *base_dto.HistoryDTO) error
	ExistedFromDate(ctx context.Context, customerID uint64, fromDate time.Time) (bool, error)
}
