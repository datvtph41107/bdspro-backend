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
	logActivityUsecase GroupLogActivityUsecase
}

func NewGroupUsecase(groupRepository repository.GroupRepository, logActivityUsecase GroupLogActivityUsecase) GroupUsecase {
	return &groupUsecase{
		groupRepository: groupRepository,
		logActivityUsecase: logActivityUsecase,
	}
}

func (u *groupUsecase) CreateGroup(ctx context.Context, group *entity.Group) (*entity.Group, error) {
	createdGroup, err := u.groupRepository.Create(ctx, group)
	if err != nil {
		return nil, err
	}

	go func() {
		log := &entity.GroupLogActivity{
			GroupId: createdGroup.Id,
			ActorId: createdGroup.CreatedBy,
			LogType: "CREATE",
			LogData: "Created new group",
		}
		u.logActivityUsecase.CreateGroupLogActivity(context.Background(), log)
	}()

	return createdGroup, nil
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

	updatedGroup, err := u.groupRepository.Update(ctx, group)
	if err != nil {
		return nil, err
	}

	go func() {
		log := &entity.GroupLogActivity{
			GroupId: updatedGroup.Id,
			ActorId: updatedGroup.CreatedBy,
			LogType: "UPDATE",
			LogData: "Updated group information",
		}
		u.logActivityUsecase.CreateGroupLogActivity(context.Background(), log)
	}()

	return updatedGroup, nil
}

func (u *groupUsecase) DeleteGroup(ctx context.Context, id uint32) error {
	group, err := u.groupRepository.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if group == nil {
		return custom_error.RecordNotFound("group not found")
	}

	err = u.groupRepository.Delete(ctx, id)
	if err != nil {
		return err
	}

	go func() {
		log := &entity.GroupLogActivity{
			GroupId: id,
			ActorId: group.CreatedBy,
			LogType: "DELETE",
			LogData: "Deleted group",
		}
		u.logActivityUsecase.CreateGroupLogActivity(context.Background(), log)
	}()

	return nil
}

func (u *groupUsecase) GetGroupByID(ctx context.Context, id uint32) (*entity.Group, error) {
	return u.groupRepository.GetByID(ctx, id)
}

func (u *groupUsecase) GetGroups(ctx context.Context, offset, limit int) ([]*entity.Group, error) {
	return u.groupRepository.List(ctx, offset, limit)
}
