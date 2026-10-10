package repository

import (
	"context"
	"fmt"
	"time"

	"userv4/modules/seo/dto"
	"userv4/modules/seo/entity"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// ── SeoDomainRepository ───────────────────────────────────────────────────────

type SeoDomainRepository interface {
	List(ctx context.Context, filter dto.SeoDomainListFilter) (*dto.SeoDomainListResult, error)
	GetByID(ctx context.Context, id uint64) (*entity.SeoDomain, error)
	GetByRef(ctx context.Context, refType uint32, refID uint64) (*entity.SeoDomain, error)
	ListByRef(ctx context.Context, refType uint32, refID uint64) ([]entity.SeoDomain, error)
	ListByLinkedRef(ctx context.Context, filter dto.SeoDomainRefsFilter) (*dto.SeoDomainListResult, error)
	GetByCanonicalURL(ctx context.Context, canonicalURL string) (*entity.SeoDomain, error)
	UpdateLink(ctx context.Context, id uint64, updates map[string]any) (*entity.SeoLink, error)
	TouchByRef(ctx context.Context, refType uint32, refID uint64) (int64, error)
	GetBySlug(ctx context.Context, slug string) (*entity.SeoDomain, error)
	Create(ctx context.Context, seo *entity.SeoDomain) (*entity.SeoDomain, error)
	Update(ctx context.Context, id uint64, updates map[string]any) (*entity.SeoDomain, error)
	Delete(ctx context.Context, id uint64) error
	Publish(ctx context.Context, id uint64) (*entity.SeoDomain, error)
	Unpublish(ctx context.Context, id uint64) (*entity.SeoDomain, error)
	Archive(ctx context.Context, id uint64) error
	Restore(ctx context.Context, id uint64) error
	GetClickStats(ctx context.Context, seoDomainID uint64) (*dto.ClickStatsItem, error)
	ListGenerationLogs(ctx context.Context, filter dto.GenerationLogFilter) ([]dto.GenerationLogItem, int64, error)
	ListForSitemap(ctx context.Context, limit, offset int) ([]entity.SeoDomain, error)
	ListPublishedStaticPaths(ctx context.Context, limit, offset int) ([]entity.SeoDomain, int64, error)
	UpdateByIDs(ctx context.Context, ids []uint64, updates map[string]any) error
	CountByRef(ctx context.Context, refType uint32, refID uint64) (int64, error)
	CountRecentByRefType(ctx context.Context, days uint32) ([]dto.SeoRefTypeCountItem, error)
	CountClicksFromTime(ctx context.Context, fromTime time.Time) (uint64, error)
	GetTopClicked(ctx context.Context, fromTime time.Time, limit int) ([]dto.SeoTopClickedItem, error)
	GetClickDailyStats(ctx context.Context, fromTime time.Time) ([]dto.SeoClickDailyItem, error)
	// Links
	CreateLink(ctx context.Context, link *entity.SeoLink) (*entity.SeoLink, error)
	ListLinks(ctx context.Context, parentSeoID uint64) ([]entity.SeoLink, error)
	DeleteLink(ctx context.Context, id uint64) error
	// Relatives
	CreateRelative(ctx context.Context, rel *entity.SeoRelative) (*entity.SeoRelative, error)
	ListRelatives(ctx context.Context, parentSeoID uint64) ([]entity.SeoRelative, error)
	DeleteRelative(ctx context.Context, id uint64) error
	// Upsert from ref (generate)
	UpsertFromRef(ctx context.Context, seo *entity.SeoDomain) (*entity.SeoDomain, bool, error)
}

type seoDomainPostgres struct{ db *gorm.DB }

func NewSeoDomainPostgres(db *gorm.DB) SeoDomainRepository { return &seoDomainPostgres{db: db} }

func (r *seoDomainPostgres) GetByID(ctx context.Context, id uint64) (*entity.SeoDomain, error) {
	var seo entity.SeoDomain
	err := r.db.WithContext(ctx).
		Preload("InternalLinks", func(db *gorm.DB) *gorm.DB { return db.Where("deleted_at IS NULL").Order("priority DESC") }).
		Preload("Relatives", func(db *gorm.DB) *gorm.DB { return db.Where("deleted_at IS NULL").Order("priority DESC") }).
		Where("id = ? AND deleted_at IS NULL", id).
		First(&seo).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	seo.NormalizeLifecycle()
	return &seo, nil
}

func (r *seoDomainPostgres) GetByRef(ctx context.Context, refType uint32, refID uint64) (*entity.SeoDomain, error) {
	var seo entity.SeoDomain
	err := r.db.WithContext(ctx).
		Where("ref_type = ? AND ref_id = ? AND deleted_at IS NULL", refType, refID).
		First(&seo).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	seo.NormalizeLifecycle()
	return &seo, nil
}

func (r *seoDomainPostgres) ListByRef(ctx context.Context, refType uint32, refID uint64) ([]entity.SeoDomain, error) {
	var list []entity.SeoDomain
	err := r.db.WithContext(ctx).
		Where("ref_type = ? AND ref_id = ? AND deleted_at IS NULL AND published = TRUE AND scope = ?",
			refType, refID, entity.ScopePublic).
		Order("published_at DESC NULLS LAST, id DESC").
		Find(&list).Error
	return list, err
}

func (r *seoDomainPostgres) ListByLinkedRef(ctx context.Context, filter dto.SeoDomainRefsFilter) (*dto.SeoDomainListResult, error) {
	page, size := normalizePage(filter.Page, filter.Size)

	q := r.db.WithContext(ctx).Table("seo_domain sd").
		Select(`sd.id, sd.slug, sd.canonical_url, sd.origin_url, sd.title, sd.description, sd.summary,
		        sd.scope, sd.page_status, sd.published, sd.published_at,
		        sd.ref_type, sd.ref_id, sd.ref_source, sd.ref_label, sd.source_status,
		        sd.category_id, sd.province_id, sd.image_url, sd.template_key,
		        sd.render_status, sd.need_generate, sd.is_site_map, sd.is_index,
		        sd.visible_on, sd.deep_link, sd.created_at, sd.updated_at,
		        COALESCE(sc.click_count, 0) AS click_count,
		        COALESCE(c.name, '') AS category_name`).
		Joins("LEFT JOIN seo_count sc ON sc.seo_domain_id = sd.id").
		Joins("LEFT JOIN seo_category c ON c.id = sd.category_id AND c.deleted_at IS NULL").
		Where("sd.deleted_at IS NULL").
		Where(`(
			(sd.ref_type = ? AND sd.ref_id = ?)
			OR EXISTS (
				SELECT 1 FROM seo_refs sr
				WHERE sr.seo_domain_id = sd.id
				  AND sr.ref_type = ?
				  AND sr.ref_id = ?
			)
		)`, filter.RefType, filter.RefID, filter.RefType, filter.RefID)

	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, fmt.Errorf("count seo domains by linked ref: %w", err)
	}
	var rows []dto.SeoDomainItem
	if err := q.Order("sd.published_at DESC NULLS LAST, sd.updated_at DESC, sd.id DESC").
		Offset((page - 1) * size).Limit(size).Scan(&rows).Error; err != nil {
		return nil, fmt.Errorf("list seo domains by linked ref: %w", err)
	}
	for i := range rows {
		rows[i].RefTypeName = entity.SEORefType(rows[i].RefType).Name()
	}
	return &dto.SeoDomainListResult{Data: rows, Total: total, Page: page, Size: size}, nil
}

