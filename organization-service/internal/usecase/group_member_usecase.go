package usecase

import (
	"context"

	"github.com/tuandq2112/organization-service/internal/custom_error"
	"github.com/tuandq2112/organization-service/internal/domain/entity"
	"github.com/tuandq2112/organization-service/internal/domain/repository"
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
}

func NewGroupMemberUsecase(groupMemberRepository repository.GroupMemberRepository) GroupMemberUsecase {
	return &groupMemberUsecase{groupMemberRepository: groupMemberRepository}
}

func (u *groupMemberUsecase) CreateGroupMember(ctx context.Context, member *entity.GroupMember) (*entity.GroupMember, error) {
	return u.groupMemberRepository.Create(ctx, member)
}

func (u *groupMemberUsecase) UpdateGroupMember(ctx context.Context, member *entity.GroupMember) (*entity.GroupMember, error) {
	existedMember, err := u.groupMemberRepository.GetByID(ctx, member.Id)
	if err != nil {
		return nil, err
	}

	if existedMember == nil {
		return nil, custom_error.RecordNotFound("group member not found")
	}

	if member.CreatedBy != existedMember.CreatedBy {
		return nil, custom_error.Forbidden("you are not allowed to update this group member")
	}

	return u.groupMemberRepository.Update(ctx, member)
}

func (u *groupMemberUsecase) DeleteGroupMember(ctx context.Context, id uint32) error {
	return u.groupMemberRepository.Delete(ctx, id)
}

func (u *groupMemberUsecase) GetGroupMemberByID(ctx context.Context, id uint32) (*entity.GroupMember, error) {
	return u.groupMemberRepository.GetByID(ctx, id)
}

func (u *groupMemberUsecase) GetGroupMembersByGroupID(ctx context.Context, groupID uint32) ([]*entity.GroupMember, error) {
	return u.groupMemberRepository.GetByGroupID(ctx, groupID)
}
