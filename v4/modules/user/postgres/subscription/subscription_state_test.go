package subscriptionpostgres

import (
	"testing"
	"time"
)

func TestCurrentRowAggregateMapsCanonicalSubscriptionState(t *testing.T) {
	start := time.Date(2026, 8, 1, 0, 0, 0, 0, time.UTC)
	end := start.AddDate(0, 1, 0)
	aggregate := (currentRow{
		ID: 9, SubscriptionKey: "sub_9", SubjectKind: "profile", SubjectID: "42",
		ProductID: 1, ProductCode: "qhpro", PlanVersionID: 3, PlanCode: "qhpro.pro",
		PlanVersion: "1.0.0", PlanTermsChecksum: "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
		PlanStatus: "active", Status: "active", StartedAt: start, CurrentPeriodStart: start,
		CurrentPeriodEnd: end, AccessUntil: &end, OrderReference: "88", CreatedAt: start, UpdatedAt: start,
	}).aggregate()

	if aggregate.SubjectID != "42" || aggregate.PlanCode != "qhpro.pro" ||
		aggregate.PlanVersion != "1.0.0" || aggregate.OrderReference != "88" {
		t.Fatalf("aggregate = %+v", aggregate)
	}
	if aggregate.AccessUntil == nil || !aggregate.AccessUntil.Equal(end) {
		t.Fatalf("access until = %v", aggregate.AccessUntil)
	}
}
