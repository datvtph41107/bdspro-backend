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
	logActivityUsecase GroupLogActivityUsecase
}

func NewGroupMemberUsecase(groupMemberRepository repository.GroupMemberRepository, logActivityUsecase GroupLogActivityUsecase) GroupMemberUsecase {
	return &groupMemberUsecase{
		groupMemberRepository: groupMemberRepository,
		logActivityUsecase: logActivityUsecase,
	}
}

func (u *groupMemberUsecase) CreateGroupMember(ctx context.Context, member *entity.GroupMember) (*entity.GroupMember, error) {
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

	if member == nil {
		return custom_error.RecordNotFound("group member not found")
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
	return u.groupMemberRepository.GetByID(ctx, id)
}

func (u *groupMemberUsecase) GetGroupMembersByGroupID(ctx context.Context, groupID uint32) ([]*entity.GroupMember, error) {
	return u.groupMemberRepository.GetByGroupID(ctx, groupID)
}
