package iusecase

import (
	"context"
	"organization/internal/dto"
)

type INotiClient interface {
	SendNotification(ctx context.Context, notification *dto.NotificationDTO) error
	SendBatchNotification(ctx context.Context, notifications []*dto.NotificationDTO) error

	HistoryDeal(ctx context.Context, history *dto.DealHistoryDTO) error
}
