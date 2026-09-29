package usecase

import (
	"context"

	"organization/internal/domain/entity"
	"organization/internal/domain/repository"
)

type GroupLogActivityUsecase interface {
	CreateGroupLogActivity(ctx context.Context, log *entity.GroupLogActivity) (*entity.GroupLogActivity, error)
	GetGroupLogActivityByID(ctx context.Context, id uint32) (*entity.GroupLogActivity, error)
	GetGroupLogActivitiesByGroupID(ctx context.Context, groupID uint32, page, size int) ([]*entity.GroupLogActivity, uint32, error)
}

type groupLogActivityUsecase struct {
	groupLogActivityRepository repository.GroupLogActivityRepository
}

func NewGroupLogActivityUsecase(groupLogActivityRepository repository.GroupLogActivityRepository) GroupLogActivityUsecase {
	return &groupLogActivityUsecase{groupLogActivityRepository: groupLogActivityRepository}
}

func (u *groupLogActivityUsecase) CreateGroupLogActivity(ctx context.Context, log *entity.GroupLogActivity) (*entity.GroupLogActivity, error) {
	return u.groupLogActivityRepository.Create(ctx, log)
}

func (u *groupLogActivityUsecase) GetGroupLogActivityByID(ctx context.Context, id uint32) (*entity.GroupLogActivity, error) {
	return u.groupLogActivityRepository.GetByID(ctx, id)
}

func (u *groupLogActivityUsecase) GetGroupLogActivitiesByGroupID(ctx context.Context, groupID uint32, page, size int) ([]*entity.GroupLogActivity, uint32, error) {
	return u.groupLogActivityRepository.GetByGroupID(ctx, groupID, page, size)
}
