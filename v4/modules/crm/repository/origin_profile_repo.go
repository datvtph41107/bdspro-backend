package repository

import (
	"context"
	domain "userv4/modules/crm/entity"
)

type OriginProfileRepo interface {
	GetByIDs(ctx context.Context, originIDs []uint64) ([]*domain.OriginProfile, error)
	GetByProfileID(ctx context.Context, profileID uint64) (*domain.OriginProfile, error)
}