package client

import (
	"bdspro/env"
	"pb/clients"
)

type PaymentClient struct {
	*clients.PaymentClient
}

func NewPaymentClient() *PaymentClient {
	return &PaymentClient{
		PaymentClient: clients.MakePaymentClient(env.ENV_RUNTIME),
	}
}
