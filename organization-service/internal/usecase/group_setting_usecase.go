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
}

func NewGroupSettingUsecase(groupSettingRepository repository.GroupSettingRepository) GroupSettingUsecase {
	return &groupSettingUsecase{groupSettingRepository: groupSettingRepository}
}

func (u *groupSettingUsecase) CreateGroupSetting(ctx context.Context, setting *entity.GroupSetting) (*entity.GroupSetting, error) {
	return u.groupSettingRepository.Create(ctx, setting)
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
	
	return u.groupSettingRepository.Update(ctx, setting)
}

func (u *groupSettingUsecase) DeleteGroupSetting(ctx context.Context, id uint32) error {
	return u.groupSettingRepository.Delete(ctx, id)
}

func (u *groupSettingUsecase) GetGroupSettingByID(ctx context.Context, id uint32) (*entity.GroupSetting, error) {
	return u.groupSettingRepository.GetByID(ctx, id)
}

func (u *groupSettingUsecase) GetGroupSettingsByGroupID(ctx context.Context, groupID uint32) ([]*entity.GroupSetting, error) {
	return u.groupSettingRepository.GetByGroupID(ctx, groupID)
}