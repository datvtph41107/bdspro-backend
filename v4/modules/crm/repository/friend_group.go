package repository

import (
	"context"
	domain "userv4/modules/crm/entity"
)

type FriendGroupRepo interface {
	CreateGroup(c context.Context, group *domain.GroupEntity) error
	UpdateGroup(c context.Context, id uint64, group *domain.GroupEntity) error
	DeleteGroup(c context.Context, profileId *uint64, id uint64) error
	ListGroups(c context.Context, createdBy uint64) ([]domain.GroupEntity, error)
	GetGroupByID(c context.Context, id uint64) (*domain.GroupEntity, error)
}