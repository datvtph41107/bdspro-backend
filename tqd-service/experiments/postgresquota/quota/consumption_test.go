package quota

import (
	commonaccess "common/access"
	"testing"
	"time"
)

func validConsumption() Consumption {
	start := time.Date(2026, 8, 1, 0, 0, 0, 0, time.UTC)
	return Consumption{
		Subject:     commonaccess.Subject{Type: commonaccess.SubjectProfile, ID: "42"},
		Meter:       "workspace.report_generation.accepted",
		Amount:      1,
		Limit:       100,
		PeriodStart: start,
		PeriodEnd:   start.AddDate(0, 1, 0),
	}
}

func TestConsumptionValidation(t *testing.T) {
	if !validConsumption().IsValid() {
		t.Fatal("valid limited consumption rejected")
	}
	invalid := validConsumption()
	invalid.Amount = 101
	if invalid.IsValid() {
		t.Fatal("amount above limit must be invalid")
	}
}
