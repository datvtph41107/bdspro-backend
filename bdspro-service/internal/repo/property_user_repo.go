package repo

import (
	"bdspro/internal/domain"
	"context"
)

type PropertyUserRepo interface {
	CreateOwner(ctx context.Context, propertyId uint64, profileID uint64) (*domain.PropertyUser, error)
}
