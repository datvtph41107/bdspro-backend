package postgre

import (
	"context"
	"crm/infra/impl"
	seo_domain "crm/internal/domain/seo"
	"crm/internal/enums"
	"time"

	"gorm.io/gorm"
)

func (r *SeoDomainPostgres) GetBySlug(ctx context.Context, slug string) (*seo_domain.SeoDomain, error) {
	var seo seo_domain.SeoDomain

	err := impl.GetDB(ctx, r.db).
		WithContext(ctx).
		Preload("InternalLinks", func(db *gorm.DB) *gorm.DB {
			return db.Where("deleted_at IS NULL").Order("priority DESC, updated_at DESC")
		}).
		Preload("InternalLinks.ChildSeo").
		Preload("Relatives", func(db *gorm.DB) *gorm.DB {
			return db.Where("deleted_at IS NULL").Order("priority DESC, updated_at DESC")
		}).
		Preload("Relatives.ChildSeo").
		Where("slug = ?", slug).
		First(&seo).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &seo, nil
}

func (r *SeoDomainPostgres) GetPublishedStaticPaths(ctx context.Context, page uint32, size uint32) ([]seo_domain.SeoDomain, int64, error) {
	var items []seo_domain.SeoDomain
	var total int64

	if page == 0 {
		page = 1
	}
	if size == 0 {
		size = 500
	}
	if size > 1000 {
		size = 1000
	}
	offset := int((page - 1) * size)

	query := impl.GetDB(ctx, r.db).
		WithContext(ctx).
		Model(&seo_domain.SeoDomain{}).
		Where("published = ?", true).
		Where("is_index = ?", true).
		Where("scope = ?", seo_domain.SeoScopePublic).
		Where("classify <> ?", enums.ESeoClassifyN)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Order("updated_at DESC").Offset(offset).Limit(int(size)).Find(&items).Error; err != nil {
		return nil, 0, err
	}
	return items, total, nil
}

func (r *SeoDomainPostgres) GetPublishedSitemapItems(ctx context.Context) ([]seo_domain.SeoDomain, error) {
	var items []seo_domain.SeoDomain
	err := impl.GetDB(ctx, r.db).
		WithContext(ctx).
		Where("published = ?", true).
		Where("is_site_map = ?", true).
		Where("is_index = ?", true).
		Where("scope = ?", seo_domain.SeoScopePublic).
		Where("classify <> ?", enums.ESeoClassifyN).
		Order("updated_at DESC").
		Find(&items).Error
	return items, err
}

func (r *SeoDomainPostgres) GetNeedGenerateItems(ctx context.Context, limit int) ([]seo_domain.SeoDomain, error) {
	var items []seo_domain.SeoDomain
	if limit <= 0 {
		limit = 200
	}
	if limit > 1000 {
		limit = 1000
	}
	err := impl.GetDB(ctx, r.db).
		WithContext(ctx).
		Where("need_generate = ?", true).
		Where("published = ?", true).
		Order("source_updated_at ASC NULLS FIRST, updated_at ASC").
		Limit(limit).
		Find(&items).Error
	return items, err
}

func (r *SeoDomainPostgres) MarkGenerated(ctx context.Context, id uint64, staticPath string, staticHash string, generatedAt time.Time) error {
	return impl.GetDB(ctx, r.db).
		WithContext(ctx).
		Model(&seo_domain.SeoDomain{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"need_generate":      false,
			"generated_at":       generatedAt,
			"site_map_lasted_at": generatedAt,
			"static_html_path":   staticPath,
			"static_html_hash":   staticHash,
			"updated_at":         gorm.Expr("NOW()"),
		}).Error
}
