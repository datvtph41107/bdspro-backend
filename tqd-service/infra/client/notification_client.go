package client

import (
	"context"
	clients "pb/clients"
	"tqd/env"
	"tqd/internal/dto"
	"tqd/internal/interface/provider"
)

// NotificationClient implements NotificationProvider interface
type NotificationClient struct {
	*clients.NotificationClient
}

// @bind: tqd/internal/interface/provider.NotificationProvider
func NewNotificationClient() provider.NotificationProvider {
	return &NotificationClient{
		NotificationClient: clients.MakeNotificationClient(env.ENV_RUNTIME),
	}
}

func (c *NotificationClient) SendNotification(ctx context.Context, notification *dto.NotificationDTO) error {
	// Return nil for now, implement later when needed
	return nil
}

func (c *NotificationClient) SendBatchNotification(ctx context.Context, notifications *dto.NotificationBatchDTO) error {
	// Return nil for now, implement later when needed
	return nil
}

func (c *NotificationClient) CreateHistory(ctx context.Context, history *dto.HistoryDTO) error {
	// Return nil for now, implement later when needed
	return nil
}

func (c *NotificationClient) GetNotificationHistory(ctx context.Context, userID uint64, page, size int) ([]*dto.HistoryDTO, int64, error) {
	// Return empty for now, implement later when needed
	return []*dto.HistoryDTO{}, 0, nil
}

func (c *NotificationClient) MarkAsRead(ctx context.Context, notificationID uint64, userID uint64) error {
	// Return nil for now, implement later when needed
	return nil
}

func (c *NotificationClient) GetUnreadCount(ctx context.Context, userID uint64) (int64, error) {
	// Return 0 for now, implement later when needed
	return 0, nil
}
