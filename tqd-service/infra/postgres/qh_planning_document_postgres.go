package postgres

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	_db "common/db"
	qh_domain "tqd/internal/domain/qh"
	qh_dto "tqd/internal/dto/qh"
	"tqd/internal/enums"
	"tqd/internal/interface/repo"
	tqdutils "tqd/internal/utils"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type QHPlanningDocumentRepo struct {
	db *_db.TransactionRepo
}

// NewQHPlanningDocumentRepo @bind: internal/interface/repo.IQHPlanningDocumentRepo
func NewQHPlanningDocumentRepo(db *_db.TransactionRepo) repo.IQHPlanningDocumentRepo {
	return &QHPlanningDocumentRepo{db: db}
}

func (r *QHPlanningDocumentRepo) Create(ctx context.Context, entity *qh_domain.QHPlanningDocument) error {
	normalizePlanningDocumentExt(entity)
	return r.db.GetDB(ctx).Create(entity).Error
}

func (r *QHPlanningDocumentRepo) Update(ctx context.Context, id uint64, entity *qh_domain.QHPlanningDocument) error {
	normalizePlanningDocumentExt(entity)
	status := entity.Status
	if status == 0 {
		status = enums.PlanningDocumentStatusNormal
	}
	updates := map[string]interface{}{
		"planning_project_id": entity.PlanningProjectID,
		"document_type":       entity.DocumentType,
		"status":              status,
		"code":                entity.Code,
		"title":               entity.Title,
		"description":         entity.Description,
		"filepath":            entity.Filepath,
		"ext":                 entity.Ext,
		"thumbnail":           entity.Thumbnail,
		"version_no":          entity.VersionNo,
		"validity_status":     entity.ValidityStatus,
		"issue_date":          entity.IssueDate,
		"effective_date":      entity.EffectiveDate,
		"metadata":            postgresJSONValue(entity.Metadata, "{}"),
		"indicators":          postgresJSONValue(entity.Indicators, "[]"),
		"authority":           entity.Authority,
		"sector":              entity.Sector,
		"field":               entity.Field,
		"signer":              entity.Signer,
		"relative_path":       entity.RelativePath,
	}
	return r.db.GetDB(ctx).Model(&qh_domain.QHPlanningDocument{}).Where("id = ? AND deleted_at IS NULL", id).Updates(updates).Error
}

func normalizePlanningDocumentExt(entity *qh_domain.QHPlanningDocument) {
	if entity == nil {
		return
	}
	entity.Ext = tqdutils.ResolveDocumentExt(entity.RelativePath, entity.Filepath)
}

func postgresJSONValue(raw, fallback string) clause.Expr {
	value := strings.TrimSpace(raw)
	if value == "" || !json.Valid([]byte(value)) {
		value = fallback
	}
	return gorm.Expr("?::jsonb", value)
}

func (r *QHPlanningDocumentRepo) Delete(ctx context.Context, id uint64) error {
	return r.db.GetDB(ctx).Where("id = ?", id).Delete(&qh_domain.QHPlanningDocument{}).Error
}

func (r *QHPlanningDocumentRepo) GetByID(ctx context.Context, id uint64) (*qh_domain.QHPlanningDocument, error) {
	var entity qh_domain.QHPlanningDocument
	err := r.db.GetDB(ctx).Where("id = ? AND deleted_at IS NULL", id).First(&entity).Error
	if err != nil {
		return nil, err
	}
	return &entity, nil
}

func (r *QHPlanningDocumentRepo) GetDetail(ctx context.Context, id uint64) (*qh_domain.QHPlanningDocument, error) {
	return r.GetByID(ctx, id)
}

func (r *QHPlanningDocumentRepo) GetAll(ctx context.Context) ([]qh_domain.QHPlanningDocument, error) {
	var entities []qh_domain.QHPlanningDocument
	err := r.db.GetDB(ctx).Where("deleted_at IS NULL").Find(&entities).Error
	return entities, err
}

