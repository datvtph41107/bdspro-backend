package publiccontent

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"tqd/internal/domain/publiccontent/model"
	"tqd/internal/enums"

	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

type parcelRow struct {
	ParcelID      uint64
	HasParcelInfo bool
	MapNumber     string
	LandNumber    string
	PropertyCode  string
	Address       string
	TotalAreaSqm  float64
	ProvinceCode  string
	WardCode      string
	IsVerified    bool
	UpdatedAt     *time.Time
	GeometryType  string
	CenterLat     float64
	CenterLon     float64
	MinLon        float64
	MinLat        float64
	MaxLon        float64
	MaxLat        float64
}

type relatedProjectRow struct {
	ID               uint64
	Code             string
	Name             string
	PlanningType     uint32
	LegalStatus      uint32
	ValidityStatus   string
	TotalArea        float64
	Metadata         string
	UpdatedAt        *time.Time
	IntersectAreaSqm float64
	ParcelAreaSqm    float64
}

type projectRow struct {
	ID               uint64
	Code             string
	Slug             string
	Name             string
	PlanningType     uint32
	PlanningLevel    uint32
	TotalArea        float64
	JurisdictionID   *uint64
	JurisdictionName string
	LegalStatus      uint32
	ValidityStatus   string
	Authority        string
	DecisionNumber   string
	Summary          string
	ResearchScope    string
	Indicators       string
	ApprovalDate     *time.Time
	EffectiveDate    *time.Time
	ExpiryDate       *time.Time
	CurrentVersion   string
	Metadata         string
	UpdatedAt        *time.Time
	MinLon           float64
	MinLat           float64
	MaxLon           float64
	MaxLat           float64
	CenterLat        float64
	CenterLon        float64
	GeometryType     string
}

type eventRow struct {
	ID               uint64
	EventType        string
	Name             string
	Description      string
	EventDate        *time.Time
	DocumentID       *uint64
	DocumentCode     string
	DocumentTitle    string
	IssuingAuthority string
	SourceURL        string
}

type documentRow struct {
	ID             uint64
	DocumentType   uint32
	Status         uint32
	Code           string
	Title          string
	Description    string
	Filepath       string
	Thumbnail      string
	VersionNo      string
	ValidityStatus string
	IssueDate      *time.Time
	EffectiveDate  *time.Time
}

type relationRow struct {
	ID       uint64
	Code     string
	Name     string
	Metadata string
}

type layerRow struct {
	ID                uint64
	PlanningProjectID uint64
	Name              string
	DisplayName       string
	Description       string
	Type              uint32
	Status            uint32
	Visible           bool
	DefaultVisible    bool
	MinZoom           uint32
	MaxZoom           uint32
	ImageURL          string
	ThumbnailURL      string
	EffectiveDate     *time.Time
	ExpiryDate        *time.Time
	UpdatedAt         *time.Time
}

var nonSlugCharacters = regexp.MustCompile(`[^a-z0-9]+`)

func slugify(value string) string {
	value = strings.ToLower(strings.TrimSpace(value))
	replacer := strings.NewReplacer(
		"à", "a", "á", "a", "ạ", "a", "ả", "a", "ã", "a",
		"â", "a", "ầ", "a", "ấ", "a", "ậ", "a", "ẩ", "a", "ẫ", "a",
		"ă", "a", "ằ", "a", "ắ", "a", "ặ", "a", "ẳ", "a", "ẵ", "a",
		"è", "e", "é", "e", "ẹ", "e", "ẻ", "e", "ẽ", "e",
		"ê", "e", "ề", "e", "ế", "e", "ệ", "e", "ể", "e", "ễ", "e",
		"ì", "i", "í", "i", "ị", "i", "ỉ", "i", "ĩ", "i",
		"ò", "o", "ó", "o", "ọ", "o", "ỏ", "o", "õ", "o",
		"ô", "o", "ồ", "o", "ố", "o", "ộ", "o", "ổ", "o", "ỗ", "o",
		"ơ", "o", "ờ", "o", "ớ", "o", "ợ", "o", "ở", "o", "ỡ", "o",
		"ù", "u", "ú", "u", "ụ", "u", "ủ", "u", "ũ", "u",
		"ư", "u", "ừ", "u", "ứ", "u", "ự", "u", "ử", "u", "ữ", "u",
		"ỳ", "y", "ý", "y", "ỵ", "y", "ỷ", "y", "ỹ", "y",
		"đ", "d",
	)
	value = replacer.Replace(value)
	value = strings.Trim(nonSlugCharacters.ReplaceAllString(value, "-"), "-")
	return value
}

