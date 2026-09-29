package client

import (
	"context"
	notification_usecase "notification/internal/usecase"
	clients "pb/clients"
	sharepb "pb/types/shared"
	userpb "pb/types/user"
)

type UserClient struct {
	Client *clients.UserGrpcClient
}

func NewUserClient(rpcClient *clients.UserGrpcClient) notification_usecase.UserProfileReader {
	return &UserClient{Client: rpcClient}
}

func (c *UserClient) GetProfileByIds(ctx context.Context, in *sharepb.GetProfileByIdsRequest) (*sharepb.GetProfileByIdsResponse, error) {
	return c.Client.GetProfileByIds(ctx, in.Ids)
}

func (c *UserClient) GetMapProfileByIds(ctx context.Context, in *sharepb.GetProfileByIdsRequest) (map[uint64]*sharepb.ProfileItem, error) {
	profileIdSet := make(map[uint64]struct{})
	for _, id := range in.Ids {
		profileIdSet[id] = struct{}{}
	}
	return c.Client.GetMapByIDs(ctx, profileIdSet)
}

func (c *UserClient) GetProfileByID(ctx context.Context, in *userpb.GetProfileByIDRequest) (*userpb.ProfileResponse, error) {
	return c.Client.InternalClient.GetProfileByID(ctx, in)
}
