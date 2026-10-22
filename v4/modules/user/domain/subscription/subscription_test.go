package subscription

import (
	"testing"
	"time"

	catalogdomain "mdn4/modules/user/domain/plan"
	"mdn4/modules/user/entity"
)

func validAggregate() Aggregate {
	start := time.Date(2026, 8, 1, 0, 0, 0, 0, time.UTC)
	end := start.AddDate(0, 1, 0)
	return Aggregate{
		ID:                 1,
		SubscriptionKey:    "sub_1",
		SubjectKind:        entity.SubscriptionSubjectProfile,
		SubjectID:          "42",
		ProductID:          2,
		ProductCode:        "qhpro",
		PlanVersionID:      3,
		PlanCode:           "qhpro.pro",
		PlanVersion:        "1.0.0",
		PlanTermsChecksum:  "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
		PlanStatus:         catalogdomain.StatusActive,
		Status:             entity.SubscriptionActive,
		StartedAt:          start,
		CurrentPeriodStart: start,
		CurrentPeriodEnd:   end,
		AccessUntil:        &end,
		CreatedAt:          start.Add(-time.Hour),
		UpdatedAt:          start,
	}
}

func TestAggregateValidate(t *testing.T) {
	aggregate := validAggregate()
	if err := aggregate.Validate(); err != nil {
		t.Fatalf("Validate() error = %v", err)
	}

	aggregate.CurrentPeriodEnd = aggregate.CurrentPeriodStart
	if err := aggregate.Validate(); err == nil {
		t.Fatal("Validate() accepted empty period")
	}
}

func TestAggregateCopyDoesNotShareTimePointers(t *testing.T) {
	aggregate := validAggregate()
	copied := aggregate.Copy()
	copied.AccessUntil = ptrTime(copied.AccessUntil.Add(time.Hour))
	if copied.AccessUntil.Equal(*aggregate.AccessUntil) {
		t.Fatal("copy changed original access window")
	}
}

func ptrTime(value time.Time) *time.Time {
	return &value
}
