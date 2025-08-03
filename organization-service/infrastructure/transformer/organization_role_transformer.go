package transformer

import (
	organizationpb "pb/types/organization"

	"organization/internal/domain/entity"
	"organization/internal/enums"
)

type OrganizationRoleTransformer interface {
	CreateRoleRequestToEntity(req *organizationpb.CreateRoleRequest) *entity.OrganizationRole
	EntityToCreateRoleResponse(entity *entity.OrganizationRole) *organizationpb.CreateRoleResponse
	EntityToGetRolesResponse(entities []*entity.OrganizationRole, total uint32) *organizationpb.GetRolesResponse
	EntityToUpdateRoleResponse(entity *entity.OrganizationRole) *organizationpb.UpdateRoleResponse
	EntityToDeleteRoleResponse(entity *entity.OrganizationRole) *organizationpb.DeleteRoleResponse
	EntityToGetRoleResponse(entity *entity.OrganizationRole) *organizationpb.GetRoleResponse
	UpdateRoleRequestToEntity(req *organizationpb.UpdateRoleRequest) *entity.OrganizationRole
	CreateNewRoleWithMembersRequestToEntity(req *organizationpb.CreateNewRoleWithMembersRequest) *entity.OrganizationRole
	EntityToCreateNewRoleWithMembersResponse(entity *entity.OrganizationRole, memberIds []uint32, permissionIds []uint64) *organizationpb.CreateNewRoleWithMembersResponse
	UpdateRoleWithMembersRequestToEntity(req *organizationpb.UpdateRoleWithMembersRequest) *entity.OrganizationRole
	EntityToUpdateRoleWithMembersResponse(entity *entity.OrganizationRole) *organizationpb.UpdateRoleWithMembersResponse
}

type organizationRoleTransformer struct{}

func NewOrganizationRoleTransformer() OrganizationRoleTransformer {
	return &organizationRoleTransformer{}
}

func (t *organizationRoleTransformer) CreateRoleRequestToEntity(req *organizationpb.CreateRoleRequest) *entity.OrganizationRole {
	return &entity.OrganizationRole{
		Name:           req.Name,
		Key:            req.Key,
		OrganizationId: req.OrganizationId,
		IsDefault:      req.IsDefault,
		DomainType:     enums.DomainType(req.DomainType),
		PermissionIds:  req.PermissionIds,
	}
}

func (t *organizationRoleTransformer) UpdateRoleRequestToEntity(req *organizationpb.UpdateRoleRequest) *entity.OrganizationRole {
	return &entity.OrganizationRole{
		Id:            req.Id,
		Name:          req.Name,
		Key:           req.Key,
		PermissionIds: req.PermissionIds,
		IsDefault:     req.IsDefault,
		DomainType:    enums.DomainType(req.DomainType),
	}
}

func (t *organizationRoleTransformer) EntityToCreateRoleResponse(entity *entity.OrganizationRole) *organizationpb.CreateRoleResponse {
	return &organizationpb.CreateRoleResponse{
		Id: entity.Id,
	}
}

func (t *organizationRoleTransformer) EntityToGetRolesResponse(roles []*entity.OrganizationRole, total uint32) *organizationpb.GetRolesResponse {
	rolesResponse := make([]*organizationpb.Role, len(roles))
	for i, role := range roles {
		permissions := make([]*organizationpb.Permission, len(role.Permissions))
		for j, permission := range role.Permissions {
			permissions[j] = &organizationpb.Permission{
				Id:   permission.Id,
				Name: permission.Name,
				Key:  permission.Key,
			}
		}
		rolesResponse[i] = &organizationpb.Role{
			Id:          role.Id,
			Name:        role.Name,
			Key:         role.Key,
			Permissions: permissions,
		}
	}
	return &organizationpb.GetRolesResponse{
		Data:  rolesResponse,
		Total: total,
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

func (t *organizationRoleTransformer) EntityToGetRoleResponse(entity *entity.OrganizationRole) *organizationpb.GetRoleResponse {
	permissions := make([]*organizationpb.Permission, len(entity.Permissions))
	for i, permission := range entity.Permissions {
		permissions[i] = &organizationpb.Permission{
			Id:   permission.Id,
			Name: permission.Name,
			Key:  permission.Key,
		}
	}

	return &organizationpb.GetRoleResponse{
		Role: &organizationpb.Role{
			Id:          entity.Id,
			Name:        entity.Name,
			Key:         entity.Key,
			Permissions: permissions,
		},
	}
}

func (t *organizationRoleTransformer) CreateNewRoleWithMembersRequestToEntity(req *organizationpb.CreateNewRoleWithMembersRequest) *entity.OrganizationRole {
	return &entity.OrganizationRole{
		Name: req.Name,
		Key:  req.Key,
	}
}

func (t *organizationRoleTransformer) EntityToCreateNewRoleWithMembersResponse(entity *entity.OrganizationRole, memberIds []uint32, permissionIds []uint64) *organizationpb.CreateNewRoleWithMembersResponse {
	return &organizationpb.CreateNewRoleWithMembersResponse{
		Id:            entity.Id,
		Name:          entity.Name,
		Key:           entity.Key,
		MemberIds:     memberIds,
		PermissionIds: permissionIds,
	}
}

func (t *organizationRoleTransformer) UpdateRoleWithMembersRequestToEntity(req *organizationpb.UpdateRoleWithMembersRequest) *entity.OrganizationRole {
	return &entity.OrganizationRole{
		Id:            req.Id,
		Name:          req.Name,
		Key:           req.Key,
		PermissionIds: req.PermissionIds,
		IsDefault:     req.IsDefault,
		DomainType:    enums.DomainType(req.DomainType),
	}
}

func (t *organizationRoleTransformer) EntityToUpdateRoleWithMembersResponse(entity *entity.OrganizationRole) *organizationpb.UpdateRoleWithMembersResponse {
	return &organizationpb.UpdateRoleWithMembersResponse{
		Id: entity.Id,
	}
}
