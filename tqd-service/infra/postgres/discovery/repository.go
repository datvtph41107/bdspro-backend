package discovery

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"tqd/infra/postgres/entitycandidate"
	"tqd/internal/domain/discovery/model"
	"tqd/internal/usecase/discovery/ports"

	"gorm.io/gorm"
)

type Repository struct{ db *gorm.DB }

func NewDiscoveryRepository(db *gorm.DB) ports.Repository {
	return &Repository{db: db}
}

type spatialRow = entitycandidate.Row

func candidate(kind domain.EntityKind, row spatialRow, matchType string) domain.EntityCandidate {
	return entitycandidate.Build(kind, row, entitycandidate.BuildOptions{
		MatchType: matchType,
	})
}

// detailURL remains as a package-level compatibility seam for existing tests
// and repository code. Canonical URL/capability rules are owned by the shared
// entitycandidate projection package.
func detailURL(kind domain.EntityKind, row spatialRow) string {
	return entitycandidate.DetailURL(kind, row)
}

func limit(reqLimit int) int {
	if reqLimit <= 0 {
		return 20
	}
	if reqLimit > 50 {
		return 50
	}
	return reqLimit
}

func (r *Repository) IdentifyParcels(ctx context.Context, req domain.IdentifyRequest) ([]domain.EntityCandidate, error) {
	var rows []spatialRow
	geometry := "NULL::text AS geo_json"
	if req.Options.IncludeGeometryPreview {
		geometry = "ST_AsGeoJSON(ST_SimplifyPreserveTopology(p.geometry, 0.000002)) AS geo_json"
	}
	query := fmt.Sprintf(`
		SELECT p.id::text AS id,
		       CONCAT(
		         'Thửa ', COALESCE(NULLIF(pi.land_number,''), p.id::text),
		         CASE WHEN COALESCE(pi.map_number,'') <> '' THEN CONCAT(', tờ ', pi.map_number) ELSE '' END
		       ) AS title,
		       COALESCE(pi.address_text, '') AS subtitle,
		       '' AS description,
		       COALESCE(pi.address_text, '') AS address,
		       COALESCE(prov.full_name, '') AS province,
		       COALESCE(pi.province_code, '') AS province_code,
		       COALESCE(w.full_name, '') AS ward,
		       COALESCE(pi.ward_code, '') AS ward_code,
		       GeometryType(p.geometry) AS geometry_type,
		       ST_Y(ST_PointOnSurface(p.geometry)) AS center_lat,
		       ST_X(ST_PointOnSurface(p.geometry)) AS center_lon,
		       ST_XMin(Box2D(p.geometry)) AS min_lon,
		       ST_YMin(Box2D(p.geometry)) AS min_lat,
		       ST_XMax(Box2D(p.geometry)) AS max_lon,
		       ST_YMax(Box2D(p.geometry)) AS max_lat,
		       %s,
		       COALESCE(pi.map_number,'') AS map_number,
		       COALESCE(pi.land_number,'') AS land_number,
		       COALESCE(NULLIF(pi.total_area_sqm,0), ST_Area(p.geometry::geography), 0) AS area_sqm,
		       30::float8 AS score
		FROM parcels p
		LEFT JOIN qh_parcel_info pi ON pi.parcel_id = p.id
		LEFT JOIN province_v2 prov ON prov.code = pi.province_code
		LEFT JOIN ward_v2 w ON w.code = pi.ward_code
		WHERE p.deleted_at IS NULL
		  AND p.geometry IS NOT NULL
		  AND p.geometry && ST_SetSRID(ST_Point(?, ?), 4326)
		  AND ST_Covers(p.geometry, ST_SetSRID(ST_Point(?, ?), 4326))
		ORDER BY COALESCE(NULLIF(pi.total_area_sqm,0), ST_Area(p.geometry::geography), 0) ASC, p.id DESC
		LIMIT ?`, geometry)
	if err := r.db.WithContext(ctx).Raw(
		query,
		req.Point.Longitude, req.Point.Latitude,
		req.Point.Longitude, req.Point.Latitude,
		limit(req.Options.Limit),
	).Scan(&rows).Error; err != nil {
		return nil, fmt.Errorf("identify parcels: %w", err)
	}
	out := make([]domain.EntityCandidate, 0, len(rows))
	for _, row := range rows {
		out = append(out, candidate(domain.EntityKindParcel, row, "point_containment"))
	}
	return out, nil
}

