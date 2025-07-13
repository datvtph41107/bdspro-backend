package repo

import (
	"bdspro/internal/domain"
	"common/case/crud3"
	"context"
)

type ProductPrivateRepo interface {
	crud3.CrudRepo[domain.ProductPrivate]
	UpdateInfo(c context.Context, productId uint64, entity *domain.ProductPrivate) error
	FindByProductId(c context.Context, productId *uint64) *domain.ProductPrivate
	GetPrivateFields(c context.Context, profileId uint64, product *domain.Product) (*domain.ProductPrivate, error)
}
