package relatedentity

import (
	"context"
	"fmt"

	discoverydomain "tqd/internal/domain/discovery/model"
	relateddomain "tqd/internal/domain/relatedentity/model"
)

func (r *Repository) FindPlanningEntities(ctx context.Context, query relateddomain.Query) ([]relateddomain.Candidate, error) {
	switch query.Primary.Kind {
	case discoverydomain.EntityKindParcel:
		return r.findPlanningRegionsForParcel(ctx, query)
	case discoverydomain.EntityKindPlanningRegion:
		return r.findPlanningEntitiesForRegion(ctx, query)
	case discoverydomain.EntityKindPlanningProject:
		return r.findPlanningRegionsForProject(ctx, query)
	default:
		return r.findPlanningRegionsAroundPoint(ctx, query)
	}
}

func planningGeometryPreview(include bool) string {
	if !include {
		return "NULL::text AS geo_json"
	}
	return "ST_AsGeoJSON(ST_SimplifyPreserveTopology(region.geometry, 0.00001)) AS geo_json"
}

func (r *Repository) findPlanningRegionsForParcel(ctx context.Context, query relateddomain.Query) ([]relateddomain.Candidate, error) {
	// A parcel's core planning relations are direct containment/intersection.
	// Nearby planning polygons are a separate exploratory concern and should not
	// make the primary Related query pay geography distance/buffer costs.
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
		FROM qh_regions region
		JOIN qh_layers layer
		  ON layer.id = region.layer_id
		 AND layer.deleted_at IS NULL
		 AND layer.status = 10
		CROSS JOIN primary_entity
		WHERE region.deleted_at IS NULL
		  AND region.status = 10
		  AND region.is_latest = true
		  AND region.geometry IS NOT NULL
		  AND region.geometry && primary_entity.geometry
		  AND ST_Intersects(region.geometry, primary_entity.geometry)
		ORDER BY score DESC, region.id DESC
		LIMIT ?`, planningSelectColumns(
		planningGeometryPreview(query.IncludeGeometryPreview),
		`CASE
		   WHEN ST_Covers(region.geometry, primary_entity.geometry) THEN 'contains_primary'
		   ELSE 'intersects_primary'
		 END`,
		`0::float8`,
		`CASE
		   WHEN ST_Covers(region.geometry, primary_entity.geometry) THEN 100
		   ELSE 82
		 END`,
	))

	return r.scanPlanningRegions(ctx, "find planning regions for parcel", sql,
		query.Primary.ID,
		query.Limit,
	)
}

func (r *Repository) findPlanningEntitiesForRegion(ctx context.Context, query relateddomain.Query) ([]relateddomain.Candidate, error) {
	projects, err := r.findPlanningProjectForRegion(ctx, query)
	if err != nil {
		return nil, err
	}
	if len(projects) >= query.Limit {
		return projects[:query.Limit], nil
	}

	remaining := query.Limit - len(projects)
	if remaining <= 0 {
		return projects, nil
	}
	nearbyQuery := query
	nearbyQuery.Limit = remaining
	regions, err := r.findPlanningRegionsAroundRegion(ctx, nearbyQuery)
	if err != nil {
		return nil, err
	}
	return append(projects, regions...), nil
}

func planningProjectGeometryPreview(include bool) string {
	if !include {
		return "NULL::text AS geo_json"
	}
	return "ST_AsGeoJSON(ST_SimplifyPreserveTopology(extent.geometry, 0.00001)) AS geo_json"
}

func (r *Repository) findPlanningProjectForRegion(ctx context.Context, query relateddomain.Query) ([]relateddomain.Candidate, error) {
	sql := fmt.Sprintf(`
		WITH selected_project AS (
			SELECT project.*
			FROM qh_regions selected_region
			JOIN qh_layers selected_layer
			  ON selected_layer.id = selected_region.layer_id
			 AND selected_layer.deleted_at IS NULL
			 AND selected_layer.status = 10
			JOIN qh_planning_projects project
			  ON project.id = selected_layer.planning_project_id
			 AND project.deleted_at IS NULL
			WHERE selected_region.id::text = ?
			  AND selected_region.deleted_at IS NULL
			  AND selected_region.status = 10
			  AND selected_region.is_latest = true
			LIMIT 1
		)
		SELECT project.id::text AS id,
		       project.name AS title,
		       CONCAT_WS(' · ', NULLIF(project.planning_type::text, ''), NULLIF(project.legal_status::text, '')) AS subtitle,
		       COALESCE(project.summary, '') AS description,
		       COALESCE(jurisdiction.name, '') AS address,
		       COALESCE(jurisdiction.name, '') AS province,
		       '' AS province_code,
		       '' AS ward,
		       '' AS ward_code,
		       COALESCE(GeometryType(extent.geometry), '') AS geometry_type,
		       COALESCE(ST_Y(ST_PointOnSurface(extent.geometry)), 0) AS center_lat,
		       COALESCE(ST_X(ST_PointOnSurface(extent.geometry)), 0) AS center_lon,
		       COALESCE(ST_XMin(Box2D(extent.geometry)), 0) AS min_lon,
		       COALESCE(ST_YMin(Box2D(extent.geometry)), 0) AS min_lat,
		       COALESCE(ST_XMax(Box2D(extent.geometry)), 0) AS max_lon,
		       COALESCE(ST_YMax(Box2D(extent.geometry)), 0) AS max_lat,
		       %s,
		       COALESCE(project.total_area, 0) AS area_sqm,
		       'contains_primary' AS relationship,
		       0::float8 AS distance_meters,
		       110::float8 AS score
		FROM selected_project project
		LEFT JOIN qh_jurisdictions jurisdiction
		  ON jurisdiction.id = project.jurisdiction_id
		 AND jurisdiction.deleted_at IS NULL
		LEFT JOIN LATERAL (
			SELECT ST_Union(region.geometry) AS geometry
			FROM qh_layers layer
			JOIN qh_regions region
			  ON region.layer_id = layer.id
			 AND region.deleted_at IS NULL
			 AND region.status = 10
			 AND region.is_latest = true
			WHERE layer.planning_project_id = project.id
			  AND layer.deleted_at IS NULL
			  AND layer.status = 10
		) extent ON true
		LIMIT 1`, planningProjectGeometryPreview(query.IncludeGeometryPreview))

	var rows []spatialRow
	if err := r.db.WithContext(ctx).Raw(sql, query.Primary.ID).Scan(&rows).Error; err != nil {
		return nil, fmt.Errorf("find planning project for region: %w", err)
	}
	return rowsToCandidates(discoverydomain.EntityKindPlanningProject, rows), nil
}

func (r *Repository) findPlanningRegionsAroundRegion(ctx context.Context, query relateddomain.Query) ([]relateddomain.Candidate, error) {
	sql := fmt.Sprintf(`
		WITH primary_entity AS (
			SELECT region.geometry,
			       ST_Envelope(ST_Buffer(region.geometry::geography, ?)::geometry) AS search_envelope
			FROM qh_regions region
			WHERE region.id::text = ?
			  AND region.deleted_at IS NULL
			  AND region.geometry IS NOT NULL
			LIMIT 1
		)
		%s
		FROM qh_regions region
		JOIN qh_layers layer
		  ON layer.id = region.layer_id
		 AND layer.deleted_at IS NULL
		 AND layer.status = 10
		CROSS JOIN primary_entity
		WHERE region.deleted_at IS NULL
		  AND region.status = 10
		  AND region.is_latest = true
		  AND region.geometry IS NOT NULL
		  AND region.id::text <> ?
		  AND region.geometry && primary_entity.search_envelope
		  AND (
		    ST_Intersects(region.geometry, primary_entity.geometry)
		    OR ST_DWithin(region.geometry::geography, primary_entity.geometry::geography, ?)
		  )
		ORDER BY score DESC, distance_meters ASC, region.id DESC
		LIMIT ?`, planningSelectColumns(
		planningGeometryPreview(query.IncludeGeometryPreview),
		`CASE
		   WHEN ST_Touches(region.geometry, primary_entity.geometry) THEN 'adjacent'
		   WHEN ST_Intersects(region.geometry, primary_entity.geometry) THEN 'intersects_primary'
		   ELSE 'nearby'
		 END`,
		`ST_Distance(region.geometry::geography, primary_entity.geometry::geography)`,
		`CASE
		   WHEN ST_Touches(region.geometry, primary_entity.geometry) THEN 78
		   WHEN ST_Intersects(region.geometry, primary_entity.geometry) THEN 70
		   ELSE 55
		 END`,
	))

	return r.scanPlanningRegions(ctx, "find planning regions around region", sql,
		query.RadiusMeters,
		query.Primary.ID,
		query.Primary.ID,
		query.RadiusMeters,
		query.Limit,
	)
}

func planningRegionsForProjectSQL(includeGeometryPreview bool) string {
	// Membership is a direct relational fact. Radius may order members around the
	// project focus, but it must never exclude a region that belongs to project.
	return fmt.Sprintf(`
		WITH focus AS (
			SELECT ST_SetSRID(ST_Point(?, ?), 4326) AS point
		)
		%s
		FROM qh_regions region
		JOIN qh_layers layer
		  ON layer.id = region.layer_id
		 AND layer.deleted_at IS NULL
		 AND layer.status = 10
		CROSS JOIN focus
		WHERE layer.planning_project_id::text = ?
		  AND region.deleted_at IS NULL
		  AND region.status = 10
		  AND region.is_latest = true
		  AND region.geometry IS NOT NULL
		ORDER BY score DESC, distance_meters ASC, region.id DESC
		LIMIT ?`, planningSelectColumns(
		planningGeometryPreview(includeGeometryPreview),
		`'member_of_primary'`,
		`ST_Distance(region.geometry::geography, focus.point::geography)`,
		`98`,
	))
}

func (r *Repository) findPlanningRegionsForProject(ctx context.Context, query relateddomain.Query) ([]relateddomain.Candidate, error) {
	return r.scanPlanningRegions(
		ctx,
		"find planning regions for project",
		planningRegionsForProjectSQL(query.IncludeGeometryPreview),
		query.Focus.Longitude,
		query.Focus.Latitude,
		query.Primary.ID,
		query.Limit,
	)
}

func (r *Repository) findPlanningRegionsAroundPoint(ctx context.Context, query relateddomain.Query) ([]relateddomain.Candidate, error) {
	pointCanBeContained := query.Primary.Kind == discoverydomain.EntityKindPOI
	sql := fmt.Sprintf(`
		WITH focus AS (
			SELECT point,
			       ST_Envelope(ST_Buffer(point::geography, ?)::geometry) AS search_envelope
			FROM (SELECT ST_SetSRID(ST_Point(?, ?), 4326) AS point) input
		)
		%s
		FROM qh_regions region
		JOIN qh_layers layer
		  ON layer.id = region.layer_id
		 AND layer.deleted_at IS NULL
		 AND layer.status = 10
		CROSS JOIN focus
		WHERE region.deleted_at IS NULL
		  AND region.status = 10
		  AND region.is_latest = true
		  AND region.geometry IS NOT NULL
		  AND region.geometry && focus.search_envelope
		  AND (
		    ST_Covers(region.geometry, focus.point)
		    OR ST_DWithin(region.geometry::geography, focus.point::geography, ?)
		  )
		ORDER BY score DESC, distance_meters ASC, region.id DESC
		LIMIT ?`, planningSelectColumns(
		planningGeometryPreview(query.IncludeGeometryPreview),
		`CASE WHEN ? AND ST_Covers(region.geometry, focus.point) THEN 'contains_primary' ELSE 'nearby' END`,
		`ST_Distance(region.geometry::geography, focus.point::geography)`,
		`CASE WHEN ? AND ST_Covers(region.geometry, focus.point) THEN 100 ELSE 55 END`,
	))

	return r.scanPlanningRegions(ctx, "find planning regions around point", sql,
		query.RadiusMeters,
		query.Focus.Longitude,
		query.Focus.Latitude,
		pointCanBeContained,
		pointCanBeContained,
		query.RadiusMeters,
		query.Limit,
	)
}

func planningSelectColumns(geometryPreview, relationship, distance, score string) string {
	return fmt.Sprintf(`SELECT region.id::text AS id,
	       COALESCE(NULLIF(region.display_name, ''), NULLIF(region.name, ''), CONCAT('Vùng quy hoạch #', region.id)) AS title,
	       COALESCE(NULLIF(layer.display_name, ''), layer.name, '') AS subtitle,
	       COALESCE(region.description, '') AS description,
	       '' AS address,
	       '' AS province,
	       '' AS province_code,
	       '' AS ward,
	       '' AS ward_code,
	       GeometryType(region.geometry) AS geometry_type,
	       ST_Y(ST_PointOnSurface(region.geometry)) AS center_lat,
	       ST_X(ST_PointOnSurface(region.geometry)) AS center_lon,
	       ST_XMin(Box2D(region.geometry)) AS min_lon,
	       ST_YMin(Box2D(region.geometry)) AS min_lat,
	       ST_XMax(Box2D(region.geometry)) AS max_lon,
	       ST_YMax(Box2D(region.geometry)) AS max_lat,
	       %s,
	       region.layer_id::text AS layer_id,
	       COALESCE(region.area_sqm, ST_Area(region.geometry::geography), 0) AS area_sqm,
	       %s AS relationship,
	       %s AS distance_meters,
	       (%s)::float8 AS score`,
		geometryPreview,
		relationship,
		distance,
		score,
	)
}

func (r *Repository) scanPlanningRegions(ctx context.Context, operation, sql string, args ...any) ([]relateddomain.Candidate, error) {
	var rows []spatialRow
	if err := r.db.WithContext(ctx).Raw(sql, args...).Scan(&rows).Error; err != nil {
		return nil, fmt.Errorf("%s: %w", operation, err)
	}
	return rowsToCandidates(discoverydomain.EntityKindPlanningRegion, rows), nil
}