func metadataString(metadata string, keys ...string) string {
	if strings.TrimSpace(metadata) == "" {
		return ""
	}
	var data map[string]any
	if err := json.Unmarshal([]byte(metadata), &data); err != nil {
		return ""
	}
	for _, key := range keys {
		if value, ok := data[key]; ok {
			text := strings.TrimSpace(fmt.Sprint(value))
			if text != "" && text != "<nil>" {
				return text
			}
		}
	}
	return ""
}

func metadataIndicators(metadata string) []domain.PlanningIndicator {
	if strings.TrimSpace(metadata) == "" {
		return []domain.PlanningIndicator{}
	}

	var root interface{}
	if err := json.Unmarshal([]byte(metadata), &root); err != nil {
		return []domain.PlanningIndicator{}
	}

	var items []interface{}

	switch value := root.(type) {
	case []interface{}:
		items = value

	case map[string]interface{}:
		for _, key := range []string{
			"indicators",
			"planningIndicators",
			"planning_indicators",
			"targets",
			"planningTargets",
		} {
			raw, exists := value[key]
			if !exists {
				continue
			}

			list, ok := raw.([]interface{})
			if ok {
				items = list
			}
			break
		}
	}

	if len(items) == 0 {
		return []domain.PlanningIndicator{}
	}

	result := make([]domain.PlanningIndicator, 0, len(items))

	for _, rawItem := range items {
		entry, ok := rawItem.(map[string]interface{})
		if !ok {
			continue
		}

		read := func(keys ...string) string {
			for _, key := range keys {
				value, exists := entry[key]
				if !exists {
					continue
				}

				rendered := strings.TrimSpace(fmt.Sprint(value))
				if rendered != "" && rendered != "<nil>" {
					return rendered
				}
			}
			return ""
		}

		indicator := domain.PlanningIndicator{
			Code:        read("code", "key", "id"),
			Name:        read("name", "label", "title"),
			Value:       read("value", "target", "amount"),
			Unit:        read("unit", "uom"),
			Description: read("description", "note"),
		}

		if indicator.Name == "" || indicator.Value == "" {
			continue
		}

		result = append(result, indicator)
	}

	return result
}

func projectSlug(id uint64, persistedSlug, code, name, metadata string) string {
	if value := strings.TrimSpace(persistedSlug); value != "" {
		return slugify(value)
	}
	if value := metadataString(metadata, "slug", "seoSlug", "seo_slug"); value != "" {
		return slugify(value)
	}
	if value := slugify(code); value != "" {
		return value
	}
	if value := slugify(name); value != "" {
		return value
	}
	return strconv.FormatUint(id, 10)
}

func preview(
	geometryType string,
	centerLat, centerLon, minLon, minLat, maxLon, maxLat float64,
) domain.PreviewContext {
	result := domain.PreviewContext{
		GeometryType: geometryType,
		PaddingRatio: 0.12,
		FitMode:      "contain",
		Completeness: "missing",
	}

	if centerLat != 0 || centerLon != 0 {
		result.Centroid = &domain.Point{
			Latitude:  centerLat,
			Longitude: centerLon,
		}
	}

	if minLon != 0 || minLat != 0 || maxLon != 0 || maxLat != 0 {
		result.Bounds = &domain.Bounds{
			MinLongitude: minLon,
			MinLatitude:  minLat,
			MaxLongitude: maxLon,
			MaxLatitude:  maxLat,
		}
		result.Completeness = "complete"
	}

	return result
}

