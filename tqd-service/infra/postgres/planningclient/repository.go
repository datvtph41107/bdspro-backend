package planningclient

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"tqd/internal/enums"
	planningdomain "tqd/internal/modules/planningclient/domain"
	tqdutils "tqd/internal/utils"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func normalizeJSON(raw json.RawMessage, fallback string) json.RawMessage {
	if len(raw) == 0 || !json.Valid(raw) {
		return json.RawMessage(fallback)
	}
	return raw
}

func normalizePage(page, size int) (int, int) {
	if page < 0 {
		page = 0
	}
	if size <= 0 {
		size = 20
	}
	if size > 100 {
		size = 100
	}
	return page, size
}

func parseCSV(value string) []string {
	parts := strings.Split(value, ",")
	result := make([]string, 0, len(parts))
	for _, part := range parts {
		if item := strings.TrimSpace(part); item != "" {
			result = append(result, item)
		}
	}
	return result
}

func parseUintCSV(value string) []uint64 {
	parts := parseCSV(value)
	result := make([]uint64, 0, len(parts))
	for _, part := range parts {
		parsed, err := strconv.ParseUint(part, 10, 64)
		if err == nil {
			result = append(result, parsed)
		}
	}
	return result
}

func lowerSearchPattern(value string) string {
	value = strings.ToLower(strings.TrimSpace(value))
	if value == "" {
		return ""
	}
	return "%" + value + "%"
}

func unaccentLike(column string) string {
	return "REPLACE(LOWER(unaccent(COALESCE(" + column + ",''))), 'đ', 'd') LIKE REPLACE(LOWER(unaccent(?)), 'đ', 'd')"
}

func slugSearchLike() string {
	return "REPLACE(REPLACE(LOWER(unaccent(COALESCE(p.slug,''))), '-', ' '), 'đ', 'd') LIKE REPLACE(REPLACE(LOWER(unaccent(?)), '-', ' '), 'đ', 'd')"
}

func applyProjectSearch(db *gorm.DB, search string, includeDocuments bool) *gorm.DB {
	pattern := lowerSearchPattern(search)
	if pattern == "" {
		return db
	}
	clause := "(" +
		unaccentLike("p.name") + " OR " +
		unaccentLike("p.code") + " OR " +
		unaccentLike("p.decision_number") + " OR " +
		slugSearchLike() + " OR " +
		unaccentLike("j.name")
	args := []any{pattern, pattern, pattern, pattern, pattern}
	if includeDocuments {
		clause += ` OR EXISTS (
			SELECT 1 FROM qh_planning_documents d
			WHERE d.planning_project_id = p.id AND d.deleted_at IS NULL
			  AND (` + unaccentLike("d.title") + ` OR ` + unaccentLike("d.code") + `)
		)`
		args = append(args, pattern, pattern)
	}
	clause += ")"
	return db.Where(clause, args...)
}

func applyProjectFilters(db *gorm.DB, filter planningdomain.ProjectListFilter) *gorm.DB {
	db = applyProjectSearch(db, filter.Search, true)
	if value := strings.TrimSpace(filter.Decision); value != "" {
		pattern := "%" + strings.ToLower(value) + "%"
		db = db.Where(`(
			LOWER(COALESCE(p.decision_number,'')) LIKE ? OR
			EXISTS (
				SELECT 1 FROM qh_planning_documents d
				WHERE d.planning_project_id = p.id AND d.deleted_at IS NULL
				  AND (LOWER(COALESCE(d.title,'')) LIKE ? OR LOWER(COALESCE(d.code,'')) LIKE ?)
			)
		)`, pattern, pattern, pattern)
	}
	if values := parseUintCSV(filter.PlanningType); len(values) > 0 {
		db = db.Where("p.planning_type IN ?", values)
	}
	if values := parseUintCSV(filter.PlanningLevel); len(values) > 0 {
		db = db.Where("p.planning_level IN ?", values)
	}
	if values := parseCSV(filter.ValidityStatus); len(values) > 0 {
		db = db.Where("p.validity_status IN ?", values)
	}
	if values := parseUintCSV(filter.ProcessStatus); len(values) > 0 {
		db = db.Where("p.process_status IN ?", values)
	}
	if values := parseUintCSV(filter.LegalStatus); len(values) > 0 {
		db = db.Where("p.legal_status IN ?", values)
	}
	if filter.JurisdictionID > 0 {
		db = db.Where("p.jurisdiction_id = ?", filter.JurisdictionID)
	}
	if value := strings.TrimSpace(filter.Area); value != "" {
		if id, err := strconv.ParseUint(value, 10, 64); err == nil && id > 0 {
			db = db.Where("p.jurisdiction_id = ?", id)
		} else {
			db = db.Where("LOWER(COALESCE(j.name,'')) LIKE ?", "%"+strings.ToLower(value)+"%")
		}
	}
	if filter.HasLayer != nil {
		clauseText := `EXISTS (
			SELECT 1 FROM qh_layers l
			WHERE l.planning_project_id = p.id AND l.deleted_at IS NULL
		)`
		if *filter.HasLayer {
			db = db.Where(clauseText)
		} else {
			db = db.Where("NOT " + clauseText)
		}
	}
	if filter.UpdatedFrom != nil {
		db = db.Where("p.updated_at >= ?", *filter.UpdatedFrom)
	}
	if filter.UpdatedTo != nil {
		db = db.Where("p.updated_at <= ?", *filter.UpdatedTo)
	}
	if values := parseUintCSV(filter.TagIDs); len(values) > 0 {
		db = db.Where(`EXISTS (
			SELECT 1 FROM qh_planning_project_tags ppt
			WHERE ppt.planning_project_id = p.id
			  AND ppt.deleted_at IS NULL
			  AND ppt.tag_id IN ?
		)`, values)
	}
	if values := parseCSV(filter.ProvinceIDs); len(values) > 0 {
		db = db.Where(`EXISTS (
			SELECT 1 FROM qh_planning_project_provinces ppp
			WHERE ppp.planning_project_id = p.id
			  AND ppp.deleted_at IS NULL
			  AND ppp.province_id IN ?
		)`, values)
	}
	if clause, arg, ok := tqdutils.PlanningProjectAdmJoinClause(filter.AdmID, "p.id"); ok {
		db = db.Where(clause, arg)
	}
	return db
}