func (r *Repository) IdentifyRegions(ctx context.Context, req domain.IdentifyRequest) ([]domain.EntityCandidate, error) {
	var rows []spatialRow
	geometry := "NULL::text AS geo_json"
	if req.Options.IncludeGeometryPreview {
		geometry = "ST_AsGeoJSON(ST_SimplifyPreserveTopology(r.geometry, 0.00001)) AS geo_json"
	}
	query := fmt.Sprintf(`
        SELECT r.id::text AS id,
               COALESCE(NULLIF(r.display_name,''), NULLIF(r.name,''), CONCAT('Vùng quy hoạch #', r.id)) AS title,
               COALESCE(NULLIF(l.display_name,''), l.name, '') AS subtitle,
               COALESCE(r.description,'') AS description,
               '' AS address, '' AS province, '' AS province_code, '' AS ward, '' AS ward_code,
               GeometryType(r.geometry) AS geometry_type,
               ST_Y(ST_PointOnSurface(r.geometry)) AS center_lat, ST_X(ST_PointOnSurface(r.geometry)) AS center_lon,
               ST_XMin(Box2D(r.geometry)) AS min_lon, ST_YMin(Box2D(r.geometry)) AS min_lat,
               ST_XMax(Box2D(r.geometry)) AS max_lon, ST_YMax(Box2D(r.geometry)) AS max_lat,
               %s,
               r.layer_id::text AS layer_id,
               COALESCE(r.area_sqm, ST_Area(r.geometry::geography), 0) AS area_sqm,
               (25 + COALESCE(l.display_order,0))::float8 AS score
        FROM qh_regions r
        JOIN qh_layers l ON l.id = r.layer_id AND l.deleted_at IS NULL AND l.status = 10
        WHERE r.deleted_at IS NULL AND r.status = 10 AND r.is_latest = true
          AND r.geometry && ST_SetSRID(ST_Point(?, ?), 4326)
          AND ST_Covers(r.geometry, ST_SetSRID(ST_Point(?, ?), 4326))
        ORDER BY COALESCE(l.display_order,0) DESC, COALESCE(r.area_sqm, ST_Area(r.geometry::geography),0) ASC, r.id DESC
        LIMIT ?`, geometry)
	if err := r.db.WithContext(ctx).Raw(query, req.Point.Longitude, req.Point.Latitude, req.Point.Longitude, req.Point.Latitude, limit(req.Options.Limit)).Scan(&rows).Error; err != nil {
		return nil, fmt.Errorf("identify regions: %w", err)
	}
	out := make([]domain.EntityCandidate, 0, len(rows))
	for _, row := range rows {
		out = append(out, candidate(domain.EntityKindPlanningRegion, row, "point_containment"))
	}
	return out, nil
}

func (r *Repository) IdentifyAdministrativeUnits(ctx context.Context, req domain.IdentifyRequest) ([]domain.EntityCandidate, error) {
	// Current schema stores administrative centroids, not boundaries. Mark results as nearest_center, never as containment.
	var rows []spatialRow
	query := `
      (SELECT CONCAT('ward:', w.id) AS id, w.full_name AS title, COALESCE(p.full_name,'') AS subtitle,
              '' AS description, '' AS address, COALESCE(p.full_name,'') AS province, COALESCE(p.code,'') AS province_code,
              w.full_name AS ward, COALESCE(w.code,'') AS ward_code, 'Point' AS geometry_type,
              w.lat AS center_lat, w.lng AS center_lon,
              0::float8 AS min_lon,0::float8 AS min_lat,0::float8 AS max_lon,0::float8 AS max_lat,
              '' AS geo_json,
              ST_DistanceSphere(ST_SetSRID(ST_Point(w.lng,w.lat),4326), ST_SetSRID(ST_Point(?,?),4326)) AS distance_meters,
              8::float8 AS score
       FROM ward_v2 w LEFT JOIN province_v2 p ON p.id = w.province_id
       WHERE w.lat IS NOT NULL AND w.lng IS NOT NULL
       ORDER BY distance_meters ASC LIMIT 1)
      UNION ALL
      (SELECT CONCAT('province:', p.id) AS id, p.full_name AS title, '' AS subtitle,
              '' AS description, '' AS address, p.full_name AS province, COALESCE(p.code,'') AS province_code,
              '' AS ward, '' AS ward_code, 'Point' AS geometry_type,
              p.lat AS center_lat, p.lng AS center_lon,
              0::float8,0::float8,0::float8,0::float8,'' AS geo_json,
              ST_DistanceSphere(ST_SetSRID(ST_Point(p.lng,p.lat),4326), ST_SetSRID(ST_Point(?,?),4326)) AS distance_meters,
              5::float8 AS score
       FROM province_v2 p WHERE p.lat IS NOT NULL AND p.lng IS NOT NULL
       ORDER BY distance_meters ASC LIMIT 1)`
	if err := r.db.WithContext(ctx).Raw(query, req.Point.Longitude, req.Point.Latitude, req.Point.Longitude, req.Point.Latitude).Scan(&rows).Error; err != nil {
		return nil, fmt.Errorf("identify administrative units: %w", err)
	}
	out := []domain.EntityCandidate{}
	for _, row := range rows {
		c := candidate(domain.EntityKindAdministrativeUnit, row, "nearest_center")
		c.Source.DataQuality = "approximate"
		out = append(out, c)
	}
	return out, nil
}

