package usage

import (
	commonaccess "common/access"
	"common/businessoperation"
	"testing"
	"time"
)

func TestEventValidityRequiresReservationAndPeriodEvidence(t *testing.T) {
	start := time.Date(2026, 8, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2026, 9, 1, 0, 0, 0, 0, time.UTC)
	event := NewEvent(
		commonaccess.Subject{Type: commonaccess.SubjectProfile, ID: "42"},
		businessoperation.WorkspaceReportGenerate,
		"op_usage_1",
		"idem_usage_1",
		"idem_usage_1",
		"reservation_1",
		1,
		start,
		end,
		start.Add(time.Hour),
	)
	if !event.IsValid() {
		t.Fatal("complete durable usage evidence must be valid")
	}

	event.ReservationID = ""
	if event.IsValid() {
		t.Fatal("usage without reservation evidence must be invalid")
	}
}

func TestEventValidityRejectsTamperedUsageKey(t *testing.T) {
	start := time.Date(2026, 8, 1, 0, 0, 0, 0, time.UTC)
	event := NewEvent(
		commonaccess.Subject{Type: commonaccess.SubjectProfile, ID: "42"},
		businessoperation.WorkspaceReportGenerate,
		"op_usage_2",
		"",
		"op_usage_2",
		"reservation_2",
		1,
		start,
		start.AddDate(0, 1, 0),
		start.Add(time.Hour),
	)
	event.UsageKey = "wrong"
	if event.IsValid() {
		t.Fatal("tampered usage key must be invalid")
	}
}
