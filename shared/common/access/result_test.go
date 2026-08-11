package access

import (
	"common/businessoperation"
	"testing"
	"time"
)

func TestResultUsesQuota(t *testing.T) {
	start := time.Date(2026, 8, 1, 0, 0, 0, 0, time.UTC)
	end := start.AddDate(0, 1, 0)

	result := Result{
		Subject:     Subject{Type: SubjectProfile, ID: "42"},
		Operation:   businessoperation.WorkspaceReportGenerate,
		Allowed:     true,
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
