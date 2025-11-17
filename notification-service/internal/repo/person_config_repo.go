package repo

import (
	"context"

	"notification/internal/domain"
)

type PersonConfigRepo interface {
	BatchUpsert(ctx context.Context, configs []*domain.PersonConfigEntity) ([]*domain.PersonConfigEntity, error)
	ListByUserID(ctx context.Context, userID uint64) ([]*domain.PersonConfigEntity, error)
}
