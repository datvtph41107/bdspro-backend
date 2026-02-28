package repo

import (
	"context"
	"hub/internal/domain"
)

type RegionRepository interface {
	GetByIDs(ctx context.Context, ids []uint64) ([]domain.Region, error)
}
