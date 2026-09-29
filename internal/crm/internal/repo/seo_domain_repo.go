package repo

import (
	"context"
	seo_domain "crm/internal/domain/seo"
	"crm/internal/dto"
	"time"
)

type SeoDomainRepo interface {
	Create(ctx context.Context, seo *seo_domain.SeoDomain) (*seo_domain.SeoDomain, error)
	GetByID(ctx context.Context, id uint64) (*seo_domain.SeoDomain, error)
	GetByRef(ctx context.Context, refType uint32, refID uint64) (*seo_domain.SeoDomain, error)
	GetByRefSource(ctx context.Context, refType uint32, refSource string, refID uint64) (*seo_domain.SeoDomain, error)
	GetByRefIdentity(ctx context.Context, refType uint32, refSource string, refURL string) (*seo_domain.SeoDomain, error)
	GetBySlug(ctx context.Context, slug string) (*seo_domain.SeoDomain, error)
	GetByCanonicalURL(ctx context.Context, canonicalURL string) (*seo_domain.SeoDomain, error)
	GetList(ctx context.Context, req *dto.SeoDomainListRequest) ([]seo_domain.SeoDomain, int64, error)
	GetListByLinkedRef(ctx context.Context, req *dto.GetSeoRefsRequest) ([]seo_domain.SeoDomain, int64, error)
	GetSitemapList(ctx context.Context) ([]seo_domain.SeoDomain, error)
	Update(ctx context.Context, id uint64, updates map[string]interface{}) error
	Delete(ctx context.Context, id uint64) error

	ExistByCanonicalURL(ctx context.Context, canonicalURL string, ignoreID *uint64) (bool, error)

	MarkNeedGenerateByRef(ctx context.Context, refType uint32, refID uint64, sourceUpdatedAt time.Time) error
	MarkNeedGenerateByRefSource(ctx context.Context, refType uint32, refSource string, refID uint64, sourceUpdatedAt time.Time) error
	UpdateSourceStateByRefSource(ctx context.Context, refType uint32, refSource string, refID uint64, updates map[string]interface{}) error

	GetSitemapItems(ctx context.Context, scope *string) ([]dto.SeoSitemapItem, error)
	GetNeedGenerateList(ctx context.Context, limit int) ([]seo_domain.SeoDomain, error)
	GetPublishedStaticPaths(ctx context.Context, page uint32, size uint32) ([]seo_domain.SeoDomain, int64, error)
	GetPublishedSitemapItems(ctx context.Context) ([]seo_domain.SeoDomain, error)
	GetNeedGenerateItems(ctx context.Context, limit int) ([]seo_domain.SeoDomain, error)

	MarkGenerated(ctx context.Context, id uint64, staticPath string, staticHash string, generatedAt time.Time) error

	// V2 render lifecycle helpers.
	MarkRenderPending(ctx context.Context, id uint64) error
	MarkRenderRunning(ctx context.Context, id uint64) error
	MarkRenderSuccess(ctx context.Context, id uint64, renderedHTML string, staticPath string, staticHash string, generatedAt time.Time) error
	MarkRenderFailed(ctx context.Context, id uint64, errorMessage string) error
}

type SeoGenerationLogRepo interface {
	Create(ctx context.Context, log *seo_domain.SeoGenerationLog) (*seo_domain.SeoGenerationLog, error)
	GetByID(ctx context.Context, id uint64) (*seo_domain.SeoGenerationLog, error)
	GetList(ctx context.Context, req *dto.SeoGenerationLogListRequest) ([]seo_domain.SeoGenerationLog, int64, error)
	Update(ctx context.Context, id uint64, updates map[string]interface{}) error
	MarkRunning(ctx context.Context, id uint64, startedAt time.Time) error
	MarkSuccess(ctx context.Context, id uint64, staticPath string, staticHash string, finishedAt time.Time) error
	MarkFailed(ctx context.Context, id uint64, errorMessage string, finishedAt time.Time) error
}

type SeoInternalLinkRepo interface {
	Create(ctx context.Context, link *seo_domain.SeoInternalLink) (*seo_domain.SeoInternalLink, error)
	GetByID(ctx context.Context, id uint64) (*seo_domain.SeoInternalLink, error)
	GetList(ctx context.Context, req *dto.SeoInternalLinkListRequest) ([]seo_domain.SeoInternalLink, int64, error)
	Update(ctx context.Context, id uint64, updates map[string]interface{}) error
	Delete(ctx context.Context, id uint64) error
	DeleteBySeoDomainID(ctx context.Context, seoDomainID uint64) error
	Exist(ctx context.Context, parentSeoID uint64, childSeoID uint64, linkType string, ignoreID *uint64) (bool, error)
}

type SeoRelativeRepo interface {
	Create(ctx context.Context, relative *seo_domain.SeoRelative) (*seo_domain.SeoRelative, error)
	GetByID(ctx context.Context, id uint64) (*seo_domain.SeoRelative, error)
	GetList(ctx context.Context, req *dto.SeoRelativeListRequest) ([]seo_domain.SeoRelative, int64, error)
	Delete(ctx context.Context, id uint64) error
	DeleteBySeoDomainID(ctx context.Context, seoDomainID uint64) error
	Exist(ctx context.Context, parentSeoID uint64, childSeoID uint64, relationType string) (bool, error)
}