func (r *Repository) GetParcelQuickView(
	ctx context.Context,
	parcelID uint64,
) (*domain.ParcelQuickView, error) {
	var row parcelRow
	err := r.db.WithContext(ctx).Raw(`
		SELECT p.id AS parcel_id,
		       (pi.parcel_id IS NOT NULL) AS has_parcel_info,
		       COALESCE(NULLIF(pi.map_number,''),NULLIF(p.map_number,''),'') AS map_number,
		       COALESCE(NULLIF(pi.land_number,''),NULLIF(p.land_number,''),'') AS land_number,
		       COALESCE(pi.property_code,'') AS property_code,
		       COALESCE(NULLIF(pi.address_text,''),NULLIF(p.address_text,''),'') AS address,
		       COALESCE(NULLIF(pi.total_area_sqm,0),p.total_area_sqm,0) AS total_area_sqm,
		       COALESCE(pi.province_code,'') AS province_code,
		       COALESCE(pi.ward_code,'') AS ward_code,
		       COALESCE(pi.is_verified,false) AS is_verified,
		       COALESCE(pi.updated_at,p.updated_at) AS updated_at,
		       GeometryType(p.geometry) AS geometry_type,
		       ST_Y(ST_PointOnSurface(p.geometry)) AS center_lat,
		       ST_X(ST_PointOnSurface(p.geometry)) AS center_lon,
		       ST_XMin(Box2D(p.geometry)) AS min_lon,
		       ST_YMin(Box2D(p.geometry)) AS min_lat,
		       ST_XMax(Box2D(p.geometry)) AS max_lon,
		       ST_YMax(Box2D(p.geometry)) AS max_lat
		FROM parcels p
		LEFT JOIN qh_parcel_info pi ON pi.parcel_id = p.id
		WHERE p.id = ?
		  AND p.deleted_at IS NULL
		  AND p.geometry IS NOT NULL
		LIMIT 1`, parcelID).Scan(&row).Error
	if err != nil {
		return nil, fmt.Errorf("load parcel quick view: %w", err)
	}
	if row.ParcelID == 0 {
		return nil, gorm.ErrRecordNotFound
	}

	var projectRows []relatedProjectRow
	if err := r.db.WithContext(ctx).Raw(`
		SELECT DISTINCT pp.id,
		       COALESCE(pp.code,'') AS code,
		       pp.name,
		       pp.planning_type,
		       pp.legal_status,
		       COALESCE(pp.validity_status,'') AS validity_status,
		       COALESCE(pp.total_area,0) AS total_area,
		       COALESCE(pp.metadata::text,'{}') AS metadata,
		       pp.updated_at,
		       ST_Area(ST_Intersection(p.geometry, r.geometry)::geography) AS intersect_area_sqm,
		       NULLIF(ST_Area(p.geometry::geography),0) AS parcel_area_sqm
		FROM parcels p
		JOIN qh_regions r
		  ON r.deleted_at IS NULL
		 AND r.status = 10
		 AND r.is_latest = true
		 AND ST_Intersects(p.geometry, r.geometry)
		JOIN qh_layers l
		  ON l.id = r.layer_id
		 AND l.deleted_at IS NULL
		 AND l.status = 10
		 AND l.planning_project_id IS NOT NULL
		JOIN qh_planning_projects pp
		  ON pp.id = l.planning_project_id
		 AND pp.deleted_at IS NULL
		WHERE p.id = ?
		ORDER BY intersect_area_sqm DESC, pp.id DESC
		LIMIT 12`, parcelID).Scan(&projectRows).Error; err != nil {
		return nil, fmt.Errorf("load parcel related planning projects: %w", err)
	}

	related := make([]domain.RelatedPlanningProject, 0, len(projectRows))
	for _, project := range projectRows {
		slug := projectSlug(project.ID, "", project.Code, project.Name, project.Metadata)
		ratio := 0.0
		if project.ParcelAreaSqm > 0 {
			ratio = project.IntersectAreaSqm / project.ParcelAreaSqm * 100
		}
		related = append(related, buildRelatedPlanningProject(project, slug, ratio))
	}

	return buildParcelQuickView(row, related), nil
}

func buildRelatedPlanningProject(project relatedProjectRow, slug string, ratio float64) domain.RelatedPlanningProject {
	id := strconv.FormatUint(project.ID, 10)
	return domain.RelatedPlanningProject{
		ID:                project.ID,
		Code:              project.Code,
		Slug:              slug,
		Name:              project.Name,
		PlanningType:      project.PlanningType,
		PlanningTypeName:  enums.GetPlanningTypeLabel(project.PlanningType),
		LegalStatus:       project.LegalStatus,
		LegalStatusName:   enums.LegalStatus(project.LegalStatus).DisplayName(),
		ValidityStatus:    project.ValidityStatus,
		TotalArea:         project.TotalArea,
		IntersectAreaSqm:  project.IntersectAreaSqm,
		IntersectionRatio: ratio,
		UpdatedAt:         project.UpdatedAt,
		PublicPath:        "/do-an-quy-hoach/" + slug,
		DossierPath:       "/ban-do/do-an/" + id,
		MapPath:           "/ban-do?projectId=" + id,
	}
}

