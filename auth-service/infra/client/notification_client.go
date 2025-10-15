package client

import (
	"context"
	"pb/clients"

	"auth/env"
)

// @bind: auth/internal/interface/provider.NotificationProvider
type NotificationClient struct {
	client *clients.NotificationClient
}

func NewNotificationClient() *NotificationClient {
	notificationClient := clients.MakeNotificationClient(env.ENV_RUNTIME)
	return &NotificationClient{
		client: notificationClient,
	}
}

// LogAdminLogin ghi log khi admin đăng nhập thành công
func (c *NotificationClient) LogAdminLogin(ctx context.Context, adminID uint64, adminRole string, ipAddress string, userAgent string) error {
	return c.client.LogAdminLogin(ctx, adminID, adminRole, ipAddress, userAgent)
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
	return c.client.CreateAdminHistory(ctx, adminID, targetID, targetType, actionType, title, notes, preStage, afterStage, ownerID, ownerType, adminRole, ipAddress, userAgent)
}

// CreateNotification tạo thông báo cho user
func (c *NotificationClient) CreateNotification(ctx context.Context,
	title string,
	message string,
	notificationType int32,
	link string,
	ownerID uint64,
) error {
	return c.client.CreateNotification(ctx, title, message, notificationType, link, ownerID)
}
