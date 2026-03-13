package repo

import (
	"bdspro/internal/domain"
	"context"
)

type PropertyUserRepo interface {
	CreateOwner(ctx context.Context, propertyId uint64, profileID uint64) (*domain.PropertyUser, error)
	GetByLineageAndOwner(ctx context.Context, propertyId uint64, ownerOriginID uint64) (*domain.PropertyUser, error)
	UpdateFieldsByLineageAndOwner(ctx context.Context, propertyId uint64, ownerOriginID uint64, fields map[string]any) error
}