func projectOrder(sortBy, sortOrder string) string {
	column := map[string]string{
		"updatedat":     "p.updated_at",
		"updated_at":    "p.updated_at",
		"createdat":     "p.created_at",
		"created_at":    "p.created_at",
		"approvaldate":  "p.approval_date",
		"approval_date": "p.approval_date",
		"time":          "p.updated_at",
		"date":          "p.updated_at",
		"name":          "p.name",
		"totalarea":     "p.total_area",
		"total_area":    "p.total_area",
	}[strings.ToLower(strings.TrimSpace(sortBy))]
	if column == "" {
		column = "p.updated_at"
	}
	order := strings.ToUpper(strings.TrimSpace(sortOrder))
	if order != "ASC" {
		order = "DESC"
	}
	if column == "p.name" {
		return column + " " + order + ", p.id " + order
	}
	return column + " " + order + " NULLS LAST, p.id " + order
}

const projectSelect = `
	p.id,
	COALESCE(p.slug,'') AS slug,
	COALESCE(p.code,'') AS code,
	p.name,
	p.planning_type,
	p.planning_level,
	COALESCE(p.total_area,0) AS total_area,
	p.jurisdiction_id,
	COALESCE(j.name,'') AS jurisdiction_name,
	COALESCE(p.legal_status,800) AS legal_status,
	COALESCE(p.validity_status,'') AS validity_status,
	COALESCE(p.authority,'') AS authority,
	COALESCE(p.decision_number,'') AS decision_number,
	COALESCE(p.summary,'') AS summary,
	COALESCE(p.research_scope,'') AS research_scope,
	COALESCE(p.indicators,'[]'::jsonb) AS indicators,
	p.approval_date,
	p.effective_date,
	p.expiry_date,
	COALESCE(p.current_version,'') AS current_version,
	COALESCE(p.metadata,'{}'::jsonb) AS metadata,
	COALESCE(p.process_status,10) AS process_status,
	p.created_at,
	p.updated_at,
	EXISTS(SELECT 1 FROM qh_layers l WHERE l.planning_project_id=p.id AND l.deleted_at IS NULL) AS has_layer,
	(SELECT COUNT(*) FROM qh_layers l WHERE l.planning_project_id=p.id AND l.deleted_at IS NULL) AS layer_count`

func (r *Repository) ListProjects(ctx context.Context, filter planningdomain.ProjectListFilter) (*planningdomain.Page[planningdomain.Project], error) {
	page, size := normalizePage(filter.Page, filter.Size)
	base := r.db.WithContext(ctx).
		Table("qh_planning_projects p").
		Joins("LEFT JOIN qh_jurisdictions j ON j.id = p.jurisdiction_id AND j.deleted_at IS NULL").
		Where("p.deleted_at IS NULL AND p.client_view = ?", true)
	base = applyProjectFilters(base, filter)

	var total int64
	if err := base.Count(&total).Error; err != nil {
		return nil, fmt.Errorf("count planning projects: %w", err)
	}

	items := []planningdomain.Project{}
	if err := base.Select(projectSelect).
		Order(projectOrder(filter.SortBy, filter.SortOrder)).
		Offset(page * size).
		Limit(size).
		Scan(&items).Error; err != nil {
		return nil, fmt.Errorf("list planning projects: %w", err)
	}
	for i := range items {
		r.decorateProject(&items[i])
	}
	if err := r.attachProjectRelations(ctx, items); err != nil {
		return nil, err
	}
	return &planningdomain.Page[planningdomain.Project]{Data: items, Total: total, Page: page, Size: size}, nil
}

