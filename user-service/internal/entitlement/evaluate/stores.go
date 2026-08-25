package evaluate

import (
	"common/operation"
	"context"
	"time"
)

/**
 * SubscriptionStore đọc subscription đang dùng của một subject.
 */
type SubscriptionStore interface {
	FindSubscriptionForOperation(
		ctx context.Context,
		subjectType string,
		subjectID string,
		operationCode operation.Code,
		now time.Time,
	) (Subscription, bool, error)
}

/**
 * PlanAccessStore đọc quyền của operation trong một plan version.
 */
type PlanAccessStore interface {
	FindPlanAccess(
		ctx context.Context,
		planVersionID uint64,
		operationCode operation.Code,
	) (PlanAccess, bool, error)
}
