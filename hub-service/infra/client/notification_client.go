package client

import (
	_enum "common/domain/enum"
	provider_grpc "common/provider/grpc"
	"context"
	"hub/env"
	notificationpb "pb/types/notification"
	sharepb "pb/types/shared"
)

// @bind: hub/internal/interface.INotificationClient
type NotificationClient struct {
	*provider_grpc.NotificationGrpcClient
}

func NewNotificationClient() *NotificationClient {
	return &NotificationClient{
		NotificationGrpcClient: provider_grpc.MakeNotificationGrpcClient(env.ENV_RUNTIME),
	}
}

func (c *NotificationClient) CreateToOwner(
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
) error {
	_, err := c.Client.CreateToOwner(ctx, &sharepb.NotificationRequest{
		Avatar:           avatar,
		Title:            title,
		Message:          message,
		NotificationType: uint32(notificationType),
		OwnerId:          ownerId,
		TargetId:         *targetId,
		OwnerOf:          uint32(ownerOf),
		AttachData:       attachData,
		SendToDevice:     true,
		IsMerge:          isMerge,
	})
	return err
}

func (c *NotificationClient) SendBatch(
	ctx context.Context,
	requests []*sharepb.NotificationRequest,
) error {
	_, err := c.Client.SendBatch(ctx, &notificationpb.NotiBatchRequest{
		Datas:        requests,
		SendToDevice: true,
	})
	return err
}
