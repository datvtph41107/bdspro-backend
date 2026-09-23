package repo

import (
	"context"
	"tqd/internal/domain"
	"tqd/internal/dto"
)

type LocationRepository interface {
	// Province operations
	GetProvinceByID(ctx context.Context, id uint64) (*domain.Province, error)
	GetProvinceDetail(ctx context.Context, id uint64) (*domain.Province, error)
	GetProvinceByCode(ctx context.Context, code string) (*domain.Province, error)
	GetProvinceByCompositeId(ctx context.Context, compositeId string) (*domain.Province, error)
	GetNearestProvince(ctx context.Context, lat, lng float64, maxDistanceKm float64) (*domain.Province, float64, error)
	SearchProvinces(ctx context.Context, query string, limit int) ([]*domain.Province, error)
	ListAllProvinces(ctx context.Context) ([]*domain.Province, error)
	AdminListProvinces(ctx context.Context, req *dto.AdminListProvincesRequest) ([]*domain.Province, int64, error)
	ListAdministrativeUnitRefs(ctx context.Context, req *dto.ListAdministrativeUnitsRequest) ([]dto.AdministrativeUnitRef, int64, error)
	GetProvincesByIDs(ctx context.Context, ids []uint64) ([]*domain.Province, error)
	GetWardsByIDs(ctx context.Context, ids []uint64) ([]*domain.Ward, error)
	FollowAdministrativeUnit(ctx context.Context, userID, compositeID uint64) error
	UnfollowAdministrativeUnit(ctx context.Context, userID, compositeID uint64) error
	ListFollowedAdministrativeUnitIDs(ctx context.Context, userID uint64, offset, limit int) ([]uint64, int64, error)
	CreateProvince(ctx context.Context, province *domain.Province) error
	UpdateProvince(ctx context.Context, province *domain.Province) error
	SoftDeleteProvince(ctx context.Context, id uint64) error
	ExistsProvinceCode(ctx context.Context, code string, excludeID uint64) (bool, error)

	// Ward operations
	GetWardByID(ctx context.Context, id uint64) (*domain.Ward, error)
	GetWardDetail(ctx context.Context, id uint64) (*domain.Ward, error)
	GetWardByCode(ctx context.Context, code string) (*domain.Ward, error)
	GetWardByCompositeId(ctx context.Context, compositeId string) (*domain.Ward, error)
	GetNearestWard(ctx context.Context, lat, lng float64, provinceID uint64, maxDistanceKm float64) (*domain.Ward, float64, error)
	SearchWards(ctx context.Context, query string, provinceID uint64, limit int) ([]*domain.Ward, error)
	ListWardsByProvince(ctx context.Context, provinceID uint64, page, pageSize int) ([]*domain.Ward, int64, error)
	ListAllWards(ctx context.Context, page, pageSize int) ([]*domain.Ward, int64, error)
	AdminListWards(ctx context.Context, req *dto.AdminListWardsRequest) ([]*domain.Ward, int64, error)
	CreateWard(ctx context.Context, ward *domain.Ward) error
	UpdateWard(ctx context.Context, ward *domain.Ward) error
	SoftDeleteWard(ctx context.Context, id uint64) error
	SoftDeleteWardsByProvinceID(ctx context.Context, provinceID uint64) error
	ExistsWardCode(ctx context.Context, code string, excludeID uint64) (bool, error)
	// Ward neighbor operations
	CreateWardNeighbor(ctx context.Context, wardID uint64, neighborWardID uint64) error
	DeleteWardNeighbor(ctx context.Context, wardID uint64, neighborWardID uint64) error
	GetWardNeighbors(ctx context.Context, wardID uint64) ([]*domain.Ward, error)

	// Batch operations
	GetNearestLocations(ctx context.Context, coordinates []*dto.Coordinate, maxDistanceKm float64) ([]*dto.LocationResult, error)
	GetLocationByParcelID(ctx context.Context, parcelID uint64) (*domain.Province, *domain.Ward, error)
	SearchTxtClient(ctx context.Context, query string, limit int) ([]*domain.SearchItem, error)

	// SEO operations
	GetAdministrativeUnitSeoSource(ctx context.Context, identity uint64) (*dto.AdministrativeUnitSeoSource, error)
	GetAdministrativeUnitSeoSourcesByIdentities(ctx context.Context, identities []uint64) ([]dto.AdministrativeUnitSeoSource, error)
	UpdateAdministrativeUnitSeoID(ctx context.Context, identity uint64, seoID uint64) error

	// Area sync operations
	SyncProvinceAreaFromGeometry(ctx context.Context, provinceID uint64) (bool, error)
	SyncWardAreaFromGeometry(ctx context.Context, wardID uint64) (bool, error)
	SyncAllProvincesAreaFromGeometry(ctx context.Context) (int64, error)
	SyncAllWardsAreaFromGeometry(ctx context.Context) (int64, error)
}
