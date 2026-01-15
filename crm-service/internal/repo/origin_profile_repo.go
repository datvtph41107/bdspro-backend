package repo

import (
	"context"
	"crm/internal/domain"
)

type OriginProfileRepo interface {
	GetByIDs(ctx context.Context, originIDs []uint64) ([]*domain.OriginProfile, error)
}
