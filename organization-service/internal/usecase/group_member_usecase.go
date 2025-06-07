package usecase

import (
	"context"
	"fmt"

	"github.com/tuandq2112/organization-service/appconfig"
	"github.com/tuandq2112/organization-service/internal/custom_error"
	"github.com/tuandq2112/organization-service/internal/domain/entity"
	"github.com/tuandq2112/organization-service/internal/domain/repository"
	"github.com/tuandq2112/organization-service/pkg/utils"
)

type GroupMemberUsecase interface {
	CreateGroupMember(ctx context.Context, member *entity.GroupMember) (*entity.GroupMember, error)
	UpdateGroupMember(ctx context.Context, member *entity.GroupMember) (*entity.GroupMember, error)
	DeleteGroupMember(ctx context.Context, id uint32) error
	GetGroupMemberByID(ctx context.Context, id uint32) (*entity.GroupMember, error)
	GetGroupMembersByGroupID(ctx context.Context, groupID uint32) ([]*entity.GroupMember, error)
}

type groupMemberUsecase struct {
	groupMemberRepository repository.GroupMemberRepository
	groupRepository       repository.GroupRepository
	logActivityUsecase    GroupLogActivityUsecase
	groupAuthUsecase      GroupAuthUsecase
	groupChatUsecase      GroupChatUsecase
}

func NewGroupMemberUsecase(groupMemberRepository repository.GroupMemberRepository,
	groupRepository repository.GroupRepository,
	logActivityUsecase GroupLogActivityUsecase,
	groupAuthUsecase GroupAuthUsecase,
	groupChatUsecase GroupChatUsecase,
) GroupMemberUsecase {
	return &groupMemberUsecase{
		groupMemberRepository: groupMemberRepository,
		groupRepository:       groupRepository,
		logActivityUsecase:    logActivityUsecase,
		groupAuthUsecase:      groupAuthUsecase,
		groupChatUsecase:      groupChatUsecase,
	}
}

func (u *groupMemberUsecase) CreateGroupMember(ctx context.Context, member *entity.GroupMember) (*entity.GroupMember, error) {
	currentUserId := utils.GetUserID(ctx, appconfig.USER_CONTEXT)
	err := u.groupAuthUsecase.IsLeaderGroup(ctx, member.GroupId, currentUserId)
	if err != nil {
		return nil, err
	}

	existedMember, err := u.groupMemberRepository.GetByGroupIDAndUserID(ctx, member.GroupId, member.UserId)
	if err != nil {
		return nil, err
	}
	if existedMember != nil {
		return nil, custom_error.RecordNotFound(fmt.Sprintf("user %d is already a member of group %d", member.UserId, member.GroupId))
	}
	err = u.groupChatUsecase.AddUserToGroupChat(ctx, member.GroupId, member.UserId)
	if err != nil {
		return nil, err
	}
	member.CreatedBy = currentUserId
	member.UpdatedBy = currentUserId
	createdMember, err := u.groupMemberRepository.Create(ctx, member)
	if err != nil {
		return nil, err
	}

	go func() {
		log := &entity.GroupLogActivity{
			GroupId: createdMember.GroupId,
			ActorId: createdMember.CreatedBy,
			LogType: "CREATE_MEMBER",
			LogData: "Added new member to group",
		}
		u.logActivityUsecase.CreateGroupLogActivity(context.Background(), log)
	}()

	return createdMember, nil
}

func (u *groupMemberUsecase) UpdateGroupMember(ctx context.Context, member *entity.GroupMember) (*entity.GroupMember, error) {
	currentUserId := utils.GetUserID(ctx, appconfig.USER_CONTEXT)
	err := u.groupAuthUsecase.IsLeaderGroup(ctx, member.GroupId, currentUserId)
	if err != nil {
		return nil, err
	}

	member.UpdatedBy = currentUserId
	updatedMember, err := u.groupMemberRepository.Update(ctx, member)
	if err != nil {
		return nil, err
	}

	go func() {
		log := &entity.GroupLogActivity{
			GroupId: updatedMember.GroupId,
			ActorId: updatedMember.CreatedBy,
			LogType: "UPDATE_MEMBER",
			LogData: "Updated member information",
		}
		u.logActivityUsecase.CreateGroupLogActivity(context.Background(), log)
	}()

	return updatedMember, nil
}

func (u *groupMemberUsecase) DeleteGroupMember(ctx context.Context, id uint32) error {
	member, err := u.groupMemberRepository.GetByID(ctx, id)
	if err != nil {
		return err
	}

	currentUserId := utils.GetUserID(ctx, appconfig.USER_CONTEXT)
	err = u.groupAuthUsecase.IsLeaderGroup(ctx, member.GroupId, currentUserId)
	if err != nil {
		return err
	}

	err = u.groupChatUsecase.RemoveUserFromGroupChat(ctx, member.GroupId, member.UserId)
	if err != nil {
		return err
	}

	err = u.groupMemberRepository.Delete(ctx, id)
	if err != nil {
		return err
	}

	go func() {
		log := &entity.GroupLogActivity{
			GroupId: member.GroupId,
			ActorId: member.CreatedBy,
			LogType: "DELETE_MEMBER",
			LogData: "Removed member from group",
		}
		u.logActivityUsecase.CreateGroupLogActivity(context.Background(), log)
	}()

	return nil
}

func (u *groupMemberUsecase) GetGroupMemberByID(ctx context.Context, id uint32) (*entity.GroupMember, error) {
	currentUserId := utils.GetUserID(ctx, appconfig.USER_CONTEXT)
	err := u.groupAuthUsecase.IsMemberGroup(ctx, id, currentUserId)
	if err != nil {
		return nil, err
	}

	return u.groupMemberRepository.GetByID(ctx, id)
}

func (u *groupMemberUsecase) GetGroupMembersByGroupID(ctx context.Context, groupID uint32) ([]*entity.GroupMember, error) {
	currentUserId := utils.GetUserID(ctx, appconfig.USER_CONTEXT)
	err := u.groupAuthUsecase.IsMemberGroup(ctx, groupID, currentUserId)
	if err != nil {
		return nil, err
	}

	return u.groupMemberRepository.GetByGroupID(ctx, groupID)
}
