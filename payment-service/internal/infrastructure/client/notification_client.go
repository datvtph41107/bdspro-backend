package client

import (
	"context"
	"fmt"

	pb_noti "pb/types/notification"

	"github.com/tuandq2112/payment-service/appconfig"
	"google.golang.org/grpc"
)

type NotificationClient interface {
	Create(ctx context.Context, req *pb_noti.NotiNewRequest) (*pb_noti.NotificationDTO, error)
	SendBatch(ctx context.Context, req *pb_noti.NotiBatchRequest) (*pb_noti.ListResponse, error)
}

type notificationClient struct {
	client pb_noti.NotificationServiceClient
}

func NewNotificationClient() NotificationClient {
	conn, err := grpc.NewClient(fmt.Sprintf("%s:%s", appconfig.NOTIFICATION_SERVICE_HOST, appconfig.NOTIFICATION_SERVICE_PORT), grpc.WithInsecure())
	if err != nil {
		panic(fmt.Sprintf("Failed to connect to notification service: %v", err))
	}
	return &notificationClient{
		client: pb_noti.NewNotificationServiceClient(conn),
	}
}

func (c *notificationClient) Create(ctx context.Context, req *pb_noti.NotiNewRequest) (*pb_noti.NotificationDTO, error) {
	return c.client.Create(ctx, req)
}

func (c *notificationClient) SendBatch(ctx context.Context, req *pb_noti.NotiBatchRequest) (*pb_noti.ListResponse, error) {
	return c.client.SendBatch(ctx, req)
}
