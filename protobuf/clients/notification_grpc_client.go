package clients

import (
	_middleware "common/middleware"
	"context"
	"fmt"
	"log"
	shared_dto "pb/dto"
	notificationpb "pb/types/notification"
	sharepb "pb/types/shared"
	"time"

	"google.golang.org/grpc"
)

type NotificationClient struct {
	client        notificationpb.NotificationServiceClient
	historyClient notificationpb.HistoryServiceClient
}

func MakeNotificationClient(env_runtime string) *NotificationClient {
	notificationClient := &NotificationClient{}
	go func() {
		url := "localhost"
		name := "notification"
		if env_runtime == "release" {
			url = name
		}

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

		client := notificationpb.NewNotificationServiceClient(conn)
		historyClient := notificationpb.NewHistoryServiceClient(conn)
		notificationClient.client = client
		notificationClient.historyClient = historyClient
	}()
	return notificationClient
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

func (s *NotificationClient) SendBatch(
	ctx context.Context,
	notifications []*sharepb.NotificationRequest,
) error {
	_, err := s.client.SendBatch(ctx, &notificationpb.NotiBatchRequest{
		Datas:        notifications,
		SendToDevice: false,
	})
	return err
}

func (s *NotificationClient) CreateHistory(ctx context.Context, history *shared_dto.HistoryDTO) error {
	_, err := s.historyClient.Create(ctx, &notificationpb.HistoryDTO{
		ActionType: int32(history.ActionType),
		TargetId:   history.TargetId,
		OwnerId:    history.OwnerID,
		Note:       history.Note,
		Title:      history.Title,
	})
	return err
}
