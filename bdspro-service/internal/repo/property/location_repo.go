// repo/property/location_repo.go
package property_repo

import (
	"bdspro/internal/domain"
	"context"
)

type PropertyLocationRepository interface {
	GetByID(ctx context.Context, id uint64) (*domain.PropertyLocation, error)
	Create(ctx context.Context, location *domain.PropertyLocation) error
	UpdateFields(ctx context.Context, id uint64, fields map[string]any) error
}

type ProvinceV2Repository interface {
	GetByID(ctx context.Context, id uint64) (*domain.ProvinceV2, error)
	GetByCode(ctx context.Context, code string) (*domain.ProvinceV2, error)
	Create(ctx context.Context, province *domain.ProvinceV2) error
	UpdateFields(ctx context.Context, id uint64, fields map[string]any) error
	FirstOrCreate(ctx context.Context, province *domain.ProvinceV2) error
	GetByFullName(ctx context.Context, name string) (*domain.ProvinceV2, error)
	GetByTQDID(ctx context.Context, tqdID string) (*domain.ProvinceV2, error)
}

type WardV2Repository interface {
	GetByID(ctx context.Context, id uint64) (*domain.WardV2, error)
	GetByCode(ctx context.Context, code string) (*domain.WardV2, error)
	Create(ctx context.Context, ward *domain.WardV2) error
	UpdateFields(ctx context.Context, id uint64, fields map[string]any) error
	FirstOrCreate(ctx context.Context, ward *domain.WardV2) error
	GetByTQDID(ctx context.Context, tqdID string) (*domain.WardV2, error)
}
