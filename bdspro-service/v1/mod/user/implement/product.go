package product_impl

import (
	"bdspro/v1/models"
	"context"
)

type ProductUpdateImpl struct {
}

func (impl ProductUpdateImpl) Save(c context.Context, price models.Product, id *uint64) error {
	return nil
}
func (impl ProductUpdateImpl) UpdatePrice(c context.Context, price models.ProductPrice) (*uint64, error) {
	return nil, nil
}
