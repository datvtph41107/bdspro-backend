package repository

import (
	"context"
	"fmt"

	"userv4/infra/postgres"
	"userv4/modules/tqd/dto"
	"userv4/modules/tqd/entity"

	"gorm.io/gorm"
)

// LocationRepository is the TQD administrative location store on the TQD DSN.
type LocationRepository interface {
	// Province operations
	GetProvinceByID(ctx context.Context, id uint64) (*entity.Province, error)
	GetProvinceDetail(ctx context.Context, id uint64) (*entity.Province, error)
	GetProvinceByCode(ctx context.Context, code string) (*entity.Province, error)
	GetProvinceByCompositeId(ctx context.Context, compositeId string) (*entity.Province, error)
	GetNearestProvince(ctx context.Context, lat, lng float64, maxDistanceKm float64) (*entity.Province, float64, error)
	SearchProvinces(ctx context.Context, query string, limit int) ([]*entity.Province, error)
	ListAllProvinces(ctx context.Context) ([]*entity.Province, error)
	AdminListProvinces(ctx context.Context, req *dto.AdminListProvincesRequest) ([]*entity.Province, int64, error)
	ListAdministrativeUnitRefs(ctx context.Context, req *dto.ListAdministrativeUnitsRequest) ([]dto.AdministrativeUnitRef, int64, error)
	GetProvincesByIDs(ctx context.Context, ids []uint64) ([]*entity.Province, error)
	GetWardsByIDs(ctx context.Context, ids []uint64) ([]*entity.Ward, error)
	FollowAdministrativeUnit(ctx context.Context, userID, compositeID uint64) error
	UnfollowAdministrativeUnit(ctx context.Context, userID, compositeID uint64) error
	ListFollowedAdministrativeUnitIDs(ctx context.Context, userID uint64, offset, limit int) ([]uint64, int64, error)
	CreateProvince(ctx context.Context, province *entity.Province) error
	UpdateProvince(ctx context.Context, province *entity.Province) error
	SoftDeleteProvince(ctx context.Context, id uint64) error
	ExistsProvinceCode(ctx context.Context, code string, excludeID uint64) (bool, error)

	// Ward operations
	GetWardByID(ctx context.Context, id uint64) (*entity.Ward, error)
	GetWardDetail(ctx context.Context, id uint64) (*entity.Ward, error)
	GetWardByCode(ctx context.Context, code string) (*entity.Ward, error)
	GetWardByCompositeId(ctx context.Context, compositeId string) (*entity.Ward, error)
	GetNearestWard(ctx context.Context, lat, lng float64, provinceID uint64, maxDistanceKm float64) (*entity.Ward, float64, error)
	SearchWards(ctx context.Context, query string, provinceID uint64, limit int) ([]*entity.Ward, error)
	ListWardsByProvince(ctx context.Context, provinceID uint64, page, pageSize int) ([]*entity.Ward, int64, error)
	ListAllWards(ctx context.Context, page, pageSize int) ([]*entity.Ward, int64, error)
	AdminListWards(ctx context.Context, req *dto.AdminListWardsRequest) ([]*entity.Ward, int64, error)
	CreateWard(ctx context.Context, ward *entity.Ward) error
	UpdateWard(ctx context.Context, ward *entity.Ward) error
	SoftDeleteWard(ctx context.Context, id uint64) error
	SoftDeleteWardsByProvinceID(ctx context.Context, provinceID uint64) error
	ExistsWardCode(ctx context.Context, code string, excludeID uint64) (bool, error)
	CreateWardNeighbor(ctx context.Context, wardID uint64, neighborWardID uint64) error
	DeleteWardNeighbor(ctx context.Context, wardID uint64, neighborWardID uint64) error
	GetWardNeighbors(ctx context.Context, wardID uint64) ([]*entity.Ward, error)

	// Batch operations
	GetNearestLocations(ctx context.Context, coordinates []*dto.Coordinate, maxDistanceKm float64) ([]*dto.LocationResult, error)
	GetLocationByParcelID(ctx context.Context, parcelID uint64) (*entity.Province, *entity.Ward, error)
	SearchTxtClient(ctx context.Context, query string, limit int) ([]*entity.SearchItem, error)

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

type locationPostgres struct{ db *gorm.DB }

func NewLocationPostgres(tqd *postgres.TqdDB) LocationRepository {
	if tqd == nil {
		return &locationPostgres{}
	}
	return &locationPostgres{db: tqd.DB}
}

func (r *locationPostgres) active(db *gorm.DB) *gorm.DB {
	return db.Where("deleted_at IS NULL")
}

func provinceSelectSQL() string {
	return `id, full_name, short_name, lat, lng, code,
		ward_count, qh_order, zoom,
		'' AS description,
		COALESCE(official_name, '') AS official_name,
		COALESCE(division_type, '') AS division_type,
		official_operation_date,
		COALESCE(natural_area_km2, 0) AS natural_area_km2,
		COALESCE(population, 0) AS population,
		COALESCE(population_density, 0) AS population_density,
		COALESCE(unit_type, 0) AS unit_type,
		COALESCE(data_source, 0) AS data_source,
		plan_doc_id,
		COALESCE(official_role, '') AS official_role,
		COALESCE(official_address, '') AS official_address,
		seo_id, created_at, updated_at`
}

func wardSelectSQL() string {
	return `id, full_name, short_name, lat, lng, code,
		province_id, province_code,
		'' AS description, COALESCE(summary, '') AS summary,
		COALESCE(official_name, '') AS official_name,
		COALESCE(division_type, '') AS division_type,
		official_operation_date,
		COALESCE(natural_area_km2, 0) AS natural_area_km2,
		COALESCE(population, 0) AS population,
		COALESCE(population_density, 0) AS population_density,
		COALESCE(unit_type, 0) AS unit_type,
		COALESCE(data_source, 0) AS data_source,
		COALESCE(official_role, '') AS official_role,
		COALESCE(official_address, '') AS official_address,
		seo_id, created_at, updated_at`
}

func (r *locationPostgres) ListAllProvinces(ctx context.Context) ([]*entity.Province, error) {
	if r.db == nil {
		return nil, nil
	}
	var provinces []*entity.Province
	err := r.active(r.db.WithContext(ctx).Model(&entity.Province{})).
		Select(provinceSelectSQL()).
		Order("full_name").
		Find(&provinces).Error
	if err != nil {
		return nil, fmt.Errorf("list provinces: %w", err)
	}
	return provinces, nil
}

func (r *locationPostgres) ListAllWards(ctx context.Context, page, pageSize int) ([]*entity.Ward, int64, error) {
	if r.db == nil {
		return nil, 0, nil
	}
	var wards []*entity.Ward
	var total int64

	err := r.active(r.db.WithContext(ctx).Model(&entity.Ward{})).Count(&total).Error
	if err != nil {
		return nil, 0, fmt.Errorf("count wards: %w", err)
	}

	offset := (page - 1) * pageSize
	err = r.active(r.db.WithContext(ctx).Model(&entity.Ward{})).
		Select(wardSelectSQL()).
		Order("full_name").
		Offset(offset).
		Limit(pageSize).
		Find(&wards).Error
	if err != nil {
		return nil, 0, fmt.Errorf("list wards: %w", err)
	}
	return wards, total, nil
}

func (r *locationPostgres) GetLocationByParcelID(ctx context.Context, parcelID uint64) (*entity.Province, *entity.Ward, error) {
	if r.db == nil {
		return nil, nil, nil
	}
	var info struct {
		ProvinceCode string
		WardCode     string
	}
	err := r.db.WithContext(ctx).
		Table("qh_parcel_info").
		Select("COALESCE(province_code,'') AS province_code, COALESCE(ward_code,'') AS ward_code").
		Where("parcel_id = ?", parcelID).
		Limit(1).
		Scan(&info).Error
	if err != nil {
		return nil, nil, fmt.Errorf("parcel location codes: %w", err)
	}
	var province *entity.Province
	var ward *entity.Ward
	if info.ProvinceCode != "" {
		var p entity.Province
		if err := r.active(r.db.WithContext(ctx).Model(&entity.Province{})).
			Select(provinceSelectSQL()).
			Where("code = ?", info.ProvinceCode).
			First(&p).Error; err == nil {
			province = &p
		}
	}
	if info.WardCode != "" {
		var w entity.Ward
		if err := r.active(r.db.WithContext(ctx).Model(&entity.Ward{})).
			Select(wardSelectSQL()).
			Where("code = ?", info.WardCode).
			First(&w).Error; err == nil {
			ward = &w
		}
	}
	return province, ward, nil
}

func errLocationNotPorted(method string) error {
	return fmt.Errorf("location repository %s not available on phase-1 catalog store; use modules/tqd/postgres.NewLocationRepository", method)
}

func (r *locationPostgres) GetProvinceByID(context.Context, uint64) (*entity.Province, error) {
	return nil, errLocationNotPorted("GetProvinceByID")
}
func (r *locationPostgres) GetProvinceDetail(context.Context, uint64) (*entity.Province, error) {
	return nil, errLocationNotPorted("GetProvinceDetail")
}
func (r *locationPostgres) GetProvinceByCode(context.Context, string) (*entity.Province, error) {
	return nil, errLocationNotPorted("GetProvinceByCode")
}
func (r *locationPostgres) GetProvinceByCompositeId(context.Context, string) (*entity.Province, error) {
	return nil, errLocationNotPorted("GetProvinceByCompositeId")
}
func (r *locationPostgres) GetNearestProvince(context.Context, float64, float64, float64) (*entity.Province, float64, error) {
	return nil, 0, errLocationNotPorted("GetNearestProvince")
}
func (r *locationPostgres) SearchProvinces(context.Context, string, int) ([]*entity.Province, error) {
	return nil, errLocationNotPorted("SearchProvinces")
}
func (r *locationPostgres) AdminListProvinces(context.Context, *dto.AdminListProvincesRequest) ([]*entity.Province, int64, error) {
	return nil, 0, errLocationNotPorted("AdminListProvinces")
}
func (r *locationPostgres) ListAdministrativeUnitRefs(context.Context, *dto.ListAdministrativeUnitsRequest) ([]dto.AdministrativeUnitRef, int64, error) {
	return nil, 0, errLocationNotPorted("ListAdministrativeUnitRefs")
}
func (r *locationPostgres) GetProvincesByIDs(context.Context, []uint64) ([]*entity.Province, error) {
	return nil, errLocationNotPorted("GetProvincesByIDs")
}
func (r *locationPostgres) GetWardsByIDs(context.Context, []uint64) ([]*entity.Ward, error) {
	return nil, errLocationNotPorted("GetWardsByIDs")
}
func (r *locationPostgres) FollowAdministrativeUnit(context.Context, uint64, uint64) error {
	return errLocationNotPorted("FollowAdministrativeUnit")
}
func (r *locationPostgres) UnfollowAdministrativeUnit(context.Context, uint64, uint64) error {
	return errLocationNotPorted("UnfollowAdministrativeUnit")
}
func (r *locationPostgres) ListFollowedAdministrativeUnitIDs(context.Context, uint64, int, int) ([]uint64, int64, error) {
	return nil, 0, errLocationNotPorted("ListFollowedAdministrativeUnitIDs")
}
func (r *locationPostgres) CreateProvince(context.Context, *entity.Province) error {
	return errLocationNotPorted("CreateProvince")
}
func (r *locationPostgres) UpdateProvince(context.Context, *entity.Province) error {
	return errLocationNotPorted("UpdateProvince")
}
func (r *locationPostgres) SoftDeleteProvince(context.Context, uint64) error {
	return errLocationNotPorted("SoftDeleteProvince")
}
func (r *locationPostgres) ExistsProvinceCode(context.Context, string, uint64) (bool, error) {
	return false, errLocationNotPorted("ExistsProvinceCode")
}
func (r *locationPostgres) GetWardByID(context.Context, uint64) (*entity.Ward, error) {
	return nil, errLocationNotPorted("GetWardByID")
}
func (r *locationPostgres) GetWardDetail(context.Context, uint64) (*entity.Ward, error) {
	return nil, errLocationNotPorted("GetWardDetail")
}
func (r *locationPostgres) GetWardByCode(context.Context, string) (*entity.Ward, error) {
	return nil, errLocationNotPorted("GetWardByCode")
}
func (r *locationPostgres) GetWardByCompositeId(context.Context, string) (*entity.Ward, error) {
	return nil, errLocationNotPorted("GetWardByCompositeId")
}
func (r *locationPostgres) GetNearestWard(context.Context, float64, float64, uint64, float64) (*entity.Ward, float64, error) {
	return nil, 0, errLocationNotPorted("GetNearestWard")
}
func (r *locationPostgres) SearchWards(context.Context, string, uint64, int) ([]*entity.Ward, error) {
	return nil, errLocationNotPorted("SearchWards")
}
func (r *locationPostgres) ListWardsByProvince(context.Context, uint64, int, int) ([]*entity.Ward, int64, error) {
	return nil, 0, errLocationNotPorted("ListWardsByProvince")
}
func (r *locationPostgres) AdminListWards(context.Context, *dto.AdminListWardsRequest) ([]*entity.Ward, int64, error) {
	return nil, 0, errLocationNotPorted("AdminListWards")
}
func (r *locationPostgres) CreateWard(context.Context, *entity.Ward) error {
	return errLocationNotPorted("CreateWard")
}
func (r *locationPostgres) UpdateWard(context.Context, *entity.Ward) error {
	return errLocationNotPorted("UpdateWard")
}
func (r *locationPostgres) SoftDeleteWard(context.Context, uint64) error {
	return errLocationNotPorted("SoftDeleteWard")
}
func (r *locationPostgres) SoftDeleteWardsByProvinceID(context.Context, uint64) error {
	return errLocationNotPorted("SoftDeleteWardsByProvinceID")
}
func (r *locationPostgres) ExistsWardCode(context.Context, string, uint64) (bool, error) {
	return false, errLocationNotPorted("ExistsWardCode")
}
func (r *locationPostgres) CreateWardNeighbor(context.Context, uint64, uint64) error {
	return errLocationNotPorted("CreateWardNeighbor")
}
func (r *locationPostgres) DeleteWardNeighbor(context.Context, uint64, uint64) error {
	return errLocationNotPorted("DeleteWardNeighbor")
}
func (r *locationPostgres) GetWardNeighbors(context.Context, uint64) ([]*entity.Ward, error) {
	return nil, errLocationNotPorted("GetWardNeighbors")
}
func (r *locationPostgres) GetNearestLocations(context.Context, []*dto.Coordinate, float64) ([]*dto.LocationResult, error) {
	return nil, errLocationNotPorted("GetNearestLocations")
}
func (r *locationPostgres) SearchTxtClient(context.Context, string, int) ([]*entity.SearchItem, error) {
	return nil, errLocationNotPorted("SearchTxtClient")
}
func (r *locationPostgres) GetAdministrativeUnitSeoSource(context.Context, uint64) (*dto.AdministrativeUnitSeoSource, error) {
	return nil, errLocationNotPorted("GetAdministrativeUnitSeoSource")
}
func (r *locationPostgres) GetAdministrativeUnitSeoSourcesByIdentities(context.Context, []uint64) ([]dto.AdministrativeUnitSeoSource, error) {
	return nil, errLocationNotPorted("GetAdministrativeUnitSeoSourcesByIdentities")
}
func (r *locationPostgres) UpdateAdministrativeUnitSeoID(context.Context, uint64, uint64) error {
	return errLocationNotPorted("UpdateAdministrativeUnitSeoID")
}
func (r *locationPostgres) SyncProvinceAreaFromGeometry(context.Context, uint64) (bool, error) {
	return false, errLocationNotPorted("SyncProvinceAreaFromGeometry")
}
func (r *locationPostgres) SyncWardAreaFromGeometry(context.Context, uint64) (bool, error) {
	return false, errLocationNotPorted("SyncWardAreaFromGeometry")
}
func (r *locationPostgres) SyncAllProvincesAreaFromGeometry(context.Context) (int64, error) {
	return 0, errLocationNotPorted("SyncAllProvincesAreaFromGeometry")
}
func (r *locationPostgres) SyncAllWardsAreaFromGeometry(context.Context) (int64, error) {
	return 0, errLocationNotPorted("SyncAllWardsAreaFromGeometry")
}
