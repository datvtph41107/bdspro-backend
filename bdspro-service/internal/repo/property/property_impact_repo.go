package property_repo

import (
	"bdspro/internal/domain"
	"bdspro/internal/dto"
	"context"
	"time"
)

type PropertyImpactRepository interface {
	GetCurrent(ctx context.Context, entityType string, entityID uint64) (*domain.PropertyImpactVersion, error)
	GetByProperty(ctx context.Context, propertyID uint64, entityTypes []string) ([]*domain.PropertyImpactVersion, error)
	Create(ctx context.Context, snapshot *domain.PropertyImpactVersion) error
	Invalidate(ctx context.Context, snapshotID uint64, validTo time.Time) error
	GetAffectedEntities(ctx context.Context, propertyID uint64) (*dto.AffectedEntities, error)
}
