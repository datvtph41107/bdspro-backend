package client

import (
	"context"
	"notification/env"
	clients "pb/clients"
)

// @bind: notification/internal/provider.IAuthClient
type AuthClient struct {
	Client *clients.UserGrpcClient
}

func NewAuthClient() *AuthClient {
	return &AuthClient{
		Client: clients.MakeUserGrpcClient(env.ENV_RUNTIME),
	}
}

func (c *AuthClient) GetPushTokensByProfileId(ctx context.Context, profileID uint64) ([]string, error) {
	if c.Client == nil {
		return nil, nil
	}

	resp, err := c.Client.GetPushTokensByProfileId(ctx, profileID)
	if err != nil {
		return nil, err
	}

	return resp.PushTokens, nil
}
