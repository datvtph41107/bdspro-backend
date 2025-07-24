package client

import (
	"context"
	"fmt"
	"log"
	"marketing/env"

	notificationpb "pb/types/notification"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type NotificationClient struct {
	client notificationpb.HistoryServiceClient
}

func NewNotificationClient() *NotificationClient {
	notificationServiceURL := fmt.Sprintf("%s:%s", env.NOTIFICATION_SERVICE_HOST, env.NOTIFICATION_SERVICE_PORT)
	
	conn, err := grpc.Dial(notificationServiceURL, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Printf("Failed to connect to notification service: %v", err)
		return nil
	}

	client := notificationpb.NewHistoryServiceClient(conn)
	return &NotificationClient{
		client: client,
	}
}

// CreateHistory creates a new history record
func (nc *NotificationClient) CreateHistory(ctx context.Context, history *notificationpb.HistoryDTO) (*notificationpb.HistoryDTO, error) {
	if nc.client == nil {
		return nil, fmt.Errorf("notification client not initialized")
	}

	response, err := nc.client.Create(ctx, history)
	if err != nil {
		return nil, fmt.Errorf("failed to create history: %w", err)
	}

	return response, nil
}

// CreateInternalHistory creates an internal history record
func (nc *NotificationClient) CreateInternalHistory(ctx context.Context, history *notificationpb.HistoryDTO) (*notificationpb.HistoryDTO, error) {
	if nc.client == nil {
		return nil, fmt.Errorf("notification client not initialized")
	}

	response, err := nc.client.CreateInternal(ctx, history)
	if err != nil {
		return nil, fmt.Errorf("failed to create internal history: %w", err)
	}

	return response, nil
}

// Lưu ý: Các method GetCampaignHistory và GetPackageHistory đã được xóa
// Sử dụng trực tiếp API từ notification service:
// GET /v2/notification/history/campaign/{id}
// GET /v2/notification/history/package/{id} 