func (r *QHPlanningDocumentRepo) GetList(ctx context.Context, req *qh_dto.ListPlanningDocumentsRequest) ([]qh_domain.QHPlanningDocument, int64, error) {
	if req == nil {
		req = &qh_dto.ListPlanningDocumentsRequest{}
	}
	if req.ParcelID > 0 {
		return r.getListByParcelRelation(ctx, req)
	}

	var documents []qh_domain.QHPlanningDocument
	var total int64

	db := r.db.GetDB(ctx).Model(&qh_domain.QHPlanningDocument{}).Where("qh_planning_documents.deleted_at IS NULL")

	if req.PlanningProjectID != 0 {
		db = db.Where("qh_planning_documents.planning_project_id = ?", req.PlanningProjectID)
	}
	if req.DocumentType != "" {
		db = db.Where("qh_planning_documents.document_type = ?", req.DocumentType)
	}
	if req.ValidityStatus != "" {
		db = db.Where("qh_planning_documents.validity_status = ?", req.ValidityStatus)
	}
	if keyword := strings.TrimSpace(req.Keyword); keyword != "" {
		pattern := "%" + keyword + "%"
		db = db.Where(`(
			REPLACE(LOWER(unaccent(COALESCE(qh_planning_documents.title,''))), 'đ', 'd') LIKE REPLACE(LOWER(unaccent(?)), 'đ', 'd') OR
			REPLACE(LOWER(unaccent(COALESCE(qh_planning_documents.code,''))), 'đ', 'd') LIKE REPLACE(LOWER(unaccent(?)), 'đ', 'd') OR
			REPLACE(LOWER(unaccent(COALESCE(qh_planning_documents.description,''))), 'đ', 'd') LIKE REPLACE(LOWER(unaccent(?)), 'đ', 'd')
		)`, pattern, pattern, pattern)
	}
	layerIDs := tqdutils.MergePositiveUint64(req.LayerID, req.LayerIDs)
	if len(layerIDs) > 0 {
		// Lọc tài liệu gắn trực tiếp với layer qua layer_plan_doc_relation.
		db = db.Where(`EXISTS (
			SELECT 1 FROM layer_plan_doc_relation rel
			WHERE rel.planning_document_id = qh_planning_documents.id
			  AND rel.deleted_at IS NULL
			  AND rel.layer_id IN ?
		)`, layerIDs)
	}

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	order := "qh_planning_documents.updated_at DESC NULLS LAST, qh_planning_documents.id DESC"

	err := db.Offset(req.GetOffset()).
		Limit(req.GetLimit()).
		Order(order).
		Find(&documents).Error

	return documents, total, err
}

