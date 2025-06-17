package usecase

import (
	"context"

	"github.com/hyperledger/fabric/common/flogging"
	"github.com/tuandq2112/organization-service/appconfig"
	"github.com/tuandq2112/organization-service/internal/domain/entity"
	"github.com/tuandq2112/organization-service/internal/domain/repository"
	"github.com/tuandq2112/organization-service/pkg/utils"
)

type GroupUsecase interface {
	CreateGroup(ctx context.Context, group *entity.Group) (*entity.Group, error)
	UpdateGroup(ctx context.Context, group *entity.Group) (*entity.Group, error)
	DeleteGroup(ctx context.Context, id uint32) error
	GetGroupByID(ctx context.Context, id uint32) (*entity.Group, error)
	GetGroups(ctx context.Context, offset, limit int) ([]*entity.Group, uint32, error)
	GetGroupByUserId(ctx context.Context, userId uint32) ([]*entity.Group, error)
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
	currentUserId := utils.GetUserID(ctx, appconfig.USER_CONTEXT)
	group.CreatedBy = currentUserId
	group.UpdatedBy = currentUserId

	createdGroup, err := u.groupRepository.Create(ctx, group)
	if err != nil {
		return nil, err
	}

	_, err = u.groupChatUsecase.CreateGroupChat(context.Background(), createdGroup.Id, []uint32{currentUserId}, group.Name, group.AvatarUrl)
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
	currentUserId := utils.GetUserID(ctx, appconfig.USER_CONTEXT)

	err := u.groupAuthUsecase.IsLeaderGroup(ctx, group.Id, currentUserId)
	if err != nil {
		return nil, err
	}
	group.UpdatedBy = currentUserId
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
	currentUserId := utils.GetUserID(ctx, appconfig.USER_CONTEXT)
	err := u.groupAuthUsecase.IsLeaderGroup(ctx, id, currentUserId)
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
	currentUserId := utils.GetUserID(ctx, appconfig.USER_CONTEXT)
	err := u.groupAuthUsecase.IsMemberGroup(ctx, id, currentUserId)
	if err != nil {
		return nil, err
	}

	return u.groupRepository.GetByID(ctx, id)
}

func (u *groupUsecase) GetGroups(ctx context.Context, page, size int) ([]*entity.Group, uint32, error) {
	currentUserId := utils.GetUserID(ctx, appconfig.USER_CONTEXT)

	return u.groupRepository.List(ctx, page, size, currentUserId)
}

func (u *groupUsecase) GetGroupByUserId(ctx context.Context, userId uint32) ([]*entity.Group, error) {
	return u.groupRepository.FindByUserId(ctx, userId)
}
