package client

import (
	_middleware "common/middleware"
	"context"
	"fmt"
	"log"
	pb_notification "pb/types/notification"
	sharepb "pb/types/shared"
	"social/env"
	"time"

	"google.golang.org/grpc"
)

// @bind: social/internal/interface.NotiClient
type NotificationClient struct {
	client pb_notification.NotificationServiceClient
}

func NewNotificationClient() *NotificationClient {
	notiClient := &NotificationClient{client: nil}
	go func() {
		url := "localhost"
		name := "notification"
		if env.ENV_RUNTIME == "release" {
			url = name
		}

		conn, err := grpc.Dial(fmt.Sprintf("%s:8004", url),
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

		client := pb_notification.NewNotificationServiceClient(conn)
		notiClient.client = client
	}()

	return notiClient
}

func (s *NotificationClient) NotifyToUser(
	ctx context.Context,
	typeNoti int32,
	userId *uint64,
	attachData []string,
	targetID uint64,
	isMerge bool,
) error {
	_, err := s.client.Create(ctx, &sharepb.NotificationRequest{
		AttachData: attachData,
		Link:       fmt.Sprintf("/news-feed/%d", targetID),
		Type:       typeNoti,
		OwnerId:    *userId,
		TargetId:   targetID,
		TargetType: typeNoti,
		IsMerge:    isMerge,
	})
	return err
}
