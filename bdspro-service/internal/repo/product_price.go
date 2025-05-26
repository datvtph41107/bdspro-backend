package repo

import (
	"bdspro/internal/domain"
	"context"
)

type ProductPriceRepo interface {
	FindByProductId(productId *uint64) *domain.ProductPrice
	GetByID(c context.Context, id uint64) (*domain.ProductPrice, error)
	UpdateProductID(c context.Context, productId *uint64, priceId *uint64) error
	UpdatePrice(c context.Context, priceId uint64, entity *domain.ProductPrice) error
}
