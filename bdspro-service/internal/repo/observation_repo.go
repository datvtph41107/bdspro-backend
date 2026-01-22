package repo

import (
	"bdspro/internal/domain"
	"bdspro/internal/dto"
	"context"
)

type ObservationRepository interface {
	Insert(ctx context.Context, o *domain.Observation) error
	GetProductStats(ctx context.Context, productID uint64) (*dto.ProductStats, error)
}
