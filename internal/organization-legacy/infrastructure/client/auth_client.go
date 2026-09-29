package client

import (
	"context"
	"organization/internal/custom_error"
	"organization/internal/dto"
	"pb/clients"
	authpb "pb/types/auth"
	organizationpb "pb/types/organization"
	sharepb "pb/types/shared"
)

type AuthClient struct {
	*clients.AuthGrpcClient
	userClient *clients.UserGrpcClient
}

func NewAuthClient(authClient *clients.AuthGrpcClient, userClient *clients.UserGrpcClient) *AuthClient {
	return &AuthClient{AuthGrpcClient: authClient, userClient: userClient}
}

func (c *AuthClient) MapRoleToGroups(ctx context.Context, groups []*dto.GroupWithDetails) {
	roleIds := make([]uint64, len(groups))
	for i, r := range groups {
		if r.RoleId != nil {
			roleIds[i] = *r.RoleId
		}
	}

	roles, err := c.userClient.GetRolesByIds(ctx, roleIds)
	if err != nil {
		return
	}

	roleMap := make(map[uint64]string)
	for _, r := range roles.Data {
		roleMap[r.Id] = r.RoleName
	}

	for i, r := range groups {
		if r.RoleId != nil {
			groups[i].UserRole = roleMap[*r.RoleId]
		} else {
			groups[i].UserRole = "Thành viên"
		}
	}
}

func toAuthRole(r *sharepb.Role) *authpb.Role {
	if r == nil {
		return nil
	}
	role := &authpb.Role{
		Id:              r.Id,
		RoleName:        r.RoleName,
		RoleKey:         r.RoleKey,
		RoleDescription: r.RoleDescription,
		Key:             r.Key,
		PermissionIds:   r.PermissionIds,
		OrganizationId:  r.OrganizationId,
		IsDefault:       r.IsDefault,
		DomainType:      r.DomainType,
		AllowAssign:     r.AllowAssign,
		CreatedAt:       r.CreatedAt,
		UpdatedAt:       r.UpdatedAt,
	}
	if r.Color != nil {
		role.Color = r.Color
		role.ColorId = &r.Color.Id
	}
	return role
}

func (c *AuthClient) GetRoleBySystemKey(ctx context.Context, systemKey uint32) (*authpb.Role, error) {
	response, err := c.userClient.GetRolesByRoleKeys(ctx, []uint32{systemKey})
	if err != nil {
		return nil, err
	}
	if len(response.Data) == 0 {
		return nil, custom_error.RecordNotFound("Role not found")
	}
	return toAuthRole(response.Data[0]), nil
}

func (c *AuthClient) MapRoleToDealInvitationPb(ctx context.Context, invitations []*organizationpb.DealMember) {
	roleIds := make([]uint32, len(invitations))
	for i, r := range invitations {
		if r.RoleKey != 0 {
			roleIds[i] = r.RoleKey
		}
	}

	roles, err := c.userClient.GetRolesByRoleKeys(ctx, roleIds)
	if err != nil {
		return
	}

	roleMap := make(map[uint32]*authpb.Role)
	for _, r := range roles.Data {
		roleMap[r.RoleKey] = toAuthRole(r)
	}

	for i, r := range invitations {
		if r.RoleKey != 0 {
			invitations[i].Role = roleMap[r.RoleKey]
		}
	}
}

func (c *AuthClient) MapRoleToOrganizationMemberPb(ctx context.Context, orgMembers []*organizationpb.OrganizationMember) error {
	roleIds := make([]uint32, len(orgMembers))
	for i, r := range orgMembers {
		if r.GetRoleKey() != 0 {
			roleIds[i] = r.RoleKey
		}
	}

	roles, err := c.userClient.GetRolesByRoleKeys(ctx, roleIds)
	if err != nil {
		return err
	}

	roleMap := make(map[uint32]*authpb.Role)
	for _, r := range roles.Data {
		roleMap[r.RoleKey] = toAuthRole(r)
	}

	for i, r := range orgMembers {
		if r.GetRoleKey() != 0 {
			role := roleMap[r.RoleKey]
			if role == nil {
				continue
			}
			orgMembers[i].RoleName = role.RoleName
			orgMembers[i].RoleKey = role.RoleKey
			orgMembers[i].RoleId = role.Id
		}
	}
	return nil
}

func (c *AuthClient) GetAdminRoleIdByGroupKey(ctx context.Context, groupKey int) (*authpb.Role, error) {
	return c.GetRoleBySystemKey(ctx, uint32(groupKey))
}
