package repo

import (
	"context"
	"crm/internal/domain"
)

type OriginProfileRepo interface {
	GetByIDs(ctx context.Context, originIDs []uint64) ([]*domain.OriginProfile, error)
	GetByProfileID(ctx context.Context, profileID uint64) (*domain.OriginProfile, error)
}