func (r *seoDomainPostgres) TouchByRef(ctx context.Context, refType uint32, refID uint64) (int64, error) {
	now := time.Now()
	res := r.db.WithContext(ctx).Model(&entity.SeoDomain{}).
		Where("ref_type = ? AND ref_id = ? AND deleted_at IS NULL", refType, refID).
		Updates(map[string]any{
			"source_updated_at": now,
			"need_generate":     true,
			"source_status":     entity.SourceStatusStale,
			"updated_at":        now,
		})
	return res.RowsAffected, res.Error
}

func (r *seoDomainPostgres) UpdateLink(ctx context.Context, id uint64, updates map[string]any) (*entity.SeoLink, error) {
	if err := r.db.WithContext(ctx).Model(&entity.SeoLink{}).Where("id = ?", id).Updates(updates).Error; err != nil {
		return nil, err
	}
	var link entity.SeoLink
	if err := r.db.WithContext(ctx).First(&link, id).Error; err != nil {
		return nil, err
	}
	return &link, nil
}

func (r *seoDomainPostgres) GetByCanonicalURL(ctx context.Context, canonicalURL string) (*entity.SeoDomain, error) {
	var seo entity.SeoDomain
	err := r.db.WithContext(ctx).
		Where("canonical_url = ? AND deleted_at IS NULL", canonicalURL).
		First(&seo).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	seo.NormalizeLifecycle()
	return &seo, nil
}

