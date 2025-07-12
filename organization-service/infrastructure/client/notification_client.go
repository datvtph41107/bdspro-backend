package client

import (
	"context"
	"fmt"
	"organization/env"
	"organization/internal/dto"
	clients "pb/clients"
	sharepb "pb/types/shared"
)

// @bind: organization/internal/interface.INotiClient
type NotificationClient struct {
	client *clients.NotificationClient
}

func NewNotificationClient() *NotificationClient {
	return &NotificationClient{
		client: clients.MakeNotificationClient(env.ENV_RUNTIME),
	}
}

func (c *NotificationClient) SendNotification(ctx context.Context, notification *dto.NotificationDTO) error {
	return c.client.NotifyToUser(ctx,
		int32(notification.Type),
		notification.UserId,
		notification.AttachData,
		notification.TargetID,
		notification.IsMerge,
	)
}

func (c *NotificationClient) SendBatchNotification(ctx context.Context, notifications []*dto.NotificationDTO) error {
	notiRequests := make([]*sharepb.NotificationRequest, 0)
	for _, notification := range notifications {
		notiRequests = append(notiRequests, &sharepb.NotificationRequest{
			Type:     int32(notification.Type),
			OwnerId:  *notification.UserId,
			TargetId: notification.TargetID,
			// TargetType: int32(notification.Type),
			IsMerge:    notification.IsMerge,
			AttachData: notification.AttachData,
			Link:       fmt.Sprintf("/news-feed/%d", notification.TargetID),
		})
	}
	return c.client.SendBatch(ctx, notiRequests)
}

func (c *NotificationClient) SendInvestmentNotification(ctx context.Context, dealID uint64, investmentID uint64, message string) error {
	// TODO: Implement investment notification logic
	// For now, return nil to allow build to pass
	return nil
}
