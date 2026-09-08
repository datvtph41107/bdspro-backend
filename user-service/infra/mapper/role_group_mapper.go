package mapper

import (
	_utils "common/utils"
	"fmt"
	authpb "pb/types/auth"
	"user/internal/domain/access"
	"user/internal/dto"
	"user/internal/enums"
)

type RoleGroupMapper struct{}

func NewRoleGroupMapper() *RoleGroupMapper {
	return &RoleGroupMapper{}
}

// MapToDomain chuyển DTO sang domain
func (m *RoleGroupMapper) MapToDomain(req *authpb.RoleGroupRequest) *access.RoleGroup {
	roleGroup := &access.RoleGroup{
		GroupName:   req.GroupName,
		Description: req.Description,
		Key:         enums.GroupRoleKeyEnum(req.Key),
		// Code is a persisted unique identity used by the legacy role schema.
		// The public contract exposes Key as that identity, so derive Code from
		// Key instead of writing the same empty string for every new group.
		Code: fmt.Sprintf("RG_%d", req.Key),
	}

	if req.Id != 0 {
		roleGroup.ID = req.Id
	}

	return roleGroup
}

// MapFromDomain chuyển domain sang DTO response
func (m *RoleGroupMapper) MapFromDomain(roleGroup *access.RoleGroup) *dto.RoleGroupResponse {
	if roleGroup == nil {
		return nil
	}

	response := &dto.RoleGroupResponse{
		ID:          roleGroup.ID,
		GroupName:   roleGroup.GroupName,
		Description: roleGroup.Description,
		Key:         uint32(roleGroup.Key),
	}

	if roleGroup.CreatedAt != nil {
		response.CreatedAt = _utils.FormatTimeToString(roleGroup.CreatedAt)
	}
	if roleGroup.UpdatedAt != nil {
		response.UpdatedAt = _utils.FormatTimeToString(roleGroup.UpdatedAt)
	}

	return response
}

// MapListFromDomain chuyển danh sách domain sang DTO response
func (m *RoleGroupMapper) MapListFromDomain(roleGroups []*access.RoleGroup) []*dto.RoleGroupResponse {
	if roleGroups == nil {
		return nil
	}

	response := make([]*dto.RoleGroupResponse, len(roleGroups))
	for i, roleGroup := range roleGroups {
		response[i] = m.MapFromDomain(roleGroup)
	}

	return response
}

func (m *RoleGroupMapper) MapToPb(roleGroup *access.RoleGroup) *authpb.RoleGroup {
	return &authpb.RoleGroup{
		Id:          roleGroup.ID,
		GroupName:   roleGroup.GroupName,
		Description: roleGroup.Description,
		Key:         uint32(roleGroup.Key),
		CreatedAt:   _utils.FormatTimeToString(roleGroup.CreatedAt),
		UpdatedAt:   _utils.FormatTimeToString(roleGroup.UpdatedAt),
	}
}

func (m *RoleGroupMapper) MapToPbList(roleGroups []access.RoleGroup) []*authpb.RoleGroup {
	if roleGroups == nil {
		return nil
	}

	response := make([]*authpb.RoleGroup, len(roleGroups))
	for i, roleGroup := range roleGroups {
		response[i] = m.MapToPb(&roleGroup)
	}

	return response
}
