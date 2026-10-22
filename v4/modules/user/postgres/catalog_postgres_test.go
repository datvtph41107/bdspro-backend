package postgres

import (
	"testing"
	"time"

	catalogdomain "mdn4/modules/user/domain/plan"
	"mdn4/modules/user/entity"
)

func TestToPlanVersionAggregateMapsOptionalFields(t *testing.T) {
	feature := "workspace.report.generate"
	meter := "workspace.report_generation.accepted"
	from := time.Date(2026, 8, 5, 0, 0, 0, 0, time.UTC)
	published := from.Add(-time.Hour)
	tier := int32(20)
	term := int32(30)

	aggregate := toPlanVersionAggregate(
		entity.CatalogProduct{ID: 1, Code: "qhpro", DisplayName: "QHPro", Status: catalogdomain.StatusActive},
		entity.CatalogPlan{ID: 2, ProductID: 1, Code: "qhpro.pro", TierRank: &tier, Status: catalogdomain.StatusActive},
		entity.CatalogPlanVersion{
			ID:                   3,
			PlanID:               2,
			Version:              "1.0.0",
			DisplayName:          "Pro",
			Status:               catalogdomain.StatusActive,
			SubjectScope:         catalogdomain.SubjectScopeProfile,
			SubscriptionTermDays: &term,
			EffectiveFrom:        &from,
			PublishedAt:          &published,
			TermsChecksum:        "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
		},
		[]entity.CatalogPlanEntitlement{
			{
				Code:        "workspace.report.feature",
				Kind:        catalogdomain.EntitlementFeatureAccess,
				FeatureCode: &feature,
				Period:      catalogdomain.PeriodNone,
			},
		},
		[]entity.CatalogPlanOperationPolicy{
			{
				OperationCode:  "workspace.report.generate",
				FeatureCode:    feature,
				MeterCode:      &meter,
				UnitsPerAction: 1,
			},
		},
		[]entity.CatalogPriceItem{
			{
				Code:        "workspace.report.overage",
				Kind:        catalogdomain.PriceOverage,
				Currency:    "VND",
				AmountMinor: 1000,
				BillingUnit: "report.accepted",
				MeterCode:   &meter,
				Quantity:    1,
			},
		},
	)

	if aggregate.ProductCode != "qhpro" || aggregate.PlanCode != "qhpro.pro" ||
		aggregate.ProductStatus != catalogdomain.StatusActive ||
		aggregate.PlanStatus != catalogdomain.StatusActive ||
		aggregate.TierRank != 20 || aggregate.SubscriptionTermDays != 30 {
		t.Fatal("catalog identity was not mapped")
	}
	if aggregate.Entitlements[0].FeatureCode != feature {
		t.Fatal("feature code was not mapped")
	}
	if aggregate.Prices[0].MeterCode != meter {
		t.Fatal("price meter was not mapped")
	}
	if len(aggregate.Operations) != 1 ||
		aggregate.Operations[0].Code != "workspace.report.generate" ||
		aggregate.Operations[0].UnitsPerAction != 1 {
		t.Fatalf("operation policy was not mapped: %#v", aggregate.Operations)
	}
	*aggregate.EffectiveFrom = aggregate.EffectiveFrom.Add(time.Hour)
	if !from.Equal(time.Date(2026, 8, 5, 0, 0, 0, 0, time.UTC)) {
		t.Fatal("aggregate retained model time pointer")
	}
}
