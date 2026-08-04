package provider_grpc

import (
	"common"
	_enum "common/domain/enum"
	_middleware "common/middleware"
	"context"
	"fmt"
	"log"
	notificationpb "pb/types/notification"
	sharepb "pb/types/shared"
	"time"

	"google.golang.org/grpc"
)

type NotificationGrpcClient struct {
	Client notificationpb.InternalNotificationServiceClient
}

func MakeNotificationGrpcClient(env_runtime string) *NotificationGrpcClient {
	notificationClient := &NotificationGrpcClient{}
	go func() {
		name := "notification"
		url := common.GetPrefixProtobufUrl(env_runtime, name)

		conn, err := grpc.Dial(fmt.Sprintf("%s:8204", url),
			grpc.WithInsecure(),
			grpc.WithBackoffMaxDelay(5*time.Second),
			grpc.WithUnaryInterceptor(_middleware.UnaryClientInterceptor),
			grpc.WithStreamInterceptor(_middleware.StreamClientInterceptor),
		)
		if err != nil {
			log.Fatalf("failed to dial: %v", err)
		}

		ctx, _ := context.WithCancel(context.Background())

		// Khởi chạy goroutine để theo dõi trạng thái kết nối
		go _middleware.MonitorConnection(ctx, name, conn)

		client := notificationpb.NewInternalNotificationServiceClient(conn)
		// historyClient := notificationpb.NewHistoryServiceClient(conn)
		notificationClient.Client = client
		// notificationClient.historyClient = historyClient
	}()
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