func (r *Repository) GetProject(ctx context.Context, id uint64) (*planningdomain.Project, error) {
	var item planningdomain.Project
	err := r.db.WithContext(ctx).
		Table("qh_planning_projects p").
		Joins("LEFT JOIN qh_jurisdictions j ON j.id = p.jurisdiction_id AND j.deleted_at IS NULL").
		Where("p.deleted_at IS NULL AND p.client_view = ? AND p.id = ?", true, id).
		Select(projectSelect).
		Take(&item).Error
	if err != nil {
		return nil, err
	}

	var stats struct {
		DocumentCount          int64
		ImportantDocumentCount int64
		LegalEventCount        int64
		LatestUpdatedAt        *time.Time
	}
	if err := r.db.WithContext(ctx).Raw(`
		SELECT
		  (SELECT COUNT(*) FROM qh_planning_documents d WHERE d.planning_project_id=? AND d.deleted_at IS NULL) AS document_count,
		  (SELECT COUNT(*) FROM qh_planning_documents d WHERE d.planning_project_id=? AND d.deleted_at IS NULL AND d.status=20) AS important_document_count,
		  (SELECT COUNT(*) FROM qh_planning_events e WHERE e.planning_project_id=? AND e.deleted_at IS NULL) AS legal_event_count,
		  GREATEST(
		    p.updated_at,
		    COALESCE((SELECT MAX(d.updated_at) FROM qh_planning_documents d WHERE d.planning_project_id=p.id AND d.deleted_at IS NULL), p.updated_at),
		    COALESCE((SELECT MAX(e.updated_at) FROM qh_planning_events e WHERE e.planning_project_id=p.id AND e.deleted_at IS NULL), p.updated_at)
		  ) AS latest_updated_at
		FROM qh_planning_projects p WHERE p.id=?`, id, id, id, id).Scan(&stats).Error; err != nil {
		return nil, fmt.Errorf("load planning project overview: %w", err)
	}

	r.decorateProject(&item)
	item.Overview = &planningdomain.ProjectOverview{
		Description:            item.Summary,
		TotalArea:              item.TotalArea,
		ResearchScope:          item.ResearchScope,
		Indicators:             normalizeJSON(item.Indicators, "[]"),
		DocumentCount:          stats.DocumentCount,
		ImportantDocumentCount: stats.ImportantDocumentCount,
		MapLayerCount:          item.LayerCount,
		LegalEventCount:        stats.LegalEventCount,
		LatestUpdatedAt:        stats.LatestUpdatedAt,
		HasLayer:               item.HasLayer,
	}
	projects := []planningdomain.Project{item}
	if err := r.attachProjectRelations(ctx, projects); err != nil {
		return nil, err
	}
	item = projects[0]
	return &item, nil
}

func (r *Repository) decorateProject(item *planningdomain.Project) {
	item.Metadata = normalizeJSON(item.Metadata, "{}")
	item.Indicators = normalizeJSON(item.Indicators, "[]")
	item.PlanningTypeName = enums.GetPlanningTypeLabel(item.PlanningType)
	item.PlanningLevelName = enums.GetPlanningLevelLabel(item.PlanningLevel)
	item.LegalStatusName = enums.LegalStatus(item.LegalStatus).DisplayName()
	item.Tags = []planningdomain.ProjectTag{
		{Key: "planningType", Label: "Loại quy hoạch", Value: item.PlanningTypeName},
		{Key: "totalArea", Label: "Diện tích", Value: formatArea(item.TotalArea)},
		{Key: "legalStatus", Label: "Trạng thái", Value: item.LegalStatusName},
		{Key: "jurisdiction", Label: "Địa bàn", Value: item.JurisdictionName},
		{Key: "approvalDate", Label: "Ngày phê duyệt", Value: formatDate(item.ApprovalDate)},
		{Key: "updatedAt", Label: "Ngày cập nhật", Value: formatDate(item.UpdatedAt)},
	}
}

