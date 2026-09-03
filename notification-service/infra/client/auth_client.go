package client

import (
	"context"
	"errors"
	clients "pb/clients"
)

// @bind: notification/internal/usecase.PushTokenResolver
type AuthClient struct {
	Client *clients.UserGrpcClient
}

func NewAuthClient(rpcClient *clients.UserGrpcClient) *AuthClient {
	return &AuthClient{Client: rpcClient}
}

func (c *AuthClient) GetPushTokensByProfileId(ctx context.Context, profileID uint64) ([]string, error) {
	if c == nil || c.Client == nil {
		return nil, errors.New("notification user gRPC client is not configured")
	}

	resp, err := c.Client.GetPushTokensByProfileId(ctx, profileID)
	if err != nil {
		return nil, err
	}

	return resp.PushTokens, nil
}
