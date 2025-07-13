package repo

import (
	"bdspro/internal/domain"
	"context"
)

type ProductPriceRepo interface {
	FindByProductId(ctx context.Context, productId *uint64) *domain.ProductPrice
	GetByID(c context.Context, id uint64) (*domain.ProductPrice, error)
	UpdateProductID(c context.Context, productId *uint64, priceId *uint64) error
	UpdatePrice(ctx context.Context, priceId *uint64, entity *domain.ProductPrice) error
	Create(c context.Context, entity *domain.ProductPrice) error
}
