package client

import (
	"context"
	clients "pb/clients"
	"tqd/env"
	"tqd/internal/dto"
	"tqd/internal/interface/provider"
)

// NotificationClient implements NotificationProvider interface
type NotificationClient struct {
	*clients.NotificationClient
}

// @bind: tqd/internal/interface/provider.NotificationProvider
func NewNotificationClient() provider.NotificationProvider {
	return &NotificationClient{
		NotificationClient: clients.MakeNotificationClient(env.ENV_RUNTIME),
	}
}

func (c *NotificationClient) SendNotification(ctx context.Context, notification *dto.NotificationDTO) error {
	// Return nil for now, implement later when needed
	return nil
}

func (c *NotificationClient) SendBatchNotification(ctx context.Context, notifications *dto.NotificationBatchDTO) error {
	// Return nil for now, implement later when needed
	return nil
}

func (c *NotificationClient) CreateHistory(ctx context.Context, history *dto.HistoryDTO) error {
	// Return nil for now, implement later when needed
	return nil
}

func (c *NotificationClient) GetNotificationHistory(ctx context.Context, userID uint64, page, size int) ([]*dto.HistoryDTO, int64, error) {
	// Return empty for now, implement later when needed
	return []*dto.HistoryDTO{}, 0, nil
}

func (c *NotificationClient) MarkAsRead(ctx context.Context, notificationID uint64, userID uint64) error {
	// Return nil for now, implement later when needed
	return nil
}

func (c *NotificationClient) GetUnreadCount(ctx context.Context, userID uint64) (int64, error) {
	// Return 0 for now, implement later when needed
	return 0, nil
}

// SendWatchAlert gửi thông báo khi có thay đổi quy hoạch
func (c *NotificationClient) SendWatchAlert(ctx context.Context, userID uint64, title string, message []string, targetID uint64, attachData []string) error {
	// Gọi CreateNotification từ notification-service
	return c.CreateNotification(ctx, title, message, 10, "", userID, attachData)
}

// SendReportReady gửi thông báo khi báo cáo hoàn thành
func (c *NotificationClient) SendReportReady(ctx context.Context, userID uint64, reportID uint64, fileURL string) error {
	message := []string{"Báo cáo của bạn đã sẵn sàng", "Nhấn để tải xuống"}
	return c.CreateNotification(ctx, "Báo cáo đã hoàn thành", message, 20, fileURL, userID, []string{})
}

// SendSystemNotification gửi thông báo hệ thống
func (c *NotificationClient) SendSystemNotification(ctx context.Context, userID uint64, title string, message []string) error {
	return c.CreateNotification(ctx, title, message, 30, "", userID, []string{})
}

// CreateNotification là wrapper gọi sang notification-service
func (c *NotificationClient) CreateNotification(ctx context.Context, title string, message []string, notificationType int32, link string, ownerID uint64, attachData []string) error {
	return c.NotificationClient.CreateNotification(ctx, title, message, notificationType, link, ownerID, attachData)
}
