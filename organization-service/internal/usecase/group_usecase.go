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
}

type groupUsecase struct {
	groupRepository    repository.GroupRepository
	logActivityUsecase GroupLogActivityUsecase
	groupAuthUsecase   GroupAuthUsecase
	groupChatUsecase   GroupChatUsecase
	logger             *flogging.FabricLogger
}

func NewGroupUsecase(
	groupRepository repository.GroupRepository,
	logActivityUsecase GroupLogActivityUsecase,
	groupAuthUsecase GroupAuthUsecase,
	groupChatUsecase GroupChatUsecase,
) GroupUsecase {
	return &groupUsecase{
		groupRepository:    groupRepository,
		logActivityUsecase: logActivityUsecase,
		groupAuthUsecase:   groupAuthUsecase,
		groupChatUsecase:   groupChatUsecase,
		logger:             flogging.MustGetLogger("group_usecase"),
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

	go func() {
		timeoutCtx, cancel := context.WithTimeout(context.Background(), appconfig.SUB_PROCESS_CONTEXT_TIMEOUT)
		defer cancel()
		_, err = u.groupChatUsecase.CreateGroupChat(timeoutCtx, createdGroup.Id, []uint32{currentUserId}, group.Name, group.AvatarUrl)
		if err != nil {
			u.logger.Errorf("Failed to create group chat: %v", err)
		}
	}()

	go func() {
		timeoutCtx, cancel := context.WithTimeout(context.Background(), appconfig.SUB_PROCESS_CONTEXT_TIMEOUT)
		defer cancel()
		log := &entity.GroupLogActivity{
			GroupId: createdGroup.Id,
			ActorId: createdGroup.CreatedBy,
			LogType: "CREATE",
			LogData: "Created new group",
		}
		u.logActivityUsecase.CreateGroupLogActivity(timeoutCtx, log)
	}()

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