func (r *Repository) IdentifyPOIs(ctx context.Context, req domain.IdentifyRequest) ([]domain.EntityCandidate, error) {
	var rows []spatialRow
	tolerance := req.Options.ToleranceMeters
	if tolerance <= 0 {
		tolerance = 25
	}
	query := `SELECT p.id::text AS id, p.name AS title, COALESCE(p.address,'') AS subtitle, COALESCE(p.description,'') AS description,
                     COALESCE(p.address,'') AS address, '' AS province, '' AS province_code, '' AS ward, '' AS ward_code,
                     'Point' AS geometry_type, p.latitude AS center_lat, p.longitude AS center_lon,
                     0::float8 AS min_lon,0::float8 AS min_lat,0::float8 AS max_lon,0::float8 AS max_lat,'' AS geo_json,
                     ST_DistanceSphere(ST_SetSRID(ST_Point(p.longitude,p.latitude),4326), ST_SetSRID(ST_Point(?,?),4326)) AS distance_meters,
                     CASE WHEN p.is_featured THEN 20 ELSE 10 END::float8 AS score
              FROM pois p WHERE p.deleted_at IS NULL AND p.is_active = true
                AND ST_DWithin(ST_SetSRID(ST_Point(p.longitude,p.latitude),4326)::geography, ST_SetSRID(ST_Point(?,?),4326)::geography, ?)
              ORDER BY distance_meters ASC, p.is_featured DESC, p.rating DESC LIMIT ?`
	if err := r.db.WithContext(ctx).Raw(query, req.Point.Longitude, req.Point.Latitude, req.Point.Longitude, req.Point.Latitude, tolerance, limit(req.Options.Limit)).Scan(&rows).Error; err != nil {
		return nil, fmt.Errorf("identify pois: %w", err)
	}
	out := []domain.EntityCandidate{}
	for _, row := range rows {
		out = append(out, candidate(domain.EntityKindPOI, row, "point_nearby"))
	}
	return out, nil
}

func searchPattern(q string) string { return "%" + strings.ToLower(strings.TrimSpace(q)) + "%" }

func viewportGeometryClause(req domain.SearchRequest, geometryExpression string) (string, []any) {
	if req.MapContext.Viewport == nil {
		return "", nil
	}

	viewport := req.MapContext.Viewport
	return " AND ST_Intersects(" + geometryExpression + ", ST_MakeEnvelope(?, ?, ?, ?, 4326))",
		[]any{viewport.MinLongitude, viewport.MinLatitude, viewport.MaxLongitude, viewport.MaxLatitude}
}

func viewportPointClause(req domain.SearchRequest, longitudeExpression, latitudeExpression string) (string, []any) {
	if req.MapContext.Viewport == nil {
		return "", nil
	}

	viewport := req.MapContext.Viewport
	return " AND " + longitudeExpression + " BETWEEN ? AND ? AND " + latitudeExpression + " BETWEEN ? AND ?",
		[]any{viewport.MinLongitude, viewport.MaxLongitude, viewport.MinLatitude, viewport.MaxLatitude}
}

func (r *Repository) SearchParcels(ctx context.Context, req domain.SearchRequest) ([]domain.EntityCandidate, error) {
	if strings.TrimSpace(req.Query) == "" {
		return nil, nil
	}

	var rows []spatialRow
	viewportClause, viewportArgs := viewportGeometryClause(req, "p.geometry")

	query := `
		SELECT p.id::text AS id,
		       CONCAT(
		         'Thửa ', COALESCE(NULLIF(pi.land_number,''), p.id::text),
		         CASE WHEN COALESCE(pi.map_number,'') <> '' THEN CONCAT(', tờ ', pi.map_number) ELSE '' END
		       ) AS title,
		       COALESCE(pi.address_text,'') AS subtitle,
		       '' AS description,
		       COALESCE(pi.address_text,'') AS address,
		       COALESCE(prov.full_name,'') AS province,
		       COALESCE(pi.province_code,'') AS province_code,
		       COALESCE(w.full_name,'') AS ward,
		       COALESCE(pi.ward_code,'') AS ward_code,
		       GeometryType(p.geometry) AS geometry_type,
		       ST_Y(ST_PointOnSurface(p.geometry)) AS center_lat,
		       ST_X(ST_PointOnSurface(p.geometry)) AS center_lon,
		       ST_XMin(Box2D(p.geometry)) AS min_lon,
		       ST_YMin(Box2D(p.geometry)) AS min_lat,
		       ST_XMax(Box2D(p.geometry)) AS max_lon,
		       ST_YMax(Box2D(p.geometry)) AS max_lat,
		       '' AS geo_json,
		       COALESCE(pi.map_number,'') AS map_number,
		       COALESCE(pi.land_number,'') AS land_number,
		       COALESCE(pi.total_area_sqm,0) AS area_sqm,
		       CASE
		         WHEN LOWER(COALESCE(pi.property_code,'')) = LOWER(?) THEN 100
		         WHEN LOWER(COALESCE(pi.map_number,'') || ' ' || COALESCE(pi.land_number,'')) = LOWER(?) THEN 90
		         WHEN LOWER(COALESCE(pi.adr_search,'')) LIKE ? THEN 35
		         ELSE 20
		       END::float8 AS score
		FROM parcels p
		JOIN qh_parcel_info pi ON pi.parcel_id = p.id
		LEFT JOIN province_v2 prov ON prov.code = pi.province_code
		LEFT JOIN ward_v2 w ON w.code = pi.ward_code
		WHERE p.deleted_at IS NULL
		  AND p.geometry IS NOT NULL
		  AND (
		    LOWER(COALESCE(pi.adr_search,'')) LIKE ?
		    OR LOWER(COALESCE(pi.address_text,'')) LIKE ?
		    OR LOWER(COALESCE(pi.property_code,'')) LIKE ?
		    OR LOWER(COALESCE(pi.map_number,'') || ' ' || COALESCE(pi.land_number,'')) LIKE ?
		  )`

	args := []any{
		req.Query,
		strings.ToLower(strings.TrimSpace(req.Query)),
		searchPattern(req.Query),
		searchPattern(req.Query),
		searchPattern(req.Query),
		searchPattern(req.Query),
		searchPattern(req.Query),
	}
	query += viewportClause
	args = append(args, viewportArgs...)

	if req.AdministrativeCode != "" {
		query += " AND (pi.ward_code = ? OR pi.province_code = ?)"
		args = append(args, req.AdministrativeCode, req.AdministrativeCode)
	}

	query += " ORDER BY score DESC, p.id DESC LIMIT ?"
	args = append(args, limit(req.Limit))

	if err := r.db.WithContext(ctx).Raw(query, args...).Scan(&rows).Error; err != nil {
		return nil, fmt.Errorf("search parcels: %w", err)
	}

	out := make([]domain.EntityCandidate, 0, len(rows))
	for _, row := range rows {
		out = append(out, candidate(domain.EntityKindParcel, row, "full_text"))
	}

	return out, nil
}

