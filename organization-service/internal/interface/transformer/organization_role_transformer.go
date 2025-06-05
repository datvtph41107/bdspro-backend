package transformer

import (
	organizationpb "pb/types/organization"

	"github.com/tuandq2112/organization-service/internal/domain/entity"
)

type OrganizationRoleTransformer interface {
	CreateRoleRequestToEntity(req *organizationpb.CreateRoleRequest) *entity.OrganizationRole
	EntityToCreateRoleResponse(entity *entity.OrganizationRole) *organizationpb.CreateRoleResponse
	EntityToGetRolesResponse(entities []*entity.OrganizationRole) *organizationpb.GetRolesResponse
	EntityToUpdateRoleResponse(entity *entity.OrganizationRole) *organizationpb.UpdateRoleResponse
	EntityToDeleteRoleResponse(entity *entity.OrganizationRole) *organizationpb.DeleteRoleResponse
	EntityToGetRoleResponse(entity *entity.OrganizationRole) *organizationpb.GetRolesResponse
	UpdateRoleRequestToEntity(req *organizationpb.UpdateRoleRequest) *entity.OrganizationRole
}

type organizationRoleTransformer struct {
}

func NewOrganizationRoleTransformer() OrganizationRoleTransformer {
	return &organizationRoleTransformer{}
}

func (t *organizationRoleTransformer) CreateRoleRequestToEntity(req *organizationpb.CreateRoleRequest) *entity.OrganizationRole {
	return &entity.OrganizationRole{
		Name: req.Name,
		Key:  req.Key,
		OrganizationId: req.OrganizationId,
	}
}

func (t *organizationRoleTransformer) UpdateRoleRequestToEntity(req *organizationpb.UpdateRoleRequest) *entity.OrganizationRole {
	return &entity.OrganizationRole{
		Id: req.Id,
		Name: req.Name,
		Key: req.Key,
		PermissionIds: req.PermissionIds,
	}
}

func (t *organizationRoleTransformer) EntityToCreateRoleResponse(entity *entity.OrganizationRole) *organizationpb.CreateRoleResponse {
	return &organizationpb.CreateRoleResponse{
		Id: entity.Id,
	}
}

func (t *organizationRoleTransformer) EntityToGetRolesResponse(entities []*entity.OrganizationRole) *organizationpb.GetRolesResponse {
	roles := make([]*organizationpb.Role, 0)
	for _, entity := range entities {
		roles = append(roles, &organizationpb.Role{
			Id: entity.Id,
			Name: entity.Name,
			Key: entity.Key,
		})
	}
	return &organizationpb.GetRolesResponse{
		Roles: roles,
	}
}

func (t *organizationRoleTransformer) EntityToUpdateRoleResponse(entity *entity.OrganizationRole) *organizationpb.UpdateRoleResponse {
	return &organizationpb.UpdateRoleResponse{
		Id: entity.Id,
	}
}

func (t *organizationRoleTransformer) EntityToDeleteRoleResponse(entity *entity.OrganizationRole) *organizationpb.DeleteRoleResponse {
	return &organizationpb.DeleteRoleResponse{
		Id: entity.Id,
	}
}

func (t *organizationRoleTransformer) EntityToGetRoleResponse(entity *entity.OrganizationRole) *organizationpb.GetRolesResponse {
	roles := make([]*organizationpb.Role, 0)
	roles = append(roles, &organizationpb.Role{
		Id: entity.Id,
		Name: entity.Name,
		Key: entity.Key,
	})	
	return &organizationpb.GetRolesResponse{
		Roles: roles,
	}
}