func buildParcelQuickView(
	row parcelRow,
	related []domain.RelatedPlanningProject,
) *domain.ParcelQuickView {
	completeness := "complete"
	warnings := []string{}
	if !row.HasParcelInfo {
		completeness = "partial"
		warnings = append(warnings, "parcel_info_missing")
	}
	if len(related) == 0 {
		warnings = append(warnings, "planning_context_empty")
	}

	return &domain.ParcelQuickView{
		ParcelID:        row.ParcelID,
		CanonicalKey:    "parcel:" + strconv.FormatUint(row.ParcelID, 10),
		MapNumber:       row.MapNumber,
		LandNumber:      row.LandNumber,
		PropertyCode:    row.PropertyCode,
		Address:         row.Address,
		TotalAreaSqm:    row.TotalAreaSqm,
		ProvinceCode:    row.ProvinceCode,
		WardCode:        row.WardCode,
		IsVerified:      row.IsVerified,
		Completeness:    completeness,
		Warnings:        warnings,
		Preview:         preview(row.GeometryType, row.CenterLat, row.CenterLon, row.MinLon, row.MinLat, row.MaxLon, row.MaxLat),
		RelatedProjects: related,
		UpdatedAt:       row.UpdatedAt,
	}
}

func (r *Repository) getPlanningProjectRow(
	ctx context.Context,
	identity string,
) (*projectRow, error) {
	var row projectRow

	/*
		Resolve exactly one planning project before calculating spatial statistics.

		The previous query joined a correlated LATERAL ST_Union before the final
		WHERE/LIMIT could reliably reduce the candidate set. Besides creating a
		full topological union only to read bounds, that plan could aggregate a
		large number of region geometries and take tens of seconds.

		This query has two explicit stages:

		1. selected_project resolves one project using only lightweight columns.
		2. The LATERAL subquery calculates ST_Extent for that selected project.

		ST_Extent is sufficient for the public projection because it only needs
		map bounds and a camera center. It avoids constructing a merged polygon.
	*/
	query := `
		WITH selected_project AS MATERIALIZED (
			SELECT pp.id,
			       COALESCE(pp.code,'') AS code,
			       COALESCE(pp.slug,'') AS slug,
			       pp.name,
			       pp.planning_type,
			       pp.planning_level,
			       COALESCE(pp.total_area,0) AS total_area,
			       pp.jurisdiction_id,
			       COALESCE(j.name,'') AS jurisdiction_name,
			       pp.legal_status,
			       COALESCE(pp.validity_status,'') AS validity_status,
			       COALESCE(pp.authority,'') AS authority,
			       COALESCE(pp.decision_number,'') AS decision_number,
			       COALESCE(pp.summary,'') AS summary,
			       COALESCE(pp.research_scope,'') AS research_scope,
			       COALESCE(pp.indicators::text,'[]') AS indicators,
			       pp.approval_date,
			       pp.effective_date,
			       pp.expiry_date,
			       COALESCE(pp.current_version,'') AS current_version,
			       COALESCE(pp.metadata::text,'{}') AS metadata,
			       pp.updated_at
			FROM qh_planning_projects pp
			LEFT JOIN qh_jurisdictions j
			  ON j.id = pp.jurisdiction_id
			 AND j.deleted_at IS NULL
			WHERE pp.deleted_at IS NULL
			  AND (
			    pp.id::text = ?
			    OR LOWER(COALESCE(pp.code,'')) = LOWER(?)
			    OR LOWER(COALESCE(pp.slug,'')) = LOWER(?)
			    OR LOWER(COALESCE(pp.metadata->>'slug','')) = LOWER(?)
			    OR LOWER(REGEXP_REPLACE(pp.name, '[^[:alnum:]]+', '-', 'g')) = LOWER(?)
			  )
			LIMIT 1
		)
		SELECT selected.id,
		       selected.code,
		       selected.slug,
		       selected.name,
		       selected.planning_type,
		       selected.planning_level,
		       selected.total_area,
		       selected.jurisdiction_id,
		       selected.jurisdiction_name,
		       selected.legal_status,
		       selected.validity_status,
		       selected.authority,
		       selected.decision_number,
		       selected.summary,
		       selected.research_scope,
		       selected.indicators,
		       selected.approval_date,
		       selected.effective_date,
		       selected.expiry_date,
		       selected.current_version,
		       selected.metadata,
		       selected.updated_at,
		       COALESCE(extent.geometry_type,'') AS geometry_type,
		       COALESCE(
		         (ST_YMin(extent.bounds) + ST_YMax(extent.bounds)) / 2.0,
		         0
		       ) AS center_lat,
		       COALESCE(
		         (ST_XMin(extent.bounds) + ST_XMax(extent.bounds)) / 2.0,
		         0
		       ) AS center_lon,
		       COALESCE(ST_XMin(extent.bounds),0) AS min_lon,
		       COALESCE(ST_YMin(extent.bounds),0) AS min_lat,
		       COALESCE(ST_XMax(extent.bounds),0) AS max_lon,
		       COALESCE(ST_YMax(extent.bounds),0) AS max_lat
		FROM selected_project selected
		LEFT JOIN LATERAL (
			SELECT
			  MIN(GeometryType(region.geometry)) AS geometry_type,
			  ST_Extent(region.geometry) AS bounds
			FROM qh_layers layer
			JOIN qh_regions region
			  ON region.layer_id = layer.id
			 AND region.deleted_at IS NULL
			 AND region.status = 10
			 AND region.is_latest = true
			 AND region.geometry IS NOT NULL
			WHERE layer.planning_project_id = selected.id
			  AND layer.deleted_at IS NULL
			  AND layer.status = 10
		) extent ON true`

	normalizedIdentity := strings.TrimSpace(identity)
	if err := r.db.WithContext(ctx).
		Raw(
			query,
			normalizedIdentity,
			normalizedIdentity,
			normalizedIdentity,
			normalizedIdentity,
			normalizedIdentity,
		).
		Scan(&row).Error; err != nil {
		return nil, fmt.Errorf(
			"query planning project %q: %w",
			normalizedIdentity,
			err,
		)
	}
	if row.ID == 0 {
		return nil, gorm.ErrRecordNotFound
	}
	return &row, nil
}