func (r *Repository) SearchRegions(ctx context.Context, req domain.SearchRequest) ([]domain.EntityCandidate, error) {
	if strings.TrimSpace(req.Query) == "" {
		return nil, nil
	}

	var rows []spatialRow
	viewportClause, viewportArgs := viewportGeometryClause(req, "r.geometry")

	query := `
		SELECT r.id::text AS id,
		       COALESCE(NULLIF(r.display_name,''),NULLIF(r.name,''),CONCAT('Vùng quy hoạch #',r.id)) AS title,
		       COALESCE(NULLIF(l.display_name,''),l.name,'') AS subtitle,
		       COALESCE(r.description,'') AS description,
		       '' AS address,
		       '' AS province,
		       '' AS province_code,
		       '' AS ward,
		       '' AS ward_code,
		       GeometryType(r.geometry) AS geometry_type,
		       ST_Y(ST_PointOnSurface(r.geometry)) AS center_lat,
		       ST_X(ST_PointOnSurface(r.geometry)) AS center_lon,
		       ST_XMin(Box2D(r.geometry)) AS min_lon,
		       ST_YMin(Box2D(r.geometry)) AS min_lat,
		       ST_XMax(Box2D(r.geometry)) AS max_lon,
		       ST_YMax(Box2D(r.geometry)) AS max_lat,
		       '' AS geo_json,
		       r.layer_id::text AS layer_id,
		       COALESCE(r.area_sqm,0) AS area_sqm,
		       CASE
		         WHEN LOWER(COALESCE(r.display_name,r.name,'')) = LOWER(?) THEN 85
		         ELSE 25
		       END::float8 AS score
		FROM qh_regions r
		JOIN qh_layers l
		  ON l.id = r.layer_id
		 AND l.deleted_at IS NULL
		 AND l.status = 10
		WHERE r.deleted_at IS NULL
		  AND r.status = 10
		  AND r.is_latest = true
		  AND (
		    LOWER(COALESCE(r.display_name,'')) LIKE ?
		    OR LOWER(COALESCE(r.name,'')) LIKE ?
		    OR LOWER(COALESCE(r.description,'')) LIKE ?
		    OR LOWER(COALESCE(l.display_name,'')) LIKE ?
		    OR LOWER(COALESCE(l.name,'')) LIKE ?
		  )`

	pattern := searchPattern(req.Query)
	args := []any{req.Query, pattern, pattern, pattern, pattern, pattern}

	query += viewportClause
	args = append(args, viewportArgs...)

	if len(req.MapContext.ActiveLayerIDs) > 0 {
		query += " AND r.layer_id::text IN ?"
		args = append(args, req.MapContext.ActiveLayerIDs)
	}

	query += " ORDER BY score DESC, COALESCE(l.display_order,0) DESC, r.id DESC LIMIT ?"
	args = append(args, limit(req.Limit))

	if err := r.db.WithContext(ctx).Raw(query, args...).Scan(&rows).Error; err != nil {
		return nil, fmt.Errorf("search regions: %w", err)
	}

	out := make([]domain.EntityCandidate, 0, len(rows))
	for _, row := range rows {
		out = append(out, candidate(domain.EntityKindPlanningRegion, row, "full_text"))
	}

	return out, nil
}

