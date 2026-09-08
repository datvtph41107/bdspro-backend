package client

import (
	"pb/clients"
)

type AuthClient struct {
	*clients.AuthGrpcClient
	UserClient *clients.UserGrpcClient
}

func NewAuthClient(
	authClient *clients.AuthGrpcClient,
	userClient *clients.UserGrpcClient,
) *AuthClient {
	return &AuthClient{
		AuthGrpcClient: authClient,
		UserClient:     userClient,
	}
}
