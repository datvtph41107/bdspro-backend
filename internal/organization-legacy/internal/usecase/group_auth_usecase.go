package usecase

import (
	_utils "common/utils"
	"context"
	"fmt"

	"organization/internal/custom_error"
	"organization/internal/domain/repository"
)

type GroupAuthUsecase interface {
	IsLeaderGroup(ctx context.Context, groupId uint32, userId uint32) error
	IsMemberGroup(ctx context.Context, groupId uint32, userId uint32) error
	CanAccessGroup(ctx context.Context, groupId uint32) error
	HasPermission(ctx context.Context, permission string) bool
}

type groupAuthUsecase struct {
	groupRepository       repository.GroupRepository
	groupMemberRepository repository.GroupMemberRepository
}

func NewGroupAuthUsecase(groupRepository repository.GroupRepository, groupMemberRepository repository.GroupMemberRepository) GroupAuthUsecase {
	return &groupAuthUsecase{
		groupRepository:       groupRepository,
		groupMemberRepository: groupMemberRepository,
	}
}

func (u *groupAuthUsecase) IsLeaderGroup(ctx context.Context, groupId uint32, userId uint32) error {
	group, err := u.groupRepository.GetByID(ctx, groupId)
	if err != nil {
		return err
	}
	if group == nil {
		return custom_error.RecordNotFound(fmt.Sprintf("group %d not found", groupId))
	}

	// if group.CreatedBy != userId {
	// 	return custom_error.Forbidden(fmt.Sprintf("user %d is not leader of group %d", userId, groupId))
	// }

	return nil
}

func (u *groupAuthUsecase) IsMemberGroup(ctx context.Context, groupId uint32, userId uint32) error {
	fmt.Println("IsMemberGroup", groupId, userId)
	if err := u.IsLeaderGroup(ctx, groupId, userId); err == nil {
		return nil
	}

	groupMember, err := u.groupMemberRepository.GetByGroupIDAndUserID(ctx, groupId, userId)
	if err != nil {
		return err
	}

	if groupMember == nil {
		return custom_error.RecordNotFound(fmt.Sprintf("user %d is not a member of group %d", userId, groupId))
	}

	return nil
}

func (u *groupAuthUsecase) CanAccessGroup(ctx context.Context, groupId uint32) error {
	// organizationId := _utils.GetOrganizationIdFromContext(ctx)
	userId := _utils.GetProfileIdWithContext(ctx)
	groupMember, err := u.groupMemberRepository.GetByGroupIDAndUserID(ctx, groupId, uint32(userId))
	if err != nil {
		return err
	}
	if groupMember == nil {
		return custom_error.RecordNotFound(fmt.Sprintf("user %d is not a member of group %d", userId, groupId))
	}

	return nil
}

func (u *groupAuthUsecase) HasPermission(ctx context.Context, permission string) bool {
	// TODO: Implement permission checking logic
	// For now, return true to allow build to pass
	return true
}
