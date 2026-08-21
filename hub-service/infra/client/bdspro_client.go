package client

import (
	"pb/clients"
)

type BDSProClient struct {
	*clients.BdsproGrpcClient
}

func NewBDSProClient(rpcClient *clients.BdsproGrpcClient) *BDSProClient {
	return &BDSProClient{BdsproGrpcClient: rpcClient}
}
