package client

import (
	"context"
	"organization/env"
	"organization/internal/custom_error"
	"organization/internal/dto"
	"pb/clients"
	authpb "pb/types/auth"
	organizationpb "pb/types/organization"
	sharepb "pb/types/shared"
)

type AuthClient struct {
	*clients.AuthGrpcClient
}

func NewAuthClient() *AuthClient {
	return &AuthClient{
		AuthGrpcClient: clients.MakeAuthGrpcClient(env.ENV_RUNTIME),
	}
}

func (c *AuthClient) MapRoleToGroups(ctx context.Context, groups []*dto.GroupWithDetails) {
	roleIds := make([]uint64, len(groups))
	for i, r := range groups {
		if r.RoleId != nil {
			roleIds[i] = *r.RoleId
		}
	}

	roles, err := c.GetRolesByIds(ctx, &sharepb.IdRequest{Ids: roleIds})
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

// GetAdminRoleIdByGroupKey lấy admin role theo group key, ưu tiên role có auto_assign = true
func (c *AuthClient) GetRoleBySystemKey(ctx context.Context, systemKey uint32) (*authpb.Role, error) {
	// Gọi sang auth service để lấy admin role
	response, err := c.AuthGrpcClient.GetRolesByRoleKeys(ctx, &authpb.GetRolesByRoleKeysRequest{
		RoleKeys: []uint32{systemKey},
	})
	if err != nil {
		return nil, err
	}
	if len(response.Data) == 0 {
		return nil, custom_error.RecordNotFound("Role not found")
	}
	return response.Data[0], nil
}

func (c *AuthClient) MapRoleToDealInvitationPb(ctx context.Context, invitations []*organizationpb.DealMember) {
	roleIds := make([]uint32, len(invitations))
	for i, r := range invitations {
		if r.RoleKey != 0 {
			roleIds[i] = r.RoleKey
		}
	}

	roles, err := c.GetRolesByRoleKeys(ctx, &authpb.GetRolesByRoleKeysRequest{RoleKeys: roleIds})
	if err != nil {
		return
	}

	roleMap := make(map[uint32]*authpb.Role)
	for _, r := range roles.Data {
		roleMap[r.RoleKey] = r
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

	roles, err := c.GetRolesByRoleKeys(ctx, &authpb.GetRolesByRoleKeysRequest{RoleKeys: roleIds})
	if err != nil {
		return err
	}

	roleMap := make(map[uint32]*authpb.Role)
	for _, r := range roles.Data {
		roleMap[r.RoleKey] = r
	}

	for i, r := range orgMembers {
		if r.GetRoleKey() != 0 {
			orgMembers[i].RoleName = roleMap[r.RoleKey].RoleName
			orgMembers[i].RoleKey = roleMap[r.RoleKey].RoleKey
			orgMembers[i].RoleId = roleMap[r.RoleKey].Id
		}
	}
	return nil
}
