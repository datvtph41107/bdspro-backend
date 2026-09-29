package repository

import (
	"context"

	"organization/internal/domain/entity"
)

type GroupMemberRepository interface {
	Create(ctx context.Context, member *entity.GroupMember) (*entity.GroupMember, error)
	CreateBatch(ctx context.Context, members []*entity.GroupMember) ([]*entity.GroupMember, error)
	Update(ctx context.Context, member *entity.GroupMember) (*entity.GroupMember, error)
	Delete(ctx context.Context, id uint32) error
	GetByID(ctx context.Context, id uint32) (*entity.GroupMember, error)
	GetByGroupID(ctx context.Context, groupID uint32) ([]*entity.GroupMember, error)
	GetByUserID(ctx context.Context, userID uint32) ([]*entity.GroupMember, error)
	GetByGroupIDAndUserID(ctx context.Context, groupID uint32, userID uint32) (*entity.GroupMember, error)
	GetByGroupIDAndUserIDs(ctx context.Context, groupID uint32, userIDs []uint64) ([]*entity.GroupMember, error)
}
