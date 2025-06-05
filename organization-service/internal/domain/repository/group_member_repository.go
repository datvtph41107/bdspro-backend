package repository

import (
	"context"

	"github.com/tuandq2112/organization-service/internal/domain/entity"
)

type GroupMemberRepository interface {
	Create(ctx context.Context, member *entity.GroupMember) (*entity.GroupMember, error)
	Update(ctx context.Context, member *entity.GroupMember) (*entity.GroupMember, error)
	Delete(ctx context.Context, id uint32) error
	GetByID(ctx context.Context, id uint32) (*entity.GroupMember, error)
	GetByGroupID(ctx context.Context, groupID uint32) ([]*entity.GroupMember, error)
	GetByUserID(ctx context.Context, userID uint32) ([]*entity.GroupMember, error)
} 