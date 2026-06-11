package postgre

import (
	"context"
	"crm/infra/impl"
	seo_domain "crm/internal/domain/seo"
	"crm/internal/dto"
	"crm/internal/repo"

	"gorm.io/gorm"
)

type SeoDomainPostgres struct {
	db *gorm.DB
}

func NewSeoDomainPostgres(db *gorm.DB) repo.SeoDomainRepo {
	return &SeoDomainPostgres{db: db}
}

func (r *SeoDomainPostgres) Create(ctx context.Context, seo *seo_domain.SeoDomain) (*seo_domain.SeoDomain, error) {
	if err := impl.GetDB(ctx, r.db).Create(seo).Error; err != nil {
		return nil, err
	}
	return r.GetByID(ctx, seo.ID)
}

func (r *SeoDomainPostgres) GetByID(ctx context.Context, id uint64) (*seo_domain.SeoDomain, error) {
	var seo seo_domain.SeoDomain
	err := r.db.WithContext(ctx).
		Preload("InternalLinks", "deleted_at IS NULL").
		Preload("InternalLinks.ChildSeo").
		Preload("Relatives", "deleted_at IS NULL").
		Preload("Relatives.ChildSeo").
		Where("id = ? AND deleted_at IS NULL", id).
		First(&seo).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &seo, nil
}

func (r *SeoDomainPostgres) GetList(ctx context.Context, req *dto.SeoDomainListRequest) ([]seo_domain.SeoDomain, int64, error) {
	var items []seo_domain.SeoDomain
	var total int64

	query := r.db.WithContext(ctx).
		Model(&seo_domain.SeoDomain{}).
		Where("deleted_at IS NULL")

	if req.Keyword != nil && *req.Keyword != "" {
		keyword := "%" + *req.Keyword + "%"
		query = query.Where("origin_url ILIKE ? OR canonical_url ILIKE ? OR note ILIKE ?", keyword, keyword, keyword)
	}
	if req.Classify != nil {
		query = query.Where("classify = ?", *req.Classify)
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
	if req.RefType != nil {
		query = query.Where("ref_type = ?", *req.RefType)
	}
	if req.RefID != nil {
		query = query.Where("ref_id = ?", *req.RefID)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.
		Preload("InternalLinks", "deleted_at IS NULL").
		Preload("Relatives", "deleted_at IS NULL").
		Order("updated_at DESC").
		Offset(req.GetOffset()).
		Limit(req.GetLimit()).
		Find(&items).Error; err != nil {
		return nil, 0, err
	}

	return items, total, nil
}

func (r *SeoDomainPostgres) Update(ctx context.Context, id uint64, updates map[string]interface{}) error {
	return impl.GetDB(ctx, r.db).
		Model(&seo_domain.SeoDomain{}).
		Where("id = ? AND deleted_at IS NULL", id).
		Updates(updates).Error
}

func (r *SeoDomainPostgres) Delete(ctx context.Context, id uint64) error {
	return impl.GetDB(ctx, r.db).
		Model(&seo_domain.SeoDomain{}).
		Where("id = ? AND deleted_at IS NULL", id).
		Update("deleted_at", gorm.Expr("NOW()")).Error
}

func (r *SeoDomainPostgres) ExistByCanonicalURL(ctx context.Context, canonicalURL string, ignoreID *uint64) (bool, error) {
	var count int64
	query := r.db.WithContext(ctx).
		Model(&seo_domain.SeoDomain{}).
		Where("canonical_url = ? AND deleted_at IS NULL", canonicalURL)
	if ignoreID != nil {
		query = query.Where("id <> ?", *ignoreID)
	}
	err := query.Count(&count).Error
	return count > 0, err
}

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
		Where("id = ? AND deleted_at IS NULL", id).
		First(&link).Error
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
	query := r.db.WithContext(ctx).
		Model(&seo_domain.SeoInternalLink{}).
		Where("deleted_at IS NULL")

	if req.ParentSeoID != nil {
		query = query.Where("parent_seo_id = ?", *req.ParentSeoID)
	}
	if req.ChildSeoID != nil {
		query = query.Where("child_seo_id = ?", *req.ChildSeoID)
	}
	if req.Keyword != nil && *req.Keyword != "" {
		keyword := "%" + *req.Keyword + "%"
		query = query.Where("link ILIKE ? OR title ILIKE ?", keyword, keyword)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if err := query.
		Preload("ParentSeo").
		Preload("ChildSeo").
		Order("updated_at DESC").
		Offset(req.GetOffset()).
		Limit(req.GetLimit()).
		Find(&items).Error; err != nil {
		return nil, 0, err
	}
	return items, total, nil
}

func (r *SeoInternalLinkPostgres) Update(ctx context.Context, id uint64, updates map[string]interface{}) error {
	return impl.GetDB(ctx, r.db).
		Model(&seo_domain.SeoInternalLink{}).
		Where("id = ? AND deleted_at IS NULL", id).
		Updates(updates).Error
}

func (r *SeoInternalLinkPostgres) Delete(ctx context.Context, id uint64) error {
	return impl.GetDB(ctx, r.db).
		Model(&seo_domain.SeoInternalLink{}).
		Where("id = ? AND deleted_at IS NULL", id).
		Update("deleted_at", gorm.Expr("NOW()")).Error
}

func (r *SeoInternalLinkPostgres) Exist(ctx context.Context, parentSeoID uint64, childSeoID uint64, ignoreID *uint64) (bool, error) {
	var count int64
	query := r.db.WithContext(ctx).
		Model(&seo_domain.SeoInternalLink{}).
		Where("parent_seo_id = ? AND child_seo_id = ? AND deleted_at IS NULL", parentSeoID, childSeoID)
	if ignoreID != nil {
		query = query.Where("id <> ?", *ignoreID)
	}
	err := query.Count(&count).Error
	return count > 0, err
}

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
		Where("id = ? AND deleted_at IS NULL", id).
		First(&relative).Error
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
	query := r.db.WithContext(ctx).
		Model(&seo_domain.SeoRelative{}).
		Where("deleted_at IS NULL")

	if req.ParentSeoID != nil {
		query = query.Where("parent_seo_id = ?", *req.ParentSeoID)
	}
	if req.ChildSeoID != nil {
		query = query.Where("child_seo_id = ?", *req.ChildSeoID)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if err := query.
		Preload("ParentSeo").
		Preload("ChildSeo").
		Order("updated_at DESC").
		Offset(req.GetOffset()).
		Limit(req.GetLimit()).
		Find(&items).Error; err != nil {
		return nil, 0, err
	}
	return items, total, nil
}

func (r *SeoRelativePostgres) Delete(ctx context.Context, id uint64) error {
	return impl.GetDB(ctx, r.db).
		Model(&seo_domain.SeoRelative{}).
		Where("id = ? AND deleted_at IS NULL", id).
		Update("deleted_at", gorm.Expr("NOW()")).Error
}

func (r *SeoRelativePostgres) Exist(ctx context.Context, parentSeoID uint64, childSeoID uint64) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&seo_domain.SeoRelative{}).
		Where("parent_seo_id = ? AND child_seo_id = ? AND deleted_at IS NULL", parentSeoID, childSeoID).
		Count(&count).Error
	return count > 0, err
}
