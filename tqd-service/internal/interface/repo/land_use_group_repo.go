package repo

import (
	"context"

	qh_domain "tqd/internal/domain/qh"
)

type LandUseGroupRepository interface {
	GetByID(ctx context.Context, id uint64) (*qh_domain.LandUseGroup, error)
}
