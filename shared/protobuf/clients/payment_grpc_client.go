package clients

import (
	"context"
	"errors"

	paymentpb "pb/types/payment"
	sharepb "pb/types/shared"

	"google.golang.org/grpc"
)

type PaymentClient struct {
	client        paymentpb.InternalServiceClient
	ServiceClient paymentpb.PaymentServiceClient
}

func BindPaymentClient(conn grpc.ClientConnInterface) *PaymentClient {
	return &PaymentClient{
		client:        paymentpb.NewInternalServiceClient(conn),
		ServiceClient: paymentpb.NewPaymentServiceClient(conn),
	}
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
