package dto

import (
	_dto "common/domain/dto"
	_utils "common/utils"
	"user/internal/domain/access"
)

// RoleGroupRequest DTO cho request tạo/sửa role group
type RoleGroupRequest struct {
	ID          *uint64 `json:"id,omitempty" validate:"omitempty"`
	GroupName   string  `json:"groupName" validate:"required,min=1,max=100"`
	Description string  `json:"description" validate:"max=500"`
	Key         uint32  `json:"key" validate:"required,min=1,max=50"`
}

// RoleGroupResponse DTO cho response role group
type RoleGroupResponse struct {
	ID          uint64 `json:"id"`
	GroupName   string `json:"groupName"`
	Description string `json:"description"`
	Key         uint32 `json:"key"`
	CreatedAt   string `json:"createdAt"`
	UpdatedAt   string `json:"updatedAt"`
}

// RoleGroupListResponse DTO cho response danh sách role group
type RoleGroupListResponse struct {
	_dto.Pagable
	Data []*RoleGroupResponse `json:"data"`
}

// AddRoleToGroupRequest DTO cho request thêm role vào group
type AddRoleToGroupRequest struct {
	GroupID uint64   `json:"groupId" validate:"required"`
	RoleIDs []uint64 `json:"roleIds" validate:"required,min=1"`
}

// RemoveRoleFromGroupRequest DTO cho request xóa role khỏi group
type RemoveRoleFromGroupRequest struct {
	GroupID uint64   `json:"groupId" validate:"required"`
	RoleIDs []uint64 `json:"roleIds" validate:"required,min=1"`
}

// AddPermissionToGroupRequest DTO cho request thêm permission vào group
type AddPermissionToGroupRequest struct {
	GroupID       uint64   `json:"groupId" validate:"required"`
	PermissionIDs []uint64 `json:"permissionIds" validate:"required,min=1"`
}

// RemovePermissionFromGroupRequest DTO cho request xóa permission khỏi group
type RemovePermissionFromGroupRequest struct {
	GroupID       uint64   `json:"groupId" validate:"required"`
	PermissionIDs []uint64 `json:"permissionIds" validate:"required,min=1"`
}

// Convert domain to DTO
func RoleGroupToResponse(roleGroup *access.RoleGroup) *RoleGroupResponse {
	if roleGroup == nil {
		return nil
	}

	response := &RoleGroupResponse{
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
