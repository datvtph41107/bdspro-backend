package providers

import (
	_dto "common/domain/dto"
	"context"
	sharepb "pb/types/shared"
)

// NotificationProvider định nghĩa interface cho việc tương tác với notification-service
type NotificationProvider interface {
	// SendWarningNotification gửi cảnh báo cho user
	// SendWarningNotification(ctx context.Context, userID uint64, warningType int32, title, content string) (*notificationpb.HistoryDTO, error)

	// CreateHistory tạo lịch sử hoạt động
	CreateHistory(ctx context.Context, history *_dto.HistoryDTO) error

	// CreateInternalHistory tạo lịch sử nội bộ
	// CreateInternalHistory(ctx context.Context, history *notificationpb.HistoryDTO) (*notificationpb.HistoryDTO, error)

	BatchPersonConfig(ctx context.Context, configs []*sharepb.PersonConfigV3Proto) error
}
