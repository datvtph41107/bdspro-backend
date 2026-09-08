package usecase

import (
	_utils "common/utils"
	"context"
	"fmt"
	"time"

	sharepb "pb/types/shared"

	"organization/env"
	"organization/internal/custom_error"
	"organization/internal/domain/entity"
	"organization/internal/domain/repository"
	"organization/internal/dto"
	iusecase "organization/internal/interface"
	"organization/pkg/utils"
)

type GroupMemberUsecase interface {
	CreateGroupMember(ctx context.Context, member *entity.GroupMember) (*entity.GroupMember, error)
	CreateGroupMemberBatch(ctx context.Context, groupId uint64, userIds []uint64) ([]*entity.GroupMember, error)
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
	userGrpcClient        iusecase.IUserClient
}

func NewGroupMemberUsecase(groupMemberRepository repository.GroupMemberRepository,
	groupRepository repository.GroupRepository,
	logWorker *LogWorker,
	groupAuthUsecase GroupAuthUsecase,
	groupChatUsecase GroupChatUsecase,
	userGrpcClient iusecase.IUserClient,
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
	go func() {
		newCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		err = u.groupChatUsecase.AddUserToGroupChat(newCtx, member.GroupId, member.UserId)
		if err != nil {
			fmt.Printf("error add user to group chat: %v", err)
		}
	}()
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

func (u *groupMemberUsecase) CreateGroupMemberBatch(ctx context.Context, groupId uint64, userIds []uint64) ([]*entity.GroupMember, error) {
	currentUserId := _utils.GetProfileIdWithContext(ctx)
	err := u.groupAuthUsecase.IsLeaderGroup(ctx, uint32(groupId), uint32(currentUserId))
	if err != nil {
		return nil, err
	}

	existedMember, err := u.groupMemberRepository.GetByGroupIDAndUserIDs(ctx, uint32(groupId), userIds)
	if err != nil {
		return nil, err
	}
	if len(existedMember) > 0 {
		return nil, custom_error.RecordNotFound(fmt.Sprintf("some users %v are already a member of group %d", userIds, groupId))
	}
	go func() {
		newCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		for _, userId := range userIds {
			err = u.groupChatUsecase.AddUserToGroupChat(newCtx, uint32(groupId), uint32(userId))
			if err != nil {
				fmt.Printf("error add user to group chat: %v", err)
			}
		}
	}()
	members := make([]*entity.GroupMember, len(userIds))
	for i, userId := range userIds {
		members[i] = &entity.GroupMember{
			GroupId:   uint32(groupId),
			UserId:    uint32(userId),
			CreatedBy: uint32(currentUserId),
			UpdatedBy: uint32(currentUserId),
		}
	}
	createdMembers, err := u.groupMemberRepository.CreateBatch(ctx, members)
	if err != nil {
		return nil, err
	}

	// u.logWorker.Push(LogEvent{
	// 	GroupId: createdMember.GroupId,
	// 	ActorId: createdMember.CreatedBy,
	// 	LogType: "CREATE_MEMBER",
	// 	LogData: "Added new member to group",
	// })

	return createdMembers, nil
}

func (u *groupMemberUsecase) UpdateGroupMember(ctx context.Context, member *entity.GroupMember) (*entity.GroupMember, error) {
	currentUserId := utils.GetUserID(ctx, env.USER_CONTEXT)
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

	currentUserId := utils.GetUserID(ctx, env.USER_CONTEXT)
	err = u.groupAuthUsecase.IsLeaderGroup(ctx, member.GroupId, currentUserId)
	if err != nil {
		return err
	}
	go func() {
		newCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		err = u.groupChatUsecase.RemoveUserFromGroupChat(newCtx, member.GroupId, member.UserId)
		if err != nil {
			fmt.Printf("error remove user from group chat: %v", err)
		}
	}()

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
	currentUserId := utils.GetUserID(ctx, env.USER_CONTEXT)
	err := u.groupAuthUsecase.IsMemberGroup(ctx, id, currentUserId)
	if err != nil {
		return nil, err
	}

	return u.groupMemberRepository.GetByID(ctx, id)
}

func (u *groupMemberUsecase) GetGroupMembersByGroupID(ctx context.Context, groupID uint32) ([]*dto.GroupMemberWithProfile, error) {
	currentUserId := utils.GetUserID(ctx, env.USER_CONTEXT)
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
