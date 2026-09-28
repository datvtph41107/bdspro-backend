package services

import (
	"context"

	notificationpb "pb/types/notification"

	"userv4/modules/notification/dto"
)

// NotificationService owns the user inbox flows migrated into modules/notification.
type NotificationService interface {
	Count(ctx context.Context, req dto.CountRequest) (*notificationpb.CountResponse, error)
	List(ctx context.Context, req dto.ListRequest) (*notificationpb.ListResponse, error)
	Read(ctx context.Context, id uint64) (*notificationpb.StatusResponse, error)
	ReadMany(ctx context.Context, ids []uint64) (*notificationpb.StatusResponse, error)
	ReadAll(ctx context.Context) (*notificationpb.StatusResponse, error)
	Remove(ctx context.Context, id uint64) (*notificationpb.StatusResponse, error)
}