func (r *Repository) attachProjectProvinces(ctx context.Context, items []planningdomain.Project) error {
	if len(items) == 0 {
		return nil
	}
	projectIDs := make([]uint64, 0, len(items))
	indexByID := make(map[uint64]int, len(items))
	for i := range items {
		if items[i].ID == 0 {
			continue
		}
		indexByID[items[i].ID] = i
		projectIDs = append(projectIDs, items[i].ID)
	}
	if len(projectIDs) == 0 {
		return nil
	}

	type projectProvinceRow struct {
		PlanningProjectID uint64
		ID                string
		FullName          string
		ShortName         string
		Code              string
		QHOrder           int32
		Zoom              uint32
		Lat               float64
		Lng               float64
	}
	rows := []projectProvinceRow{}
	// province_v2.id là bigint; qh_planning_project_provinces.province_id vẫn VARCHAR(50)
	// (lưu dạng text của id). Cast để tránh SQLSTATE 42883 bigint = character varying.
	err := r.db.WithContext(ctx).
		Table("qh_planning_project_provinces ppp").
		Joins("JOIN province_v2 pv ON pv.id::text = ppp.province_id").
		Where("ppp.deleted_at IS NULL AND ppp.planning_project_id IN ?", projectIDs).
		Select(`ppp.planning_project_id,
			pv.id::text AS id,
			pv.full_name,
			COALESCE(pv.short_name,'') AS short_name,
			COALESCE(pv.code,'') AS code,
			COALESCE(pv.qh_order,0) AS qh_order,
			COALESCE(pv.zoom,10) AS zoom,
			COALESCE(pv.lat,0) AS lat,
			COALESCE(pv.lng,0) AS lng`).
		Order("pv.qh_order DESC, pv.full_name ASC, pv.id ASC").
		Scan(&rows).Error
	if err != nil {
		return fmt.Errorf("load planning project provinces: %w", err)
	}
	for _, row := range rows {
		index, ok := indexByID[row.PlanningProjectID]
		if !ok {
			continue
		}
		items[index].Provinces = append(items[index].Provinces, planningdomain.ProjectProvince{
			ID:        row.ID,
			FullName:  row.FullName,
			ShortName: row.ShortName,
			Code:      row.Code,
			QHOrder:   row.QHOrder,
			Zoom:      row.Zoom,
			Lat:       row.Lat,
			Lng:       row.Lng,
		})
	}
	return nil
}

func (r *Repository) attachProjectRelations(ctx context.Context, items []planningdomain.Project) error {
	return r.attachProjectProvinces(ctx, items)
}

func formatArea(value float64) string {
	if value <= 0 {
		return "Chưa cập nhật"
	}
	return strconv.FormatFloat(value, 'f', -1, 64) + " ha"
}

func formatDate(value *time.Time) string {
	if value == nil || value.IsZero() {
		return "Chưa cập nhật"
	}
	return value.Format("02/01/2006")
}

func documentStatusName(status uint32) string {
	if status == planningdomain.DocumentStatusImportant {
		return "Quan trọng"
	}
	return "Thường"
}

func (r *Repository) ListDocuments(ctx context.Context, projectID uint64, page, size int, keyword, documentType, status string) (*planningdomain.Page[planningdomain.Document], error) {
	page, size = normalizePage(page, size)
	base := r.db.WithContext(ctx).
		Table("qh_planning_documents d").
		Joins("JOIN qh_planning_projects p ON p.id = d.planning_project_id AND p.deleted_at IS NULL AND p.client_view = ?", true).
		Where("d.deleted_at IS NULL AND d.planning_project_id = ?", projectID)
	if value := strings.TrimSpace(keyword); value != "" {
		pattern := "%" + strings.ToLower(value) + "%"
		base = base.Where(`LOWER(COALESCE(d.title,'')) LIKE ? OR LOWER(COALESCE(d.code,'')) LIKE ? OR LOWER(COALESCE(d.description,'')) LIKE ?`, pattern, pattern, pattern)
	}
	if values := parseCSV(documentType); len(values) > 0 {
		base = base.Where("d.document_type IN ?", values)
	}
	if values := parseCSV(status); len(values) > 0 {
		base = base.Where("d.status IN ?", values)
	}
	var total int64
	if err := base.Count(&total).Error; err != nil {
		return nil, err
	}
	items := []planningdomain.Document{}
	if err := base.Select(documentSelect).
		Order("CASE WHEN d.status=20 THEN 0 ELSE 1 END, d.issue_date DESC NULLS LAST, d.updated_at DESC, d.id DESC").
		Offset(page * size).Limit(size).Scan(&items).Error; err != nil {
		return nil, err
	}
	for i := range items {
		r.decorateDocument(&items[i])
	}
	return &planningdomain.Page[planningdomain.Document]{Data: items, Total: total, Page: page, Size: size}, nil
}

const documentSelect = `d.id,d.planning_project_id,d.document_type,COALESCE(d.status,10) AS status,
		COALESCE(d.code,'') AS code,d.title,COALESCE(d.description,'') AS description,
		COALESCE(d.filepath,'') AS filepath,COALESCE(d.relative_path,'') AS relative_path,
		COALESCE(d.ext,'') AS ext,
		COALESCE(d.thumbnail,'') AS thumbnail,
		COALESCE(d.version_no,'') AS version_no,COALESCE(d.validity_status,'') AS validity_status,
		d.issue_date,d.effective_date,
		COALESCE(d.indicators,'[]'::jsonb) AS indicators,
		COALESCE(d.authority,'') AS authority,
		COALESCE(d.sector,'') AS sector,
		COALESCE(d.field,'') AS field,
		COALESCE(d.signer,'') AS signer,
		COALESCE(d.metadata,'{}'::jsonb) AS metadata,d.created_at,d.updated_at`

