package repo

import (
	"context"
	"tqd/internal/domain"
)

type NotificationRepository interface {
	Create(ctx context.Context, notif *domain.UserNotification) error
	MarkAsRead(ctx context.Context, userID uint64, notificationID uint64) error
	MarkAllAsRead(ctx context.Context, userID uint64) error
	ListByUser(ctx context.Context, userID uint64, notifType *uint32, read *bool, page, limit int) ([]domain.UserNotification, int64, error)
	CountUnread(ctx context.Context, userID uint64) (int64, error)
	DeleteExpired(ctx context.Context) error
}
