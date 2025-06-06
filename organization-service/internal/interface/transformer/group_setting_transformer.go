package transformer

import (
	grouppb "pb/types/group"

	"github.com/tuandq2112/organization-service/internal/domain/entity"
)

type GroupSettingTransformer interface {
	CreateGroupSettingRequestToEntity(request *grouppb.CreateGroupSettingRequest) *entity.GroupSetting
	UpdateGroupSettingRequestToEntity(request *grouppb.UpdateGroupSettingRequest) *entity.GroupSetting
	EntityToCreateGroupSettingResponse(setting *entity.GroupSetting) *grouppb.CreateGroupSettingResponse
	EntityToUpdateGroupSettingResponse(setting *entity.GroupSetting) *grouppb.UpdateGroupSettingResponse
	EntityToGetGroupSettingResponse(setting *entity.GroupSetting) *grouppb.GetGroupSettingResponse
	EntityToGetGroupSettingsResponse(settings []*entity.GroupSetting, total uint32) *grouppb.GetGroupSettingsResponse
}

type groupSettingTransformer struct{}

func NewGroupSettingTransformer() GroupSettingTransformer {
	return &groupSettingTransformer{}
}

func (t *groupSettingTransformer) CreateGroupSettingRequestToEntity(request *grouppb.CreateGroupSettingRequest) *entity.GroupSetting {
	return &entity.GroupSetting{
		GroupId:     request.GroupId,
		ConfigKey:   request.ConfigKey,
		ConfigValue: request.ConfigValue,
	}
}

func (t *groupSettingTransformer) UpdateGroupSettingRequestToEntity(request *grouppb.UpdateGroupSettingRequest) *entity.GroupSetting {
	return &entity.GroupSetting{
		Id:          request.Id,
		ConfigKey:   request.ConfigKey,
		ConfigValue: request.ConfigValue,
	}
}

func (t *groupSettingTransformer) EntityToCreateGroupSettingResponse(setting *entity.GroupSetting) *grouppb.CreateGroupSettingResponse {
	return &grouppb.CreateGroupSettingResponse{
		Id: setting.Id,
	}
}

func (t *groupSettingTransformer) EntityToUpdateGroupSettingResponse(setting *entity.GroupSetting) *grouppb.UpdateGroupSettingResponse {
	return &grouppb.UpdateGroupSettingResponse{
		Id: setting.Id,
	}
}

func (t *groupSettingTransformer) EntityToGetGroupSettingResponse(setting *entity.GroupSetting) *grouppb.GetGroupSettingResponse {
	return &grouppb.GetGroupSettingResponse{
		Setting: &grouppb.GroupSetting{
			Id:          setting.Id,
			ConfigKey:   setting.ConfigKey,
			ConfigValue: setting.ConfigValue,
		},
	}
}

func (t *groupSettingTransformer) EntityToGetGroupSettingsResponse(settings []*entity.GroupSetting, total uint32) *grouppb.GetGroupSettingsResponse {
	settingsResponse := make([]*grouppb.GroupSetting, len(settings))
	for i, setting := range settings {
		settingsResponse[i] = &grouppb.GroupSetting{
			Id:          setting.Id,
			ConfigKey:   setting.ConfigKey,
			ConfigValue: setting.ConfigValue,
		}
	}
	return &grouppb.GetGroupSettingsResponse{
		Settings: settingsResponse,
		Total:    total,
	}
} 