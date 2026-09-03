package access

import (
	"common/operation"
	"testing"
	"time"
)

func TestResultUsesQuota(t *testing.T) {
	start := time.Date(2026, 8, 1, 0, 0, 0, 0, time.UTC)
	end := start.AddDate(0, 1, 0)

	result := Result{
		Subject:   Subject{Type: SubjectProfile, ID: "42"},
		Operation: operation.Code("workspace.report.generate"),
		Allowed:   true,
		Metering: Metering{
			FeatureCode:    "workspace.report.generate",
			MeterCode:      "workspace.report_generation.accepted",
			UnitsPerAction: 1,
			PolicyVersion:  "1.0.0",
		},
		Limit:       100,
		Period:      PeriodSubscriptionCycle,
		PeriodStart: start,
		PeriodEnd:   end,
	}

	if !result.IsValid() {
		t.Fatal("expected valid access result")
	}
	if !result.UsesQuota() {
		t.Fatal("expected quota to be required")
	}

	result.Unlimited = true
	result.Limit = 0
	if result.UsesQuota() {
		t.Fatal("unlimited access should not use quota")
	}
}
