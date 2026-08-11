package access

import (
	"common/businessoperation"
	"context"
	"time"
)

/**
 * SubscriptionStore đọc subscription đang dùng của một subject.
 */
type SubscriptionStore interface {
	FindActiveSubscription(
		ctx context.Context,
		subjectType string,
		subjectID string,
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
		operation businessoperation.Code,
	) (PlanAccess, bool, error)
}
