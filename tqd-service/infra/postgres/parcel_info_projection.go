package postgres

// Parcel geometry/identity belongs to parcels. Descriptive attributes belong
// to qh_parcel_info. All overview/detail reads use this projection; adding
// shadow columns to parcels would reintroduce conflicting metadata owners.
func parcelInfoSelect(includeGeometry bool) string {
	geometry := ""
	if includeGeometry {
		geometry = ", ST_AsGeoJSON(p.geometry) AS geometry"
	}
	return `SELECT p.id AS parcel_id,
		COALESCE(NULLIF(info.lat, 0), NULLIF(p.lat, 0), ST_Y(ST_PointOnSurface(p.geometry)), 0) AS lat,
		COALESCE(NULLIF(info.lon, 0), NULLIF(p.lon, 0), ST_X(ST_PointOnSurface(p.geometry)), 0) AS lon,
		COALESCE(NULLIF(info.total_area_sqm, 0), ST_Area(p.geometry::geography), 0) AS area_sqm,
		COALESCE(info.address_text, '') AS address_text,
		COALESCE(info.map_number, '') AS map_number,
		COALESCE(info.land_number, '') AS land_number,
		COALESCE(info.property_code, '') AS property_code,
		COALESCE(info.property_uuid::text, '') AS property_uuid,
		'' AS shape_type, '' AS direction,
		COALESCE(pu.planning_type_code, '') AS land_use_code,
		COALESCE(pu.planning_type_name, 'Chưa phân loại') AS land_use_name,
		'#CCCCCC' AS land_use_color,
		COALESCE(prov.full_name, '') AS province,
		COALESCE(info.province_code, '') AS province_code,
		COALESCE(info.ward_code, '') AS ward_code,
		COALESCE(p.seo_id, 0) AS seo_id,
		COALESCE(p.is_seo, false) AS is_seo` + geometry + `
		FROM parcels p
		LEFT JOIN LATERAL (
			SELECT * FROM qh_parcel_info WHERE parcel_id = p.id
			ORDER BY updated_at DESC, id DESC LIMIT 1
		) info ON true
		LEFT JOIN qh_planning_land_use pu ON info.land_type_id = pu.id
		LEFT JOIN province_v2 prov ON info.province_code = prov.code`
}
