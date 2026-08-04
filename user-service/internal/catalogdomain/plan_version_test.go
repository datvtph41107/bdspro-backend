package catalogdomain

import (
	"testing"
	"time"

	"common/productcatalog"
)

func TestContractReturnsDefensiveCopy(t *testing.T) {
	until := time.Date(2027, 1, 1, 0, 0, 0, 0, time.UTC)
	aggregate := PlanVersionAggregate{
		ProductCode:  "qhpro",
		PlanCode:     "qhpro.pro",
		Version:      "1.0.0",
		DisplayName:  "Pro",
		SubjectScope: productcatalog.SubjectScopeProfile,
		Entitlements: []productcatalog.Entitlement{{
			Code:        "workspace.report.feature",
			Kind:        productcatalog.EntitlementFeatureAccess,
			FeatureCode: "workspace.report.generate",
			Period:      productcatalog.PeriodNone,
		}},
	}

	contract := aggregate.Contract(
		productcatalog.StatusActive,
		time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
		&until,
	)
	contract.Entitlements[0].Code = "changed"
	*contract.EffectiveUntil = contract.EffectiveUntil.Add(time.Hour)

	if aggregate.Entitlements[0].Code != "workspace.report.feature" {
		t.Fatal("aggregate entitlement was mutated")
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
		Status:       productcatalog.StatusDraft,
		SubjectScope: productcatalog.SubjectScopeProfile,
	}

	contract := aggregate.CurrentContract()
	if !contract.EffectiveFrom.IsZero() {
		t.Fatal("draft effective_from must remain zero")
	}
}
