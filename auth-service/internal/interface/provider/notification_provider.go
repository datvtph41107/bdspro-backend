package provider

import "context"

type NotificationProvider interface {
	LogAdminLogin(ctx context.Context, adminID uint64, adminRole string, ipAddress string, userAgent string) error
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
	CreateAdminHistoryBatch(ctx context.Context, requests []interface{}) error
}
