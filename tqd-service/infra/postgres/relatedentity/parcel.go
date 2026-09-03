package relatedentity

import (
	"context"
	"fmt"

	discoverydomain "tqd/internal/domain/discovery/model"
	relateddomain "tqd/internal/domain/relatedentity/model"
)

func (r *Repository) FindParcels(ctx context.Context, query relateddomain.Query) ([]relateddomain.Candidate, error) {
	switch query.Primary.Kind {
	case discoverydomain.EntityKindParcel:
		return r.findParcelsAroundParcel(ctx, query)
	case discoverydomain.EntityKindPlanningRegion:
		return r.findParcelsInsidePlanningRegion(ctx, query)
	default:
		return r.findParcelsAroundPoint(ctx, query)
	}
}

func parcelGeometryPreview(include bool) string {
	if !include {
		return "NULL::text AS geo_json"
	}
	return "ST_AsGeoJSON(ST_SimplifyPreserveTopology(parcel.geometry, 0.000002)) AS geo_json"
}

func (r *Repository) findParcelsAroundParcel(ctx context.Context, query relateddomain.Query) ([]relateddomain.Candidate, error) {
	// The Parcel provider represents adjacency. Nearby exploration is a
	// different intent and must not make this core relation pay geography
	// buffer/distance costs.
	sql := fmt.Sprintf(`
		WITH primary_entity AS (
			SELECT parcel.geometry
			FROM parcels parcel
			WHERE parcel.id::text = ?
			  AND parcel.deleted_at IS NULL
			  AND parcel.geometry IS NOT NULL
			LIMIT 1
		)
		%s
		FROM parcels parcel
		CROSS JOIN primary_entity
		LEFT JOIN qh_parcel_info info ON info.parcel_id = parcel.id
		LEFT JOIN province_v2 province ON province.code = info.province_code
		LEFT JOIN ward_v2 ward ON ward.code = info.ward_code
		WHERE parcel.deleted_at IS NULL
		  AND parcel.geometry IS NOT NULL
		  AND parcel.id::text <> ?
		  AND parcel.geometry && primary_entity.geometry
		  AND ST_Touches(parcel.geometry, primary_entity.geometry)
		ORDER BY parcel.id DESC
		LIMIT ?`, parcelSelectColumns(
		parcelGeometryPreview(query.IncludeGeometryPreview),
		`'adjacent'`,
		`0::float8`,
		`68`,
	))

	return r.scanParcels(ctx, "find parcels around parcel", sql,
		query.Primary.ID,
		query.Primary.ID,
		query.Limit,
	)
}

func (r *Repository) findParcelsInsidePlanningRegion(ctx context.Context, query relateddomain.Query) ([]relateddomain.Candidate, error) {
	// Small planning regions may expose contained parcels, but the relation must
	// represent the region geometry itself—not an undocumented centroid sample.
	sql := fmt.Sprintf(`
		WITH primary_entity AS (
			SELECT region.geometry
			FROM qh_regions region
			WHERE region.id::text = ?
			  AND region.deleted_at IS NULL
			  AND region.status = 10
			  AND region.is_latest = true
			  AND region.geometry IS NOT NULL
			LIMIT 1
		)
		%s
		FROM parcels parcel
		CROSS JOIN primary_entity
		LEFT JOIN qh_parcel_info info ON info.parcel_id = parcel.id
		LEFT JOIN province_v2 province ON province.code = info.province_code
		LEFT JOIN ward_v2 ward ON ward.code = info.ward_code
		WHERE parcel.deleted_at IS NULL
		  AND parcel.geometry IS NOT NULL
		  AND parcel.geometry && primary_entity.geometry
		  AND ST_Covers(primary_entity.geometry, parcel.geometry)
		ORDER BY parcel.id DESC
		LIMIT ?`, parcelSelectColumns(
		parcelGeometryPreview(query.IncludeGeometryPreview),
		`'contained_by_primary'`,
		`0::float8`,
		`76`,
	))

	return r.scanParcels(ctx, "find parcels inside planning region", sql,
		query.Primary.ID,
		query.Limit,
	)
}

