package client

import (
	"context"
	"pb/clients"
	sharepb "pb/types/shared"
)

type AuthClient struct {
	*clients.AuthGrpcClient
	userClient *clients.UserGrpcClient
}

func NewAuthClient(
	authClient *clients.AuthGrpcClient,
	userClient *clients.UserGrpcClient,
) *AuthClient {
	return &AuthClient{
		AuthGrpcClient: authClient,
		userClient:     userClient,
	}
}

func (c *AuthClient) MapRoleToDealInvitationPb(ctx context.Context, invitations []*sharepb.DealMember) {
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

	roleMap := make(map[uint32]*sharepb.Role)
	for _, r := range roles.Data {
		role := &sharepb.Role{
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
			role.Color = &sharepb.Color{
				Id:              r.Color.Id,
				Name:            r.Color.Name,
				ContentColor:    r.Color.ContentColor,
				BackgroundColor: r.Color.BackgroundColor,
			}
			role.ColorId = &r.Color.Id
		}

		roleMap[r.RoleKey] = role
	}

	for i, r := range invitations {
		if r.RoleKey != 0 {
			invitations[i].Role = roleMap[r.RoleKey]
		}
	}
}
