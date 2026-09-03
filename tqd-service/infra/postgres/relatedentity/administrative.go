package relatedentity

import (
	"context"
	"fmt"

	discoverydomain "tqd/internal/domain/discovery/model"
	relateddomain "tqd/internal/domain/relatedentity/model"
)

func (r *Repository) FindAdministrativeUnits(ctx context.Context, query relateddomain.Query) ([]relateddomain.Candidate, error) {
	// Administrative boundaries are not available in the current schema. Exact
	// ward/province codes from the canonical primary are authoritative. Nearest
	// centre fallbacks remain radius-bounded and are explicitly marked approximate.
	sql := `
		WITH focus AS (
			SELECT point,
			       ST_Envelope(ST_Buffer(point::geography, ?)::geometry) AS search_envelope
			FROM (SELECT ST_SetSRID(ST_Point(?, ?), 4326) AS point) input
		), candidates AS (
			SELECT CONCAT('ward:', ward.id) AS id,
			       ward.full_name AS title,
			       COALESCE(province.full_name, '') AS subtitle,
			       '' AS description,
			       '' AS address,
			       COALESCE(province.full_name, '') AS province,
			       COALESCE(province.code, '') AS province_code,
			       ward.full_name AS ward,
			       COALESCE(ward.code, '') AS ward_code,
			       'Point' AS geometry_type,
			       ward.lat AS center_lat,
			       ward.lng AS center_lon,
			       0::float8 AS min_lon,
			       0::float8 AS min_lat,
			       0::float8 AS max_lon,
			       0::float8 AS max_lat,
			       NULL::text AS geo_json,
			       CASE WHEN ward.code = ? THEN 'administrative_context' ELSE 'nearby' END AS relationship,
			       ST_DistanceSphere(ST_SetSRID(ST_Point(ward.lng, ward.lat), 4326), focus.point) AS distance_meters,
			       CASE WHEN ward.code = ? THEN 95 ELSE 48 END::float8 AS score
			FROM ward_v2 ward
			LEFT JOIN province_v2 province ON province.id = ward.province_id
			CROSS JOIN focus
			WHERE ward.lat IS NOT NULL
			  AND ward.lng IS NOT NULL
			  AND (
			    (? <> '' AND ward.code = ?)
			    OR (
			      ward.lng BETWEEN ST_XMin(focus.search_envelope) AND ST_XMax(focus.search_envelope)
			      AND ward.lat BETWEEN ST_YMin(focus.search_envelope) AND ST_YMax(focus.search_envelope)
			      AND ST_DWithin(
			        ST_SetSRID(ST_Point(ward.lng, ward.lat), 4326)::geography,
			        focus.point::geography,
			        ?
			      )
			    )
			  )
			  AND NOT (? = 'administrative_unit' AND CONCAT('ward:', ward.id) = ?)

			UNION ALL

			SELECT CONCAT('province:', province.id) AS id,
			       province.full_name AS title,
			       '' AS subtitle,
			       '' AS description,
			       '' AS address,
			       province.full_name AS province,
			       COALESCE(province.code, '') AS province_code,
			       '' AS ward,
			       '' AS ward_code,
			       'Point' AS geometry_type,
			       province.lat AS center_lat,
			       province.lng AS center_lon,
			       0::float8 AS min_lon,
			       0::float8 AS min_lat,
			       0::float8 AS max_lon,
			       0::float8 AS max_lat,
			       NULL::text AS geo_json,
			       CASE WHEN province.code = ? THEN 'administrative_context' ELSE 'nearby' END AS relationship,
			       ST_DistanceSphere(ST_SetSRID(ST_Point(province.lng, province.lat), 4326), focus.point) AS distance_meters,
			       CASE WHEN province.code = ? THEN 90 ELSE 42 END::float8 AS score
			FROM province_v2 province
			CROSS JOIN focus
			WHERE province.lat IS NOT NULL
			  AND province.lng IS NOT NULL
			  AND (
			    (? <> '' AND province.code = ?)
			    OR (
			      province.lng BETWEEN ST_XMin(focus.search_envelope) AND ST_XMax(focus.search_envelope)
			      AND province.lat BETWEEN ST_YMin(focus.search_envelope) AND ST_YMax(focus.search_envelope)
			      AND ST_DWithin(
			        ST_SetSRID(ST_Point(province.lng, province.lat), 4326)::geography,
			        focus.point::geography,
			        ?
			      )
			    )
			  )
			  AND NOT (? = 'administrative_unit' AND CONCAT('province:', province.id) = ?)
		)
		SELECT *
		FROM candidates
		ORDER BY score DESC, distance_meters ASC, id ASC
		LIMIT ?`

	args := []any{
		query.RadiusMeters,
		query.Focus.Longitude,
		query.Focus.Latitude,

		query.WardCode,
		query.WardCode,
		query.WardCode,
		query.WardCode,
		query.RadiusMeters,
		string(query.Primary.Kind),
		query.Primary.ID,

		query.ProvinceCode,
		query.ProvinceCode,
		query.ProvinceCode,
		query.ProvinceCode,
		query.RadiusMeters,
		string(query.Primary.Kind),
		query.Primary.ID,

		query.Limit,
	}

	var rows []spatialRow
	if err := r.db.WithContext(ctx).Raw(sql, args...).Scan(&rows).Error; err != nil {
		return nil, fmt.Errorf("find related administrative units: %w", err)
	}
	out := rowsToCandidates(discoverydomain.EntityKindAdministrativeUnit, rows)
	for index := range out {
		if out[index].Relationship == relateddomain.RelationshipNearby {
			out[index].Entity.Source.DataQuality = "approximate"
			out[index].ReasonCodes = append(out[index].ReasonCodes, "administrative_centroid_only")
		}
	}
	return out, nil
}
