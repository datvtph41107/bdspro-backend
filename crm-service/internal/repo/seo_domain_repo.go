package repo

import (
	"context"
	seo_domain "crm/internal/domain/seo"
	"crm/internal/dto"
)

type SeoDomainRepo interface {
	Create(ctx context.Context, seo *seo_domain.SeoDomain) (*seo_domain.SeoDomain, error)
	GetByID(ctx context.Context, id uint64) (*seo_domain.SeoDomain, error)
	GetList(ctx context.Context, req *dto.SeoDomainListRequest) ([]seo_domain.SeoDomain, int64, error)
	Update(ctx context.Context, id uint64, updates map[string]interface{}) error
	Delete(ctx context.Context, id uint64) error
	ExistByCanonicalURL(ctx context.Context, canonicalURL string, ignoreID *uint64) (bool, error)
}

type SeoInternalLinkRepo interface {
	Create(ctx context.Context, link *seo_domain.SeoInternalLink) (*seo_domain.SeoInternalLink, error)
	GetByID(ctx context.Context, id uint64) (*seo_domain.SeoInternalLink, error)
	GetList(ctx context.Context, req *dto.SeoInternalLinkListRequest) ([]seo_domain.SeoInternalLink, int64, error)
	Update(ctx context.Context, id uint64, updates map[string]interface{}) error
	Delete(ctx context.Context, id uint64) error
	Exist(ctx context.Context, parentSeoID uint64, childSeoID uint64, ignoreID *uint64) (bool, error)
}

type SeoRelativeRepo interface {
	Create(ctx context.Context, relative *seo_domain.SeoRelative) (*seo_domain.SeoRelative, error)
	GetByID(ctx context.Context, id uint64) (*seo_domain.SeoRelative, error)
	GetList(ctx context.Context, req *dto.SeoRelativeListRequest) ([]seo_domain.SeoRelative, int64, error)
	Delete(ctx context.Context, id uint64) error
	Exist(ctx context.Context, parentSeoID uint64, childSeoID uint64) (bool, error)
}
