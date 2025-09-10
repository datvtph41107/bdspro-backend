package client

import (
	"context"
	"organization/env"
	"organization/internal/custom_error"
	"organization/internal/dto"
	"pb/clients"
	authpb "pb/types/auth"
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
	response, err := c.AuthGrpcClient.GetRolesBySystemKeys(ctx, &authpb.GetRolesBySystemKeysRequest{
		SystemKeys: []uint32{systemKey},
	})
	if err != nil {
		return nil, err
	}
	if len(response.Data) == 0 {
		return nil, custom_error.RecordNotFound("Role not found")
	}
	return response.Data[0], nil
}
