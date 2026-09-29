package repository

import (
	"context"

	"organization/internal/domain/entity"
)

type GroupNotificationRepository interface {
	Create(ctx context.Context, notification *entity.GroupNotification) (*entity.GroupNotification, error)
	Update(ctx context.Context, notification *entity.GroupNotification) (*entity.GroupNotification, error)
	Delete(ctx context.Context, id uint32) error
	GetByID(ctx context.Context, id uint32) (*entity.GroupNotification, error)
	GetByGroupID(ctx context.Context, groupID uint32) ([]*entity.GroupNotification, error)
}
