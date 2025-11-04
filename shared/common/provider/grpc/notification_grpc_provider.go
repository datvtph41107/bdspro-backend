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
