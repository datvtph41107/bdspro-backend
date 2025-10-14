package providers

import (
	"context"
	"user/internal/dto"
)

type IPaymentProvider interface {
	CreateOrder(ctx context.Context, req dto.OrderRequest) (*dto.OrderResponse, error)
}