func (r *QHPlanningDocumentRepo) getListByParcelRelation(ctx context.Context, req *qh_dto.ListPlanningDocumentsRequest) ([]qh_domain.QHPlanningDocument, int64, error) {
	page := int(req.GetPage())
	size := int(req.GetLimit())
	if page < 0 {
		page = 0
	}
	if size <= 0 {
		size = 20
	}

	layerIDs := tqdutils.MergePositiveUint64(req.LayerID, req.LayerIDs)

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
	`
	args := []any{req.ParcelID}
	if len(layerIDs) > 0 {
		baseSQL += ` AND r.layer_id IN ?`
		args = append(args, layerIDs)
	}
	baseSQL += `
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
		WHERE 1=1
	`
	if req.PlanningProjectID != 0 {
		baseSQL += ` AND d.planning_project_id = ?`
		args = append(args, req.PlanningProjectID)
	}
	if req.DocumentType != "" {
		baseSQL += ` AND d.document_type = ?`
		args = append(args, req.DocumentType)
	}
	if req.ValidityStatus != "" {
		baseSQL += ` AND d.validity_status = ?`
		args = append(args, req.ValidityStatus)
	}
	if keyword := strings.TrimSpace(req.Keyword); keyword != "" {
		pattern := "%" + keyword + "%"
		baseSQL += ` AND (
			REPLACE(LOWER(unaccent(COALESCE(d.title,''))), 'đ', 'd') LIKE REPLACE(LOWER(unaccent(?)), 'đ', 'd') OR
			REPLACE(LOWER(unaccent(COALESCE(d.code,''))), 'đ', 'd') LIKE REPLACE(LOWER(unaccent(?)), 'đ', 'd') OR
			REPLACE(LOWER(unaccent(COALESCE(d.description,''))), 'đ', 'd') LIKE REPLACE(LOWER(unaccent(?)), 'đ', 'd')
		)`
		args = append(args, pattern, pattern, pattern)
	}

	countSQL := fmt.Sprintf(baseSQL, "COUNT(*)")
	var total int64
	if err := r.db.GetDB(ctx).Raw(countSQL, args...).Scan(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("count planning documents by parcel failed: %w", err)
	}
	if total == 0 {
		return []qh_domain.QHPlanningDocument{}, 0, nil
	}

	listSQL := fmt.Sprintf(baseSQL, "d.*") + `
		ORDER BY rr.sort_order ASC, rr.rel_id ASC
		OFFSET ? LIMIT ?
	`
	listArgs := append(append([]any{}, args...), page*size, size)
	var documents []qh_domain.QHPlanningDocument
	if err := r.db.GetDB(ctx).Raw(listSQL, listArgs...).Scan(&documents).Error; err != nil {
		return nil, 0, fmt.Errorf("list planning documents by parcel failed: %w", err)
	}
	return documents, total, nil
}

func (r *QHPlanningDocumentRepo) CountByProject(ctx context.Context, projectID uint64) (int64, error) {
	var total int64
	q := r.db.GetDB(ctx).Model(&qh_domain.QHPlanningDocument{}).Where("deleted_at IS NULL")
	if projectID != 0 {
		q = q.Where("planning_project_id = ?", projectID)
	}
	if err := q.Count(&total).Error; err != nil {
		return 0, err
	}
	return total, nil
}

// ListByProcessStatus lấy các tài liệu ở 1 trạng thái xử lý, dùng cho job classify quét theo batch.
func (r *QHPlanningDocumentRepo) ListByProcessStatus(ctx context.Context, status uint32, limit int) ([]qh_domain.QHPlanningDocument, error) {
	var documents []qh_domain.QHPlanningDocument
	err := r.db.GetDB(ctx).
		Where("deleted_at IS NULL AND process_status = ?", status).
		Order("id ASC").
		Limit(limit).
		Find(&documents).Error
	return documents, err
}

// CountByProjectAndStatuses đếm tài liệu của 1 đồ án theo danh sách trạng thái.
func (r *QHPlanningDocumentRepo) CountByProjectAndStatuses(ctx context.Context, projectID uint64, statuses []uint32) (int64, error) {
	var total int64
	err := r.db.GetDB(ctx).Model(&qh_domain.QHPlanningDocument{}).
		Where("deleted_at IS NULL AND planning_project_id = ? AND process_status IN ?", projectID, statuses).
		Count(&total).Error
	return total, err
}

// ApproveClassifiedByProject chuyển toàn bộ tài liệu đã Classified của đồ án sang Approved.
func (r *QHPlanningDocumentRepo) ApproveClassifiedByProject(ctx context.Context, projectID uint64) error {
	return r.db.GetDB(ctx).Model(&qh_domain.QHPlanningDocument{}).
		Where("deleted_at IS NULL AND planning_project_id = ? AND process_status = ?", projectID, enums.PlanningProcessStatusClassified).
		Update("process_status", enums.PlanningProcessStatusApproved).Error
}

// ResetClassifyForRetry đưa tài liệu Failed về Pending và xóa lỗi classify.
func (r *QHPlanningDocumentRepo) ResetClassifyForRetry(ctx context.Context, id uint64) error {
	return r.db.GetDB(ctx).Model(&qh_domain.QHPlanningDocument{}).
		Where("id = ? AND deleted_at IS NULL AND process_status = ?", id, enums.PlanningProcessStatusFailed).
		Updates(map[string]interface{}{
			"process_status": enums.PlanningProcessStatusPending,
			"classify_error": "",
		}).Error
}

// ResetFailedClassifyByProject đưa tất cả tài liệu Failed của đồ án về Pending.
func (r *QHPlanningDocumentRepo) ResetFailedClassifyByProject(ctx context.Context, projectID uint64) (int64, error) {
	tx := r.db.GetDB(ctx).Model(&qh_domain.QHPlanningDocument{}).
		Where("deleted_at IS NULL AND planning_project_id = ? AND process_status = ?", projectID, enums.PlanningProcessStatusFailed).
		Updates(map[string]interface{}{
			"process_status": enums.PlanningProcessStatusPending,
			"classify_error": "",
		})
	return tx.RowsAffected, tx.Error
}

func (r *QHPlanningDocumentRepo) AssignToLayers(ctx context.Context, layerIDs, docIDs []uint64) error {
	now := time.Now()
	rows := make([]map[string]any, 0, len(layerIDs)*len(docIDs))
	for _, layerID := range layerIDs {
		for i, docID := range docIDs {
			rows = append(rows, map[string]any{
				"layer_id":             layerID,
				"planning_document_id": docID,
				"sort_order":           i,
				"created_at":           now,
				"updated_at":           now,
				"deleted_at":           nil,
			})
		}
	}
	return r.db.GetDB(ctx).Table("layer_plan_doc_relation").
		Clauses(clause.OnConflict{
			Columns: []clause.Column{{Name: "layer_id"}, {Name: "planning_document_id"}},
			DoUpdates: clause.Assignments(map[string]any{
				"sort_order": gorm.Expr("EXCLUDED.sort_order"),
				"updated_at": now,
				"deleted_at": nil,
			}),
		}).Create(&rows).Error
}

func (r *QHPlanningDocumentRepo) UnassignFromLayers(ctx context.Context, layerIDs, docIDs []uint64) error {
	return r.db.GetDB(ctx).
		Unscoped().
		Where("layer_id IN ? AND planning_document_id IN ?", layerIDs, docIDs).
		Delete(&qh_domain.LayerPlanDocRelation{}).Error
}

func (r *QHPlanningDocumentRepo) UnassignLayer(ctx context.Context, planningID, layerID uint64) error {
	return r.db.GetDB(ctx).
		Where("layer_id = ? AND deleted_at IS NULL AND planning_document_id IN (?)", layerID,
			r.db.GetDB(ctx).Model(&qh_domain.QHPlanningDocument{}).
				Select("id").
				Where("planning_project_id = ? AND deleted_at IS NULL", planningID),
		).
		Delete(&qh_domain.LayerPlanDocRelation{}).Error
}

func (r *QHPlanningDocumentRepo) GetPlanningDocumentSeoSource(ctx context.Context, documentID uint64) (*qh_dto.PlanningDocumentSeoSource, error) {
	var source qh_dto.PlanningDocumentSeoSource
	err := r.db.GetDB(ctx).
		Raw(`
			SELECT
				id AS planning_document_id,
				COALESCE(NULLIF(title, ''), '') AS title,
				COALESCE(NULLIF(code, ''), NULLIF(title, ''), '') AS slug,
				seo_id AS seo_id,
				COALESCE(filepath, '') AS filepath,
				COALESCE(thumbnail, '') AS thumbnail,
				COALESCE(ext, '') AS ext
			FROM qh_planning_documents
			WHERE id = ? AND deleted_at IS NULL
			LIMIT 1
		`, documentID).
		Scan(&source).Error
	if err != nil {
		return nil, fmt.Errorf("get planning document seo source failed: %w", err)
	}
	if source.PlanningDocumentID == 0 {
		return nil, nil
	}
	return &source, nil
}

func (r *QHPlanningDocumentRepo) GetPlanningDocumentSeoSourcesByIDs(ctx context.Context, documentIDs []uint64) ([]qh_dto.PlanningDocumentSeoSource, error) {
	ids := dedupeNonZeroUint64(documentIDs)
	if len(ids) == 0 {
		return nil, nil
	}

	var sources []qh_dto.PlanningDocumentSeoSource
	err := r.db.GetDB(ctx).
		Raw(`
			SELECT
				id AS planning_document_id,
				COALESCE(NULLIF(title, ''), '') AS title,
				COALESCE(NULLIF(code, ''), NULLIF(title, ''), '') AS slug,
				seo_id AS seo_id,
				COALESCE(filepath, '') AS filepath,
				COALESCE(thumbnail, '') AS thumbnail,
				COALESCE(ext, '') AS ext
			FROM qh_planning_documents
			WHERE id IN ? AND deleted_at IS NULL
		`, ids).
		Scan(&sources).Error
	if err != nil {
		return nil, fmt.Errorf("get planning document seo sources failed: %w", err)
	}
	return sources, nil
}

func (r *QHPlanningDocumentRepo) UpdatePlanningDocumentSeoID(ctx context.Context, documentID uint64, seoID uint64) error {
	if documentID == 0 || seoID == 0 {
		return errors.New("planning_document_id và seo_id là bắt buộc")
	}
	return r.db.GetDB(ctx).
		Model(&qh_domain.QHPlanningDocument{}).
		Where("id = ? AND deleted_at IS NULL", documentID).
		Update("seo_id", seoID).Error
}

func (r *QHPlanningDocumentRepo) CreateHistory(ctx context.Context, entity *qh_domain.QHPlanningDocumentEvent) error {
	normalizePlanningDocumentHistory(entity)
	return r.db.GetDB(ctx).Create(entity).Error
}

func (r *QHPlanningDocumentRepo) DeleteHistory(ctx context.Context, id uint64) error {
	return r.db.GetDB(ctx).Unscoped().Where("id = ?", id).Delete(&qh_domain.QHPlanningDocumentEvent{}).Error
}

func (r *QHPlanningDocumentRepo) GetHistoryByID(ctx context.Context, id uint64) (*qh_domain.QHPlanningDocumentEvent, error) {
	var entity qh_domain.QHPlanningDocumentEvent
	err := r.db.GetDB(ctx).Where("id = ? AND deleted_at IS NULL", id).First(&entity).Error
	if err != nil {
		return nil, err
	}
	return &entity, nil
}

func (r *QHPlanningDocumentRepo) ListHistories(ctx context.Context, req *qh_dto.ListPlanningDocumentHistoriesRequest) ([]qh_domain.QHPlanningDocumentEvent, int64, error) {
	if req == nil {
		req = &qh_dto.ListPlanningDocumentHistoriesRequest{}
	}

	var entities []qh_domain.QHPlanningDocumentEvent
	var total int64

	db := r.db.GetDB(ctx).Model(&qh_domain.QHPlanningDocumentEvent{}).Where("deleted_at IS NULL")
	if req.PlanningDocumentID != 0 {
		db = db.Where("planning_document_id = ?", req.PlanningDocumentID)
	}

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := db.
		Offset(req.GetOffset()).
		Limit(req.GetLimit()).
		Order("event_date DESC NULLS LAST, created_at DESC").
		Find(&entities).Error; err != nil {
		return nil, 0, err
	}
	return entities, total, nil
}

func normalizePlanningDocumentHistory(entity *qh_domain.QHPlanningDocumentEvent) {
	if entity == nil {
		return
	}
	if len(entity.Metadata) == 0 {
		entity.Metadata = []byte("{}")
	}
}
