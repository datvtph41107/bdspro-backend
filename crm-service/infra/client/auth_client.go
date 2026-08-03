package client

import (
	"crm/env"
	"pb/clients"
)

type AuthClient struct {
	*clients.AuthGrpcClient
	UserClient *clients.UserGrpcClient
}

func NewAuthClient() *AuthClient {
	return &AuthClient{
		AuthGrpcClient: clients.MakeAuthGrpcClient(env.ENV_RUNTIME),
		UserClient:     clients.MakeUserGrpcClient(env.ENV_RUNTIME),
	}
}
