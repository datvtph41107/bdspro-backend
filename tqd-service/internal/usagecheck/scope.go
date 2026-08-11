package usagecheck

import (
	commonaccess "common/access"
	"common/businessoperation"
	"context"
	"time"
)

/**
 * Scope xác định đúng counter cần kiểm tra giữa PostgreSQL và Redis.
 */
type Scope struct {
	Subject     commonaccess.Subject
	Operation   businessoperation.Code
	PeriodStart time.Time
	PeriodEnd   time.Time
}

/**
 * ScopeStore trả các usage scope đã có dữ liệu bền vững.
 */
type ScopeStore interface {
	ListUsageScopes(ctx context.Context, offset int, limit int) ([]Scope, error)
}