func (r *Repository) SearchPlanningProjects(ctx context.Context, req domain.SearchRequest) ([]domain.EntityCandidate, error) {
	if strings.TrimSpace(req.Query) == "" {
		return nil, nil
	}

	var rows []spatialRow
	pattern := searchPattern(req.Query)
	viewportClause, viewportArgs := viewportGeometryClause(req, "extent.geometry")

	query := `
		SELECT pp.id::text AS id,
		       pp.name AS title,
		       CONCAT_WS(' · ', NULLIF(pp.planning_type::text,''), NULLIF(pp.legal_status::text,'')) AS subtitle,
		       COALESCE(pp.summary,'') AS description,
		       COALESCE(j.name,'') AS address,
		       COALESCE(j.name,'') AS province,
		       '' AS province_code,
		       '' AS ward,
		       '' AS ward_code,
		       COALESCE(GeometryType(extent.geometry),'') AS geometry_type,
		       COALESCE(ST_Y(ST_PointOnSurface(extent.geometry)),0) AS center_lat,
		       COALESCE(ST_X(ST_PointOnSurface(extent.geometry)),0) AS center_lon,
		       COALESCE(ST_XMin(Box2D(extent.geometry)),0) AS min_lon,
		       COALESCE(ST_YMin(Box2D(extent.geometry)),0) AS min_lat,
		       COALESCE(ST_XMax(Box2D(extent.geometry)),0) AS max_lon,
		       COALESCE(ST_YMax(Box2D(extent.geometry)),0) AS max_lat,
		       '' AS geo_json,
		       COALESCE(pp.total_area,0) AS area_sqm,
		       CASE
		         WHEN LOWER(COALESCE(pp.code,'')) = LOWER(?) THEN 100
		         WHEN LOWER(COALESCE(pp.name,'')) = LOWER(?) THEN 95
		         WHEN LOWER(COALESCE(pp.metadata->>'slug','')) = LOWER(?) THEN 90
		         ELSE 35
		       END::float8 AS score
		FROM qh_planning_projects pp
		LEFT JOIN qh_jurisdictions j
		  ON j.id = pp.jurisdiction_id
		 AND j.deleted_at IS NULL
		LEFT JOIN LATERAL (
		  SELECT ST_Union(r.geometry) AS geometry
		  FROM qh_layers l
		  JOIN qh_regions r
		    ON r.layer_id = l.id
		   AND r.deleted_at IS NULL
		   AND r.status = 10
		   AND r.is_latest = true
		  WHERE l.planning_project_id = pp.id
		    AND l.deleted_at IS NULL
		    AND l.status = 10
		) extent ON true
		WHERE pp.deleted_at IS NULL
		  AND (
		    LOWER(COALESCE(pp.code,'')) LIKE ?
		    OR LOWER(COALESCE(pp.name,'')) LIKE ?
		    OR LOWER(COALESCE(pp.summary,'')) LIKE ?
		    OR LOWER(COALESCE(pp.authority,'')) LIKE ?
		    OR LOWER(COALESCE(pp.metadata->>'slug','')) LIKE ?
		  )`

	args := []any{
		req.Query,
		req.Query,
		req.Query,
		pattern,
		pattern,
		pattern,
		pattern,
		pattern,
	}
	query += viewportClause
	args = append(args, viewportArgs...)
	query += " ORDER BY score DESC, pp.updated_at DESC NULLS LAST, pp.id DESC LIMIT ?"
	args = append(args, limit(req.Limit))

	if err := r.db.WithContext(ctx).Raw(query, args...).Scan(&rows).Error; err != nil {
		return nil, fmt.Errorf("search planning projects: %w", err)
	}

	out := make([]domain.EntityCandidate, 0, len(rows))
	for _, row := range rows {
		item := candidate(domain.EntityKindPlanningProject, row, "full_text")
		item.Links.CanonicalURL = "/do-an-quy-hoach/" + row.ID
		item.Links.DetailURL = item.Links.CanonicalURL
		item.Capabilities.CanCreateReport = true
		item.Capabilities.CanCompare = true
		out = append(out, item)
	}
	return out, nil
}

