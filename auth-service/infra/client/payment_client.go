package client

import (
	"auth/internal/dto"
	"context"
	"net/http"
)

// @bind: auth/internal/interface/provider.PaymentProvider
type PaymentClient struct {
	HttpClient *http.Client
}

func NewPaymentClient(httpClient *http.Client) *PaymentClient {
	return &PaymentClient{HttpClient: httpClient}
}

func (c *PaymentClient) GetPayment(ctx context.Context, paymentID string) (*dto.PaymentDTO, error) {
	return nil, nil
}

func (c *PaymentClient) CreatePayment(ctx context.Context, payment *dto.PaymentDTO) (*dto.PaymentDTO, error) {
	return nil, nil
}

func (c *PaymentClient) UpdatePayment(ctx context.Context, payment *dto.PaymentDTO) (*dto.PaymentDTO, error) {
	return nil, nil
}

func (c *PaymentClient) DeletePayment(ctx context.Context, paymentID string) error {
	return nil
}
