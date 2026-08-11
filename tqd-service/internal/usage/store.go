package usage

import (
	commonaccess "common/access"
	"common/businessoperation"
	"context"
	"time"
)

/**
 * Store lưu và đọc usage events.
 */
type Store interface {
	SaveUsageEvent(ctx context.Context, event Event) (created bool, err error)

	SumUsage(
		ctx context.Context,
		subject commonaccess.Subject,
		operation businessoperation.Code,
		periodStart time.Time,
		periodEnd time.Time,
	) (int64, error)
}
