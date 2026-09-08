package providers

import (
	_dto "common/domain/dto"
	_enum "common/domain/enum"
	"context"
	sharepb "pb/types/shared"
	"user/internal/dto"
)

// NotificationProvider định nghĩa interface cho việc tương tác với notification-service
type NotificationProvider interface {
	// _provider.NotificationProvider
	// _provider.NotificationProvider
	// SendWarningNotification gửi cảnh báo cho user
	// SendWarningNotification(ctx context.Context, userID uint64, warningType int32, title, content string) (*notificationpb.HistoryDTO, error)

	// CreateHistory tạo lịch sử hoạt động
	CreateHistory(ctx context.Context, history *_dto.HistoryDTO) error

	// CreateInternalHistory tạo lịch sử nội bộ
	// CreateInternalHistory(ctx context.Context, history *notificationpb.HistoryDTO) (*notificationpb.HistoryDTO, error)
	CreateNotification(
		ctx context.Context,
		avatar string,
		title string,
		message []string,
		notificationType _enum.ENotificationType,
		targetId *uint64,
		ownerId uint64,
		ownerOf _enum.EOwnerOf,
		attachData []string,
	) error

	BatchPersonConfig(ctx context.Context, configs []*sharepb.PersonConfigV3Proto) error
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

	CreateHistoryAuth(ctx context.Context, payload *dto.HistoryAuthCreateDTO) error

	GetPersonConfigs(ctx context.Context, userID uint64) ([]*dto.PersonConfigDTO, error)
}
