package client

import (
	"pb/clients"
)

type PaymentClient struct {
	*clients.PaymentClient
}

func NewPaymentClient(rpcClient *clients.PaymentClient) *PaymentClient {
	return &PaymentClient{PaymentClient: rpcClient}
}
