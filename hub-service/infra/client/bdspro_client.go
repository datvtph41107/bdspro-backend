package client

import (
	"hub/env"

	"pb/clients"
)

type BDSProClient struct {
	*clients.BdsproGrpcClient
}

func NewBDSProClient() *BDSProClient {
	return &BDSProClient{
		BdsproGrpcClient: clients.MakeBdsproGrpcClient(env.ENV_RUNTIME),
	}
}
