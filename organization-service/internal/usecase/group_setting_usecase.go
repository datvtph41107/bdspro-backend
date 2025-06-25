package usecase

import (
	"context"

	"organization/env"
	"organization/internal/custom_error"
	"organization/internal/domain/entity"
	"organization/internal/domain/repository"
	"organization/pkg/utils"
)

type GroupSettingUsecase interface {
	CreateGroupSetting(ctx context.Context, setting *entity.GroupSetting) (*entity.GroupSetting, error)
	UpdateGroupSetting(ctx context.Context, setting *entity.GroupSetting) (*entity.GroupSetting, error)
	DeleteGroupSetting(ctx context.Context, id uint32) error
	GetGroupSettingByID(ctx context.Context, id uint32) (*entity.GroupSetting, error)
	GetGroupSettingsByGroupID(ctx context.Context, groupID uint32, page, size int) ([]*entity.GroupSetting, uint32, error)
}

type groupSettingUsecase struct {
	groupSettingRepository repository.GroupSettingRepository
	logWorker              *LogWorker
	groupAuthUsecase       GroupAuthUsecase
}

func NewGroupSettingUsecase(groupSettingRepository repository.GroupSettingRepository, logWorker *LogWorker, groupAuthUsecase GroupAuthUsecase) GroupSettingUsecase {
	return &groupSettingUsecase{
		groupSettingRepository: groupSettingRepository,
		logWorker:              logWorker,
		groupAuthUsecase:       groupAuthUsecase,
	}
}

func (u *groupSettingUsecase) CreateGroupSetting(ctx context.Context, setting *entity.GroupSetting) (*entity.GroupSetting, error) {
	currentUserId := utils.GetUserID(ctx, env.USER_CONTEXT)
	err := u.groupAuthUsecase.IsLeaderGroup(ctx, setting.GroupId, currentUserId)
	if err != nil {
		return nil, custom_error.Forbidden("you are not allowed to create group setting")
	}
	setting.CreatedBy = currentUserId
	setting.UpdatedBy = currentUserId

	createdSetting, err := u.groupSettingRepository.Create(ctx, setting)
	if err != nil {
		return nil, err
	}

	u.logWorker.Push(LogEvent{
		GroupId: createdSetting.GroupId,
		ActorId: createdSetting.CreatedBy,
		LogType: "CREATE_SETTING",
		LogData: "Created new group setting",
	})

	return createdSetting, nil
}

func (u *groupSettingUsecase) UpdateGroupSetting(ctx context.Context, setting *entity.GroupSetting) (*entity.GroupSetting, error) {
	currentUserId := utils.GetUserID(ctx, env.USER_CONTEXT)

	setting.UpdatedBy = currentUserId
	existedSetting, err := u.groupSettingRepository.GetByID(ctx, setting.Id)
	if err != nil {
		return nil, err
	}

	if existedSetting == nil {
		return nil, custom_error.RecordNotFound("group setting not found")
	}

	err = u.groupAuthUsecase.IsLeaderGroup(ctx, existedSetting.GroupId, currentUserId)
	if err != nil {
		return nil, custom_error.Forbidden("you are not allowed to update group setting")
	}
	updatedSetting, err := u.groupSettingRepository.Update(ctx, setting)
	if err != nil {
		return nil, err
	}

	u.logWorker.Push(LogEvent{
		GroupId: updatedSetting.GroupId,
		ActorId: updatedSetting.CreatedBy,
		LogType: "UPDATE_SETTING",
		LogData: "Updated group setting",
	})

	return updatedSetting, nil
}

func (u *groupSettingUsecase) DeleteGroupSetting(ctx context.Context, id uint32) error {

	setting, err := u.groupSettingRepository.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if setting == nil {
		return custom_error.RecordNotFound("group setting not found")
	}
	currentUserId := utils.GetUserID(ctx, env.USER_CONTEXT)
	err = u.groupAuthUsecase.IsLeaderGroup(ctx, setting.GroupId, currentUserId)
	if err != nil {
		return custom_error.Forbidden("you are not allowed to delete group setting")
	}
	err = u.groupSettingRepository.Delete(ctx, id)
	if err != nil {
		return err
	}

	u.logWorker.Push(LogEvent{
		GroupId: setting.GroupId,
		ActorId: setting.CreatedBy,
		LogType: "DELETE_SETTING",
		LogData: "Deleted group setting",
	})

	return nil
}

func (u *groupSettingUsecase) GetGroupSettingByID(ctx context.Context, id uint32) (*entity.GroupSetting, error) {
	currentUserId := utils.GetUserID(ctx, env.USER_CONTEXT)
	setting, err := u.groupSettingRepository.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if setting == nil {
		return nil, custom_error.RecordNotFound("group setting not found")
	}

	err = u.groupAuthUsecase.IsMemberGroup(ctx, setting.GroupId, currentUserId)
	if err != nil {
		return nil, custom_error.Forbidden("you are not allowed to get group setting")
	}
	if err != nil {
		return nil, err
	}

	return setting, nil
}

func (u *groupSettingUsecase) GetGroupSettingsByGroupID(ctx context.Context, groupID uint32, page, size int) ([]*entity.GroupSetting, uint32, error) {
	currentUserId := utils.GetUserID(ctx, env.USER_CONTEXT)
	err := u.groupAuthUsecase.IsMemberGroup(ctx, groupID, currentUserId)
	if err != nil {
		return nil, 0, custom_error.Forbidden("you are not allowed to get group settings")
	}
	return u.groupSettingRepository.GetByGroupID(ctx, groupID, page, size)
}