func (r *seoDomainPostgres) GetBySlug(ctx context.Context, slug string) (*entity.SeoDomain, error) {
	var seo entity.SeoDomain
	err := r.db.WithContext(ctx).
		Where("slug = ? AND deleted_at IS NULL", slug).
		First(&seo).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	seo.NormalizeLifecycle()
	return &seo, nil
}

func (r *seoDomainPostgres) List(ctx context.Context, filter dto.SeoDomainListFilter) (*dto.SeoDomainListResult, error) {
	page, size := normalizePage(filter.Page, filter.Size)

	q := r.db.WithContext(ctx).Table("seo_domain sd").
		Select(`sd.id, sd.slug, sd.canonical_url, sd.origin_url, sd.title, sd.description, sd.summary,
		        sd.scope, sd.page_status, sd.published, sd.published_at,
		        sd.ref_type, sd.ref_id, sd.ref_source, sd.ref_label, sd.source_status,
		        sd.category_id, sd.province_id, sd.image_url, sd.template_key,
		        sd.render_status, sd.need_generate, sd.is_site_map, sd.is_index,
		        sd.visible_on, sd.deep_link, sd.created_at, sd.updated_at,
		        COALESCE(sc.click_count, 0) AS click_count,
		        COALESCE(c.name, '') AS category_name`).
		Joins("LEFT JOIN seo_count sc ON sc.seo_domain_id = sd.id").
		Joins("LEFT JOIN seo_category c ON c.id = sd.category_id AND c.deleted_at IS NULL").
		Where("sd.deleted_at IS NULL")

	if kw := likeKeyword(filter.Keyword); kw != "" {
		q = q.Where("(LOWER(COALESCE(sd.title,'')) LIKE ? OR LOWER(COALESCE(sd.slug,'')) LIKE ?)", kw, kw)
	}
	if filter.PageStatus != "" {
		q = q.Where("sd.page_status = ?", filter.PageStatus)
	}
	if filter.Scope != "" {
		q = q.Where("sd.scope = ?", filter.Scope)
	}
	if filter.Published != nil {
		q = q.Where("sd.published = ?", *filter.Published)
	}
	if filter.CategoryID > 0 {
		q = q.Where("sd.category_id = ?", filter.CategoryID)
	}
	if filter.ProvinceID > 0 {
		q = q.Where("sd.province_id = ?", filter.ProvinceID)
	}
	if filter.RefType > 0 {
		q = q.Where("sd.ref_type = ?", filter.RefType)
	}

	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, fmt.Errorf("count seo domains: %w", err)
	}
	var rows []dto.SeoDomainItem
	if err := q.Order("sd.updated_at DESC, sd.id DESC").Offset((page - 1) * size).Limit(size).Scan(&rows).Error; err != nil {
		return nil, fmt.Errorf("list seo domains: %w", err)
	}
	for i := range rows {
		rows[i].RefTypeName = entity.SEORefType(rows[i].RefType).Name()
	}
	return &dto.SeoDomainListResult{Data: rows, Total: total, Page: page, Size: size}, nil
}

