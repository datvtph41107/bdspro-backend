package repo

import (
	"context"

	seo_domain "crm/internal/domain/seo"
)

type SeoCountRepo interface {
	// FlushClicks writes batch logs and increments aggregate seo_count in one transaction.
	FlushClicks(ctx context.Context, logs []*seo_domain.SeoCountLog) error
}
