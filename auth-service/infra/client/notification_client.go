package client

import (
	"context"
	"fmt"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// @bind: auth/internal/interface/provider.NotificationProvider
type NotificationClient struct {
	client *grpc.ClientConn
}

func NewNotificationClient() *NotificationClient {
	notificationClient := &NotificationClient{}
	go func() {
		url := "localhost"
		// TODO: Add environment check for production
		// if env_runtime == "release" {
		// 	url = "notification"
		// }

		conn, err := grpc.Dial(fmt.Sprintf("%s:8204", url),
			grpc.WithTransportCredentials(insecure.NewCredentials()),
			grpc.WithBackoffMaxDelay(5*time.Second),
		)
		if err != nil {
			log.Printf("Failed to connect to notification service: %v", err)
			return
		}

		notificationClient.client = conn
		log.Printf("Connected to notification service at %s:8204", url)
	}()
	return notificationClient
}

// LogAdminLogin ghi log khi admin đăng nhập thành công
func (c *NotificationClient) LogAdminLogin(ctx context.Context, adminID uint64, adminRole string, ipAddress string, userAgent string) error {
	// TODO: Implement gRPC call to notification service
	// For now, just log the action
	log.Printf("Admin login logged: AdminID=%d, Role=%s, IP=%s, UserAgent=%s",
		adminID, adminRole, ipAddress, userAgent)
	return nil
}

// CreateAdminHistory tạo lịch sử admin
func (c *NotificationClient) CreateAdminHistory(ctx context.Context,
	adminID uint64,
	targetID uint64,
	targetType int32,
	actionType int32,
	title string,
	notes []string,
	preStage string,
	afterStage string,
	ownerID *uint64,
	ownerType int32,
	adminRole string,
	ipAddress string,
	userAgent string,
) error {
	// TODO: Implement gRPC call to notification service
	log.Printf("Admin history created: AdminID=%d, TargetID=%d, Action=%s",
		adminID, targetID, title)
	return nil
}

// CreateAdminHistoryBatch tạo hàng loạt lịch sử admin
func (c *NotificationClient) CreateAdminHistoryBatch(ctx context.Context, requests []interface{}) error {
	// TODO: Implement gRPC call to notification service
	log.Printf("Admin history batch created: %d records", len(requests))
	return nil
}

// LogAdminAction ghi log hành động của admin
func (c *NotificationClient) LogAdminAction(ctx context.Context,
	adminID uint64,
	targetID uint64,
	targetType int32,
	actionType int32,
	title string,
	notes []string,
	preStage string,
	afterStage string,
	ownerID *uint64,
	ownerType int32,
	adminRole string,
	ipAddress string,
	userAgent string,
) error {
	// TODO: Implement gRPC call to notification service
	// For now, just log the action
	log.Printf("Admin action logged: AdminID=%d, TargetID=%d, Action=%s, Title=%s",
		adminID, targetID, title, title)
	return nil
}