func (r *Repository) decorateDocument(item *planningdomain.Document) {
	item.Metadata = normalizeJSON(item.Metadata, "{}")
	item.Indicators = normalizeJSON(item.Indicators, "[]")
	item.DocumentTypeName = enums.GetPlanningDocumentTypeLabel(item.DocumentType)
	item.StatusName = documentStatusName(item.Status)
	if strings.TrimSpace(item.Ext) == "" {
		item.Ext = tqdutils.ResolveDocumentExt(item.RelativePath, item.Filepath)
	}
}

func (r *Repository) ListDocumentsByParcelRelation(ctx context.Context, parcelID uint64, page, size int, docType string) (*planningdomain.Page[planningdomain.Document], error) {
	if parcelID == 0 {
		return nil, errors.New("parcel id is required")
	}
	page, size = normalizePage(page, size)
	docTypes := parseUintCSV(docType)

	// parcel → regions (intersects) → layers → layer_plan_doc_relation → planning-documents
	baseSQL := `
		WITH parcel_layers AS (
			SELECT DISTINCT r.layer_id
			FROM parcels p
			JOIN qh_regions r
			  ON r.geometry && p.geometry
			 AND ST_Intersects(r.geometry, p.geometry)
			JOIN qh_layers l
			  ON l.id = r.layer_id
			 AND l.deleted_at IS NULL
			 AND l.status = 10
			WHERE p.id = ?
			  AND p.deleted_at IS NULL
			  AND p.geometry IS NOT NULL
			  AND r.deleted_at IS NULL
			  AND r.status = 10
			  AND r.is_latest = true
			  AND r.geometry IS NOT NULL
		),
		ranked AS (
			SELECT
				rel.planning_document_id,
				MIN(rel.sort_order) AS sort_order,
				MIN(rel.id) AS rel_id
			FROM layer_plan_doc_relation rel
			JOIN parcel_layers pl ON pl.layer_id = rel.layer_id
			WHERE rel.deleted_at IS NULL
			GROUP BY rel.planning_document_id
		)
		SELECT %s
		FROM ranked rr
		JOIN qh_planning_documents d ON d.id = rr.planning_document_id AND d.deleted_at IS NULL
		JOIN qh_planning_projects p ON p.id = d.planning_project_id AND p.deleted_at IS NULL AND p.client_view = true
	`
	args := []any{parcelID}
	if len(docTypes) > 0 {
		baseSQL += ` WHERE d.document_type IN ?`
		args = append(args, docTypes)
	}

	countSQL := fmt.Sprintf(baseSQL, "COUNT(*)")
	var total int64
	if err := r.db.WithContext(ctx).Raw(countSQL, args...).Scan(&total).Error; err != nil {
		return nil, fmt.Errorf("count layer plan doc relations for parcel: %w", err)
	}
	if total == 0 {
		return &planningdomain.Page[planningdomain.Document]{Data: []planningdomain.Document{}, Total: 0, Page: page, Size: size}, nil
	}

	listSQL := fmt.Sprintf(baseSQL, documentSelect) + `
		ORDER BY rr.sort_order ASC, rr.rel_id ASC
		OFFSET ? LIMIT ?
	`
	listArgs := append(append([]any{}, args...), page*size, size)
	items := []planningdomain.Document{}
	if err := r.db.WithContext(ctx).Raw(listSQL, listArgs...).Scan(&items).Error; err != nil {
		return nil, fmt.Errorf("list layer plan doc relations for parcel: %w", err)
	}
	for i := range items {
		r.decorateDocument(&items[i])
	}
	return &planningdomain.Page[planningdomain.Document]{Data: items, Total: total, Page: page, Size: size}, nil
}

func (r *Repository) GetDocument(ctx context.Context, documentID uint64) (*planningdomain.Document, error) {
	var row struct {
		planningdomain.Document
		ProjectRefID   uint64 `gorm:"column:project_ref_id"`
		ProjectRefName string `gorm:"column:project_ref_name"`
	}
	err := r.db.WithContext(ctx).Table("qh_planning_documents d").
		Joins("LEFT JOIN qh_planning_projects p ON p.id = d.planning_project_id AND p.deleted_at IS NULL AND p.client_view = ?", true).
		Where("d.id = ? AND d.deleted_at IS NULL", documentID).
		Select(documentSelect+`, p.id AS project_ref_id, COALESCE(p.name, '') AS project_ref_name`).
		Take(&row).Error
	if err != nil {
		return nil, err
	}
	r.decorateDocument(&row.Document)
	if row.ProjectRefID > 0 {
		row.Document.PlanningProject = &planningdomain.DocumentPlanningProject{
			ID:   row.ProjectRefID,
			Name: row.ProjectRefName,
		}
	} else if row.Document.PlanningProjectID > 0 {
		row.Document.PlanningProject = &planningdomain.DocumentPlanningProject{
			ID: row.Document.PlanningProjectID,
		}
	}
	return &row.Document, nil
}

