package client

import (
	"organization/env"
	"pb/clients"
)

type TransactionClient struct {
	*clients.TransactionGrpcClient
}

func NewTransactionClient() *TransactionClient {
	return &TransactionClient{
		TransactionGrpcClient: clients.NewTransactionGrpcClient(env.ENV_RUNTIME),
	}
}
