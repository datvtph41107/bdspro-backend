package usecase

import (
	"context"
	"fmt"

	"github.com/tuandq2112/organization-service/internal/custom_error"
	"github.com/tuandq2112/organization-service/internal/domain/repository"
)

type GroupAuthUsecase interface {
	IsLeaderGroup(ctx context.Context, groupId uint32, userId uint32) error
	IsMemberGroup(ctx context.Context, groupId uint32, userId uint32) error
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

	if group.CreatedBy != userId {
		return custom_error.Forbidden(fmt.Sprintf("user %d is not leader of group %d", userId, groupId))
	}

	return nil
}

func (u *groupAuthUsecase) IsMemberGroup(ctx context.Context, groupId uint32, userId uint32) error {
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
