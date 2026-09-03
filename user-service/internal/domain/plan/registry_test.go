package plan

import (
	"testing"
	"time"
)

func TestNewBuildsImmutableCatalog(t *testing.T) {
	start := time.Date(2026, 8, 1, 0, 0, 0, 0, time.UTC)
	catalog := sampleCatalog(start)

	registry, err := NewRegistry(catalog)
	if err != nil {
		t.Fatalf("NewRegistry() error = %v", err)
	}

	catalog.Products[0].Plans[0].Entitlements[0].FeatureCode = "mutated.feature"
	catalog.Products[0].Plans[0].Operations[0].UnitsPerAction = 999
	catalog.Products[0].Plans[0].Prices[0].AmountMinor = 999

	plan, ok := registry.Plan("qhpro.pro", "1.0.0")
	if !ok {
		t.Fatal("Plan() missing approved version")
	}
	if plan.Entitlements[0].FeatureCode != "workspace.report.generate" {
		t.Fatalf("feature code mutated: %q", plan.Entitlements[0].FeatureCode)
	}
	if plan.Prices[0].AmountMinor != 9900000 {
		t.Fatalf("price mutated: %d", plan.Prices[0].AmountMinor)
	}
	if plan.Operations[0].UnitsPerAction != 1 {
		t.Fatalf("operation policy mutated: %d", plan.Operations[0].UnitsPerAction)
	}

	plan.Entitlements[0].FeatureCode = "caller.mutation"
	plan.Operations[0].UnitsPerAction = 99
	again, _ := registry.Plan("qhpro.pro", "1.0.0")
	if again.Entitlements[0].FeatureCode != "workspace.report.generate" {
		t.Fatal("Plan() did not return a defensive copy")
	}
	if again.Operations[0].UnitsPerAction != 1 {
		t.Fatal("Plan() leaked an operation-policy mutation")
	}
}

func TestEffectivePlanUsesExclusiveEnd(t *testing.T) {
	start := time.Date(2026, 8, 1, 0, 0, 0, 0, time.UTC)
	end := start.AddDate(0, 1, 0)
	catalog := sampleCatalog(start)
	catalog.Products[0].Plans[0].EffectiveUntil = &end

	registry, err := NewRegistry(catalog)
	if err != nil {
		t.Fatalf("NewRegistry() error = %v", err)
	}

	if _, ok := registry.EffectivePlan("qhpro.pro", start); !ok {
		t.Fatal("plan should be active at effective_from")
	}
	if _, ok := registry.EffectivePlan("qhpro.pro", end); ok {
		t.Fatal("plan should be inactive at exclusive effective_until")
	}
}

func TestNewRejectsOverlappingActiveVersions(t *testing.T) {
	start := time.Date(2026, 8, 1, 0, 0, 0, 0, time.UTC)
	catalog := sampleCatalog(start)
	catalog.Products[0].Plans = append(
		catalog.Products[0].Plans,
		PlanVersion{
			ProductCode:          "qhpro",
			PlanCode:             "qhpro.pro",
			Version:              "2.0.0",
			DisplayName:          "Pro V2",
			Status:               StatusActive,
			SubjectScope:         SubjectScopeProfile,
			SubscriptionTermDays: 30,
			EffectiveFrom:        start.AddDate(0, 0, 10),
		},
	)

	if _, err := NewRegistry(catalog); err == nil {
		t.Fatal("NewRegistry() accepted overlapping active versions")
	}
}

func TestNewRejectsDuplicateEntitlementTarget(t *testing.T) {
	start := time.Date(2026, 8, 1, 0, 0, 0, 0, time.UTC)
	catalog := sampleCatalog(start)
	catalog.Products[0].Plans[0].Entitlements = append(
		catalog.Products[0].Plans[0].Entitlements,
		Entitlement{
			Code:      "workspace.report.allowance.duplicate",
			Kind:      EntitlementUsageAllowance,
			MeterCode: "workspace.report_generation.accepted",
			Amount:    50,
			Period:    PeriodSubscriptionCycle,
		},
	)

	if _, err := NewRegistry(catalog); err == nil {
		t.Fatal("NewRegistry() accepted duplicate entitlement target")
	}
}

func TestNewRejectsFloatStyleAndTransportCodesByContract(t *testing.T) {
	start := time.Date(2026, 8, 1, 0, 0, 0, 0, time.UTC)
	catalog := sampleCatalog(start)
	catalog.Products[0].Plans[0].PlanCode = "FindParcelByLocation"

	if _, err := NewRegistry(catalog); err == nil {
		t.Fatal("NewRegistry() accepted a transport-style plan code")
	}
}

