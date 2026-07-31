package providers

import (
	"context"
	"user/internal/dto"
)

type IPaymentProvider interface {
	CreateOrder(ctx context.Context, req dto.OrderRequest) (*dto.OrderResponse, error)
	// GetPayment(ctx context.Context, paymentID string) (*dto.PaymentDTO, error)
	// CreatePayment(ctx context.Context, payment *dto.PaymentDTO) (*dto.PaymentDTO, error)
	// UpdatePayment(ctx context.Context, payment *dto.PaymentDTO) (*dto.PaymentDTO, error)
	// DeletePayment(ctx context.Context, paymentID string) error
}
