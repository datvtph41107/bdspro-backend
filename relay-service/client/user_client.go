package client

import (
	"context"
	"fmt"
	"pb/clients"
	sharepb "pb/types/shared"
	userpb "pb/types/user"
)

type UserClient struct {
	InternalClient userpb.InternalUserServiceClient
}

func NewUserClient(rpcClient *clients.UserGrpcClient) *UserClient {
	return &UserClient{InternalClient: rpcClient.InternalClient}
}

func (c *UserClient) UpdateLastSeen(ctx context.Context, profileID uint64) error {
	if c.InternalClient == nil {
		return fmt.Errorf("user client not available")
	}

	_, err := c.InternalClient.UpdateLastSeen(ctx, &sharepb.IdRequest{
		Id: profileID,
	})
	return err
}
