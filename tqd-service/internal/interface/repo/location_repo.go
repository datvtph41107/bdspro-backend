// internal/interface/repo/location_repo.go
package repo

import (
	"context"
	"tqd/internal/domain"
	"tqd/internal/dto"
)

type LocationRepository interface {
	// Province operations
	GetProvinceByID(ctx context.Context, id string) (*domain.Province, error)
	GetProvinceByCode(ctx context.Context, code string) (*domain.Province, error)
	GetNearestProvince(ctx context.Context, lat, lng float64, maxDistanceKm float64) (*domain.Province, float64, error)
	SearchProvinces(ctx context.Context, query string, limit int) ([]*domain.Province, error)
	ListAllProvinces(ctx context.Context) ([]*domain.Province, error)
	AdminListProvinces(ctx context.Context, req *dto.AdminListProvincesRequest) ([]*domain.Province, int64, error)

	// Ward operations
	GetWardByID(ctx context.Context, id string) (*domain.Ward, error)
	GetWardByCode(ctx context.Context, code string) (*domain.Ward, error)
	GetNearestWard(ctx context.Context, lat, lng float64, provinceID string, maxDistanceKm float64) (*domain.Ward, float64, error)
	SearchWards(ctx context.Context, query string, provinceID string, limit int) ([]*domain.Ward, error)
	ListWardsByProvince(ctx context.Context, provinceID string, page, pageSize int) ([]*domain.Ward, int64, error)
	AdminListWards(ctx context.Context, req *dto.AdminListWardsRequest) ([]*domain.Ward, int64, error)

	// Batch operations
	GetNearestLocations(ctx context.Context, coordinates []*dto.Coordinate, maxDistanceKm float64) ([]*dto.LocationResult, error)
	GetLocationByParcelID(ctx context.Context, parcelID uint64) (*domain.Province, *domain.Ward, error)
	SearchTxtClient(ctx context.Context, query string, limit int) ([]*domain.SearchItem, error)

	// SEO operations
	GetAdministrativeUnitSeoSource(ctx context.Context, identity string) (*dto.AdministrativeUnitSeoSource, error)
	GetAdministrativeUnitSeoSourcesByIdentities(ctx context.Context, identities []string) ([]dto.AdministrativeUnitSeoSource, error)
	UpdateAdministrativeUnitSeoID(ctx context.Context, identity string, seoID uint64) error
}