func (r *seoDomainPostgres) Create(ctx context.Context, seo *entity.SeoDomain) (*entity.SeoDomain, error) {
	if seo.Scope == "" {
		seo.Scope = entity.ScopePublic
	}
	if seo.PageStatus == "" {
		seo.PageStatus = entity.PageStatusDraft
	}
	if seo.RenderStatus == "" {
		seo.RenderStatus = entity.RenderStatusNone
	}
	// GORM zero-value datatypes.JSON inserts SQL NULL; columns are NOT NULL without reliable DEFAULT on older DBs.
	if len(seo.RefSnapshotJSON) == 0 {
		seo.RefSnapshotJSON = datatypes.JSON([]byte("{}"))
	}
	if len(seo.Metadata) == 0 {
		seo.Metadata = datatypes.JSON([]byte("{}"))
	}
	if err := r.db.WithContext(ctx).Create(seo).Error; err != nil {
		return nil, err
	}
	return r.GetByID(ctx, seo.ID)
}

func (r *seoDomainPostgres) Update(ctx context.Context, id uint64, updates map[string]any) (*entity.SeoDomain, error) {
	updates["updated_at"] = time.Now()
	if err := r.db.WithContext(ctx).Model(&entity.SeoDomain{}).
		Where("id = ? AND deleted_at IS NULL", id).Updates(updates).Error; err != nil {
		return nil, err
	}
	return r.GetByID(ctx, id)
}

func (r *seoDomainPostgres) Delete(ctx context.Context, id uint64) error {
	return r.db.WithContext(ctx).Delete(&entity.SeoDomain{}, id).Error
}

func (r *seoDomainPostgres) Publish(ctx context.Context, id uint64) (*entity.SeoDomain, error) {
	now := time.Now()
	updates := map[string]any{
		"published":    true,
		"published_at": now,
		"page_status":  entity.PageStatusPublished,
		"updated_at":   now,
	}
	if err := r.db.WithContext(ctx).Model(&entity.SeoDomain{}).
		Where("id = ? AND deleted_at IS NULL", id).Updates(updates).Error; err != nil {
		return nil, err
	}
	return r.GetByID(ctx, id)
}

func (r *seoDomainPostgres) Unpublish(ctx context.Context, id uint64) (*entity.SeoDomain, error) {
	updates := map[string]any{
		"published":   false,
		"page_status": entity.PageStatusDraft,
		"updated_at":  time.Now(),
	}
	if err := r.db.WithContext(ctx).Model(&entity.SeoDomain{}).
		Where("id = ? AND deleted_at IS NULL", id).Updates(updates).Error; err != nil {
		return nil, err
	}
	return r.GetByID(ctx, id)
}

func (r *seoDomainPostgres) Archive(ctx context.Context, id uint64) error {
	updates := map[string]any{
		"published":   false,
		"page_status": entity.PageStatusArchived,
		"updated_at":  time.Now(),
	}
	return r.db.WithContext(ctx).Model(&entity.SeoDomain{}).
		Where("id = ? AND deleted_at IS NULL", id).Updates(updates).Error
}

func (r *seoDomainPostgres) Restore(ctx context.Context, id uint64) error {
	updates := map[string]any{
		"page_status": entity.PageStatusDraft,
		"updated_at":  time.Now(),
	}
	return r.db.WithContext(ctx).Model(&entity.SeoDomain{}).
		Where("id = ? AND deleted_at IS NULL", id).Updates(updates).Error
}

