package plan

import (
	"testing"
	"time"
)

func TestContractReturnsDefensiveCopy(t *testing.T) {
	until := time.Date(2027, 1, 1, 0, 0, 0, 0, time.UTC)
	aggregate := PlanVersionAggregate{
		ProductCode:  "qhpro",
		PlanCode:     "qhpro.pro",
		Version:      "1.0.0",
		DisplayName:  "Pro",
		SubjectScope: SubjectScopeProfile,
		Entitlements: []Entitlement{{
			Code:        "workspace.report.feature",
			Kind:        EntitlementFeatureAccess,
			FeatureCode: "workspace.report.generate",
			Period:      PeriodNone,
		}},
		Operations: []OperationBinding{{
			Code:           "workspace.report.generate",
			FeatureCode:    "workspace.report.generate",
			MeterCode:      "workspace.report_generation.accepted",
			UnitsPerAction: 1,
		}},
	}

	contract := aggregate.Contract(
		StatusActive,
		time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
		&until,
	)
	contract.Entitlements[0].Code = "changed"
	contract.Operations[0].UnitsPerAction = 99
	*contract.EffectiveUntil = contract.EffectiveUntil.Add(time.Hour)

	if aggregate.Entitlements[0].Code != "workspace.report.feature" {
		t.Fatal("aggregate entitlement was mutated")
	}
	if aggregate.Operations[0].UnitsPerAction != 1 {
		t.Fatal("aggregate operation policy was mutated")
	}
	if !until.Equal(time.Date(2027, 1, 1, 0, 0, 0, 0, time.UTC)) {
		t.Fatal("effective_until input was mutated")
	}
}

func TestCurrentContractSupportsDraft(t *testing.T) {
	aggregate := PlanVersionAggregate{
		ProductCode:  "qhpro",
		PlanCode:     "qhpro.free",
		Version:      "1.0.0",
		DisplayName:  "Free",
		Status:       StatusDraft,
		SubjectScope: SubjectScopeProfile,
	}

	contract := aggregate.CurrentContract()
	if !contract.EffectiveFrom.IsZero() {
		t.Fatal("draft effective_from must remain zero")
	}
}
