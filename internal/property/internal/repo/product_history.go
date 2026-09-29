package repo

import (
	"bdspro/internal/domain"
	"bdspro/internal/dto"
	"bdspro/internal/enums"
	"context"
)

type ProductHistoryRepo interface {
	History(
		c context.Context,
		profileId uint64,
		productId uint64,
		dto dto.ProductHistorySearch,
	) (*[]domain.ProductHistory, int64, error)
	CreateHistory(
		c context.Context,
		action enums.EProductHistory,
		productID *uint64,
		parentID uint64,
		childIDs []uint64,
		areas []float64,
	) error
}