func (r *seoDomainPostgres) GetClickStats(ctx context.Context, seoDomainID uint64) (*dto.ClickStatsItem, error) {
	var item dto.ClickStatsItem
	err := r.db.WithContext(ctx).Table("seo_count").
		Where("seo_domain_id = ?", seoDomainID).
		Select("seo_domain_id, click_count, view_count").
		Scan(&item).Error
	if err != nil {
		return nil, err
	}
	item.SeoDomainID = seoDomainID
	return &item, nil
}

func (r *seoDomainPostgres) ListGenerationLogs(ctx context.Context, filter dto.GenerationLogFilter) ([]dto.GenerationLogItem, int64, error) {
	page, size := normalizePage(filter.Page, filter.Size)
	q := r.db.WithContext(ctx).Table("seo_generation_log")
	if filter.SeoDomainID > 0 {
		q = q.Where("seo_domain_id = ?", filter.SeoDomainID)
	}
	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var logs []dto.GenerationLogItem
	if err := q.Select("id, seo_domain_id, status, message, started_at, finished_at, created_at").
		Order("id DESC").Offset((page - 1) * size).Limit(size).Scan(&logs).Error; err != nil {
		return nil, 0, err
	}
	return logs, total, nil
}

func (r *seoDomainPostgres) ListForSitemap(ctx context.Context, limit, offset int) ([]entity.SeoDomain, error) {
	var list []entity.SeoDomain
	err := r.db.WithContext(ctx).
		Where("deleted_at IS NULL AND published = TRUE AND is_site_map = TRUE AND is_index = TRUE AND scope = 'public' AND page_status = 'published'").
		Order("published_at DESC NULLS LAST, id DESC").
		Limit(limit).Offset(offset).
		Find(&list).Error
	return list, err
}

func (r *seoDomainPostgres) ListPublishedStaticPaths(ctx context.Context, limit, offset int) ([]entity.SeoDomain, int64, error) {
	q := r.db.WithContext(ctx).Model(&entity.SeoDomain{}).
		Where("deleted_at IS NULL AND published = TRUE AND is_index = TRUE AND scope = ? AND page_status = ?",
			entity.ScopePublic, entity.PageStatusPublished)
	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var list []entity.SeoDomain
	if err := q.Order("updated_at DESC").Limit(limit).Offset(offset).Find(&list).Error; err != nil {
		return nil, 0, err
	}
	for i := range list {
		list[i].NormalizeLifecycle()
	}
	return list, total, nil
}

func (r *seoDomainPostgres) UpdateByIDs(ctx context.Context, ids []uint64, updates map[string]any) error {
	if len(ids) == 0 || len(updates) == 0 {
		return nil
	}
	updates["updated_at"] = time.Now()
	return r.db.WithContext(ctx).Model(&entity.SeoDomain{}).
		Where("id IN ? AND deleted_at IS NULL", ids).
		Updates(updates).Error
}

func (r *seoDomainPostgres) CountByRef(ctx context.Context, refType uint32, refID uint64) (int64, error) {
	var total int64
	err := r.db.WithContext(ctx).Model(&entity.SeoDomain{}).
		Where("ref_type = ? AND ref_id = ? AND deleted_at IS NULL", refType, refID).
		Count(&total).Error
	return total, err
}

func (r *seoDomainPostgres) CountRecentByRefType(ctx context.Context, days uint32) ([]dto.SeoRefTypeCountItem, error) {
	if days == 0 {
		days = 5
	}
	type row struct {
		RefType uint32 `gorm:"column:ref_type"`
		Total   int64  `gorm:"column:total"`
	}
	since := time.Now().AddDate(0, 0, -int(days))
	var rows []row
	err := r.db.WithContext(ctx).Model(&entity.SeoDomain{}).
		Select("ref_type, COUNT(*) AS total").
		Where("updated_at >= ? AND deleted_at IS NULL", since).
		Group("ref_type").
		Scan(&rows).Error
	if err != nil {
		return nil, err
	}
	totals := map[uint32]int64{}
	for _, item := range rows {
		refType := item.RefType
		if refType == 0 {
			refType = uint32(entity.SEORefTypeLegacy)
		}
		totals[refType] += item.Total
	}
	allTypes := []entity.SEORefType{
		entity.SEORefTypeAdmUnit, entity.SEORefTypeParcel, entity.SEORefTypeRegion,
		entity.SEORefTypeProject, entity.SEORefTypeDocument, entity.SEORefTypeMap,
		entity.SEORefTypeReport, entity.SEORefTypeLegacy,
	}
	result := make([]dto.SeoRefTypeCountItem, 0, len(allTypes))
	for _, refType := range allTypes {
		total := totals[uint32(refType)]
		if total == 0 {
			continue
		}
		result = append(result, dto.SeoRefTypeCountItem{
			RefType:     uint32(refType),
			RefTypeName: refType.Name(),
			Total:       total,
		})
	}
	return result, nil
}

