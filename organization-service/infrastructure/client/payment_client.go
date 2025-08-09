package client

import (
	"organization/env"
	clients "pb/clients"
)

type PaymentClient struct {
	*clients.PaymentClient
}

func NewPaymentClient() *PaymentClient {
	return &PaymentClient{
		PaymentClient: clients.MakePaymentClient(env.ENV_RUNTIME),
	}
}
