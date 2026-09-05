package subscription

import "time"

// AdminProjection is the User-owned operational view of one durable subscription.
type AdminProjection struct {
	SubscriptionID       uint64
	SubjectKind          string
	SubjectID            string
	ProductCode          string
	ProductDisplayName   string
	PlanCode             string
	PlanDisplayName      string
	PlanVersionID        uint64
	PlanVersion          string
	Status               string
	StartedAt            time.Time
	CurrentPeriodStart   time.Time
	CurrentPeriodEnd     time.Time
	AccessUntil          *time.Time
	PendingPlanVersionID *uint64
	PendingEffectiveAt   *time.Time
	OrderReference       string
	UpdatedAt            time.Time
}

type AdminUserProjection struct {
	ProfileID             uint64
	EffectiveSubscription *AdminProjection
}
