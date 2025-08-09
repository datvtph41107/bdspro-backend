package client

import (
	"context"
	"pb/clients"
	sharepb "pb/types/shared"

	"github.com/tuandq2112/appointment-service/internal/env"
)

type UserClient struct {
	*clients.UserGrpcClient
}

func NewUserClient() *UserClient {
	return &UserClient{
		UserGrpcClient: clients.MakeUserGrpcClient(env.ENV_RUNTIME),
	}
}

func (c *UserClient) GetProfileByIds(ctx context.Context, ids []uint64) ([]*sharepb.ProfileItem, error) {
	request := &sharepb.GetProfileByIdsRequest{
		Ids: ids,
	}
	response, err := c.Client.GetProfileByIds(ctx, request)
	if err != nil {
		return nil, err
	}

	return response.Profiles, nil
}
