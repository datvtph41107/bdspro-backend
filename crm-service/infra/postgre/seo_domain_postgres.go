package postgre

import (
	"context"
	"crm/infra/impl"
	seo_domain "crm/internal/domain/seo"
	"crm/internal/dto"
	"crm/internal/repo"
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
	return &seo, nil
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

	query := r.db.WithContext(ctx).Model(&seo_domain.SeoDomain{})

	if req.Keyword != nil && strings.TrimSpace(*req.Keyword) != "" {
		keyword := "%" + strings.TrimSpace(*req.Keyword) + "%"
		query = query.Where(
			"slug ILIKE ? OR origin_url ILIKE ? OR canonical_url ILIKE ? OR title ILIKE ? OR description ILIKE ? OR note ILIKE ? OR ref_label ILIKE ?",
			keyword, keyword, keyword, keyword, keyword, keyword, keyword,
		)
	}
	if req.RefType != nil {
		query = query.Where("ref_type = ?", *req.RefType)
	}
	if req.RefID != nil {
		query = query.Where("ref_id = ?", *req.RefID)
	}
	if req.RefSource != nil && strings.TrimSpace(*req.RefSource) != "" {
		query = query.Where("ref_source = ?", strings.TrimSpace(*req.RefSource))
	}
	if req.SourceStatus != nil && strings.TrimSpace(*req.SourceStatus) != "" {
		query = query.Where("source_status = ?", strings.TrimSpace(*req.SourceStatus))
	}
	if req.RefMissing != nil {
		query = query.Where("ref_missing = ?", *req.RefMissing)
	}
	if req.Scope != nil && strings.TrimSpace(*req.Scope) != "" {
		query = query.Where("scope = ?", strings.TrimSpace(*req.Scope))
	}
	if req.Published != nil {
		query = query.Where("published = ?", *req.Published)
	}
	if req.IsSiteMap != nil {
		query = query.Where("is_site_map = ?", *req.IsSiteMap)
	}
	if req.IsIndex != nil {
		query = query.Where("is_index = ?", *req.IsIndex)
	}
	if req.IsRobot != nil {
		query = query.Where("is_robot = ?", *req.IsRobot)
	}
	if req.NeedGenerate != nil {
		query = query.Where("need_generate = ?", *req.NeedGenerate)
	}

	// V2 lifecycle/render filters.
	if req.PageStatus != nil && strings.TrimSpace(*req.PageStatus) != "" {
		query = query.Where("page_status = ?", strings.TrimSpace(*req.PageStatus))
	}
	if req.RenderStatus != nil && strings.TrimSpace(*req.RenderStatus) != "" {
		query = query.Where("render_status = ?", strings.TrimSpace(*req.RenderStatus))
	}
	if req.TemplateKey != nil && strings.TrimSpace(*req.TemplateKey) != "" {
		query = query.Where("template_key = ?", strings.TrimSpace(*req.TemplateKey))
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.
		Order("updated_at DESC").
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

func (r *SeoDomainPostgres) GetListByLinkedRef(ctx context.Context, req *dto.GetSeoRefsRequest) ([]seo_domain.SeoDomain, int64, error) {
	var items []seo_domain.SeoDomain
	var total int64

	if req == nil {
		req = &dto.GetSeoRefsRequest{}
	}
	req.Normalize()

	buildQuery := func(db *gorm.DB) *gorm.DB {
		return db.WithContext(ctx).
			Model(&seo_domain.SeoDomain{}).
			Joins("INNER JOIN seo_refs sr ON sr.seo_domain_id = seo_domain.id").
			Where("sr.ref_type = ? AND sr.ref_id = ?", req.RefType, req.RefID).
			Where("seo_domain.published = ?", true).
			Where("(seo_domain.scope = ? OR seo_domain.scope = '')", seo_domain.SeoScopePublic)
	}

	db := impl.GetDB(ctx, r.db)

	if err := buildQuery(db).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := buildQuery(db).
		Select("seo_domain.*").
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

func (r *SeoDomainPostgres) Delete(ctx context.Context, id uint64) error {
	return impl.GetDB(ctx, r.db).Delete(&seo_domain.SeoDomain{}, id).Error
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
		Where("page_status = ?", seo_domain.SeoPageStatusPublished).
		Where("render_status IN ?", []string{seo_domain.SeoRenderStatusPending, seo_domain.SeoRenderStatusFailed}).
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
	return impl.GetDB(ctx, r.db).
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
	return impl.GetDB(ctx, r.db).
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
