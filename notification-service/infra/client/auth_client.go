package client

import (
	"context"
	"notification/env"
	sharepb "pb/types/shared"
	clients "pb/clients"
)

// @bind: notification/internal/interface.IAuthClient
type AuthClient struct {
	Client *clients.AuthGrpcClient
}

func NewAuthClient() *AuthClient {
	return &AuthClient{
		Client: clients.MakeAuthGrpcClient(env.ENV_RUNTIME),
	}
}

func (c *AuthClient) GetProfileByIds(ctx context.Context, in *sharepb.GetProfileByIdsRequest) (*sharepb.GetProfileByIdsResponse, error) {
	return c.Client.GetProfileByIds(ctx, in.Ids)
}

func (c *AuthClient) GetMapProfileByIds(ctx context.Context, in *sharepb.GetProfileByIdsRequest) (map[uint64]*sharepb.ProfileItem, error) {
	profileIdSet := make(map[uint64]struct{})
	for _, id := range in.Ids {
		profileIdSet[id] = struct{}{}
	}
	return c.Client.GetMapByIDs(ctx, profileIdSet)
}
