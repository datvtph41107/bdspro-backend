package property_repo

import (
	"bdspro/internal/domain"
	"context"
)

type PropertyLineageRepository interface {
	Create(ctx context.Context, entity *domain.PropertyLineage) error
	GetByID(ctx context.Context, id uint64) (*domain.PropertyLineage, error)
	GetByPropertyIdentifyID(ctx context.Context, propertyIdentifyID uint64) (*domain.PropertyLineage, error)

	GetByIdAndOwnerId(ctx context.Context, id uint64, ownerID uint64) (*domain.PropertyLineage, error)
	UpdateFields(ctx context.Context, id uint64, fields map[string]any) error
}
