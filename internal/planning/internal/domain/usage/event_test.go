package usage

import (
	"common/operation"
	"testing"
	"time"
	"tqd/internal/access"
)

func validEventInput() NewEventInput {
	start := time.Date(2026, 8, 1, 0, 0, 0, 0, time.UTC)
	return NewEventInput{
		Subject:        access.Subject{Type: access.SubjectProfile, ID: "42"},
		Operation:      operation.Code("workspace.report.generate"),
		MeterCode:      "workspace.report_generation.accepted",
		OperationID:    "op_usage_1",
		IdempotencyKey: "idem_usage_1",
		CommandKey:     "idem_usage_1",
		ReservationID:  "reservation_1",
		Amount:         1,
		PeriodStart:    start,
		PeriodEnd:      start.AddDate(0, 1, 0),
		CreatedAt:      start.Add(time.Hour),
		Commercial: CommercialEvidence{
			SubscriptionID: 77,
			PlanCode:       "qhpro.pro",
			PlanVersion:    "2.0.0",
			PolicyVersion:  "2.0.0",
			LimitSnapshot:  200,
		},
	}
}

func TestEventValidityDoesNotRequireReservation(t *testing.T) {
	input := validEventInput()
	input.ReservationID = ""

	event := NewEvent(input)
	if !event.IsValid() {
		t.Fatal("durable usage must be valid without Redis reservation evidence")
	}
}

func TestEventValidityAcceptsLegacyCommercialEvidenceBridge(t *testing.T) {
	input := validEventInput()
	input.Commercial = CommercialEvidence{}

	event := NewEvent(input)
	if !event.IsValid() {
		t.Fatal("legacy usage without catalog evidence must remain readable during migration")
	}
}

func TestEventValidityRejectsPartialCommercialEvidence(t *testing.T) {
	input := validEventInput()
	input.Commercial = CommercialEvidence{SubscriptionID: 77, PlanCode: "qhpro.pro"}

	if NewEvent(input).IsValid() {
		t.Fatal("partially populated commercial evidence must be rejected")
	}
}

func TestEventValidityRejectsTamperedUsageKey(t *testing.T) {
	event := NewEvent(validEventInput())
	event.UsageKey = "wrong"
	if event.IsValid() {
		t.Fatal("tampered usage key must be invalid")
	}
}