func (r *seoDomainPostgres) CountClicksFromTime(ctx context.Context, fromTime time.Time) (uint64, error) {
	var total int64
	err := r.db.WithContext(ctx).Raw(`
		SELECT COALESCE(COUNT(*), 0)
		FROM seo_count_log
		WHERE created_at >= ? AND event_type = 'click'
	`, fromTime).Scan(&total).Error
	if err != nil {
		return 0, err
	}
	return uint64(total), nil
}

func (r *seoDomainPostgres) GetTopClicked(ctx context.Context, fromTime time.Time, limit int) ([]dto.SeoTopClickedItem, error) {
	if limit <= 0 {
		limit = 15
	}
	type row struct {
		ID         uint64  `gorm:"column:id"`
		Name       string  `gorm:"column:name"`
		RefType    uint32  `gorm:"column:ref_type"`
		RefID      *uint64 `gorm:"column:ref_id"`
		ClickCount uint64  `gorm:"column:click_count"`
	}
	var rows []row
	err := r.db.WithContext(ctx).Raw(`
		SELECT
			sd.id,
			COALESCE(NULLIF(TRIM(sd.title), ''), COALESCE(sd.ref_label, '')) AS name,
			sd.ref_type,
			sd.ref_id,
			COUNT(*)::bigint AS click_count
		FROM seo_count_log scl
		INNER JOIN seo_domain sd ON sd.id = scl.seo_domain_id AND sd.deleted_at IS NULL
		WHERE scl.created_at >= ? AND scl.event_type = 'click'
		GROUP BY sd.id, sd.title, sd.ref_label, sd.ref_type, sd.ref_id
		ORDER BY click_count DESC, sd.id DESC
		LIMIT ?
	`, fromTime, limit).Scan(&rows).Error
	if err != nil {
		return nil, err
	}
	result := make([]dto.SeoTopClickedItem, 0, len(rows))
	for _, item := range rows {
		var refID uint64
		if item.RefID != nil {
			refID = *item.RefID
		}
		result = append(result, dto.SeoTopClickedItem{
			ID:   item.ID,
			Name: item.Name,
			Source: dto.SeoClickSource{
				RefType: item.RefType,
				RefID:   refID,
			},
			ClickCount: item.ClickCount,
		})
	}
	return result, nil
}

func (r *seoDomainPostgres) GetClickDailyStats(ctx context.Context, fromTime time.Time) ([]dto.SeoClickDailyItem, error) {
	type row struct {
		Date  string `gorm:"column:date"`
		Count uint64 `gorm:"column:count"`
	}
	var rows []row
	err := r.db.WithContext(ctx).Raw(`
		SELECT
			to_char((created_at AT TIME ZONE 'Asia/Ho_Chi_Minh'), 'YYYY-MM-DD') AS date,
			COUNT(*)::bigint AS count
		FROM seo_count_log
		WHERE created_at >= ? AND event_type = 'click'
		GROUP BY 1
		ORDER BY 1 ASC
	`, fromTime).Scan(&rows).Error
	if err != nil {
		return nil, err
	}
	result := make([]dto.SeoClickDailyItem, 0, len(rows))
	for _, item := range rows {
		result = append(result, dto.SeoClickDailyItem{Date: item.Date, Count: item.Count})
	}
	return result, nil
}

