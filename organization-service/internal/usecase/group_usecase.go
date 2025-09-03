package usecase

import (
	"context"

	_utils "common/utils"
	"organization/internal/domain/entity"
	"organization/internal/domain/repository"
	"organization/internal/dto"

	"github.com/hyperledger/fabric/common/flogging"
)

type GroupUsecase interface {
	CreateGroup(ctx context.Context, group *entity.Group) (*entity.Group, error)
	UpdateGroup(ctx context.Context, group *entity.Group) (*entity.Group, error)
	DeleteGroup(ctx context.Context, id uint32) error
	GetGroupByID(ctx context.Context, id uint32) (*entity.Group, error)
	GetGroups(ctx context.Context, offset, limit int) ([]*entity.Group, uint32, error)
	GetGroupByUserId(ctx context.Context) ([]*entity.Group, error)
	GetGroupByUserIdWithDetails(ctx context.Context) ([]*dto.GroupWithDetails, error)
	GetGroupsWithDetails(ctx context.Context, page, size int) ([]*dto.GroupWithDetails, uint32, error)
}

type groupUsecase struct {
	groupRepository  repository.GroupRepository
	logWorker        *LogWorker
	groupAuthUsecase GroupAuthUsecase
	groupChatUsecase GroupChatUsecase
	logger           *flogging.FabricLogger
}

func NewGroupUsecase(
	groupRepository repository.GroupRepository,
	logWorker *LogWorker,
	groupAuthUsecase GroupAuthUsecase,
	groupChatUsecase GroupChatUsecase,
) GroupUsecase {
	return &groupUsecase{
		groupRepository:  groupRepository,
		logWorker:        logWorker,
		groupAuthUsecase: groupAuthUsecase,
		groupChatUsecase: groupChatUsecase,
		logger:           flogging.MustGetLogger("group_usecase"),
	}
}

func (u *groupUsecase) CreateGroup(ctx context.Context, group *entity.Group) (*entity.Group, error) {
	currentUserId := _utils.GetProfileIdWithContext(ctx)
	group.CreatedBy = uint32(currentUserId)
	group.UpdatedBy = uint32(currentUserId)

	// Set default status if not provided
	if group.Status == "" {
		group.Status = entity.GroupStatusActive
	}

	createdGroup, err := u.groupRepository.Create(ctx, group)
	if err != nil {
		return nil, err
	}

	_, err = u.groupChatUsecase.CreateGroupChat(context.Background(), createdGroup.Id, []uint32{uint32(currentUserId)}, group.Name, group.AvatarUrl)
	if err != nil {
		u.logger.Errorf("Failed to create group chat: %v", err)
	}

	u.logWorker.Push(LogEvent{
		GroupId: createdGroup.Id,
		ActorId: createdGroup.CreatedBy,
		LogType: "CREATE_GROUP",
		LogData: "Created new group",
	})

	return createdGroup, nil
}

func (u *groupUsecase) UpdateGroup(ctx context.Context, group *entity.Group) (*entity.Group, error) {
	currentUserId := _utils.GetProfileIdWithContext(ctx)

	err := u.groupAuthUsecase.IsLeaderGroup(ctx, group.Id, uint32(currentUserId))
	if err != nil {
		return nil, err
	}
	group.UpdatedBy = uint32(currentUserId)
	updatedGroup, err := u.groupRepository.Update(ctx, group)
	if err != nil {
		return nil, err
	}

	u.logWorker.Push(LogEvent{
		GroupId: updatedGroup.Id,
		ActorId: updatedGroup.CreatedBy,
		LogType: "UPDATE",
		LogData: "Updated group information",
	})

	return updatedGroup, nil
}

func (u *groupUsecase) DeleteGroup(ctx context.Context, id uint32) error {
	currentUserId := _utils.GetProfileIdWithContext(ctx)
	err := u.groupAuthUsecase.IsLeaderGroup(ctx, id, uint32(currentUserId))
	if err != nil {
		return err
	}

	group, err := u.groupRepository.GetByID(ctx, id)
	if err != nil {
		return err
	}

	err = u.groupRepository.Delete(ctx, id)
	if err != nil {
		return err
	}

	u.logWorker.Push(LogEvent{
		GroupId: id,
		ActorId: group.CreatedBy,
		LogType: "DELETE_GROUP",
		LogData: "Deleted group",
	})

	return nil
}

func (u *groupUsecase) GetGroupByID(ctx context.Context, id uint32) (*entity.Group, error) {
	currentUserId := _utils.GetProfileIdWithContext(ctx)
	err := u.groupAuthUsecase.IsMemberGroup(ctx, id, uint32(currentUserId))
	if err != nil {
		return nil, err
	}

	return u.groupRepository.GetByID(ctx, id)
}

func (u *groupUsecase) GetGroups(ctx context.Context, page, size int) ([]*entity.Group, uint32, error) {
	currentUserId := _utils.GetProfileIdWithContext(ctx)

	return u.groupRepository.List(ctx, page, size, uint32(currentUserId))
}

func (u *groupUsecase) GetGroupByUserId(ctx context.Context) ([]*entity.Group, error) {
	currentUserId := _utils.GetProfileIdWithContext(ctx)
	return u.groupRepository.FindByUserId(ctx, uint32(currentUserId))
}

func (u *groupUsecase) GetGroupByUserIdWithDetails(ctx context.Context) ([]*dto.GroupWithDetails, error) {
	currentUserId := _utils.GetProfileIdWithContext(ctx)
	return u.groupRepository.FindByUserIdWithDetails(ctx, uint32(currentUserId))
}

func (u *groupUsecase) GetGroupsWithDetails(ctx context.Context, page, size int) ([]*dto.GroupWithDetails, uint32, error) {
	currentUserId := _utils.GetProfileIdWithContext(ctx)
	return u.groupRepository.ListWithDetails(ctx, page, size, uint32(currentUserId))
}
