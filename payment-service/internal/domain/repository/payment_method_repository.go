package repository

import (
	"context"

	paymentpb "pb/types/payment"

	"payment/internal/domain/entity"
)

type PaymentMethodRepository interface {
	GetPaymentMethods(ctx context.Context, organizationId uint32, page int, size int) ([]*paymentpb.PaymentMethod, int64, error)
	CreatePaymentMethod(ctx context.Context, method *entity.PaymentMethod) (*entity.PaymentMethod, error)
	GetPaymentMethodById(ctx context.Context, id uint32) (*entity.PaymentMethod, error)
	UpdatePaymentMethod(ctx context.Context, method *entity.PaymentMethod) (*entity.PaymentMethod, error)
	DeletePaymentMethod(ctx context.Context, id uint32) error
}
