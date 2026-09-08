package repo

import (
	"context"
)

type ProductCareRepo interface {
	Bulk(ctx context.Context, productCares []uint64, leadID uint64) error
}