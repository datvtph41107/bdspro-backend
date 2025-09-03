package client

import (
	"context"
	"organization/env"
	"organization/internal/dto"
	"pb/clients"
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
