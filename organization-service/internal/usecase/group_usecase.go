package usecase

import (
	"context"

	"github.com/tuandq2112/organization-service/internal/custom_error"
	"github.com/tuandq2112/organization-service/internal/domain/entity"
	"github.com/tuandq2112/organization-service/internal/domain/repository"
)

type GroupUsecase interface {
	CreateGroup(ctx context.Context, group *entity.Group) (*entity.Group, error)
	UpdateGroup(ctx context.Context, group *entity.Group) (*entity.Group, error)
	DeleteGroup(ctx context.Context, id uint32) error
	GetGroupByID(ctx context.Context, id uint32) (*entity.Group, error)
	GetGroups(ctx context.Context, offset, limit int) ([]*entity.Group, error)
}

type groupUsecase struct {
	groupRepository repository.GroupRepository
}

func NewGroupUsecase(groupRepository repository.GroupRepository) GroupUsecase {
	return &groupUsecase{groupRepository: groupRepository}
}

func (u *groupUsecase) CreateGroup(ctx context.Context, group *entity.Group) (*entity.Group, error) {
	return u.groupRepository.Create(ctx, group)
}

func (u *groupUsecase) UpdateGroup(ctx context.Context, group *entity.Group) (*entity.Group, error) {
	existedGroup, err := u.groupRepository.GetByID(ctx, group.Id)
	if err != nil {
		return nil, err
	}

	if existedGroup == nil {
		return nil, custom_error.RecordNotFound("group not found")
	}

	if group.CreatedBy != existedGroup.CreatedBy {
		return nil, custom_error.Forbidden("you are not allowed to update this group")
	}

	return u.groupRepository.Update(ctx, group)
}

func (u *groupUsecase) DeleteGroup(ctx context.Context, id uint32) error {
	return u.groupRepository.Delete(ctx, id)
}

func (u *groupUsecase) GetGroupByID(ctx context.Context, id uint32) (*entity.Group, error) {
	return u.groupRepository.GetByID(ctx, id)
}

func (u *groupUsecase) GetGroups(ctx context.Context, offset, limit int) ([]*entity.Group, error) {
	return u.groupRepository.List(ctx, offset, limit)
}
