package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	qh_domain "tqd/internal/domain/qh"
	"tqd/internal/interface/repo"

	"gorm.io/gorm"
)

type oneHouseRepo struct {
	db *gorm.DB
}

func NewOneHouseRepo(db *gorm.DB) repo.OneHouseRepo {
	return &oneHouseRepo{db: db}
}

func (r *oneHouseRepo) GetByID(ctx context.Context, id uint64) (*qh_domain.OneHouse, error) {
	var result qh_domain.OneHouse
	query := `
		SELECT 
			id, lat, lng, address, district, district_code, city, city_code,
			province, province_code, ward, ward_code, map_number, land_number,
			latitude, longitude, direction, shape, number_of_facades, total_area,
			planning_land_type, planning_land_type_code, planning_infos, shape_file_id,
			is_saved_by_user, property_code, property_uuid, accepted_tnc,
			ST_AsGeoJSON(geom) as geojson
		FROM one_house
		WHERE id = $1
	`
	row := r.db.WithContext(ctx).Raw(query, id).Row()
	err := r.scanOneHouse(row, &result)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("GetByID failed: %w", err)
	}
	return &result, nil
}

func (r *oneHouseRepo) GetByPropertyUUID(ctx context.Context, uuid string) (*qh_domain.OneHouse, error) {
	if uuid == "" {
		return nil, nil
	}
	var result qh_domain.OneHouse
	query := `
		SELECT 
			id, lat, lng, address, district, district_code, city, city_code,
			province, province_code, ward, ward_code, map_number, land_number,
			latitude, longitude, direction, shape, number_of_facades, total_area,
			planning_land_type, planning_land_type_code, planning_infos, shape_file_id,
			is_saved_by_user, property_code, property_uuid, accepted_tnc,
			ST_AsGeoJSON(geom) as geojson
		FROM one_house
		WHERE property_uuid = $1
	`
	row := r.db.WithContext(ctx).Raw(query, uuid).Row()
	err := r.scanOneHouse(row, &result)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("GetByPropertyUUID failed: %w", err)
	}
	return &result, nil
}

func (r *oneHouseRepo) scanOneHouse(row *sql.Row, oh *qh_domain.OneHouse) error {
	var geojsonStr string
	err := row.Scan(
		&oh.ID, &oh.Lat, &oh.Lng, &oh.Address, &oh.District, &oh.DistrictCode,
		&oh.City, &oh.CityCode, &oh.Province, &oh.ProvinceCode, &oh.Ward, &oh.WardCode,
		&oh.MapNumber, &oh.LandNumber, &oh.Latitude, &oh.Longitude, &oh.Direction,
		&oh.Shape, &oh.NumberOfFacades, &oh.TotalArea, &oh.PlanningLandType,
		&oh.PlanningLandTypeCode, &oh.PlanningInfos, &oh.ShapeFileID,
		&oh.IsSavedByUser, &oh.PropertyCode, &oh.PropertyUUID, &oh.AcceptedTnc,
		&geojsonStr,
	)
	if err != nil {
		return err
	}
	oh.GeomGeoJSON = geojsonStr
	return nil
}