func (r *Repository) GetPlanningProjectProjection(
	ctx context.Context,
	identity string,
) (*domain.PlanningProjectProjection, error) {
	row, err := r.getPlanningProjectRow(ctx, strings.TrimSpace(identity))
	if err != nil {
		return nil, fmt.Errorf("load planning project: %w", err)
	}

	slug := projectSlug(row.ID, row.Slug, row.Code, row.Name, row.Metadata)
	researchScope := strings.TrimSpace(row.ResearchScope)
	if researchScope == "" {
		researchScope = metadataString(
			row.Metadata,
			"researchScope",
			"research_scope",
			"studyScope",
			"study_scope",
		)
	}
	indicators := metadataIndicators(row.Indicators)
	if len(indicators) == 0 {
		indicators = metadataIndicators(row.Metadata)
	}

	var eventRows []eventRow
	if err := r.db.WithContext(ctx).Raw(`
		SELECT e.id,
		       COALESCE(e.event_type,'') AS event_type,
		       e.event_name AS name,
		       COALESCE(e.description,'') AS description,
		       e.event_date,
		       e.document_id,
		       COALESCE(d.code,'') AS document_code,
		       COALESCE(d.title,'') AS document_title,
		       COALESCE(
		         NULLIF(e.metadata->>'issuingAuthority',''),
		         NULLIF(e.metadata->>'issuing_authority',''),
		         ''
		       ) AS issuing_authority,
		       COALESCE(e.source_url,'') AS source_url
		FROM qh_planning_events e
		LEFT JOIN qh_planning_documents d
		  ON d.id = e.document_id
		 AND d.deleted_at IS NULL
		WHERE e.planning_project_id = ?
		  AND e.deleted_at IS NULL
		ORDER BY e.event_date DESC NULLS LAST, e.id DESC`, row.ID).Scan(&eventRows).Error; err != nil {
		return nil, fmt.Errorf("load planning project events: %w", err)
	}

	events := make([]domain.PlanningEvent, 0, len(eventRows))
	for _, item := range eventRows {
		events = append(events, domain.PlanningEvent{
			ID:               item.ID,
			EventType:        item.EventType,
			Name:             item.Name,
			Description:      item.Description,
			EventDate:        item.EventDate,
			DocumentID:       item.DocumentID,
			DocumentCode:     item.DocumentCode,
			DocumentTitle:    item.DocumentTitle,
			IssuingAuthority: item.IssuingAuthority,
			SourceURL:        item.SourceURL,
		})
	}

	var documentRows []documentRow
	if err := r.db.WithContext(ctx).Raw(`
		SELECT id,
		       document_type,
		       COALESCE(status,10) AS status,
		       COALESCE(code,'') AS code,
		       title,
		       COALESCE(description,'') AS description,
		       COALESCE(filepath,'') AS filepath,
		       COALESCE(thumbnail,'') AS thumbnail,
		       COALESCE(version_no,'') AS version_no,
		       COALESCE(validity_status,'') AS validity_status,
		       issue_date,
		       effective_date
		FROM qh_planning_documents
		WHERE planning_project_id = ?
		  AND deleted_at IS NULL
		ORDER BY issue_date DESC NULLS LAST, id DESC`, row.ID).Scan(&documentRows).Error; err != nil {
		return nil, fmt.Errorf("load planning project documents: %w", err)
	}

	documents := make([]domain.PlanningDocument, 0, len(documentRows))
	for _, item := range documentRows {
		documents = append(documents, domain.PlanningDocument{
			ID:               item.ID,
			DocumentType:     item.DocumentType,
			DocumentTypeName: enums.GetPlanningDocumentTypeLabel(item.DocumentType),
			Status:           item.Status,
			StatusName:       map[bool]string{true: "Quan trọng", false: "Thông thường"}[item.Status == 20],
			Code:             item.Code,
			Title:            item.Title,
			Description:      item.Description,
			Filepath:         item.Filepath,
			Thumbnail:        item.Thumbnail,
			VersionNo:        item.VersionNo,
			ValidityStatus:   item.ValidityStatus,
			IssueDate:        item.IssueDate,
			EffectiveDate:    item.EffectiveDate,
		})
	}

	var layerRows []layerRow
	if err := r.db.WithContext(ctx).Raw(`
		SELECT l.id,
		       l.planning_project_id,
		       COALESCE(l.name,'') AS name,
		       COALESCE(l.display_name,'') AS display_name,
		       COALESCE(l.description,'') AS description,
		       l.type,
		       l.status,
		       (l.visible = 1) AS visible,
		       COALESCE(l.default_visible,FALSE) AS default_visible,
		       l.min_zoom,
		       l.max_zoom,
		       COALESCE(l.image_url,'') AS image_url,
		       COALESCE(l.thumbnail_url,'') AS thumbnail_url,
		       l.effective_date,
		       l.expiry_date,
		       l.updated_at
		FROM qh_layers l
		WHERE l.planning_project_id = ?
		  AND l.deleted_at IS NULL
		ORDER BY l.default_visible DESC, l.display_order ASC, l.id ASC`, row.ID).Scan(&layerRows).Error; err != nil {
		return nil, fmt.Errorf("load planning project layers: %w", err)
	}

	layers := make([]domain.PlanningLayer, 0, len(layerRows))
	for _, item := range layerRows {
		layers = append(layers, domain.PlanningLayer{
			ID:                item.ID,
			PlanningProjectID: item.PlanningProjectID,
			Name:              item.Name,
			DisplayName:       item.DisplayName,
			Description:       item.Description,
			Type:              item.Type,
			Status:            item.Status,
			Visible:           item.Visible,
			DefaultVisible:    item.DefaultVisible,
			MinZoom:           item.MinZoom,
			MaxZoom:           item.MaxZoom,
			ImageURL:          item.ImageURL,
			ThumbnailURL:      item.ThumbnailURL,
			EffectiveDate:     item.EffectiveDate,
			ExpiryDate:        item.ExpiryDate,
			UpdatedAt:         item.UpdatedAt,
		})
	}

	var relationRows []relationRow
	if err := r.db.WithContext(ctx).Raw(`
		SELECT related.id,
		       COALESCE(related.code,'') AS code,
		       related.name,
		       COALESCE(related.metadata::text,'{}') AS metadata
		FROM qh_planning_relations relation
		JOIN qh_planning_projects related
		  ON related.id = relation.related_planning_project_id
		 AND related.deleted_at IS NULL
		WHERE relation.planning_project_id = ?
		  AND relation.deleted_at IS NULL
		ORDER BY related.name`, row.ID).Scan(&relationRows).Error; err != nil {
		return nil, fmt.Errorf("load related planning projects: %w", err)
	}

	relations := make([]domain.PlanningProjectRelation, 0, len(relationRows))
	for _, item := range relationRows {
		relatedSlug := projectSlug(item.ID, "", item.Code, item.Name, item.Metadata)
		relations = append(relations, domain.PlanningProjectRelation{
			ID:         item.ID,
			Slug:       relatedSlug,
			Code:       item.Code,
			Name:       item.Name,
			PublicPath: "/do-an-quy-hoach/" + relatedSlug,
			MapPath:    "/ban-do?projectId=" + strconv.FormatUint(item.ID, 10),
		})
	}

	limitations := []string{}
	dataQuality := "complete"
	if row.Summary == "" {
		dataQuality = "partial"
		limitations = append(limitations, "Đồ án chưa có nội dung tóm tắt.")
	}
	if strings.TrimSpace(researchScope) == "" {
		dataQuality = "partial"
		limitations = append(limitations, "Chưa có phạm vi nghiên cứu được cấu trúc trong nguồn hiện tại.")
	}
	if len(indicators) == 0 {
		dataQuality = "partial"
		limitations = append(limitations, "Chưa có chỉ tiêu quy hoạch được cấu trúc trong nguồn hiện tại.")
	}
	if len(events) == 0 {
		dataQuality = "partial"
		limitations = append(limitations, "Chưa có lịch sử pháp lý được công bố.")
	}
	if row.MinLon == 0 && row.MaxLon == 0 {
		dataQuality = "partial"
		limitations = append(limitations, "Chưa có geometry tổng hợp để dựng bản đồ ngữ cảnh.")
	}

	return &domain.PlanningProjectProjection{
		ID:                row.ID,
		CanonicalKey:      "planning_project:" + strconv.FormatUint(row.ID, 10),
		Slug:              slug,
		CanonicalPath:     "/do-an-quy-hoach/" + slug,
		Code:              row.Code,
		Name:              row.Name,
		PlanningType:      row.PlanningType,
		PlanningTypeName:  enums.GetPlanningTypeLabel(row.PlanningType),
		PlanningLevel:     row.PlanningLevel,
		PlanningLevelName: enums.GetPlanningLevelLabel(row.PlanningLevel),
		TotalArea:         row.TotalArea,
		LegalStatus:       row.LegalStatus,
		LegalStatusName:   enums.LegalStatus(row.LegalStatus).DisplayName(),
		ValidityStatus:    row.ValidityStatus,
		Authority:         row.Authority,
		DecisionNumber:    row.DecisionNumber,
		Summary:           row.Summary,
		ResearchScope:     researchScope,
		ApprovalDate:      row.ApprovalDate,
		EffectiveDate:     row.EffectiveDate,
		ExpiryDate:        row.ExpiryDate,
		CurrentVersion:    row.CurrentVersion,
		JurisdictionID:    row.JurisdictionID,
		JurisdictionName:  row.JurisdictionName,
		Preview:           preview(row.GeometryType, row.CenterLat, row.CenterLon, row.MinLon, row.MinLat, row.MaxLon, row.MaxLat),
		Indicators:        indicators,
		Events:            events,
		Documents:         documents,
		Layers:            layers,
		RelatedProjects:   relations,
		Source: domain.ProjectionSource{
			System:         "tqd-service",
			Dataset:        "qh_planning_projects",
			Authority:      "aggregated",
			AuthorityName:  row.Authority,
			AuthorityClass: "aggregated",
			DataQuality:    dataQuality,
			UpdatedAt:      row.UpdatedAt,
			Limitations:    limitations,
		},
	}, nil
}

func IsNotFound(err error) bool {
	return errors.Is(err, gorm.ErrRecordNotFound)
}
