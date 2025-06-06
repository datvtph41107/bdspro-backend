package usecase

import (
	"context"

	"github.com/tuandq2112/organization-service/internal/custom_error"
	"github.com/tuandq2112/organization-service/internal/domain/entity"
	"github.com/tuandq2112/organization-service/internal/domain/repository"
)

type GroupSettingUsecase interface {
	CreateGroupSetting(ctx context.Context, setting *entity.GroupSetting) (*entity.GroupSetting, error)
	UpdateGroupSetting(ctx context.Context, setting *entity.GroupSetting) (*entity.GroupSetting, error)
	DeleteGroupSetting(ctx context.Context, id uint32) error
	GetGroupSettingByID(ctx context.Context, id uint32) (*entity.GroupSetting, error)
	GetGroupSettingsByGroupID(ctx context.Context, groupID uint32) ([]*entity.GroupSetting, error)
}

type groupSettingUsecase struct {
	groupSettingRepository repository.GroupSettingRepository
	logWorker             *LogWorker
}

func NewGroupSettingUsecase(groupSettingRepository repository.GroupSettingRepository, logWorker *LogWorker) GroupSettingUsecase {
	return &groupSettingUsecase{
		groupSettingRepository: groupSettingRepository,
		logWorker:             logWorker,
	}
}

func (u *groupSettingUsecase) CreateGroupSetting(ctx context.Context, setting *entity.GroupSetting) (*entity.GroupSetting, error) {
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
	existedSetting, err := u.groupSettingRepository.GetByID(ctx, setting.Id)
	if err != nil {
		return nil, err
	}
	
	if existedSetting == nil {
		return nil, custom_error.RecordNotFound("group setting not found")
	}
	
	if setting.CreatedBy != existedSetting.CreatedBy {
		return nil, custom_error.Forbidden("you are not allowed to update this group setting")
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
	return u.groupSettingRepository.GetByID(ctx, id)
}

func (u *groupSettingUsecase) GetGroupSettingsByGroupID(ctx context.Context, groupID uint32) ([]*entity.GroupSetting, error) {
	return u.groupSettingRepository.GetByGroupID(ctx, groupID)
}