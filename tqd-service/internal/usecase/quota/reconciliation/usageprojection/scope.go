package usageprojection

import (
	commonmetering "common/metering"
	"context"
	"time"
	"tqd/internal/access"
)

/**
 * Scope xác định đúng counter cần kiểm tra giữa PostgreSQL và Redis.
 */
type Scope struct {
	Subject     access.Subject
	MeterCode   commonmetering.Code
	PeriodStart time.Time
	PeriodEnd   time.Time
}

/**
 * ScopeStore trả các usage scope đã có dữ liệu bền vững.
 */
type ScopeStore interface {
	ListUsageScopes(ctx context.Context, offset int, limit int) ([]Scope, error)
}
