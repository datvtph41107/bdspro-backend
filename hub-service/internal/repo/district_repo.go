package repo

import (
	"context"
	"hub/internal/domain"
)

type IDistrictRepo interface {
	GetAll(ctx context.Context) ([]domain.District, error)
	GetByProvinceID(ctx context.Context, provinceID string) ([]domain.District, error)
	GetByIDs(ctx context.Context, ids []uint64) ([]domain.District, error)
	Search(ctx context.Context, keyword string) ([]domain.District, error)
	SearchByProvinceID(ctx context.Context, provinceID string, keyword string) ([]domain.District, error)
	SearchLocation(ctx context.Context, keyword string) ([]domain.LocationSearchResult, error)
}
