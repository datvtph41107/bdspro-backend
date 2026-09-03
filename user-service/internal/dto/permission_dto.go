package dto

import (
	_utils "common/utils"
	"user/internal/domain/access"
)

type PermissionRequest struct {
	Name           string  `json:"name" validate:"required"`
	Key            string  `json:"key" validate:"required"`
	Description    string  `json:"description"`
	PermissionType string  `json:"permissionType"`
	Module         string  `json:"module" validate:"required"`
	ParentID       *uint64 `json:"parentId"`
}

type PermissionResponse struct {
	ID             uint64                `json:"id"`
	Name           string                `json:"name"`
	Key            string                `json:"key"`
	Description    string                `json:"description"`
	PermissionType string                `json:"permissionType"`
	Module         string                `json:"module"`
	ParentID       *uint64               `json:"parentId"`
	Children       []*PermissionResponse `json:"children"`
	CreatedAt      string                `json:"createdAt"`
	UpdatedAt      string                `json:"updatedAt"`
}

// MapToDomain chuyển DTO sang domain
func (dto *PermissionRequest) MapToDomain() *access.Permission {
	return &access.Permission{
		Name:           dto.Name,
		Key:            dto.Key,
		Description:    dto.Description,
		PermissionType: dto.PermissionType,
		Module:         dto.Module,
		ParentID:       dto.ParentID,
	}
}

// MapFromDomain chuyển domain sang DTO response
func MapPermissionFromDomain(permission *access.Permission) *PermissionResponse {
	if permission == nil {
		return nil
	}

	response := &PermissionResponse{
		ID:             permission.ID,
		Name:           permission.Name,
		Key:            permission.Key,
		Description:    permission.Description,
		PermissionType: permission.PermissionType,
		Module:         permission.Module,
		ParentID:       permission.ParentID,
	}

	if permission.CreatedAt != nil {
		response.CreatedAt = _utils.FormatTimeToString(permission.CreatedAt)
	}
	if permission.UpdatedAt != nil {
		response.UpdatedAt = _utils.FormatTimeToString(permission.UpdatedAt)
	}

	if permission.Children != nil {
		response.Children = make([]*PermissionResponse, len(permission.Children))
		for i, child := range permission.Children {
			response.Children[i] = MapPermissionFromDomain(child)
		}
	}

	return response
}