func (r *Repository) ListEvents(ctx context.Context, projectID uint64, page, size int) (*planningdomain.Page[planningdomain.Event], error) {
	page, size = normalizePage(page, size)
	base := r.db.WithContext(ctx).
		Table("qh_planning_events e").
		Joins("JOIN qh_planning_projects p ON p.id = e.planning_project_id AND p.deleted_at IS NULL AND p.client_view = ?", true).
		Where("e.deleted_at IS NULL AND e.planning_project_id = ?", projectID)
	var total int64
	if err := base.Count(&total).Error; err != nil {
		return nil, err
	}
	items := []planningdomain.Event{}
	if err := base.Select(`e.id,e.planning_project_id,COALESCE(e.event_type,'') AS event_type,e.event_name,
		COALESCE(e.description,'') AS description,e.event_date,e.document_id,COALESCE(e.source_url,'') AS source_url,
		COALESCE(e.metadata,'{}'::jsonb) AS metadata,e.created_at,e.updated_at`).
		Order("e.event_date DESC NULLS LAST, e.id DESC").Offset(page * size).Limit(size).Scan(&items).Error; err != nil {
		return nil, err
	}
	for i := range items {
		items[i].Metadata = normalizeJSON(items[i].Metadata, "{}")
		if items[i].DocumentID != nil {
			var doc planningdomain.Document
			err := r.db.WithContext(ctx).Table("qh_planning_documents d").
				Where("d.id=? AND d.deleted_at IS NULL", *items[i].DocumentID).
				Select(documentSelect).
				Take(&doc).Error
			if err == nil {
				r.decorateDocument(&doc)
				items[i].Document = &doc
			}
		}
	}
	return &planningdomain.Page[planningdomain.Event]{Data: items, Total: total, Page: page, Size: size}, nil
}

type followRow struct {
	ID                uint64
	UserID            uint64
	PlanningProjectID uint64
	Note              string
	UpdatedAt         *time.Time
}

func (r *Repository) Follow(ctx context.Context, userID, projectID uint64, note string) (*planningdomain.Follow, error) {
	if userID == 0 || projectID == 0 {
		return nil, errors.New("userID and planningProjectID are required")
	}
	var exists int64
	if err := r.db.WithContext(ctx).Table("qh_planning_projects").Where("id=? AND deleted_at IS NULL AND client_view = ?", projectID, true).Count(&exists).Error; err != nil {
		return nil, err
	}
	if exists == 0 {
		return nil, gorm.ErrRecordNotFound
	}
	row := map[string]any{
		"user_id": userID, "planning_project_id": projectID, "note": strings.TrimSpace(note),
		"created_at": time.Now(), "updated_at": time.Now(), "deleted_at": nil,
	}
	if err := r.db.WithContext(ctx).Table("user_followed_planning_projects").
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "user_id"}, {Name: "planning_project_id"}},
			DoUpdates: clause.Assignments(map[string]any{"note": strings.TrimSpace(note), "updated_at": time.Now(), "deleted_at": nil}),
		}).Create(row).Error; err != nil {
		return nil, err
	}
	return r.GetFollow(ctx, userID, projectID)
}

func (r *Repository) GetFollow(ctx context.Context, userID, projectID uint64) (*planningdomain.Follow, error) {
	var row followRow
	err := r.db.WithContext(ctx).Table("user_followed_planning_projects").
		Where("user_id=? AND planning_project_id=? AND deleted_at IS NULL", userID, projectID).
		Select("id,user_id,planning_project_id,COALESCE(note,'') AS note,updated_at").Take(&row).Error
	if err != nil {
		return nil, err
	}
	return &planningdomain.Follow{ID: row.ID, UserID: row.UserID, PlanningProjectID: row.PlanningProjectID, Note: row.Note, FollowedAt: row.UpdatedAt}, nil
}

func (r *Repository) Unfollow(ctx context.Context, userID, projectID uint64) error {
	tx := r.db.WithContext(ctx).Table("user_followed_planning_projects").
		Where("user_id=? AND planning_project_id=? AND deleted_at IS NULL", userID, projectID).
		Updates(map[string]any{"deleted_at": time.Now(), "updated_at": time.Now()})
	if tx.Error != nil {
		return tx.Error
	}
	// Unfollow is idempotent: deleting an already-unfollowed relationship keeps
	// the desired state and therefore remains a successful command.
	return nil
}

