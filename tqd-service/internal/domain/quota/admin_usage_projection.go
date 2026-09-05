package quota

import "time"

// AdminUsageProjection is TQD-owned evidence assembled from durable accounting.
// Runtime counters are attached by the usecase and never replace durable truth.
type AdminUsageProjection struct {
	SubjectType      string
	SubjectID        string
	Operation        string
	MeterCode        string
	DurableUsed      int64
	RuntimeUsed      int64
	RuntimeReserved  int64
	Remaining        int64
	LimitKnown       bool
	LimitSnapshot    int64
	PeriodStart      time.Time
	PeriodEnd        time.Time
	SubscriptionID   uint64
	PlanCode         string
	PlanVersion      string
	PolicyVersion    string
	RuntimeAvailable bool
	Reconciliation   string
	LastUsageAt      time.Time
}
