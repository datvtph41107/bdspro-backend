package transformer

import (
	organizationpb "pb/types/organization"

	"github.com/tuandq2112/organization-service/internal/domain/entity"
)

type OrganizationPermissionTransformer interface {
	CreatePermissionRequestToEntity(req *organizationpb.CreatePermissionRequest) *entity.OrganizationPermission
	EntityToCreatePermissionResponse(entity *entity.OrganizationPermission) *organizationpb.CreatePermissionResponse
	UpdatePermissionRequestToEntity(req *organizationpb.UpdatePermissionRequest) *entity.OrganizationPermission
	EntityToUpdatePermissionResponse(entity *entity.OrganizationPermission) *organizationpb.UpdatePermissionResponse
	EntityToDeletePermissionResponse(entity *entity.OrganizationPermission) *organizationpb.DeletePermissionResponse
	EntityToGetPermissionsResponse(entities []*entity.OrganizationPermission) *organizationpb.GetPermissionsResponse
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
		OrganizationId: req.OrganizationId,
	}
}

func (t *organizationPermissionTransformer) UpdatePermissionRequestToEntity(req *organizationpb.UpdatePermissionRequest) *entity.OrganizationPermission {
	return &entity.OrganizationPermission{
		Id: req.Id,
		Name: req.Name,
		Key: req.Key,
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

func (t *organizationPermissionTransformer) EntityToGetPermissionsResponse(entities []*entity.OrganizationPermission) *organizationpb.GetPermissionsResponse {
	permissions := make([]*organizationpb.Permission, 0)	
	for _, permission := range entities {
		permissions = append(permissions, &organizationpb.Permission{
			Id: permission.Id,
			Name: permission.Name,
			Key: permission.Key,
		})
	}
	return &organizationpb.GetPermissionsResponse{
		Permissions: permissions,
	}
}

