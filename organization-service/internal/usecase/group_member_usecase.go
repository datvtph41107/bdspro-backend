package usecase

import (
	_utils "common/utils"
	"context"
	"fmt"

	sharepb "pb/types/shared"

	"github.com/tuandq2112/organization-service/appconfig"
	"github.com/tuandq2112/organization-service/internal/custom_error"
	"github.com/tuandq2112/organization-service/internal/domain/entity"
	"github.com/tuandq2112/organization-service/internal/domain/repository"
	"github.com/tuandq2112/organization-service/internal/infrastructure/client"
	"github.com/tuandq2112/organization-service/internal/interface/dto"
	"github.com/tuandq2112/organization-service/pkg/utils"
)

type GroupMemberUsecase interface {
	CreateGroupMember(ctx context.Context, member *entity.GroupMember) (*entity.GroupMember, error)
	UpdateGroupMember(ctx context.Context, member *entity.GroupMember) (*entity.GroupMember, error)
	DeleteGroupMember(ctx context.Context, id uint32) error
	GetGroupMemberByID(ctx context.Context, id uint32) (*entity.GroupMember, error)
	GetGroupMembersByGroupID(ctx context.Context, groupID uint32) ([]*dto.GroupMemberWithProfile, error)
}

type groupMemberUsecase struct {
	groupMemberRepository repository.GroupMemberRepository
	groupRepository       repository.GroupRepository
	logWorker             *LogWorker
	groupAuthUsecase      GroupAuthUsecase
	groupChatUsecase      GroupChatUsecase
	userGrpcClient        client.UserClient
}

func NewGroupMemberUsecase(groupMemberRepository repository.GroupMemberRepository,
	groupRepository repository.GroupRepository,
	logWorker *LogWorker,
	groupAuthUsecase GroupAuthUsecase,
	groupChatUsecase GroupChatUsecase,
	userGrpcClient client.UserClient,
) GroupMemberUsecase {
	return &groupMemberUsecase{
		groupMemberRepository: groupMemberRepository,
		groupRepository:       groupRepository,
		logWorker:             logWorker,
		groupAuthUsecase:      groupAuthUsecase,
		groupChatUsecase:      groupChatUsecase,
		userGrpcClient:        userGrpcClient,
	}
}

func (u *groupMemberUsecase) CreateGroupMember(ctx context.Context, member *entity.GroupMember) (*entity.GroupMember, error) {
	currentUserId := _utils.GetProfileIdWithContext(ctx)
	err := u.groupAuthUsecase.IsLeaderGroup(ctx, member.GroupId, uint32(currentUserId))
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
	member.CreatedBy = uint32(currentUserId)
	member.UpdatedBy = uint32(currentUserId)
	createdMember, err := u.groupMemberRepository.Create(ctx, member)
	if err != nil {
		return nil, err
	}

	u.logWorker.Push(LogEvent{
		GroupId: createdMember.GroupId,
		ActorId: createdMember.CreatedBy,
		LogType: "CREATE_MEMBER",
		LogData: "Added new member to group",
	})

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

	u.logWorker.Push(LogEvent{
		GroupId: updatedMember.GroupId,
		ActorId: updatedMember.CreatedBy,
		LogType: "UPDATE_MEMBER",
		LogData: "Updated member information",
	})

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

	u.logWorker.Push(LogEvent{
		GroupId: member.GroupId,
		ActorId: member.CreatedBy,
		LogType: "DELETE_MEMBER",
		LogData: "Removed member from group",
	})

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

func (u *groupMemberUsecase) GetGroupMembersByGroupID(ctx context.Context, groupID uint32) ([]*dto.GroupMemberWithProfile, error) {
	currentUserId := utils.GetUserID(ctx, appconfig.USER_CONTEXT)
	err := u.groupAuthUsecase.IsMemberGroup(ctx, groupID, currentUserId)
	if err != nil {
		return nil, err
	}

	members, err := u.groupMemberRepository.GetByGroupID(ctx, groupID)
	if err != nil {
		return nil, err
	}
	profileIds := make([]uint64, len(members))
	for i, member := range members {
		profileIds[i] = uint64(member.UserId)
	}
	profiles, err := u.userGrpcClient.GetProfileByIds(ctx, &sharepb.GetProfileByIdsRequest{
		Ids: profileIds,
	})
	if err != nil {
		return nil, err
	}
	membersWithProfiles := make([]*dto.GroupMemberWithProfile, len(members))
	for i, member := range members {
		membersWithProfiles[i] = &dto.GroupMemberWithProfile{
			GroupMember: member,
			Profile:     profiles.Profiles[i],
		}
	}
	return membersWithProfiles, nil
}
