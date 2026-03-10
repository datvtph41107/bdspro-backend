package property_repo

import (
	"bdspro/internal/domain"
	"context"
)

// PropertyLocationRepo — quản lý bảng property_location.
type PropertyLocationRepository interface {
	GetByID(ctx context.Context, id uint64) (*domain.PropertyLocation, error)
	GetByLineageID(ctx context.Context, lineageID uint64) (*domain.PropertyLocation, error)
	Create(ctx context.Context, entity *domain.PropertyLocation) error
	UpdateFields(ctx context.Context, id uint64, fields map[string]any) error
}
