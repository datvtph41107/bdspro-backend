//go:build ignore

package repository

import (
	"context"

	seo_domain "mdn4/modules/crm/entity/seo"
)

type SeoCountRepo interface {
	// FlushClicks writes batch logs and increments aggregate seo_count in one transaction.
	FlushClicks(ctx context.Context, logs []*seo_domain.SeoCountLog) error
}
