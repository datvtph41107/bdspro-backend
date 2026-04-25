package provider

import (
	"context"
	"tqd/internal/dto"
)

// NotificationProvider defines the interface for notification operations
type NotificationProvider interface {
	// SendNotification sends a single notification
	SendNotification(ctx context.Context, notification *dto.NotificationDTO) error

	// SendBatchNotification sends multiple notifications
	SendBatchNotification(ctx context.Context, notifications *dto.NotificationBatchDTO) error

	// CreateHistory creates notification history
	CreateHistory(ctx context.Context, history *dto.HistoryDTO) error

	// GetNotificationHistory gets notification history for a user
	GetNotificationHistory(ctx context.Context, userID uint64, page, size int) ([]*dto.HistoryDTO, int64, error)

	// MarkAsRead marks notification as read
	MarkAsRead(ctx context.Context, notificationID uint64, userID uint64) error

	// GetUnreadCount gets unread notification count for a user
	GetUnreadCount(ctx context.Context, userID uint64) (int64, error)
	SendWatchAlert(ctx context.Context, userID uint64, title string, message []string, targetID uint64, attachData []string) error
	// Gửi thông báo khi báo cáo hoàn thành
	SendReportReady(ctx context.Context, userID uint64, reportID uint64, fileURL string) error
	// Gửi thông báo hệ thống
	SendSystemNotification(ctx context.Context, userID uint64, title string, message []string) error
}