func (r *Repository) findParcelsAroundPoint(ctx context.Context, query relateddomain.Query) ([]relateddomain.Candidate, error) {
	pointCanBeContained := query.Primary.Kind == discoverydomain.EntityKindPOI
	sql := fmt.Sprintf(`
		WITH focus AS (
			SELECT point,
			       ST_Envelope(ST_Buffer(point::geography, ?)::geometry) AS search_envelope
			FROM (SELECT ST_SetSRID(ST_Point(?, ?), 4326) AS point) input
		)
		%s
		FROM parcels parcel
		CROSS JOIN focus
		LEFT JOIN qh_parcel_info info ON info.parcel_id = parcel.id
		LEFT JOIN province_v2 province ON province.code = info.province_code
		LEFT JOIN ward_v2 ward ON ward.code = info.ward_code
		WHERE parcel.deleted_at IS NULL
		  AND parcel.geometry IS NOT NULL
		  AND parcel.geometry && focus.search_envelope
		  AND (
		    ST_Covers(parcel.geometry, focus.point)
		    OR ST_DWithin(parcel.geometry::geography, focus.point::geography, ?)
		  )
		ORDER BY score DESC, distance_meters ASC, parcel.id DESC
		LIMIT ?`, parcelSelectColumns(
		parcelGeometryPreview(query.IncludeGeometryPreview),
		`CASE WHEN ? AND ST_Covers(parcel.geometry, focus.point) THEN 'contains_primary' ELSE 'nearby' END`,
		`ST_Distance(parcel.geometry::geography, focus.point::geography)`,
		`CASE WHEN ? AND ST_Covers(parcel.geometry, focus.point) THEN 82 ELSE 45 END`,
	))

	return r.scanParcels(ctx, "find parcels around point", sql,
		query.RadiusMeters,
		query.Focus.Longitude,
		query.Focus.Latitude,
		pointCanBeContained,
		pointCanBeContained,
		query.RadiusMeters,
		query.Limit,
	)
}

func parcelSelectColumns(geometryPreview, relationship, distance, score string) string {
	return fmt.Sprintf(`SELECT parcel.id::text AS id,
	       CONCAT(
	         'Thửa ', COALESCE(NULLIF(info.land_number, ''), NULLIF(parcel.land_number, ''), parcel.id::text),
	         CASE
	           WHEN COALESCE(NULLIF(info.map_number, ''), NULLIF(parcel.map_number, ''), '') <> ''
	           THEN CONCAT(', tờ ', COALESCE(NULLIF(info.map_number, ''), NULLIF(parcel.map_number, '')))
	           ELSE ''
	         END
	       ) AS title,
	       COALESCE(NULLIF(info.address_text, ''), NULLIF(parcel.address_text, ''), '') AS subtitle,
	       '' AS description,
	       COALESCE(NULLIF(info.address_text, ''), NULLIF(parcel.address_text, ''), '') AS address,
	       COALESCE(province.full_name, '') AS province,
	       COALESCE(info.province_code, '') AS province_code,
	       COALESCE(ward.full_name, '') AS ward,
	       COALESCE(info.ward_code, '') AS ward_code,
	       GeometryType(parcel.geometry) AS geometry_type,
	       ST_Y(ST_PointOnSurface(parcel.geometry)) AS center_lat,
	       ST_X(ST_PointOnSurface(parcel.geometry)) AS center_lon,
	       ST_XMin(Box2D(parcel.geometry)) AS min_lon,
	       ST_YMin(Box2D(parcel.geometry)) AS min_lat,
	       ST_XMax(Box2D(parcel.geometry)) AS max_lon,
	       ST_YMax(Box2D(parcel.geometry)) AS max_lat,
	       %s,
	       COALESCE(NULLIF(info.map_number, ''), NULLIF(parcel.map_number, ''), '') AS map_number,
	       COALESCE(NULLIF(info.land_number, ''), NULLIF(parcel.land_number, ''), '') AS land_number,
	       COALESCE(NULLIF(info.total_area_sqm, 0), parcel.total_area_sqm, ST_Area(parcel.geometry::geography), 0) AS area_sqm,
	       %s AS relationship,
	       %s AS distance_meters,
	       (%s)::float8 AS score`,
		geometryPreview,
		relationship,
		distance,
		score,
	)
}

func (r *Repository) scanParcels(ctx context.Context, operation, sql string, args ...any) ([]relateddomain.Candidate, error) {
	var rows []spatialRow
	if err := r.db.WithContext(ctx).Raw(sql, args...).Scan(&rows).Error; err != nil {
		return nil, fmt.Errorf("%s: %w", operation, err)
	}
	return rowsToCandidates(discoverydomain.EntityKindParcel, rows), nil
}
