package property_repo

import (
	"bdspro/internal/domain"
	"context"
)

type PropertyInfoRepository interface {
	GetByID(ctx context.Context, id uint64) (*domain.PropertyInfo, error)
	GetByLineageID(ctx context.Context, lineageID uint64) (*domain.PropertyInfo, error)
	Create(ctx context.Context, entity *domain.PropertyInfo) error
	UpdateFields(ctx context.Context, id uint64, fields map[string]any) error
}
