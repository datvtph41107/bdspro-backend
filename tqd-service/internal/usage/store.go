package usage

import (
	commonaccess "common/access"
	commonmetering "common/metering"
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
		meterCode commonmetering.Code,
		periodStart time.Time,
		periodEnd time.Time,
	) (int64, error)
}
