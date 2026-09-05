package provider_grpc

import (
	_enum "common/domain/enum"
	"common/tcp/connect"
	"context"
	"fmt"
	"log"
	notificationpb "pb/types/notification"
	sharepb "pb/types/shared"
)

type NotificationGrpcClient struct {
	Client notificationpb.InternalNotificationServiceClient
}

func MakeNotificationGrpcClient(env_runtime string) *NotificationGrpcClient {
	notificationClient := &NotificationGrpcClient{}
	name := "notification"
	addr := connect.AddressOf(env_runtime, name, connect.DefaultPort(name))

	conn, err := connect.DialNamed(name, addr)
	if err != nil {
		log.Printf("failed to dial notification at %s: %v", addr, err)
		return notificationClient
	}

	notificationClient.Client = notificationpb.NewInternalNotificationServiceClient(conn)
	return notificationClient
}

func (c *NotificationGrpcClient) CreateNotification(
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
	log.Println("avatar", avatar)
	_, err := c.Client.CreateToOwner(ctx, &sharepb.NotificationRequest{
		Avatar:           avatar,
		Title:            title,
		Message:          message,
		Type:             int32(notificationType),
		OwnerId:          ownerId,
		OwnerOf:          uint32(ownerOf),
		NotificationType: uint32(notificationType),
		TargetId:         *targetId,
		AttachData:       attachData,
		SendToDevice:     false,
	})
	return err
}

func (c *NotificationGrpcClient) RemoveNotification(
	ctx context.Context,
	ownerOf _enum.EOwnerOf,
	ownerId uint64,
	targetId uint64,
	notificationType _enum.ENotificationType,
) error {
	_, err := c.Client.RemoveToOwner(ctx, &sharepb.NotificationRequest{
		OwnerOf:          uint32(ownerOf),
		OwnerId:          ownerId,
		TargetId:         targetId,
		NotificationType: uint32(notificationType),
	})
	return err
}

func (c *NotificationGrpcClient) CreateHistoryAuth(
	ctx context.Context,
	req *notificationpb.HistoryAuthCreateRequest,
) (*notificationpb.HistoryAuthDTO, error) {
	return c.Client.CreateHistoryAuth(ctx, req)
}

func (c *NotificationGrpcClient) GetPersonConfigs(ctx context.Context, userID uint64) ([]*sharepb.PersonConfigV3Proto, error) {
	if userID == 0 {
		return []*sharepb.PersonConfigV3Proto{}, nil
	}

	if c.Client == nil {
		return nil, fmt.Errorf("notification grpc client chưa được khởi tạo")
	}

	resp, err := c.Client.GetPersonConfigs(ctx, &notificationpb.GetPersonConfigsRequest{
		UserId: userID,
	})
	if err != nil {
		return nil, err
	}

	return resp.GetData(), nil
}
