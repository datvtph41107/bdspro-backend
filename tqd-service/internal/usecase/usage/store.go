package usage

import (
	commonmetering "common/metering"
	"context"
	"time"
	"tqd/internal/access"
)

/**
 * Store lưu và đọc usage events.
 */
type Store interface {
	SaveUsageEvent(ctx context.Context, event Event) (created bool, err error)

	SumUsage(
		ctx context.Context,
		subject access.Subject,
		meterCode commonmetering.Code,
		periodStart time.Time,
		periodEnd time.Time,
	) (int64, error)
}
