package postgres

import (
	_db "common/db"
	"context"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"
	qh_domain "tqd/internal/domain/qh"
	qh_dto "tqd/internal/dto/qh"
	"tqd/internal/enums"
	"tqd/internal/interface/repo"
	tqdutils "tqd/internal/utils"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type QHPlanningProjectRepo struct {
	db *_db.TransactionRepo
}

// NewQHPlanningProjectRepo @bind: internal/interface/repo.IQHPlanningProjectRepo
func NewQHPlanningProjectRepo(db *_db.TransactionRepo) repo.IQHPlanningProjectRepo {
	return &QHPlanningProjectRepo{db: db}
}

func (r *QHPlanningProjectRepo) Create(ctx context.Context, entity *qh_domain.QHPlanningProject) error {
	return r.db.GetDB(ctx).Create(entity).Error
}

func (r *QHPlanningProjectRepo) Update(ctx context.Context, id uint64, entity *qh_domain.QHPlanningProject) error {
	tx := r.db.GetDB(ctx).Model(&qh_domain.QHPlanningProject{}).Where("id = ?", id)
	if err := tx.Updates(entity).Error; err != nil {
		return err
	}
	// Updates() bỏ qua bool zero-value — ghi riêng để admin ẩn đồ án (clientView=false).
	return tx.Update("client_view", entity.ClientView).Error
}

func (r *QHPlanningProjectRepo) Delete(ctx context.Context, id uint64) error {
	return r.db.GetDB(ctx).Where("id = ?", id).Delete(&qh_domain.QHPlanningProject{}).Error
}

func (r *QHPlanningProjectRepo) GetByID(ctx context.Context, id uint64) (*qh_domain.QHPlanningProject, error) {
	var entity qh_domain.QHPlanningProject
	err := r.db.GetDB(ctx).
		Preload("Jurisdiction").
		Preload("Layers").
		Preload("LegalBasis").
		Where("id = ? AND deleted_at IS NULL", id).
		First(&entity).Error
	if err != nil {
		return nil, err
	}
	var count int64
	if err := r.db.GetDB(ctx).
		Table("qh_planning_documents").
		Where("planning_project_id = ? AND deleted_at IS NULL", id).
		Count(&count).Error; err != nil {
		return nil, err
	}
	entity.DocumentCount = count
	return &entity, nil
}

func (r *QHPlanningProjectRepo) GetDetail(ctx context.Context, id uint64) (*qh_domain.QHPlanningProject, error) {
	return r.GetByID(ctx, id)
}

func (r *QHPlanningProjectRepo) GetAll(ctx context.Context) ([]qh_domain.QHPlanningProject, error) {
	var entities []qh_domain.QHPlanningProject
	err := r.db.GetDB(ctx).Where("deleted_at IS NULL").Find(&entities).Error
	return entities, err
}

func (r *QHPlanningProjectRepo) GetList(ctx context.Context, req *qh_dto.ListPlanningProjectsRequest) ([]qh_domain.QHPlanningProject, int64, error) {
	var projects []qh_domain.QHPlanningProject
	var total int64

	db := r.db.GetDB(ctx).Model(&qh_domain.QHPlanningProject{}).Where("deleted_at IS NULL")

	if search := strings.TrimSpace(req.Search); search != "" {
		searchTerm := "%" + search + "%"
		// unaccent('Đ')→'D' nhưng unaccent('đ') giữ nguyên → sau ToLower phía Go sẽ lệch.
		// Chuẩn hóa thêm đ→d ở SQL để search không phân biệt hoa/thường + dấu.
		db = db.Where(`(
			REPLACE(LOWER(unaccent(COALESCE(name,''))), 'đ', 'd') LIKE REPLACE(LOWER(unaccent(?)), 'đ', 'd') OR
			REPLACE(LOWER(unaccent(COALESCE(code,''))), 'đ', 'd') LIKE REPLACE(LOWER(unaccent(?)), 'đ', 'd') OR
			REPLACE(REPLACE(LOWER(unaccent(COALESCE(slug,''))), '-', ' '), 'đ', 'd') LIKE REPLACE(REPLACE(LOWER(unaccent(?)), '-', ' '), 'đ', 'd')
		)`, searchTerm, searchTerm, searchTerm)
	}
	if req.PlanningType != "" {
		db = db.Where("planning_type = ?", req.PlanningType)
	}
	if req.PlanningLevel != "" {
		db = db.Where("planning_level = ?", req.PlanningLevel)
	}
	if req.ValidityStatus != "" {
		db = db.Where("validity_status = ?", req.ValidityStatus)
	}
	if req.ClientViewOnly {
		db = db.Where("client_view = ?", true)
	}
	if req.ProcessStatus != "" {
		statuses := strings.Split(req.ProcessStatus, ",")
		trimmed := make([]string, 0, len(statuses))
		for _, s := range statuses {
			s = strings.TrimSpace(s)
			if s != "" {
				trimmed = append(trimmed, s)
			}
		}
		if len(trimmed) == 1 {
			db = db.Where("process_status = ?", trimmed[0])
		} else if len(trimmed) > 1 {
			db = db.Where("process_status IN ?", trimmed)
		}
	}

	if clause, arg, ok := tqdutils.PlanningProjectAdmJoinClause(req.AdmID, "qh_planning_projects.id"); ok {
		db = db.Where(clause, arg)
	}
	if layerIDs := tqdutils.MergePositiveUint64(0, req.LayerIDs); len(layerIDs) > 0 {
		// Lọc đồ án có tài liệu gắn layer qua layer_plan_doc_relation.
		db = db.Where(`EXISTS (
			SELECT 1 FROM layer_plan_doc_relation rel
			JOIN qh_planning_documents d
			  ON d.id = rel.planning_document_id
			 AND d.deleted_at IS NULL
			WHERE rel.deleted_at IS NULL
			  AND d.planning_project_id = qh_planning_projects.id
			  AND rel.layer_id IN ?
		)`, layerIDs)
	}

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := db.
		Preload("Jurisdiction").
		Preload("Layers").
		Offset(req.GetOffset()).
		Limit(req.GetLimit()).
		Order("updated_at DESC NULLS LAST, id DESC").
		Find(&projects).Error
	if err != nil {
		return projects, total, err
	}
	if err := r.attachDocumentCounts(ctx, projects); err != nil {
		return nil, 0, err
	}

	log.Printf("projects: %+v", len(projects))
	return projects, total, nil
}

func (r *QHPlanningProjectRepo) attachDocumentCounts(ctx context.Context, projects []qh_domain.QHPlanningProject) error {
	if len(projects) == 0 {
		return nil
	}
	ids := make([]uint64, 0, len(projects))
	indexByID := make(map[uint64]int, len(projects))
	for i := range projects {
		ids = append(ids, projects[i].ID)
		indexByID[projects[i].ID] = i
	}
	type row struct {
		PlanningProjectID uint64 `gorm:"column:planning_project_id"`
		DocumentCount     int64  `gorm:"column:document_count"`
	}
	var rows []row
	if err := r.db.GetDB(ctx).
		Table("qh_planning_documents").
		Select("planning_project_id, COUNT(*) AS document_count").
		Where("planning_project_id IN ? AND deleted_at IS NULL", ids).
		Group("planning_project_id").
		Scan(&rows).Error; err != nil {
		return err
	}
	for _, item := range rows {
		if idx, ok := indexByID[item.PlanningProjectID]; ok {
			projects[idx].DocumentCount = item.DocumentCount
		}
	}
	return nil
}

func (r *QHPlanningProjectRepo) matchedLayersQuery(ctx context.Context, planningID uint64) *gorm.DB {
	return r.db.GetDB(ctx).Table("qh_layers").
		Select("qh_layers.id").
		Joins("LEFT JOIN layer_plan_doc_relation rel ON rel.layer_id = qh_layers.id AND rel.deleted_at IS NULL").
		Joins("LEFT JOIN qh_planning_documents d ON d.id = rel.planning_document_id AND d.deleted_at IS NULL").
		Where("qh_layers.deleted_at IS NULL").
		Where("qh_layers.planning_project_id = ? OR d.planning_project_id = ?", planningID, planningID).
		Group("qh_layers.id")
}

func (r *QHPlanningProjectRepo) ListLayers(ctx context.Context, req *qh_dto.ListPlanningProjectLayersRequest) ([]qh_domain.QHLayer, int64, error) {
	var layers []qh_domain.QHLayer
	var total int64

	matched := r.matchedLayersQuery(ctx, req.PlanningID)

	if err := r.db.GetDB(ctx).Table("(?) AS matched_layers", matched).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := r.db.GetDB(ctx).Model(&qh_domain.QHLayer{}).
		Where("id IN (?)", matched).
		Order("display_order ASC, id ASC").
		Offset(req.GetOffset()).
		Limit(req.GetLimit()).
		Find(&layers).Error
	if err != nil {
		return nil, 0, err
	}
	if err := r.attachLayerDocs(ctx, req.PlanningID, layers); err != nil {
		return nil, 0, err
	}
	return layers, total, nil
}

func (r *QHPlanningProjectRepo) CountLayersByProject(ctx context.Context, projectID uint64) (int64, error) {
	var total int64
	matched := r.matchedLayersQuery(ctx, projectID)
	if err := r.db.GetDB(ctx).Table("(?) AS matched_layers", matched).Count(&total).Error; err != nil {
		return 0, err
	}
	return total, nil
}

type planningLayerDocRow struct {
	LayerID      uint64 `gorm:"column:layer_id"`
	ID           uint64 `gorm:"column:id"`
	Title        string `gorm:"column:title"`
	Filepath     string `gorm:"column:filepath"`
	Thumbnail    string `gorm:"column:thumbnail"`
	Ext          string `gorm:"column:ext"`
	RelativePath string `gorm:"column:relative_path"`
	Code         string `gorm:"column:code"`
	DocumentType uint32 `gorm:"column:document_type"`
}

func (r *QHPlanningProjectRepo) attachLayerDocs(ctx context.Context, planningID uint64, layers []qh_domain.QHLayer) error {
	if len(layers) == 0 {
		return nil
	}

	layerIDs := make([]uint64, 0, len(layers))
	indexByID := make(map[uint64]int, len(layers))
	for i := range layers {
		layerIDs = append(layerIDs, layers[i].ID)
		indexByID[layers[i].ID] = i
		layers[i].Docs = []qh_domain.QHPlanningDocument{}
	}

	var rows []planningLayerDocRow
	err := r.db.GetDB(ctx).Table("layer_plan_doc_relation rel").
		Select(`rel.layer_id,
			d.id, d.title, d.filepath, d.thumbnail, d.ext, d.relative_path, d.code, d.document_type`).
		Joins("JOIN qh_planning_documents d ON d.id = rel.planning_document_id AND d.deleted_at IS NULL").
		Where("rel.deleted_at IS NULL AND rel.layer_id IN ? AND d.planning_project_id = ?", layerIDs, planningID).
		Order("rel.sort_order ASC, rel.id ASC").
		Scan(&rows).Error
	if err != nil {
		return err
	}

	for _, row := range rows {
		idx, ok := indexByID[row.LayerID]
		if !ok {
			continue
		}
		doc := qh_domain.QHPlanningDocument{
			Title:        row.Title,
			Filepath:     row.Filepath,
			Thumbnail:    row.Thumbnail,
			Ext:          row.Ext,
			RelativePath: row.RelativePath,
			Code:         row.Code,
			DocumentType: enums.PlanningDocumentType(row.DocumentType),
		}
		doc.ID = row.ID
		layers[idx].Docs = append(layers[idx].Docs, doc)
	}
	return nil
}

// UpdateProcessStatus cập nhật riêng cột process_status, dùng cho job classify.
func (r *QHPlanningProjectRepo) UpdateProcessStatus(ctx context.Context, id uint64, status uint32) error {
	return r.db.GetDB(ctx).Model(&qh_domain.QHPlanningProject{}).Where("id = ?", id).
		Update("process_status", status).Error
}

func (r *QHPlanningProjectRepo) GetPlanningProjectSeoSource(ctx context.Context, planningProjectID uint64) (*qh_dto.PlanningProjectSeoSource, error) {
	var source qh_dto.PlanningProjectSeoSource
	err := r.db.GetDB(ctx).
		Raw(`
			SELECT
				p.id AS planning_project_id,
				COALESCE(NULLIF(p.name, ''), '') AS name,
				COALESCE(NULLIF(p.slug, ''), NULLIF(p.name, ''), '') AS slug,
				p.seo_id AS seo_id,
				(
					SELECT CAST(NULLIF(TRIM(CAST(pp.province_id AS text)), '') AS bigint)
					FROM qh_planning_project_provinces pp
					WHERE pp.planning_project_id = p.id
					  AND pp.deleted_at IS NULL
					  AND TRIM(CAST(pp.province_id AS text)) ~ '^[0-9]+$'
					ORDER BY pp.id ASC
					LIMIT 1
				) AS province_id
			FROM qh_planning_projects p
			WHERE p.id = ? AND p.deleted_at IS NULL
			LIMIT 1
		`, planningProjectID).
		Scan(&source).Error
	if err != nil {
		return nil, fmt.Errorf("get planning project seo source failed: %w", err)
	}
	if source.PlanningProjectID == 0 {
		return nil, nil
	}
	return &source, nil
}

func (r *QHPlanningProjectRepo) GetPlanningProjectSeoSourcesByIDs(ctx context.Context, planningProjectIDs []uint64) ([]qh_dto.PlanningProjectSeoSource, error) {
	ids := dedupePlanningProjectIDs(planningProjectIDs)
	if len(ids) == 0 {
		return nil, nil
	}

	var sources []qh_dto.PlanningProjectSeoSource
	err := r.db.GetDB(ctx).
		Raw(`
			SELECT
				p.id AS planning_project_id,
				COALESCE(NULLIF(p.name, ''), '') AS name,
				COALESCE(NULLIF(p.slug, ''), NULLIF(p.name, ''), '') AS slug,
				p.seo_id AS seo_id,
				(
					SELECT CAST(NULLIF(TRIM(CAST(pp.province_id AS text)), '') AS bigint)
					FROM qh_planning_project_provinces pp
					WHERE pp.planning_project_id = p.id
					  AND pp.deleted_at IS NULL
					  AND TRIM(CAST(pp.province_id AS text)) ~ '^[0-9]+$'
					ORDER BY pp.id ASC
					LIMIT 1
				) AS province_id
			FROM qh_planning_projects p
			WHERE p.id IN ? AND p.deleted_at IS NULL
		`, ids).
		Scan(&sources).Error
	if err != nil {
		return nil, fmt.Errorf("get planning project seo sources failed: %w", err)
	}
	return sources, nil
}

func (r *QHPlanningProjectRepo) UpdatePlanningProjectSeoID(ctx context.Context, planningProjectID uint64, seoID uint64) error {
	if planningProjectID == 0 || seoID == 0 {
		return errors.New("planning_project_id và seo_id là bắt buộc")
	}
	return r.db.GetDB(ctx).
		Model(&qh_domain.QHPlanningProject{}).
		Where("id = ? AND deleted_at IS NULL", planningProjectID).
		Update("seo_id", seoID).Error
}

// AssignProvinces inserts relation rows (ignore conflicts)
func (r *QHPlanningProjectRepo) AssignProvinces(ctx context.Context, planningProjectID uint64, provinceIDs []uint64) error {
	if planningProjectID == 0 || len(provinceIDs) == 0 {
		return nil
	}
	rows := make([]map[string]any, 0, len(provinceIDs))
	now := time.Now()
	for _, pid := range provinceIDs {
		if pid == 0 {
			continue
		}
		rows = append(rows, map[string]any{
			"planning_project_id": planningProjectID,
			"province_id":         pid,
			"created_at":          now,
			"updated_at":          now,
			"deleted_at":          nil,
		})
	}
	if len(rows) == 0 {
		return nil
	}
	return r.db.GetDB(ctx).Table("qh_planning_project_provinces").
		Clauses(clause.OnConflict{DoNothing: true}).
		Create(rows).Error
}

// UnassignProvinces deletes matching rows
func (r *QHPlanningProjectRepo) UnassignProvinces(ctx context.Context, planningProjectID uint64, provinceIDs []uint64) error {
	if planningProjectID == 0 || len(provinceIDs) == 0 {
		return nil
	}
	ids := make([]uint64, 0, len(provinceIDs))
	for _, pid := range provinceIDs {
		if pid == 0 {
			continue
		}
		ids = append(ids, pid)
	}
	if len(ids) == 0 {
		return nil
	}
	return r.db.GetDB(ctx).Table("qh_planning_project_provinces").
		Where("planning_project_id = ? AND province_id IN ?", planningProjectID, ids).
		Delete(nil).Error
}

// AssignWards inserts relation rows (ignore conflicts)
func (r *QHPlanningProjectRepo) AssignWards(ctx context.Context, planningProjectID uint64, wardIDs []uint64) error {
	if planningProjectID == 0 || len(wardIDs) == 0 {
		return nil
	}
	rows := make([]map[string]any, 0, len(wardIDs))
	now := time.Now()
	for _, wid := range wardIDs {
		if wid == 0 {
			continue
		}
		rows = append(rows, map[string]any{
			"planning_project_id": planningProjectID,
			"ward_id":             wid,
			"created_at":          now,
			"updated_at":          now,
			"deleted_at":          nil,
		})
	}
	if len(rows) == 0 {
		return nil
	}
	return r.db.GetDB(ctx).Table("qh_planning_project_wards").
		Clauses(clause.OnConflict{DoNothing: true}).
		Create(rows).Error
}

// UnassignWards deletes matching rows
func (r *QHPlanningProjectRepo) UnassignWards(ctx context.Context, planningProjectID uint64, wardIDs []uint64) error {
	if planningProjectID == 0 || len(wardIDs) == 0 {
		return nil
	}
	ids := make([]uint64, 0, len(wardIDs))
	for _, wid := range wardIDs {
		if wid == 0 {
			continue
		}
		ids = append(ids, wid)
	}
	if len(ids) == 0 {
		return nil
	}
	return r.db.GetDB(ctx).Table("qh_planning_project_wards").
		Where("planning_project_id = ? AND ward_id IN ?", planningProjectID, ids).
		Delete(nil).Error
}

func dedupePlanningProjectIDs(ids []uint64) []uint64 {
	seen := make(map[uint64]struct{}, len(ids))
	out := make([]uint64, 0, len(ids))
	for _, id := range ids {
		if id == 0 {
			continue
		}
		if _, ok := seen[id]; ok {
			continue
		}
		seen[id] = struct{}{}
		out = append(out, id)
	}
	return out
}
