package repository

import (
	"context"

	"organization/internal/domain/entity"
)

type GroupLogActivityRepository interface {
	Create(ctx context.Context, log *entity.GroupLogActivity) (*entity.GroupLogActivity, error)
	GetByID(ctx context.Context, id uint32) (*entity.GroupLogActivity, error)
	GetByGroupID(ctx context.Context, groupID uint32, page, size int) ([]*entity.GroupLogActivity, uint32, error)
	GetByActorID(ctx context.Context, actorID uint32) ([]*entity.GroupLogActivity, error)
}
