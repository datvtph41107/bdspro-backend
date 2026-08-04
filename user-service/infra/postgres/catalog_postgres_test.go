package postgres

import (
	"testing"
	"time"

	"common/productcatalog"
	"user/internal/models"
)

func TestToPlanVersionAggregateMapsOptionalFields(t *testing.T) {
	feature := "workspace.report.generate"
	meter := "workspace.report_generation.accepted"
	from := time.Date(2026, 8, 5, 0, 0, 0, 0, time.UTC)
	published := from.Add(-time.Hour)

	aggregate := toPlanVersionAggregate(
		models.CatalogProduct{ID: 1, Code: "qhpro", DisplayName: "QHPro", Status: productcatalog.StatusActive},
		models.CatalogPlan{ID: 2, ProductID: 1, Code: "qhpro.pro", Status: productcatalog.StatusActive},
		models.CatalogPlanVersion{
			ID:            3,
			PlanID:        2,
			Version:       "1.0.0",
			DisplayName:   "Pro",
			Status:        productcatalog.StatusActive,
			SubjectScope:  productcatalog.SubjectScopeProfile,
			EffectiveFrom: &from,
			PublishedAt:   &published,
			TermsChecksum: "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
		},
		[]models.CatalogPlanEntitlement{
			{
				Code:        "workspace.report.feature",
				Kind:        productcatalog.EntitlementFeatureAccess,
				FeatureCode: &feature,
				Period:      productcatalog.PeriodNone,
			},
		},
		[]models.CatalogPriceItem{
			{
				Code:        "workspace.report.overage",
				Kind:        productcatalog.PriceOverage,
				Currency:    "VND",
				AmountMinor: 1000,
				BillingUnit: "report.accepted",
				MeterCode:   &meter,
				Quantity:    1,
			},
		},
	)

	if aggregate.ProductCode != "qhpro" || aggregate.PlanCode != "qhpro.pro" ||
		aggregate.ProductStatus != productcatalog.StatusActive ||
		aggregate.PlanStatus != productcatalog.StatusActive {
		t.Fatal("catalog identity was not mapped")
	}
	if aggregate.Entitlements[0].FeatureCode != feature {
		t.Fatal("feature code was not mapped")
	}
	if aggregate.Prices[0].MeterCode != meter {
		t.Fatal("price meter was not mapped")
	}
	*aggregate.EffectiveFrom = aggregate.EffectiveFrom.Add(time.Hour)
	if !from.Equal(time.Date(2026, 8, 5, 0, 0, 0, 0, time.UTC)) {
		t.Fatal("aggregate retained model time pointer")
	}
}

func TestToSubscriptionAggregateMapsPlanAndPeriod(t *testing.T) {
	start := time.Date(2026, 8, 1, 0, 0, 0, 0, time.UTC)
	end := start.AddDate(0, 1, 0)
	aggregate := toSubscriptionAggregate(
		models.CatalogSubscription{
			ID:                 9,
			SubscriptionKey:    "sub_legacy_membership_7",
			SubjectKind:        models.SubscriptionSubjectProfile,
			SubjectID:          "42",
			ProductID:          1,
			PlanVersionID:      3,
			Status:             models.SubscriptionActive,
			StartedAt:          start,
			CurrentPeriodStart: start,
			CurrentPeriodEnd:   end,
			AccessUntil:        &end,
			OrderReference:     "88",
		},
		models.CatalogProduct{ID: 1, Code: "qhpro"},
		models.CatalogPlan{ID: 2, Code: "qhpro.pro"},
		models.CatalogPlanVersion{
			ID:            3,
			Version:       "1.0.0",
			Status:        productcatalog.StatusActive,
			TermsChecksum: "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
		},
	)
	if aggregate.SubjectID != "42" || aggregate.PlanCode != "qhpro.pro" ||
		aggregate.PlanVersion != "1.0.0" || aggregate.OrderReference != "88" {
		t.Fatalf("aggregate = %+v", aggregate)
	}
	if aggregate.AccessUntil == nil || !aggregate.AccessUntil.Equal(end) {
		t.Fatalf("access until = %v", aggregate.AccessUntil)
	}
}
