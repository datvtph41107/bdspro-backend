package iusecase

import (
	"context"
	"organization/internal/dto"
)

type INotiClient interface {
	SendNotification(ctx context.Context, notification *dto.NotificationDTO) error
	SendBatchNotification(ctx context.Context, notifications []*dto.NotificationDTO) error
	SendInvestmentNotification(ctx context.Context, dealID uint64, investmentID uint64, message string) error
}
