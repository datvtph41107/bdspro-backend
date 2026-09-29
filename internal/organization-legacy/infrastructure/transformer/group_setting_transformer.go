package transformer

import (
	organizationpb "pb/types/organization"

	"organization/internal/domain/entity"
)

type GroupSettingTransformer interface {
	CreateGroupSettingRequestToEntity(request *organizationpb.CreateGroupSettingRequest) *entity.GroupSetting
	UpdateGroupSettingRequestToEntity(request *organizationpb.UpdateGroupSettingRequest) *entity.GroupSetting
	EntityToCreateGroupSettingResponse(setting *entity.GroupSetting) *organizationpb.CreateGroupSettingResponse
	EntityToUpdateGroupSettingResponse(setting *entity.GroupSetting) *organizationpb.UpdateGroupSettingResponse
	EntityToGetGroupSettingResponse(setting *entity.GroupSetting) *organizationpb.GetGroupSettingResponse
	EntityToGetGroupSettingsResponse(settings []*entity.GroupSetting, total uint32) *organizationpb.GetGroupSettingsResponse
}

type groupSettingTransformer struct{}

func NewGroupSettingTransformer() GroupSettingTransformer {
	return &groupSettingTransformer{}
}

func (t *groupSettingTransformer) CreateGroupSettingRequestToEntity(request *organizationpb.CreateGroupSettingRequest) *entity.GroupSetting {
	return &entity.GroupSetting{
		GroupId:     request.GroupId,
		ConfigKey:   request.ConfigKey,
		ConfigValue: request.ConfigValue,
	}
}

func (t *groupSettingTransformer) UpdateGroupSettingRequestToEntity(request *organizationpb.UpdateGroupSettingRequest) *entity.GroupSetting {
	return &entity.GroupSetting{
		Id:          request.Id,
		ConfigKey:   request.ConfigKey,
		ConfigValue: request.ConfigValue,
	}
}

func (t *groupSettingTransformer) EntityToCreateGroupSettingResponse(setting *entity.GroupSetting) *organizationpb.CreateGroupSettingResponse {
	return &organizationpb.CreateGroupSettingResponse{
		Id: setting.Id,
	}
}

func (t *groupSettingTransformer) EntityToUpdateGroupSettingResponse(setting *entity.GroupSetting) *organizationpb.UpdateGroupSettingResponse {
	return &organizationpb.UpdateGroupSettingResponse{
		Id: setting.Id,
	}
}

func (t *groupSettingTransformer) EntityToGetGroupSettingResponse(setting *entity.GroupSetting) *organizationpb.GetGroupSettingResponse {
	return &organizationpb.GetGroupSettingResponse{
		Setting: &organizationpb.GroupSetting{
			Id:          setting.Id,
			ConfigKey:   setting.ConfigKey,
			ConfigValue: setting.ConfigValue,
		},
	}
}

func (t *groupSettingTransformer) EntityToGetGroupSettingsResponse(settings []*entity.GroupSetting, total uint32) *organizationpb.GetGroupSettingsResponse {
	settingsResponse := make([]*organizationpb.GroupSetting, len(settings))
	for i, setting := range settings {
		settingsResponse[i] = &organizationpb.GroupSetting{
			Id:          setting.Id,
			ConfigKey:   setting.ConfigKey,
			ConfigValue: setting.ConfigValue,
		}
	}
	return &organizationpb.GetGroupSettingsResponse{
		Settings: settingsResponse,
		Total:    total,
	}
}
