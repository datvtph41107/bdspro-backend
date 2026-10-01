//go:build ignore

package postgres

import (
	"context"
	"crm/infra/impl"
	seo_domain "userv4/modules/crm/entity/seo"
	"userv4/modules/crm/dto"
	"userv4/modules/crm/enums"
	repo "userv4/modules/crm/repository"
	"strings"
	"time"

	"gorm.io/gorm"
)

type SeoDomainPostgres struct {
	db *gorm.DB
}

func NewSeoDomainPostgres(db *gorm.DB) repo.SeoDomainRepo {
	return &SeoDomainPostgres{db: db}
}

func (r *SeoDomainPostgres) Create(ctx context.Context, seo *seo_domain.SeoDomain) (*seo_domain.SeoDomain, error) {
	if seo != nil {
		seo.NormalizeLifecycle()
	}
	if err := impl.GetDB(ctx, r.db).Create(seo).Error; err != nil {
		return nil, err
	}
	return r.GetByID(ctx, seo.ID)
}

func (r *SeoDomainPostgres) GetByID(ctx context.Context, id uint64) (*seo_domain.SeoDomain, error) {
	var seo seo_domain.SeoDomain
	err := impl.GetDB(ctx, r.db).WithContext(ctx).
		Preload("InternalLinks", func(db *gorm.DB) *gorm.DB {
			return db.Order("priority DESC, updated_at DESC")
		}).
		Preload("InternalLinks.ChildSeo").
		Preload("Relatives", func(db *gorm.DB) *gorm.DB {
			return db.Order("priority DESC, updated_at DESC")
		}).
		Preload("Relatives.ChildSeo").
		First(&seo, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	seo.NormalizeLifecycle()
	r.fillCategoryName(ctx, &seo)
	return &seo, nil
}

func (r *SeoDomainPostgres) fillCategoryName(ctx context.Context, seo *seo_domain.SeoDomain) {
	if seo == nil || seo.CategoryID == nil || *seo.CategoryID == 0 {
		return
	}
	var name string
	_ = impl.GetDB(ctx, r.db).WithContext(ctx).
		Table("seo_category").
		Select("name").
		Where("id = ? AND deleted_at IS NULL", *seo.CategoryID).
		Limit(1).
		Scan(&name).Error
	seo.CategoryName = name
}

func (r *SeoDomainPostgres) GetByRef(ctx context.Context, refType uint32, refID uint64) (*seo_domain.SeoDomain, error) {
	var seo seo_domain.SeoDomain
	err := impl.GetDB(ctx, r.db).WithContext(ctx).
		Preload("InternalLinks", func(db *gorm.DB) *gorm.DB {
			return db.Order("priority DESC, updated_at DESC")
		}).
		Preload("InternalLinks.ChildSeo").
		Preload("Relatives", func(db *gorm.DB) *gorm.DB {
			return db.Order("priority DESC, updated_at DESC")
		}).
		Preload("Relatives.ChildSeo").
		Where("ref_type = ? AND ref_id = ?", refType, refID).
		First(&seo).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	seo.NormalizeLifecycle()
	return &seo, nil
}

func (r *SeoDomainPostgres) GetByRefSource(ctx context.Context, refType uint32, refSource string, refID uint64) (*seo_domain.SeoDomain, error) {
	var seo seo_domain.SeoDomain
	err := impl.GetDB(ctx, r.db).WithContext(ctx).
		Preload("InternalLinks", func(db *gorm.DB) *gorm.DB {
			return db.Order("priority DESC, updated_at DESC")
		}).
		Preload("InternalLinks.ChildSeo").
		Preload("Relatives", func(db *gorm.DB) *gorm.DB {
			return db.Order("priority DESC, updated_at DESC")
		}).
		Preload("Relatives.ChildSeo").
		Where("ref_type = ? AND ref_source = ? AND ref_id = ?", refType, refSource, refID).
		First(&seo).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	seo.NormalizeLifecycle()
	return &seo, nil
}

func (r *SeoDomainPostgres) GetByRefIdentity(ctx context.Context, refType uint32, refSource string, refURL string) (*seo_domain.SeoDomain, error) {
	var seo seo_domain.SeoDomain
	err := impl.GetDB(ctx, r.db).WithContext(ctx).
		Preload("InternalLinks", func(db *gorm.DB) *gorm.DB {
			return db.Order("priority DESC, updated_at DESC")
		}).
		Preload("InternalLinks.ChildSeo").
		Preload("Relatives", func(db *gorm.DB) *gorm.DB {
			return db.Order("priority DESC, updated_at DESC")
		}).
		Preload("Relatives.ChildSeo").
		Where("ref_type = ? AND ref_source = ? AND ref_url = ?", refType, refSource, strings.TrimSpace(refURL)).
		First(&seo).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	seo.NormalizeLifecycle()
	return &seo, nil
}

func (r *SeoDomainPostgres) GetByCanonicalURL(ctx context.Context, canonicalURL string) (*seo_domain.SeoDomain, error) {
	var seo seo_domain.SeoDomain
	err := impl.GetDB(ctx, r.db).WithContext(ctx).
		Where("canonical_url = ?", canonicalURL).
		First(&seo).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	seo.NormalizeLifecycle()
	return &seo, nil
}

func (r *SeoDomainPostgres) GetList(ctx context.Context, req *dto.SeoDomainListRequest) ([]seo_domain.SeoDomain, int64, error) {
	var items []seo_domain.SeoDomain
	var total int64

	if req == nil {
		req = &dto.SeoDomainListRequest{}
	}
	req.Normalize()

	query := r.db.WithContext(ctx).
		Model(&seo_domain.SeoDomain{}).
		Select("seo_domain.*, COALESCE(sc.count, 0) AS click_count, sc.updated_at AS click_updated_at, COALESCE(cat.name, '') AS category_name").
		Joins("LEFT JOIN seo_count sc ON sc.seo_id = seo_domain.id").
		Joins("LEFT JOIN seo_category cat ON cat.id = seo_domain.category_id AND cat.deleted_at IS NULL")

	if req.Keyword != nil && strings.TrimSpace(*req.Keyword) != "" {
		keyword := strings.ToLower(strings.TrimSpace(*req.Keyword))
		keyword = strings.ReplaceAll(keyword, "+", " ")
		parts := strings.Fields(keyword)
		if len(parts) > 0 {
			keyword = "%" + strings.Join(parts, "%") + "%"
			query = query.Where(
				"LOWER(seo_domain.slug) LIKE ? OR LOWER(seo_domain.origin_url) LIKE ? OR LOWER(seo_domain.canonical_url) LIKE ? OR LOWER(seo_domain.title) LIKE ? OR LOWER(seo_domain.description) LIKE ? OR LOWER(seo_domain.note) LIKE ? OR LOWER(seo_domain.ref_label) LIKE ?",
				keyword, keyword, keyword, keyword, keyword, keyword, keyword,
			)
		}
	}
	if refTypes := enums.ExpandSEORefTypes(req.RefTypes); len(refTypes) > 0 {
		query = query.Where("seo_domain.ref_type IN ?", refTypes)
	} else if req.RefType != nil {
		query = query.Where("seo_domain.ref_type = ?", *req.RefType)
	}
	if req.RefID != nil {
		query = query.Where("seo_domain.ref_id = ?", *req.RefID)
	}
	if req.RefSource != nil && strings.TrimSpace(*req.RefSource) != "" {
		query = query.Where("seo_domain.ref_source = ?", strings.TrimSpace(*req.RefSource))
	}
	if req.SourceStatus != nil && strings.TrimSpace(*req.SourceStatus) != "" {
		query = query.Where("seo_domain.source_status = ?", strings.TrimSpace(*req.SourceStatus))
	}
	if req.RefMissing != nil {
		query = query.Where("seo_domain.ref_missing = ?", *req.RefMissing)
	}
	if req.Scope != nil && strings.TrimSpace(*req.Scope) != "" {
		query = query.Where("seo_domain.scope = ?", strings.TrimSpace(*req.Scope))
	}
	if req.Published != nil {
		query = query.Where("seo_domain.published = ?", *req.Published)
	}
	if req.IsSiteMap != nil {
		query = query.Where("seo_domain.is_site_map = ?", *req.IsSiteMap)
	}
	if len(req.VisibleOns) > 0 {
		query = query.Where("seo_domain.visible_on IN ?", req.VisibleOns)
	}
	if pageStatuses := seo_domain.PageStatusesFromFilterCodes(req.Status); len(pageStatuses) > 0 {
		query = query.Where("seo_domain.page_status IN ?", pageStatuses)
	}
	if req.IsIndex != nil {
		query = query.Where("seo_domain.is_index = ?", *req.IsIndex)
	}
	if req.IsRobot != nil {
		query = query.Where("seo_domain.is_robot = ?", *req.IsRobot)
	}
	if req.NeedGenerate != nil {
		query = query.Where("seo_domain.need_generate = ?", *req.NeedGenerate)
	}

	// V2 lifecycle/render filters.
	if req.PageStatus != nil && strings.TrimSpace(*req.PageStatus) != "" {
		query = query.Where("seo_domain.page_status = ?", strings.TrimSpace(*req.PageStatus))
	}
	if req.RenderStatus != nil && strings.TrimSpace(*req.RenderStatus) != "" {
		query = query.Where("seo_domain.render_status = ?", strings.TrimSpace(*req.RenderStatus))
	}
	if req.TemplateKey != nil && strings.TrimSpace(*req.TemplateKey) != "" {
		query = query.Where("seo_domain.template_key = ?", strings.TrimSpace(*req.TemplateKey))
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.
		Order("seo_domain.updated_at DESC").
		Offset(req.GetOffset()).
		Limit(req.GetLimitAdmin()).
		Find(&items).Error; err != nil {
		return nil, 0, err
	}

	for i := range items {
		items[i].NormalizeLifecycle()
	}

	return items, total, nil
}

func (r *SeoDomainPostgres) linkedRefQuery(ctx context.Context, req *dto.GetSeoRefsRequest) *gorm.DB {
	if req == nil {
		req = &dto.GetSeoRefsRequest{}
	}
	return impl.GetDB(ctx, r.db).WithContext(ctx).
		Model(&seo_domain.SeoDomain{}).
		Where("seo_domain.published = ?", true).
		Where("(seo_domain.scope = ? OR seo_domain.scope = '')", seo_domain.SeoScopePublic).
		Where(`(
			(seo_domain.ref_type = ? AND seo_domain.ref_id = ?)
			OR EXISTS (
				SELECT 1 FROM seo_refs sr
				WHERE sr.seo_domain_id = seo_domain.id
				  AND sr.ref_type = ?
				  AND sr.ref_id = ?
			)
		)`, req.RefType, req.RefID, req.RefType, req.RefID)
}

func (r *SeoDomainPostgres) GetListByLinkedRef(ctx context.Context, req *dto.GetSeoRefsRequest) ([]seo_domain.SeoDomain, int64, error) {
	var items []seo_domain.SeoDomain
	var total int64

	if req == nil {
		req = &dto.GetSeoRefsRequest{}
	}
	req.Normalize()

	db := r.linkedRefQuery(ctx, req)

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := r.linkedRefQuery(ctx, req).
		Order("seo_domain.published_at DESC NULLS LAST, seo_domain.created_at DESC").
		Offset(req.GetOffset()).
		Limit(req.GetLimit()).
		Find(&items).Error
	if err != nil {
		return nil, 0, err
	}

	for i := range items {
		items[i].NormalizeLifecycle()
	}

	return items, total, nil
}

func (r *SeoDomainPostgres) CountByLinkedRef(ctx context.Context, req *dto.GetSeoRefsRequest) (int64, error) {
	if req == nil {
		req = &dto.GetSeoRefsRequest{}
	}

	var total int64
	if err := r.linkedRefQuery(ctx, req).Count(&total).Error; err != nil {
		return 0, err
	}
	return total, nil
}

func (r *SeoDomainPostgres) CountRecentByRefType(ctx context.Context, days uint32) ([]dto.SeoRefTypeCountItem, error) {
	if days == 0 {
		days = 5
	}

	type row struct {
		RefType uint32 `gorm:"column:ref_type"`
		Total   int64  `gorm:"column:total"`
	}

	since := time.Now().AddDate(0, 0, -int(days))
	var rows []row
	err := impl.GetDB(ctx, r.db).WithContext(ctx).
		Model(&seo_domain.SeoDomain{}).
		Select("ref_type, COUNT(*) AS total").
		Where("updated_at >= ?", since).
		Group("ref_type").
		Scan(&rows).Error
	if err != nil {
		return nil, err
	}

	totals := make(map[uint32]int64, len(rows))
	for _, item := range rows {
		totals[enums.CanonicalSEORefType(item.RefType)] += item.Total
	}

	allTypes := enums.AllSEORefTypes()
	result := make([]dto.SeoRefTypeCountItem, 0, len(allTypes))
	for _, refType := range allTypes {
		total := totals[uint32(refType)]
		if total == 0 {
			continue
		}
		result = append(result, dto.SeoRefTypeCountItem{
			RefType:     refType,
			RefTypeName: refType.Name(),
			Total:       total,
		})
	}
	return result, nil
}

func (r *SeoDomainPostgres) GetTopClicked(ctx context.Context, fromTime time.Time, limit int) ([]dto.SeoTopClickedItem, error) {
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
	err := impl.GetDB(ctx, r.db).WithContext(ctx).
		Raw(`
			SELECT
				sd.id,
				COALESCE(NULLIF(TRIM(sd.title), ''), COALESCE(sd.ref_label, '')) AS name,
				sd.ref_type,
				sd.ref_id,
				SUM(scl.clicked)::bigint AS click_count
			FROM seo_count_log scl
			INNER JOIN seo_domain sd ON sd.id = scl.seo_id AND sd.deleted_at IS NULL
			WHERE scl."timestamp" >= ?
			GROUP BY sd.id, sd.title, sd.ref_label, sd.ref_type, sd.ref_id
			ORDER BY click_count DESC, sd.id DESC
			LIMIT ?
		`, fromTime, limit).
		Scan(&rows).Error
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

func (r *SeoDomainPostgres) CountClicksFromTime(ctx context.Context, fromTime time.Time) (uint64, error) {
	var total int64
	err := impl.GetDB(ctx, r.db).WithContext(ctx).
		Raw(`
			SELECT COALESCE(SUM(clicked), 0)
			FROM seo_count_log
			WHERE "timestamp" >= ?
		`, fromTime).
		Scan(&total).Error
	if err != nil {
		return 0, err
	}
	return uint64(total), nil
}

func (r *SeoDomainPostgres) GetClickDailyStats(ctx context.Context, fromTime time.Time) ([]dto.SeoClickDailyItem, error) {
	type row struct {
		Date  string `gorm:"column:date"`
		Count uint64 `gorm:"column:count"`
	}

	var rows []row
	err := impl.GetDB(ctx, r.db).WithContext(ctx).
		Raw(`
			SELECT
				to_char(("timestamp" AT TIME ZONE 'Asia/Ho_Chi_Minh'), 'YYYY-MM-DD') AS date,
				COALESCE(SUM(clicked), 0)::bigint AS count
			FROM seo_count_log
			WHERE "timestamp" >= ?
			GROUP BY 1
			ORDER BY 1 ASC
		`, fromTime).
		Scan(&rows).Error
	if err != nil {
		return nil, err
	}

	result := make([]dto.SeoClickDailyItem, 0, len(rows))
	for _, item := range rows {
		result = append(result, dto.SeoClickDailyItem{
			Date:  item.Date,
			Count: item.Count,
		})
	}
	return result, nil
}

func (r *SeoDomainPostgres) GetSitemapList(ctx context.Context) ([]seo_domain.SeoDomain, error) {
	var items []seo_domain.SeoDomain

	err := r.db.WithContext(ctx).
		Model(&seo_domain.SeoDomain{}).
		Where("deleted_at IS NULL").
		Where("COALESCE(NULLIF(canonical_url, ''), NULLIF(origin_url, '')) IS NOT NULL").
		Order("site_map_lasted_at DESC NULLS LAST").
		Order("updated_at DESC").
		Find(&items).Error
	if err != nil {
		return nil, err
	}

	for i := range items {
		items[i].NormalizeLifecycle()
	}

	return items, nil
}

func (r *SeoDomainPostgres) Update(ctx context.Context, id uint64, updates map[string]interface{}) error {
	if len(updates) == 0 {
		return nil
	}
	updates["updated_at"] = gorm.Expr("NOW()")
	return impl.GetDB(ctx, r.db).
		Model(&seo_domain.SeoDomain{}).
		Where("id = ?", id).
		Updates(updates).Error
}

func (r *SeoDomainPostgres) UpdateByIDs(ctx context.Context, ids []uint64, updates map[string]interface{}) error {
	if len(ids) == 0 || len(updates) == 0 {
		return nil
	}
	updates["updated_at"] = gorm.Expr("NOW()")
	return impl.GetDB(ctx, r.db).
		Model(&seo_domain.SeoDomain{}).
		Where("id IN ? AND deleted_at IS NULL", ids).
		Updates(updates).Error
}

func (r *SeoDomainPostgres) Delete(ctx context.Context, id uint64) error {
	db := impl.GetDB(ctx, r.db)
	hardDelete := func(model interface{}, query string, args ...interface{}) error {
		q := db.Session(&gorm.Session{NewDB: true}).Unscoped()
		if query != "" {
			q = q.Where(query, args...)
		}
		return q.Delete(model).Error
	}

	if err := hardDelete(&seo_domain.SeoInternalLink{}, "parent_seo_id = ? OR child_seo_id = ?", id, id); err != nil {
		return err
	}
	if err := hardDelete(&seo_domain.SeoRelative{}, "parent_seo_id = ? OR child_seo_id = ?", id, id); err != nil {
		return err
	}
	if err := hardDelete(&seo_domain.SeoGenerationLog{}, "seo_domain_id = ?", id); err != nil {
		return err
	}
	if err := hardDelete(&seo_domain.SeoRef{}, "seo_domain_id = ?", id); err != nil {
		return err
	}
	if err := hardDelete(&seo_domain.SeoCount{}, "seo_id = ?", id); err != nil {
		return err
	}
	if err := hardDelete(&seo_domain.SeoCountLog{}, "seo_id = ?", id); err != nil {
		return err
	}
	return hardDelete(&seo_domain.SeoDomain{}, "id = ?", id)
}

func (r *SeoDomainPostgres) ExistByCanonicalURL(ctx context.Context, canonicalURL string, ignoreID *uint64) (bool, error) {
	var count int64
	query := r.db.WithContext(ctx).
		Model(&seo_domain.SeoDomain{}).
		Where("canonical_url = ?", canonicalURL)
	if ignoreID != nil {
		query = query.Where("id <> ?", *ignoreID)
	}
	err := query.Count(&count).Error
	return count > 0, err
}

func (r *SeoDomainPostgres) MarkNeedGenerateByRef(ctx context.Context, refType uint32, refID uint64, sourceUpdatedAt time.Time) error {
	return impl.GetDB(ctx, r.db).
		Model(&seo_domain.SeoDomain{}).
		Where("ref_type = ? AND ref_id = ?", refType, refID).
		Updates(map[string]interface{}{
			"need_generate":     true,
			"render_status":     seo_domain.SeoRenderStatusPending,
			"source_updated_at": sourceUpdatedAt,
			"updated_at":        gorm.Expr("NOW()"),
		}).Error
}

func (r *SeoDomainPostgres) MarkNeedGenerateByRefSource(ctx context.Context, refType uint32, refSource string, refID uint64, sourceUpdatedAt time.Time) error {
	return impl.GetDB(ctx, r.db).
		Model(&seo_domain.SeoDomain{}).
		Where("ref_type = ? AND ref_source = ? AND ref_id = ?", refType, refSource, refID).
		Updates(map[string]interface{}{
			"need_generate":     true,
			"render_status":     seo_domain.SeoRenderStatusPending,
			"source_updated_at": sourceUpdatedAt,
			"updated_at":        gorm.Expr("NOW()"),
		}).Error
}

func (r *SeoDomainPostgres) UpdateSourceStateByRefSource(ctx context.Context, refType uint32, refSource string, refID uint64, updates map[string]interface{}) error {
	if len(updates) == 0 {
		return nil
	}
	if _, ok := updates["need_generate"]; ok {
		if _, hasRender := updates["render_status"]; !hasRender {
			updates["render_status"] = seo_domain.SeoRenderStatusPending
		}
	}
	updates["updated_at"] = gorm.Expr("NOW()")
	return impl.GetDB(ctx, r.db).
		Model(&seo_domain.SeoDomain{}).
		Where("ref_type = ? AND ref_source = ? AND ref_id = ?", refType, refSource, refID).
		Updates(updates).Error
}

func (r *SeoDomainPostgres) GetSitemapItems(ctx context.Context, scope *string) ([]dto.SeoSitemapItem, error) {
	var rows []dto.SeoSitemapItem

	query := r.db.WithContext(ctx).
		Model(&seo_domain.SeoDomain{}).
		Select(`
			canonical_url AS canonical_url,
			COALESCE(site_map_lasted_at, updated_at) AS site_map_lasted_at,
			updated_at AS updated_at,
			sitemap_priority AS sitemap_priority,
			sitemap_change_freq AS sitemap_change_freq
		`).
		Where("published = ? AND is_site_map = ?", true, true).
		Where("page_status = ?", seo_domain.SeoPageStatusPublished).
		Where("is_index = ?", true).
		Where("scope = ?", seo_domain.SeoScopePublic)

	if scope != nil && strings.TrimSpace(*scope) != "" {
		query = query.Where("scope = ?", strings.TrimSpace(*scope))
	}

	if err := query.Order("updated_at DESC").Scan(&rows).Error; err != nil {
		return nil, err
	}
	return rows, nil
}

func (r *SeoDomainPostgres) GetNeedGenerateList(ctx context.Context, limit int) ([]seo_domain.SeoDomain, error) {
	var items []seo_domain.SeoDomain
	if limit <= 0 || limit > 1000 {
		limit = 100
	}
	err := r.db.WithContext(ctx).
		Where("need_generate = ? AND published = ?", true, true).
		// Where("page_status = ?", seo_domain.SeoPageStatusPublished).
		// Where("render_status IN ?", []string{seo_domain.SeoRenderStatusPending, seo_domain.SeoRenderStatusFailed}).
		Order("source_updated_at ASC NULLS FIRST, updated_at ASC").
		Limit(limit).
		Find(&items).Error
	return items, err
}

// -----------------------------------------------------------------------------
// Internal links
// -----------------------------------------------------------------------------

type SeoInternalLinkPostgres struct {
	db *gorm.DB
}

func NewSeoInternalLinkPostgres(db *gorm.DB) repo.SeoInternalLinkRepo {
	return &SeoInternalLinkPostgres{db: db}
}

func (r *SeoInternalLinkPostgres) Create(ctx context.Context, link *seo_domain.SeoInternalLink) (*seo_domain.SeoInternalLink, error) {
	if err := impl.GetDB(ctx, r.db).Create(link).Error; err != nil {
		return nil, err
	}
	return r.GetByID(ctx, link.ID)
}

func (r *SeoInternalLinkPostgres) GetByID(ctx context.Context, id uint64) (*seo_domain.SeoInternalLink, error) {
	var link seo_domain.SeoInternalLink
	err := r.db.WithContext(ctx).
		Preload("ParentSeo").
		Preload("ChildSeo").
		First(&link, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &link, nil
}

func (r *SeoInternalLinkPostgres) GetList(ctx context.Context, req *dto.SeoInternalLinkListRequest) ([]seo_domain.SeoInternalLink, int64, error) {
	var items []seo_domain.SeoInternalLink
	var total int64

	if req == nil {
		req = &dto.SeoInternalLinkListRequest{}
	}
	req.Normalize()

	query := r.db.WithContext(ctx).Model(&seo_domain.SeoInternalLink{})

	if req.ParentSeoID != nil {
		query = query.Where("parent_seo_id = ?", *req.ParentSeoID)
	}
	if req.ChildSeoID != nil {
		query = query.Where("child_seo_id = ?", *req.ChildSeoID)
	}
	if req.LinkType != nil && strings.TrimSpace(*req.LinkType) != "" {
		query = query.Where("link_type = ?", strings.TrimSpace(*req.LinkType))
	}
	if req.Keyword != nil && strings.TrimSpace(*req.Keyword) != "" {
		keyword := "%" + strings.TrimSpace(*req.Keyword) + "%"
		query = query.Where("link ILIKE ? OR title ILIKE ?", keyword, keyword)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if err := query.
		Preload("ParentSeo").
		Preload("ChildSeo").
		Order("priority DESC, updated_at DESC").
		Offset(req.GetOffset()).
		Limit(req.GetLimit()).
		Find(&items).Error; err != nil {
		return nil, 0, err
	}
	return items, total, nil
}

func (r *SeoInternalLinkPostgres) Update(ctx context.Context, id uint64, updates map[string]interface{}) error {
	if len(updates) == 0 {
		return nil
	}
	updates["updated_at"] = gorm.Expr("NOW()")
	return impl.GetDB(ctx, r.db).
		Model(&seo_domain.SeoInternalLink{}).
		Where("id = ?", id).
		Updates(updates).Error
}

func (r *SeoInternalLinkPostgres) Delete(ctx context.Context, id uint64) error {
	return impl.GetDB(ctx, r.db).Delete(&seo_domain.SeoInternalLink{}, id).Error
}

func (r *SeoInternalLinkPostgres) DeleteBySeoDomainID(ctx context.Context, seoDomainID uint64) error {
	return impl.GetDB(ctx, r.db).Session(&gorm.Session{NewDB: true}).Unscoped().
		Where("parent_seo_id = ? OR child_seo_id = ?", seoDomainID, seoDomainID).
		Delete(&seo_domain.SeoInternalLink{}).Error
}

func (r *SeoInternalLinkPostgres) Exist(ctx context.Context, parentSeoID uint64, childSeoID uint64, linkType string, ignoreID *uint64) (bool, error) {
	var count int64
	query := r.db.WithContext(ctx).
		Model(&seo_domain.SeoInternalLink{}).
		Where("parent_seo_id = ? AND child_seo_id = ? AND link_type = ?", parentSeoID, childSeoID, linkType)
	if ignoreID != nil {
		query = query.Where("id <> ?", *ignoreID)
	}
	err := query.Count(&count).Error
	return count > 0, err
}

// -----------------------------------------------------------------------------
// Relatives
// -----------------------------------------------------------------------------

type SeoRelativePostgres struct {
	db *gorm.DB
}

func NewSeoRelativePostgres(db *gorm.DB) repo.SeoRelativeRepo {
	return &SeoRelativePostgres{db: db}
}

func (r *SeoRelativePostgres) Create(ctx context.Context, relative *seo_domain.SeoRelative) (*seo_domain.SeoRelative, error) {
	if err := impl.GetDB(ctx, r.db).Create(relative).Error; err != nil {
		return nil, err
	}
	return r.GetByID(ctx, relative.ID)
}

func (r *SeoRelativePostgres) GetByID(ctx context.Context, id uint64) (*seo_domain.SeoRelative, error) {
	var relative seo_domain.SeoRelative
	err := r.db.WithContext(ctx).
		Preload("ParentSeo").
		Preload("ChildSeo").
		First(&relative, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &relative, nil
}

func (r *SeoRelativePostgres) GetList(ctx context.Context, req *dto.SeoRelativeListRequest) ([]seo_domain.SeoRelative, int64, error) {
	var items []seo_domain.SeoRelative
	var total int64

	if req == nil {
		req = &dto.SeoRelativeListRequest{}
	}
	req.Normalize()

	query := r.db.WithContext(ctx).Model(&seo_domain.SeoRelative{})

	if req.ParentSeoID != nil {
		query = query.Where("parent_seo_id = ?", *req.ParentSeoID)
	}
	if req.ChildSeoID != nil {
		query = query.Where("child_seo_id = ?", *req.ChildSeoID)
	}
	if req.RelationType != nil && strings.TrimSpace(*req.RelationType) != "" {
		query = query.Where("relation_type = ?", strings.TrimSpace(*req.RelationType))
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if err := query.
		Preload("ParentSeo").
		Preload("ChildSeo").
		Order("priority DESC, updated_at DESC").
		Offset(req.GetOffset()).
		Limit(req.GetLimit()).
		Find(&items).Error; err != nil {
		return nil, 0, err
	}
	return items, total, nil
}

func (r *SeoRelativePostgres) Delete(ctx context.Context, id uint64) error {
	return impl.GetDB(ctx, r.db).Delete(&seo_domain.SeoRelative{}, id).Error
}

func (r *SeoRelativePostgres) DeleteBySeoDomainID(ctx context.Context, seoDomainID uint64) error {
	return impl.GetDB(ctx, r.db).Session(&gorm.Session{NewDB: true}).Unscoped().
		Where("parent_seo_id = ? OR child_seo_id = ?", seoDomainID, seoDomainID).
		Delete(&seo_domain.SeoRelative{}).Error
}

func (r *SeoRelativePostgres) Exist(ctx context.Context, parentSeoID uint64, childSeoID uint64, relationType string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&seo_domain.SeoRelative{}).
		Where("parent_seo_id = ? AND child_seo_id = ? AND relation_type = ?", parentSeoID, childSeoID, relationType).
		Count(&count).Error
	return count > 0, err
}
