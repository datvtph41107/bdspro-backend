package clients

import (
	"common/tcp/connect"
	"context"
	"errors"
	"log"
	paymentpb "pb/types/payment"
	sharepb "pb/types/shared"
)

type PaymentClient struct {
	client paymentpb.InternalServiceClient
}

func MakePaymentClient(env_runtime string) *PaymentClient {
	paymentClient := &PaymentClient{}
	name := "payment"
	addr := connect.AddressOf(env_runtime, name, connect.DefaultPort(name))

	conn, err := connect.DialNamed(name, addr)
	if err != nil {
		log.Printf("failed to dial payment at %s: %v", addr, err)
		return paymentClient
	}

	paymentClient.client = paymentpb.NewInternalServiceClient(conn)
	return paymentClient
}

func (c *PaymentClient) GetBankById(ctx context.Context, id uint64) (*paymentpb.Bank, error) {
	if c.client == nil {
		return nil, errors.New("payment client not avaiable")
	}

	bank, err := c.client.GetBankById(ctx, &sharepb.IdRequest{Id: id})
	if err != nil {
		return nil, err
	}
	return bank, nil
}
