package repo

import (
	"context"
	"hub/internal/domain"
)

type IWardV2Repo interface {
	// GetAll returns all wards
	GetAll(ctx context.Context) ([]domain.WardV2, error)

	// GetList returns paginated list of wards
	GetList(ctx context.Context, provinceCode *int, keyword string, page, size int) ([]domain.WardV2, int64, error)

	// GetByID returns ward by ID
	GetByID(ctx context.Context, id uint64) (*domain.WardV2, error)

	// GetByCode returns ward by code
	GetByCode(ctx context.Context, code int) (*domain.WardV2, error)

	// GetByProvinceCode returns all wards of a province
	GetByProvinceCode(ctx context.Context, provinceCode int) ([]domain.WardV2, error)

	// Search searches wards by keyword
	Search(ctx context.Context, keyword string, page, size int) ([]domain.WardV2, int64, error)

	// SearchLocation searches wards with province info by joining two tables
	SearchLocation(ctx context.Context, keyword string, page, size int) ([]domain.LocationSearchResultV2, int64, error)

	InferFromText(ctx context.Context, text string) (*domain.WardV2, error)
}