func (r *Repository) ListFollowed(ctx context.Context, userID uint64, page, size int, search string) (*planningdomain.Page[planningdomain.FollowedProject], error) {
	page, size = normalizePage(page, size)
	base := r.db.WithContext(ctx).Table("user_followed_planning_projects f").
		Joins("JOIN qh_planning_projects p ON p.id=f.planning_project_id AND p.deleted_at IS NULL AND p.client_view = ?", true).
		Joins("LEFT JOIN qh_jurisdictions j ON j.id=p.jurisdiction_id AND j.deleted_at IS NULL").
		Where("f.user_id=? AND f.deleted_at IS NULL", userID)
	base = applyProjectSearch(base, search, false)
	var total int64
	if err := base.Count(&total).Error; err != nil {
		return nil, err
	}
	type followedRow struct {
		FollowID   uint64
		FollowedAt *time.Time
		Note       string
		planningdomain.Project
	}
	rows := []followedRow{}
	selectSQL := "f.id AS follow_id,f.updated_at AS followed_at,COALESCE(f.note,'') AS note," + projectSelect
	if err := base.Select(selectSQL).Order("f.updated_at DESC,f.id DESC").Offset(page * size).Limit(size).Scan(&rows).Error; err != nil {
		return nil, err
	}
	items := make([]planningdomain.FollowedProject, 0, len(rows))
	for i := range rows {
		r.decorateProject(&rows[i].Project)
		items = append(items, planningdomain.FollowedProject{FollowID: rows[i].FollowID, FollowedAt: rows[i].FollowedAt, Note: rows[i].Note, Project: rows[i].Project})
	}
	projects := make([]planningdomain.Project, len(items))
	for i := range items {
		projects[i] = items[i].Project
	}
	// if err := r.attachProjectTags(ctx, projects); err != nil {
	// 	return nil, err
	// }
	if err := r.attachProjectRelations(ctx, projects); err != nil {
		return nil, err
	}
	for i := range items {
		items[i].Project = projects[i]
	}
	return &planningdomain.Page[planningdomain.FollowedProject]{Data: items, Total: total, Page: page, Size: size}, nil
}

func (r *Repository) ListProjectsByLayerRelation(ctx context.Context, layerIDs []uint64, page, size int) (*planningdomain.Page[planningdomain.Project], error) {
	if len(layerIDs) == 0 {
		return nil, errors.New("layer ids is required")
	}
	page, size = normalizePage(page, size)

	baseSQL := `
		WITH ranked AS (
			SELECT
				d.planning_project_id,
				MIN(rel.sort_order) AS sort_order,
				MIN(rel.id) AS rel_id
			FROM layer_plan_doc_relation rel
			JOIN qh_planning_documents d
			  ON d.id = rel.planning_document_id
			 AND d.deleted_at IS NULL
			JOIN qh_layers l
			  ON l.id = rel.layer_id
			 AND l.deleted_at IS NULL
			 AND l.status = 10
			WHERE rel.deleted_at IS NULL
			  AND rel.layer_id IN ?
			GROUP BY d.planning_project_id
		)
		SELECT %s
		FROM ranked rr
		JOIN qh_planning_projects p ON p.id = rr.planning_project_id AND p.deleted_at IS NULL AND p.client_view = true
		LEFT JOIN qh_jurisdictions j ON j.id = p.jurisdiction_id AND j.deleted_at IS NULL
	`
	return r.scanPlanningProjectsPage(ctx, baseSQL, []any{layerIDs}, page, size)
}

func (r *Repository) ListProjectsByParcelRelation(ctx context.Context, parcelID uint64, page, size int) (*planningdomain.Page[planningdomain.Project], error) {
	if parcelID == 0 {
		return nil, errors.New("parcel id is required")
	}
	page, size = normalizePage(page, size)

	// parcel → regions (intersects) → layers → layer_plan_doc_relation → documents → planning-projects
	baseSQL := `
		WITH parcel_layers AS (
			SELECT DISTINCT r.layer_id
			FROM parcels p
			JOIN qh_regions r
			  ON r.geometry && p.geometry
			 AND ST_Intersects(r.geometry, p.geometry)
			JOIN qh_layers l
			  ON l.id = r.layer_id
			 AND l.deleted_at IS NULL
			 AND l.status = 10
			WHERE p.id = ?
			  AND p.deleted_at IS NULL
			  AND p.geometry IS NOT NULL
			  AND r.deleted_at IS NULL
			  AND r.status = 10
			  AND r.is_latest = true
			  AND r.geometry IS NOT NULL
		),
		ranked AS (
			SELECT
				d.planning_project_id,
				MIN(rel.sort_order) AS sort_order,
				MIN(rel.id) AS rel_id
			FROM layer_plan_doc_relation rel
			JOIN parcel_layers pl ON pl.layer_id = rel.layer_id
			JOIN qh_planning_documents d
			  ON d.id = rel.planning_document_id
			 AND d.deleted_at IS NULL
			WHERE rel.deleted_at IS NULL
			GROUP BY d.planning_project_id
		)
		SELECT %s
		FROM ranked rr
		JOIN qh_planning_projects p ON p.id = rr.planning_project_id AND p.deleted_at IS NULL AND p.client_view = true
		LEFT JOIN qh_jurisdictions j ON j.id = p.jurisdiction_id AND j.deleted_at IS NULL
	`
	return r.scanPlanningProjectsPage(ctx, baseSQL, []any{parcelID}, page, size)
}