// ── Links ──────────────────────────────────────────────────────────────────────

func (r *seoDomainPostgres) CreateLink(ctx context.Context, link *entity.SeoLink) (*entity.SeoLink, error) {
	if err := r.db.WithContext(ctx).Create(link).Error; err != nil {
		return nil, err
	}
	return link, nil
}

func (r *seoDomainPostgres) ListLinks(ctx context.Context, parentSeoID uint64) ([]entity.SeoLink, error) {
	var list []entity.SeoLink
	err := r.db.WithContext(ctx).
		Where("parent_seo_id = ? AND deleted_at IS NULL", parentSeoID).
		Order("priority DESC, updated_at DESC").
		Find(&list).Error
	return list, err
}

func (r *seoDomainPostgres) DeleteLink(ctx context.Context, id uint64) error {
	return r.db.WithContext(ctx).Delete(&entity.SeoLink{}, id).Error
}

// ── Relatives ─────────────────────────────────────────────────────────────────

func (r *seoDomainPostgres) CreateRelative(ctx context.Context, rel *entity.SeoRelative) (*entity.SeoRelative, error) {
	if err := r.db.WithContext(ctx).Create(rel).Error; err != nil {
		return nil, err
	}
	return rel, nil
}

func (r *seoDomainPostgres) ListRelatives(ctx context.Context, parentSeoID uint64) ([]entity.SeoRelative, error) {
	var list []entity.SeoRelative
	err := r.db.WithContext(ctx).
		Where("parent_seo_id = ? AND deleted_at IS NULL", parentSeoID).
		Order("priority DESC").
		Find(&list).Error
	return list, err
}

func (r *seoDomainPostgres) DeleteRelative(ctx context.Context, id uint64) error {
	return r.db.WithContext(ctx).Delete(&entity.SeoRelative{}, id).Error
}

// ── Upsert from ref ───────────────────────────────────────────────────────────

func (r *seoDomainPostgres) UpsertFromRef(ctx context.Context, seo *entity.SeoDomain) (*entity.SeoDomain, bool, error) {
	if seo.RefID == nil || *seo.RefID == 0 || seo.RefType == 0 {
		created, err := r.Create(ctx, seo)
		return created, true, err
	}
	existing, err := r.GetByRef(ctx, uint32(seo.RefType), *seo.RefID)
	if err != nil {
		return nil, false, err
	}
	if existing != nil {
		updates := buildUpdateMap(seo)
		updated, err := r.Update(ctx, existing.ID, updates)
		return updated, false, err
	}
	created, err := r.Create(ctx, seo)
	return created, true, err
}

func buildUpdateMap(seo *entity.SeoDomain) map[string]any {
	m := map[string]any{}
	if seo.Title != "" {
		m["title"] = seo.Title
	}
	if seo.Description != "" {
		m["description"] = seo.Description
	}
	if seo.Summary != "" {
		m["summary"] = seo.Summary
	}
	if seo.Slug != "" {
		m["slug"] = seo.Slug
	}
	if seo.CanonicalURL != "" {
		m["canonical_url"] = seo.CanonicalURL
	}
	if seo.ImageURL != "" {
		m["image_url"] = seo.ImageURL
	}
	if seo.DeepLink != "" {
		m["deep_link"] = seo.DeepLink
	}
	if seo.RefLabel != "" {
		m["ref_label"] = seo.RefLabel
	}
	if seo.TemplateKey != "" {
		m["template_key"] = seo.TemplateKey
	}
	now := time.Now()
	m["source_updated_at"] = now
	m["need_generate"] = true
	if seo.SourceStatus != "" {
		m["source_status"] = seo.SourceStatus
	} else {
		m["source_status"] = entity.SourceStatusLinked
	}
	return m
}
