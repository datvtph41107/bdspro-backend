package provider

import (
	"auth/internal/dto"
	"context"
)

type PaymentProvider interface {
	GetPayment(ctx context.Context, paymentID string) (*dto.PaymentDTO, error)
	CreatePayment(ctx context.Context, payment *dto.PaymentDTO) (*dto.PaymentDTO, error)
	UpdatePayment(ctx context.Context, payment *dto.PaymentDTO) (*dto.PaymentDTO, error)
	DeletePayment(ctx context.Context, paymentID string) error
}
