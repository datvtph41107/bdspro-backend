package client

import (
	_enum "common/domain/enum"
	_provider "common/domain/provider"
	provider_grpc "common/provider/grpc"
	"context"
	"feedback/env"
	sharepb "pb/types/shared"
)

type NotificationClient struct {
	*provider_grpc.NotificationGrpcClient
}

func NewNotificationClient() _provider.NotificationProvider {
	return &NotificationClient{
		NotificationGrpcClient: provider_grpc.MakeNotificationGrpcClient(env.ENV_RUNTIME),
	}
}

func (c *NotificationClient) CreateNotification(
	ctx context.Context,
	avatar string,
	title string,
	message []string,
	notificationType _enum.ENotificationType,
	targetId *uint64,
	ownerId uint64,
	ownerOf _enum.EOwnerOf,
	attachData []string,
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
	})
	return err
}

func (c *NotificationClient) RemoveNotification(
	ctx context.Context,
	ownerOf _enum.EOwnerOf,
	ownerId uint64,
	targetId uint64,
	notificationType _enum.ENotificationType,
) error {
	return c.NotificationGrpcClient.RemoveNotification(ctx, ownerOf, ownerId, targetId, notificationType)
}