func (r *Repository) scanPlanningProjectsPage(ctx context.Context, baseSQL string, args []any, page, size int) (*planningdomain.Page[planningdomain.Project], error) {
	countSQL := fmt.Sprintf(baseSQL, "COUNT(*)")
	var total int64
	if err := r.db.WithContext(ctx).Raw(countSQL, args...).Scan(&total).Error; err != nil {
		return nil, err
	}
	if total == 0 {
		return &planningdomain.Page[planningdomain.Project]{Data: []planningdomain.Project{}, Total: 0, Page: page, Size: size}, nil
	}

	listSQL := fmt.Sprintf(baseSQL, projectSelect) + `
		ORDER BY rr.sort_order ASC, rr.rel_id ASC
		OFFSET ? LIMIT ?
	`
	listArgs := append(append([]any{}, args...), page*size, size)
	items := []planningdomain.Project{}
	if err := r.db.WithContext(ctx).Raw(listSQL, listArgs...).Scan(&items).Error; err != nil {
		return nil, err
	}
	for i := range items {
		r.decorateProject(&items[i])
	}
	if err := r.attachProjectRelations(ctx, items); err != nil {
		return nil, err
	}
	return &planningdomain.Page[planningdomain.Project]{Data: items, Total: total, Page: page, Size: size}, nil
}

func (r *Repository) ListProjectLayers(ctx context.Context, projectID uint64) ([]planningdomain.ProjectLayer, error) {
	if projectID == 0 {
		return nil, errors.New("planning project id is required")
	}
	var exists int64
	if err := r.db.WithContext(ctx).Table("qh_planning_projects").
		Where("id = ? AND deleted_at IS NULL AND client_view = ?", projectID, true).
		Count(&exists).Error; err != nil {
		return nil, err
	}
	if exists == 0 {
		return nil, gorm.ErrRecordNotFound
	}
	items := make([]planningdomain.ProjectLayer, 0)
	err := r.db.WithContext(ctx).Table("qh_layers l").
		Select(`l.id, l.planning_project_id, COALESCE(l.name,'') AS name,
			COALESCE(l.display_name,'') AS display_name,
			COALESCE(l.description,'') AS description, l.type, l.status,
			(l.visible = 1) AS visible, COALESCE(l.default_visible,FALSE) AS default_visible,
			COALESCE(l.display_order,0) AS display_order, l.min_zoom, l.max_zoom,
			COALESCE(l.source_code, l.name, '') AS source_code,
			COALESCE(l.source_type,'vector') AS source_type,
			COALESCE(l.layer_url,'') AS layer_url,
			COALESCE(l.style_config,'[]'::jsonb) AS style_config,
			COALESCE(l.image_url,'') AS image_url,
			COALESCE(l.thumbnail_url,'') AS thumbnail_url, l.effective_date,
			l.expiry_date, l.updated_at`).
		Where("l.deleted_at IS NULL AND l.planning_project_id = ?", projectID).
		Order("l.default_visible DESC, l.display_order ASC, l.id ASC").
		Scan(&items).Error
	if err != nil {
		return nil, fmt.Errorf("list planning project layers: %w", err)
	}
	for i := range items {
		items[i].StyleConfig = normalizeJSON(items[i].StyleConfig, "[]")
	}
	return items, nil
}

func (r *Repository) ListProvinces(ctx context.Context, filter planningdomain.ProvinceListFilter) ([]planningdomain.ProvinceItem, error) {
	db := r.db.WithContext(ctx).Table("province_v2")
	if filter.Scope == planningdomain.ProvinceScopeQH {
		db = db.Where("qh_order > 0")
	}
	if pattern := lowerSearchPattern(filter.Search); pattern != "" {
		code := strings.ToLower(strings.TrimSpace(filter.Search))
		db = db.Where(`(
			`+unaccentLike("full_name")+` OR
			`+unaccentLike("short_name")+` OR
			LOWER(COALESCE(code,'')) = ?
		)`, pattern, pattern, code)
	}

	items := []planningdomain.ProvinceItem{}
	err := db.Select(`id,
		full_name,
		COALESCE(short_name,'') AS short_name,
		COALESCE(code,'') AS code,
		COALESCE(lat,0) AS lat,
		COALESCE(lng,0) AS lng,
		COALESCE(qh_order,0) AS qh_order,
		COALESCE(zoom,10) AS zoom,
		COALESCE(description,'') AS description,
		COALESCE(unit_type,0) AS unit_type,
		COALESCE(data_source,0) AS data_source,
		plan_doc_id`).
		Order("qh_order DESC, full_name ASC, id ASC").
		Scan(&items).Error
	if err != nil {
		return nil, fmt.Errorf("list provinces: %w", err)
	}
	return items, nil
}
