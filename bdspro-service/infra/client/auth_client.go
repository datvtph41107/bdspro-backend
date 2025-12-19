package client

import (
	"bdspro/env"
	"context"
	"pb/clients"
	authpb "pb/types/auth"
	bdspropb "pb/types/bdspro"
)

type AuthClient struct {
	*clients.AuthGrpcClient
}

func NewAuthClient() *AuthClient {
	return &AuthClient{
		clients.MakeAuthGrpcClient(env.ENV_RUNTIME),
	}
}

func (c *AuthClient) MapRoleToDealInvitationPb(ctx context.Context, invitations []*bdspropb.DealMember) {
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
