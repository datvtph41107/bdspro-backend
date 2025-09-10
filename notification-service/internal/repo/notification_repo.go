package repo

import (
	"context"
	"notification/internal/domain"
	"notification/internal/dto"
	"notification/internal/enums"
)

type NotificationRepo interface {
	Search(ctx context.Context, profileId uint64, dto *dto.SearchNotiRequest) ([]domain.NotificationEntity, int32, error)
	Create(ctx context.Context, notification *domain.NotificationEntity) error
	MarkAsRead(ctx context.Context, profileId uint64, id uint64) (int64, error)
	ReadAll(ctx context.Context, profileId uint64) (int64, error)
	Count(ctx context.Context, profileId uint64) (int64, error)
	Delete(ctx context.Context, profileId uint64, id uint64) error
	ListNotiWithTargetIdAndOwnerType(ctx context.Context, targetId uint64, ownerType int32, hour int) ([]domain.NotificationEntity, error)
	DeleteIds(ctx context.Context, ids []uint64) error
	RemoveUnique(ctx context.Context, ownerId, targetId uint64, typeNoti enums.TypeEnum, ownerType enums.OwnerType) (*domain.NotificationEntity, error)
}
