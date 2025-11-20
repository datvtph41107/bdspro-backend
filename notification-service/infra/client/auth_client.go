package client

import (
	"context"
	"notification/env"
	clients "pb/clients"
	sharepb "pb/types/shared"
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

// GetPushTokensByProfileId lấy danh sách push token theo profile ID
func (c *AuthClient) GetPushTokensByProfileId(ctx context.Context, profileID uint64) ([]string, error) {
	if c.Client == nil {
		return nil, nil
	}

	// AuthGrpcClient embed AuthInternalServiceClient, nên có thể gọi trực tiếp
	resp, err := c.Client.GetPushTokensByProfileId(ctx, &sharepb.IdRequest{
		Id: profileID,
	})
	if err != nil {
		return nil, err
	}

	return resp.PushTokens, nil
}
