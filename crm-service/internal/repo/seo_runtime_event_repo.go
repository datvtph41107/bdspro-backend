package repo

import (
	"context"

	seo_domain "crm/internal/domain/seo"
)

type SeoRuntimeEventRepo interface {
	CreateIdempotent(ctx context.Context, event *seo_domain.SeoRuntimeEvent) (created bool, err error)
}
