package relatedentity

import (
	"context"
	"fmt"

	discoverydomain "tqd/internal/domain/discovery/model"
	relateddomain "tqd/internal/domain/relatedentity/model"
)

func (r *Repository) FindPOIs(ctx context.Context, query relateddomain.Query) ([]relateddomain.Candidate, error) {
	switch query.Primary.Kind {
	case discoverydomain.EntityKindParcel:
		return r.findPOIsForPrimaryArea(ctx, query, `
			SELECT parcel.geometry
			FROM parcels parcel
			WHERE parcel.id::text = ?
			  AND parcel.deleted_at IS NULL
			  AND parcel.geometry IS NOT NULL
			LIMIT 1`)
	case discoverydomain.EntityKindPlanningRegion:
		return r.findPOIsForPrimaryArea(ctx, query, `
			SELECT region.geometry
			FROM qh_regions region
			WHERE region.id::text = ?
			  AND region.deleted_at IS NULL
			  AND region.status = 10
			  AND region.is_latest = true
			  AND region.geometry IS NOT NULL
			LIMIT 1`)
	default:
		return r.findPOIsAroundPoint(ctx, query)
	}
}

func (r *Repository) findPOIsForPrimaryArea(
	ctx context.Context,
	query relateddomain.Query,
	primaryAreaSQL string,
) ([]relateddomain.Candidate, error) {
	sql := fmt.Sprintf(`
		WITH focus AS (
			SELECT point,
			       ST_Envelope(ST_Buffer(point::geography, ?)::geometry) AS search_envelope
			FROM (SELECT ST_SetSRID(ST_Point(?, ?), 4326) AS point) input
		), primary_entity AS (
			%s
		), candidates AS (
			%s
			FROM pois poi
			CROSS JOIN focus
			CROSS JOIN primary_entity
			WHERE poi.deleted_at IS NULL
			  AND poi.is_active = true
			  AND poi.longitude BETWEEN ST_XMin(focus.search_envelope) AND ST_XMax(focus.search_envelope)
			  AND poi.latitude BETWEEN ST_YMin(focus.search_envelope) AND ST_YMax(focus.search_envelope)
			  AND ST_DWithin(
			    ST_SetSRID(ST_Point(poi.longitude, poi.latitude), 4326)::geography,
			    focus.point::geography,
			    ?
			  )
		)
		SELECT *
		FROM candidates
		ORDER BY score DESC, distance_meters ASC, id DESC
		LIMIT ?`, primaryAreaSQL, poiSelectColumns(
		`CASE
		   WHEN ST_Covers(primary_entity.geometry, ST_SetSRID(ST_Point(poi.longitude, poi.latitude), 4326))
		   THEN 'contained_by_primary'
		   ELSE 'nearby'
		 END`,
		`CASE
		   WHEN ST_Covers(primary_entity.geometry, ST_SetSRID(ST_Point(poi.longitude, poi.latitude), 4326))
		   THEN 62
		   WHEN poi.is_featured THEN 42
		   ELSE 35
		 END`,
	))

	return r.scanPOIs(ctx, "find POIs for primary area", sql,
		query.RadiusMeters,
		query.Focus.Longitude,
		query.Focus.Latitude,
		query.Primary.ID,
		query.RadiusMeters,
		query.Limit,
	)
}

func (r *Repository) findPOIsAroundPoint(ctx context.Context, query relateddomain.Query) ([]relateddomain.Candidate, error) {
	sql := fmt.Sprintf(`
		WITH focus AS (
			SELECT point,
			       ST_Envelope(ST_Buffer(point::geography, ?)::geometry) AS search_envelope
			FROM (SELECT ST_SetSRID(ST_Point(?, ?), 4326) AS point) input
		)
		%s
		FROM pois poi
		CROSS JOIN focus
		WHERE poi.deleted_at IS NULL
		  AND poi.is_active = true
		  AND poi.longitude BETWEEN ST_XMin(focus.search_envelope) AND ST_XMax(focus.search_envelope)
		  AND poi.latitude BETWEEN ST_YMin(focus.search_envelope) AND ST_YMax(focus.search_envelope)
		  AND ST_DWithin(
		    ST_SetSRID(ST_Point(poi.longitude, poi.latitude), 4326)::geography,
		    focus.point::geography,
		    ?
		  )
		  AND NOT (? = 'poi' AND poi.id::text = ?)
		ORDER BY score DESC, distance_meters ASC, poi.rating DESC, poi.id DESC
		LIMIT ?`, poiSelectColumns(`'nearby'`, `CASE WHEN poi.is_featured THEN 42 ELSE 35 END`))

	return r.scanPOIs(ctx, "find nearby POIs", sql,
		query.RadiusMeters,
		query.Focus.Longitude,
		query.Focus.Latitude,
		query.RadiusMeters,
		string(query.Primary.Kind),
		query.Primary.ID,
		query.Limit,
	)
}

func poiSelectColumns(relationship, score string) string {
	return fmt.Sprintf(`SELECT poi.id::text AS id,
	       poi.name AS title,
	       COALESCE(poi.address, '') AS subtitle,
	       COALESCE(poi.description, '') AS description,
	       COALESCE(poi.address, '') AS address,
	       '' AS province,
	       '' AS province_code,
	       '' AS ward,
	       '' AS ward_code,
	       'Point' AS geometry_type,
	       poi.latitude AS center_lat,
	       poi.longitude AS center_lon,
	       0::float8 AS min_lon,
	       0::float8 AS min_lat,
	       0::float8 AS max_lon,
	       0::float8 AS max_lat,
	       NULL::text AS geo_json,
	       %s AS relationship,
	       ST_DistanceSphere(
	         ST_SetSRID(ST_Point(poi.longitude, poi.latitude), 4326),
	         focus.point
	       ) AS distance_meters,
	       (%s)::float8 AS score`, relationship, score)
}

func (r *Repository) scanPOIs(ctx context.Context, operation, sql string, args ...any) ([]relateddomain.Candidate, error) {
	var rows []spatialRow
	if err := r.db.WithContext(ctx).Raw(sql, args...).Scan(&rows).Error; err != nil {
		return nil, fmt.Errorf("%s: %w", operation, err)
	}
	return rowsToCandidates(discoverydomain.EntityKindPOI, rows), nil
}
