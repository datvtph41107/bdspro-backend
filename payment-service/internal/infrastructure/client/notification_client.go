package client

import (
	"context"

	clients "pb/clients"
	notificationpb "pb/types/notification"
	sharepb "pb/types/shared"
)

type NotificationClient interface {
	Create(ctx context.Context, req *sharepb.NotificationRequest) (*notificationpb.NotificationDTO, error)
	SendBatch(ctx context.Context, req []*sharepb.NotificationRequest) (*notificationpb.BatchResponse, error)
}

type notificationClient struct {
	client *clients.NotificationClient
}

func NewNotificationClient() NotificationClient {
	result := &notificationClient{
		client: clients.MakeNotificationClient("release"),
	}

	return result
}

func (c *notificationClient) Create(ctx context.Context, req *sharepb.NotificationRequest) (*notificationpb.NotificationDTO, error) {
	err := c.client.NotifyToUser(ctx, int32(req.Type), &req.OwnerId, []string{req.Message}, req.TargetId, req.IsMerge)
	return nil, err
}

func (c *notificationClient) SendBatch(ctx context.Context, req []*sharepb.NotificationRequest) (*notificationpb.BatchResponse, error) {
	err := c.client.SendBatch(ctx, req)
	return nil, err
}
