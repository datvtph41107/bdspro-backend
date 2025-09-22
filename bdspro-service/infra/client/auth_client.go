package client

import (
	"bdspro/env"
	"pb/clients"
)

func NewAuthClient() *clients.AuthGrpcClient {
	return clients.MakeAuthGrpcClient(env.ENV_RELEASE)
}
