package repo

import (
	_crud "common/domain/crud"
	"context"
	"tqd/internal/domain"
)

// IPoiCategoryRepo defines the repository interface for POI category operations
type IPoiCategoryRepo interface {
	_crud.ICrudRepo[domain.PoiCategory]
	// GetList(ctx context.Context, req *dto.PoiCategoryListRequest) ([]domain.PoiCategory, int64, error)
	GetTree(ctx context.Context) ([]domain.PoiCategory, error)
}
