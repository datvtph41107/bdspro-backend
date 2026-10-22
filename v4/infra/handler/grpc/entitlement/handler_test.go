package accessgrpc

import (
	"common/operation"
	useraccess "mdn4/modules/user/domain/entitlement"
	"testing"
	"time"
)

func TestToProtoCarriesMeteringEvidence(t *testing.T) {
	start := time.Date(2026, 8, 1, 0, 0, 0, 0, time.UTC)
	end := start.AddDate(0, 1, 0)

	response := toProto(useraccess.Result{
		Subject:   useraccess.Subject{Type: useraccess.SubjectProfile, ID: "42"},
		Operation: operation.Code("workspace.report.generate"),
		Allowed:   true,
		Metering: useraccess.Metering{
			FeatureCode:    "workspace.report.generate",
			MeterCode:      "workspace.report_generation.accepted",
			UnitsPerAction: 1,
			PolicyVersion:  "1.0.0",
		},
		Limit:       100,
		Period:      useraccess.PeriodSubscriptionCycle,
		PeriodStart: start,
		PeriodEnd:   end,
	})

	if response.GetFeatureCode() != "workspace.report.generate" || response.GetUnitsPerAction() != 1 {
		t.Fatalf("response metering = %+v", response)
	}
}