func (r *Repository) SearchAdministrativeUnits(ctx context.Context, req domain.SearchRequest) ([]domain.EntityCandidate, error) {
	if strings.TrimSpace(req.Query) == "" {
		return nil, nil
	}

	var rows []spatialRow
	wardViewportClause, wardViewportArgs := viewportPointClause(req, "w.lng", "w.lat")
	provinceViewportClause, provinceViewportArgs := viewportPointClause(req, "p.lng", "p.lat")
	pattern := searchPattern(req.Query)
	l := limit(req.Limit)

	wardAdministrativeClause := ""
	wardAdministrativeArgs := []any{}
	provinceAdministrativeClause := ""
	provinceAdministrativeArgs := []any{}

	if req.AdministrativeCode != "" {
		wardAdministrativeClause = " AND (w.code = ? OR p.code = ?)"
		wardAdministrativeArgs = append(
			wardAdministrativeArgs,
			req.AdministrativeCode,
			req.AdministrativeCode,
		)
		provinceAdministrativeClause = " AND p.code = ?"
		provinceAdministrativeArgs = append(
			provinceAdministrativeArgs,
			req.AdministrativeCode,
		)
	}

	query := `
		SELECT * FROM (
		  SELECT CONCAT('ward:',w.id) AS id,
		         w.full_name AS title,
		         COALESCE(p.full_name,'') AS subtitle,
		         '' AS description,
		         '' AS address,
		         COALESCE(p.full_name,'') AS province,
		         COALESCE(p.code,'') AS province_code,
		         w.full_name AS ward,
		         COALESCE(w.code,'') AS ward_code,
		         'Point' AS geometry_type,
		         w.lat AS center_lat,
		         w.lng AS center_lon,
		         0::float8 AS min_lon,
		         0::float8 AS min_lat,
		         0::float8 AS max_lon,
		         0::float8 AS max_lat,
		         '' AS geo_json,
		         CASE
		           WHEN LOWER(COALESCE(w.full_name,'')) = LOWER(?) THEN 95
		           WHEN LOWER(COALESCE(w.code,'')) = LOWER(?) THEN 100
		           ELSE 30
		         END::float8 AS score
		  FROM ward_v2 w
		  LEFT JOIN province_v2 p ON p.id = w.province_id
		  WHERE (
		    LOWER(COALESCE(w.full_name,'')) LIKE ?
		    OR LOWER(COALESCE(w.short_name,'')) LIKE ?
		    OR LOWER(COALESCE(w.code,'')) LIKE ?
		  )`

	args := []any{req.Query, req.Query, pattern, pattern, pattern}
	query += wardViewportClause
	args = append(args, wardViewportArgs...)
	query += wardAdministrativeClause
	args = append(args, wardAdministrativeArgs...)
	query += " ORDER BY score DESC, w.full_name LIMIT ?"
	args = append(args, l)

	query += `
		) ward_results
		UNION ALL
		SELECT * FROM (
		  SELECT CONCAT('province:',p.id) AS id,
		         p.full_name AS title,
		         '' AS subtitle,
		         '' AS description,
		         '' AS address,
		         p.full_name AS province,
		         COALESCE(p.code,'') AS province_code,
		         '' AS ward,
		         '' AS ward_code,
		         'Point' AS geometry_type,
		         p.lat AS center_lat,
		         p.lng AS center_lon,
		         0::float8 AS min_lon,
		         0::float8 AS min_lat,
		         0::float8 AS max_lon,
		         0::float8 AS max_lat,
		         '' AS geo_json,
		         CASE
		           WHEN LOWER(COALESCE(p.full_name,'')) = LOWER(?) THEN 95
		           WHEN LOWER(COALESCE(p.code,'')) = LOWER(?) THEN 100
		           ELSE 30
		         END::float8 AS score
		  FROM province_v2 p
		  WHERE (
		    LOWER(COALESCE(p.full_name,'')) LIKE ?
		    OR LOWER(COALESCE(p.short_name,'')) LIKE ?
		    OR LOWER(COALESCE(p.code,'')) LIKE ?
		  )`

	args = append(args, req.Query, req.Query, pattern, pattern, pattern)
	query += provinceViewportClause
	args = append(args, provinceViewportArgs...)
	query += provinceAdministrativeClause
	args = append(args, provinceAdministrativeArgs...)
	query += " ORDER BY score DESC, p.full_name LIMIT ?"
	args = append(args, l)
	query += ") province_results"

	if err := r.db.WithContext(ctx).Raw(query, args...).Scan(&rows).Error; err != nil {
		return nil, fmt.Errorf("search administrative units: %w", err)
	}

	out := make([]domain.EntityCandidate, 0, len(rows))
	for _, row := range rows {
		out = append(out, candidate(domain.EntityKindAdministrativeUnit, row, "full_text"))
	}

	return out, nil
}

func (r *Repository) SearchPOIs(ctx context.Context, req domain.SearchRequest) ([]domain.EntityCandidate, error) {
	if strings.TrimSpace(req.Query) == "" {
		return nil, nil
	}

	var rows []spatialRow
	viewportClause, viewportArgs := viewportPointClause(req, "p.longitude", "p.latitude")
	pattern := searchPattern(req.Query)

	query := `
		SELECT p.id::text AS id,
		       p.name AS title,
		       COALESCE(p.address,'') AS subtitle,
		       COALESCE(p.description,'') AS description,
		       COALESCE(p.address,'') AS address,
		       '' AS province,
		       '' AS province_code,
		       '' AS ward,
		       '' AS ward_code,
		       'Point' AS geometry_type,
		       p.latitude AS center_lat,
		       p.longitude AS center_lon,
		       0::float8 AS min_lon,
		       0::float8 AS min_lat,
		       0::float8 AS max_lon,
		       0::float8 AS max_lat,
		       '' AS geo_json,
		       CASE
		         WHEN LOWER(p.name) = LOWER(?) THEN 90
		         WHEN LOWER(p.code) = LOWER(?) THEN 100
		         ELSE 25
		       END::float8 AS score
		FROM pois p
		WHERE p.deleted_at IS NULL
		  AND p.is_active = true
		  AND (
		    LOWER(p.name) LIKE ?
		    OR LOWER(COALESCE(p.address,'')) LIKE ?
		    OR LOWER(COALESCE(p.code,'')) LIKE ?
		  )`

	args := []any{req.Query, req.Query, pattern, pattern, pattern}
	query += viewportClause
	args = append(args, viewportArgs...)
	query += " ORDER BY score DESC, p.is_featured DESC, p.rating DESC LIMIT ?"
	args = append(args, limit(req.Limit))

	if err := r.db.WithContext(ctx).Raw(query, args...).Scan(&rows).Error; err != nil {
		return nil, fmt.Errorf("search pois: %w", err)
	}

	out := make([]domain.EntityCandidate, 0, len(rows))
	for _, row := range rows {
		out = append(out, candidate(domain.EntityKindPOI, row, "full_text"))
	}

	return out, nil
}

