package property_repo

import (
	"bdspro/internal/domain"
	"context"
)

type PropertyLineageRepository interface {
	GetByID(ctx context.Context, id uint64) (*domain.PropertyLineage, error)
	UpdateFields(ctx context.Context, id uint64, fields map[string]any) error
}
