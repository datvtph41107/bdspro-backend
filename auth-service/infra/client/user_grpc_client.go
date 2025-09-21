package client

import (
	"auth/env"
	"context"
	clients "pb/clients"
	sharepb "pb/types/shared"
)

// @bind: auth/internal/interface.IUserGrpcClient
type UserGrpcClient struct {
	Client *clients.UserGrpcClient
}

func NewUserGrpcClient() *UserGrpcClient {
	return &UserGrpcClient{
		Client: clients.MakeUserGrpcClient(env.ENV_RELEASE),
	}
}

func (c *UserGrpcClient) GetProfileByIds(ctx context.Context, in *sharepb.GetProfileByIdsRequest) (*sharepb.GetProfileByIdsResponse, error) {
	return c.Client.GetProfileByIds(ctx, in.Ids)
}

func (c *UserGrpcClient) GetMapProfileByIds(ctx context.Context, in *sharepb.GetProfileByIdsRequest) (map[uint64]*sharepb.ProfileItem, error) {
	profileIdSet := make(map[uint64]struct{})
	for _, id := range in.Ids {
		profileIdSet[id] = struct{}{}
	}
	return c.Client.GetMapByIDs(ctx, profileIdSet)
}
