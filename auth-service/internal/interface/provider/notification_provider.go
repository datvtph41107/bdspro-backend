package provider

import "context"

type NotificationProvider interface {
	CreateAdminHistory(ctx context.Context,
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
		userAgent string) error

	CreateNotification(ctx context.Context,
		title string,
		message string,
		notificationType int32,
		link string,
		ownerID uint64) error
}
