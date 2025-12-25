package repo

import (
	_dto "common/domain/dto"
	"context"
)

type ContactProductRepo interface {
	Bulk(ctx context.Context, productIds []uint64, contactID uint64) error
	GetProductIds(ctx context.Context, contactID uint64) ([]uint64, error)
	GetProductIdsWithPagination(ctx context.Context, contactID uint64, pagable _dto.Pagable) ([]uint64, int64, error)
}