func TestNewValidatesEntitlementShapes(t *testing.T) {
	start := time.Date(2026, 8, 1, 0, 0, 0, 0, time.UTC)
	tests := []struct {
		name        string
		entitlement Entitlement
	}{
		{
			name: "feature with quota fields",
			entitlement: Entitlement{
				Code:        "workspace.report.feature",
				Kind:        EntitlementFeatureAccess,
				FeatureCode: "workspace.report.generate",
				Amount:      1,
			},
		},
		{
			name: "usage without period",
			entitlement: Entitlement{
				Code:      "workspace.report.allowance",
				Kind:      EntitlementUsageAllowance,
				MeterCode: "workspace.report_generation.accepted",
				Amount:    100,
				Period:    PeriodNone,
			},
		},
		{
			name: "capacity with reset period",
			entitlement: Entitlement{
				Code:      "organization.member.capacity",
				Kind:      EntitlementCapacityLimit,
				MeterCode: "organization.member.active",
				Amount:    10,
				Period:    PeriodCalendarMonth,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			catalog := sampleCatalog(start)
			catalog.Products[0].Plans[0].Entitlements = []Entitlement{test.entitlement}
			if _, err := NewRegistry(catalog); err == nil {
				t.Fatal("NewRegistry() accepted invalid entitlement")
			}
		})
	}
}

func TestNewValidatesPriceShapes(t *testing.T) {
	start := time.Date(2026, 8, 1, 0, 0, 0, 0, time.UTC)
	catalog := sampleCatalog(start)
	catalog.Products[0].Plans[0].Prices = []PriceItem{
		{
			Code:        "workspace.report.overage",
			Kind:        PriceOverage,
			Currency:    "VND",
			AmountMinor: 500000,
			BillingUnit: "accepted.report",
			Quantity:    1,
		},
	}

	if _, err := NewRegistry(catalog); err == nil {
		t.Fatal("NewRegistry() accepted overage price without meter")
	}
}

func TestNewValidatesOperationPolicyBindings(t *testing.T) {
	start := time.Date(2026, 8, 1, 0, 0, 0, 0, time.UTC)
	tests := []struct {
		name   string
		policy OperationBinding
	}{
		{
			name: "feature is not granted",
			policy: OperationBinding{
				Code:        "workspace.report.generate",
				FeatureCode: "workspace.report.export",
			},
		},
		{
			name: "meter has no usage allowance",
			policy: OperationBinding{
				Code:           "workspace.report.generate",
				FeatureCode:    "workspace.report.generate",
				MeterCode:      "workspace.report_generation.rendered",
				UnitsPerAction: 1,
			},
		},
		{
			name: "unmetered operation has units",
			policy: OperationBinding{
				Code:           "workspace.report.generate",
				FeatureCode:    "workspace.report.generate",
				UnitsPerAction: 1,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			catalog := sampleCatalog(start)
			catalog.Products[0].Plans[0].Operations = []OperationBinding{test.policy}
			if _, err := NewRegistry(catalog); err == nil {
				t.Fatal("NewRegistry() accepted an invalid operation policy")
			}
		})
	}
}

func sampleCatalog(start time.Time) Snapshot {
	return Snapshot{
		Version: "1.0.0",
		Products: []Product{
			{
				Code:        "qhpro",
				DisplayName: "QHPRO",
				Plans: []PlanVersion{
					{
						ProductCode:          "qhpro",
						PlanCode:             "qhpro.pro",
						Version:              "1.0.0",
						DisplayName:          "QHPRO Pro",
						Status:               StatusActive,
						SubjectScope:         SubjectScopeProfile,
						SubscriptionTermDays: 30,
						EffectiveFrom:        start,
						Entitlements: []Entitlement{
							{
								Code:        "workspace.report.feature",
								Kind:        EntitlementFeatureAccess,
								FeatureCode: "workspace.report.generate",
								Period:      PeriodNone,
							},
							{
								Code:      "workspace.report.allowance",
								Kind:      EntitlementUsageAllowance,
								MeterCode: "workspace.report_generation.accepted",
								Amount:    100,
								Period:    PeriodSubscriptionCycle,
							},
						},
						Operations: []OperationBinding{
							{
								Code:           "workspace.report.generate",
								FeatureCode:    "workspace.report.generate",
								MeterCode:      "workspace.report_generation.accepted",
								UnitsPerAction: 1,
							},
						},
						Prices: []PriceItem{
							{
								Code:        "qhpro.pro.subscription",
								Kind:        PriceRecurring,
								Currency:    "VND",
								AmountMinor: 9900000,
								BillingUnit: "subscription.cycle",
								Quantity:    1,
							},
						},
					},
				},
			},
		},
	}
}

func TestNewRejectsActivePlanWithoutSubscriptionTerm(t *testing.T) {
	start := time.Date(2026, 8, 1, 0, 0, 0, 0, time.UTC)
	catalog := sampleCatalog(start)
	catalog.Products[0].Plans[0].SubscriptionTermDays = 0
	if _, err := NewRegistry(catalog); err == nil {
		t.Fatal("NewRegistry() accepted active plan without subscription_term_days")
	}
}
