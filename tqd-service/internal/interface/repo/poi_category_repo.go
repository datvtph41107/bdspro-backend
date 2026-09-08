package repo

import (
	"context"

	_crud "common/domain/crud"
	"tqd/internal/domain"
	"tqd/internal/dto"
)

// IPoiCategoryRepo defines repository interface for POI category
type IPoiCategoryRepo interface {
	_crud.ICrudRepo[domain.PoiCategory]

	// FindByCode finds category by code
	FindByCode(ctx context.Context, code string) (*domain.PoiCategory, error)

	// FindByParent finds categories by parent ID
	FindByParent(ctx context.Context, parentID *uint64) ([]domain.PoiCategory, error)

	// GetTree gets category tree
	GetTree(ctx context.Context) ([]domain.PoiCategory, error)

	// UpdatePath updates category path
	UpdatePath(ctx context.Context, id uint64, path string) error

	// UpdateLevel updates category level
	UpdateLevel(ctx context.Context, id uint64, level int) error

	// IncrementPOICount increments POI count for category
	IncrementPOICount(ctx context.Context, id uint64) error

	// DecrementPOICount decrements POI count for category
	DecrementPOICount(ctx context.Context, id uint64) error

	// ListWithFilter lists categories with filters
	ListWithFilter(ctx context.Context, filter *dto.PoiCategoryFilter) ([]domain.PoiCategory, int64, error)
}
