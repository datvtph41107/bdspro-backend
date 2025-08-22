package client

import (
	"pb/clients"
)

func NewAuthClient() *clients.AuthGrpcClient {
	return clients.MakeAuthGrpcClient("local")
}
