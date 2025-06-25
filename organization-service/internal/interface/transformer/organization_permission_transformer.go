package transformer

import (
	organizationpb "pb/types/organization"

	"organization/internal/domain/entity"
)

type OrganizationPermissionTransformer interface {
	CreatePermissionRequestToEntity(req *organizationpb.CreatePermissionRequest) *entity.OrganizationPermission
	EntityToCreatePermissionResponse(entity *entity.OrganizationPermission) *organizationpb.CreatePermissionResponse
	UpdatePermissionRequestToEntity(req *organizationpb.UpdatePermissionRequest) *entity.OrganizationPermission
	EntityToUpdatePermissionResponse(entity *entity.OrganizationPermission) *organizationpb.UpdatePermissionResponse
	EntityToDeletePermissionResponse(entity *entity.OrganizationPermission) *organizationpb.DeletePermissionResponse
	EntityToGetPermissionsResponse(entities []*entity.OrganizationPermission, total uint32) *organizationpb.GetPermissionsResponse
}
type organizationPermissionTransformer struct {
}

func NewOrganizationPermissionTransformer() OrganizationPermissionTransformer {
	return &organizationPermissionTransformer{}
}

func (t *organizationPermissionTransformer) CreatePermissionRequestToEntity(req *organizationpb.CreatePermissionRequest) *entity.OrganizationPermission {
	return &entity.OrganizationPermission{
		Name: req.Name,
		Key:  req.Key,
		// OrganizationId: req.OrganizationId,
	}
}

func (t *organizationPermissionTransformer) UpdatePermissionRequestToEntity(req *organizationpb.UpdatePermissionRequest) *entity.OrganizationPermission {
	return &entity.OrganizationPermission{
		Id:   req.Id,
		Name: req.Name,
		Key:  req.Key,
	}
}

func (t *organizationPermissionTransformer) EntityToCreatePermissionResponse(entity *entity.OrganizationPermission) *organizationpb.CreatePermissionResponse {
	return &organizationpb.CreatePermissionResponse{
		Id: entity.Id,
	}
}

func (t *organizationPermissionTransformer) EntityToUpdatePermissionResponse(entity *entity.OrganizationPermission) *organizationpb.UpdatePermissionResponse {
	return &organizationpb.UpdatePermissionResponse{
		Id: entity.Id,
	}
}

func (t *organizationPermissionTransformer) EntityToDeletePermissionResponse(entity *entity.OrganizationPermission) *organizationpb.DeletePermissionResponse {
	return &organizationpb.DeletePermissionResponse{
		Id: entity.Id,
	}
}

func (t *organizationPermissionTransformer) EntityToPb(permission *entity.OrganizationPermission) *organizationpb.Permission {
	permissionResponse := &organizationpb.Permission{
		Id:   permission.Id,
		Name: permission.Name,
		Key:  permission.Key,
	}
	if len(permission.Childrens) > 0 {
		childrens := make([]*organizationpb.Permission, len(permission.Childrens))
		for i, child := range permission.Childrens {
			childrens[i] = t.EntityToPb(child)
		}
		permissionResponse.Childrens = childrens
	}
	return permissionResponse
}
func (t *organizationPermissionTransformer) EntityToGetPermissionsResponse(permissions []*entity.OrganizationPermission, total uint32) *organizationpb.GetPermissionsResponse {
	permissionsResponse := make([]*organizationpb.Permission, len(permissions))
	for i, permission := range permissions {
		permissionsResponse[i] = t.EntityToPb(permission)
	}
	return &organizationpb.GetPermissionsResponse{
		Permissions: permissionsResponse,
		Total:       total,
	}
}