func (r *Repository) GetEntity(ctx context.Context, ref domain.EntityRef) (*domain.EntityCandidate, error) {
	var row spatialRow

	if ref.Kind == domain.EntityKindAdministrativeUnit {
		unitType, unitID, ok := strings.Cut(ref.ID, ":")
		if !ok || unitID == "" {
			return nil, fmt.Errorf("administrative entity id must be ward:<id> or province:<id>")
		}
		switch unitType {
		case "ward":
			err := r.db.WithContext(ctx).Raw(`
				SELECT CONCAT('ward:',w.id) AS id,
				       w.full_name AS title,
				       COALESCE(p.full_name,'') AS subtitle,
				       '' AS description,
				       '' AS address,
				       COALESCE(p.full_name,'') AS province,
				       COALESCE(p.code,'') AS province_code,
				       w.full_name AS ward,
				       COALESCE(w.code,'') AS ward_code,
				       'Point' AS geometry_type,
				       w.lat AS center_lat,
				       w.lng AS center_lon,
				       0::float8 AS min_lon,0::float8 AS min_lat,0::float8 AS max_lon,0::float8 AS max_lat,
				       '' AS geo_json,
				       50::float8 AS score
				FROM ward_v2 w
				LEFT JOIN province_v2 p ON p.id = w.province_id
				WHERE w.id = ?
				LIMIT 1`, unitID).Scan(&row).Error
			if err != nil {
				return nil, err
			}
		case "province":
			err := r.db.WithContext(ctx).Raw(`
				SELECT CONCAT('province:',p.id) AS id,
				       p.full_name AS title,
				       '' AS subtitle,
				       '' AS description,
				       '' AS address,
				       p.full_name AS province,
				       COALESCE(p.code,'') AS province_code,
				       '' AS ward,
				       '' AS ward_code,
				       'Point' AS geometry_type,
				       p.lat AS center_lat,
				       p.lng AS center_lon,
				       0::float8 AS min_lon,0::float8 AS min_lat,0::float8 AS max_lon,0::float8 AS max_lat,
				       '' AS geo_json,
				       50::float8 AS score
				FROM province_v2 p
				WHERE p.id = ?
				LIMIT 1`, unitID).Scan(&row).Error
			if err != nil {
				return nil, err
			}
		default:
			return nil, fmt.Errorf("unsupported administrative unit type %q", unitType)
		}
		if row.ID == "" {
			return nil, nil
		}
		c := candidate(ref.Kind, row, "entity_key")
		return &c, nil
	}

	numericID, err := strconv.ParseUint(ref.ID, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("entity id must be numeric: %w", err)
	}
	switch ref.Kind {
	case domain.EntityKindParcel:
		err = r.db.WithContext(ctx).Raw(`
			SELECT p.id::text AS id,
			       CONCAT(
			         'Thửa ',
			         COALESCE(NULLIF(pi.land_number,''),NULLIF(p.land_number,''),p.id::text),
			         CASE
			           WHEN COALESCE(NULLIF(pi.map_number,''),NULLIF(p.map_number,''),'') <> ''
			           THEN CONCAT(', tờ ',COALESCE(NULLIF(pi.map_number,''),NULLIF(p.map_number,'')))
			           ELSE ''
			         END
			       ) AS title,
			       COALESCE(NULLIF(pi.address_text,''),NULLIF(p.address_text,''),'') AS subtitle,
			       '' AS description,
			       COALESCE(NULLIF(pi.address_text,''),NULLIF(p.address_text,''),'') AS address,
			       COALESCE(prov.full_name,'') AS province,
			       COALESCE(pi.province_code,'') AS province_code,
			       COALESCE(w.full_name,'') AS ward,
			       COALESCE(pi.ward_code,'') AS ward_code,
			       GeometryType(p.geometry) AS geometry_type,
			       ST_Y(ST_PointOnSurface(p.geometry)) AS center_lat,
			       ST_X(ST_PointOnSurface(p.geometry)) AS center_lon,
			       ST_XMin(Box2D(p.geometry)) AS min_lon,
			       ST_YMin(Box2D(p.geometry)) AS min_lat,
			       ST_XMax(Box2D(p.geometry)) AS max_lon,
			       ST_YMax(Box2D(p.geometry)) AS max_lat,
			       ST_AsGeoJSON(ST_SimplifyPreserveTopology(p.geometry,0.000002)) AS geo_json,
			       COALESCE(NULLIF(pi.map_number,''),NULLIF(p.map_number,''),'') AS map_number,
			       COALESCE(NULLIF(pi.land_number,''),NULLIF(p.land_number,''),'') AS land_number,
			       COALESCE(NULLIF(pi.total_area_sqm,0),p.total_area_sqm,0) AS area_sqm,
			       50::float8 AS score
			FROM parcels p
			LEFT JOIN qh_parcel_info pi ON pi.parcel_id = p.id
			LEFT JOIN province_v2 prov ON prov.code = pi.province_code
			LEFT JOIN ward_v2 w ON w.code = pi.ward_code
			WHERE (p.id = ? OR p.ref_id = ?)
			  AND p.deleted_at IS NULL
			  AND p.geometry IS NOT NULL
			ORDER BY CASE WHEN p.id = ? THEN 0 ELSE 1 END, p.id DESC
			LIMIT 1`, numericID, numericID, numericID).Scan(&row).Error
	case domain.EntityKindPlanningRegion:
		err = r.db.WithContext(ctx).Raw(`
			SELECT r.id::text AS id,
			       COALESCE(NULLIF(r.display_name,''),NULLIF(r.name,''),CONCAT('Vùng quy hoạch #',r.id)) AS title,
			       COALESCE(NULLIF(l.display_name,''),l.name,'') AS subtitle,
			       COALESCE(r.description,'') AS description,
			       '' AS address,'' AS province,'' AS province_code,'' AS ward,'' AS ward_code,
			       GeometryType(r.geometry) AS geometry_type,
			       ST_Y(ST_PointOnSurface(r.geometry)) AS center_lat,
			       ST_X(ST_PointOnSurface(r.geometry)) AS center_lon,
			       ST_XMin(Box2D(r.geometry)) AS min_lon,
			       ST_YMin(Box2D(r.geometry)) AS min_lat,
			       ST_XMax(Box2D(r.geometry)) AS max_lon,
			       ST_YMax(Box2D(r.geometry)) AS max_lat,
			       ST_AsGeoJSON(ST_SimplifyPreserveTopology(r.geometry,0.00001)) AS geo_json,
			       r.layer_id::text AS layer_id,
			       COALESCE(r.area_sqm,0) AS area_sqm,
			       50::float8 AS score
			FROM qh_regions r
			JOIN qh_layers l ON l.id = r.layer_id
			WHERE r.deleted_at IS NULL
			  AND (
			    r.id = ?
			    OR EXISTS (
			      SELECT 1
			      FROM jsonb_each_text(COALESCE(r.original_properties,'{}'::jsonb)) AS property(key,value)
			      WHERE LOWER(REPLACE(property.key,'_','')) IN (
			        'regionid','planningregionid','areaid','objectid','gid','fid','id'
			      )
			        AND property.value = ?
			    )
			  )
			ORDER BY CASE WHEN r.id = ? THEN 0 ELSE 1 END, r.id DESC
			LIMIT 1`, numericID, ref.ID, numericID).Scan(&row).Error
	case domain.EntityKindPlanningProject:
		err = r.db.WithContext(ctx).Raw(`
			SELECT pp.id::text AS id,
			       pp.name AS title,
			       CONCAT_WS(' · ', NULLIF(pp.planning_type::text,''), NULLIF(pp.legal_status::text,'')) AS subtitle,
			       COALESCE(pp.summary,'') AS description,
			       COALESCE(j.name,'') AS address,
			       COALESCE(j.name,'') AS province,
			       '' AS province_code,'' AS ward,'' AS ward_code,
			       COALESCE(GeometryType(extent.geometry),'') AS geometry_type,
			       COALESCE(ST_Y(ST_PointOnSurface(extent.geometry)),0) AS center_lat,
			       COALESCE(ST_X(ST_PointOnSurface(extent.geometry)),0) AS center_lon,
			       COALESCE(ST_XMin(Box2D(extent.geometry)),0) AS min_lon,
			       COALESCE(ST_YMin(Box2D(extent.geometry)),0) AS min_lat,
			       COALESCE(ST_XMax(Box2D(extent.geometry)),0) AS max_lon,
			       COALESCE(ST_YMax(Box2D(extent.geometry)),0) AS max_lat,
			       '' AS geo_json,
			       COALESCE(pp.total_area,0) AS area_sqm,
			       50::float8 AS score
			FROM qh_planning_projects pp
			LEFT JOIN qh_jurisdictions j ON j.id = pp.jurisdiction_id AND j.deleted_at IS NULL
			LEFT JOIN LATERAL (
			  SELECT ST_Union(r.geometry) AS geometry
			  FROM qh_layers l
			  JOIN qh_regions r ON r.layer_id = l.id AND r.deleted_at IS NULL AND r.status = 10 AND r.is_latest = true
			  WHERE l.planning_project_id = pp.id AND l.deleted_at IS NULL AND l.status = 10
			) extent ON true
			WHERE pp.id = ? AND pp.deleted_at IS NULL
			LIMIT 1`, numericID).Scan(&row).Error
	case domain.EntityKindPOI:
		err = r.db.WithContext(ctx).Raw(`
			SELECT p.id::text AS id,p.name AS title,COALESCE(p.address,'') AS subtitle,
			       COALESCE(p.description,'') AS description,COALESCE(p.address,'') AS address,
			       '' AS province,'' AS province_code,'' AS ward,'' AS ward_code,
			       'Point' AS geometry_type,p.latitude AS center_lat,p.longitude AS center_lon,
			       0::float8 AS min_lon,0::float8 AS min_lat,0::float8 AS max_lon,0::float8 AS max_lat,
			       '' AS geo_json,50::float8 AS score
			FROM pois p
			WHERE p.id = ? AND p.deleted_at IS NULL
			LIMIT 1`, numericID).Scan(&row).Error
	default:
		return nil, fmt.Errorf("unsupported entity kind %q", ref.Kind)
	}
	if err != nil {
		return nil, err
	}
	if row.ID == "" {
		return nil, nil
	}
	c := candidate(ref.Kind, row, "entity_key")
	return &c, nil
}
