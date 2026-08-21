package migration

import (
	"math"
	"testing"

	"common/productcatalog"
)

func TestBuildDraftRequiresCompleteFieldDecisions(t *testing.T) {
	decision := sampleDecision()
	delete(decision.IgnoredFields, FieldOrgProductLimit)

	if _, err := BuildDraft(sampleLegacyPlan(), decision); err == nil {
		t.Fatal("BuildDraft() accepted an uncovered legacy field")
	}
}

func TestBuildDraftCreatesOnlyApprovedDraftTerms(t *testing.T) {
	draft, err := BuildDraft(sampleLegacyPlan(), sampleDecision())
	if err != nil {
		t.Fatalf("BuildDraft() error = %v", err)
	}

	plan := draft.Catalog.Products[0].Plans[0]
	if plan.Status != productcatalog.StatusDraft {
		t.Fatalf("status = %q, want draft", plan.Status)
	}
	if plan.EffectiveFrom.IsZero() == false {
		t.Fatal("legacy backfill draft unexpectedly has effective_from")
	}
	if len(plan.Entitlements) != 2 {
		t.Fatalf("entitlements = %d, want 2", len(plan.Entitlements))
	}
	if plan.Entitlements[1].Amount != 5 {
		t.Fatalf("post allowance = %d, want 5", plan.Entitlements[1].Amount)
	}
	if len(plan.Prices) != 1 || plan.Prices[0].AmountMinor != 9900000 {
		t.Fatalf("price = %+v", plan.Prices)
	}
	if len(draft.SourceChecksum) != 64 || len(draft.TermsChecksum) != 64 {
		t.Fatal("checksums must be sha256 hex")
	}
	if draft.SubscriptionTermDays != 365 {
		t.Fatalf("term days = %d", draft.SubscriptionTermDays)
	}
}

func TestBuildDraftRejectsNonBooleanFeatureFlag(t *testing.T) {
	snapshot := sampleLegacyPlan()
	snapshot.CreateProduct = 10

	if _, err := BuildDraft(snapshot, sampleDecision()); err == nil {
		t.Fatal("BuildDraft() accepted a non-boolean feature flag")
	}
}

func TestBuildDraftRejectsAmbiguousZeroAmount(t *testing.T) {
	decision := sampleDecision()
	decision.Entitlements[1].SourceField = FieldOrgMemberLimit
	decision.Entitlements[1].Code = "organization.member.capacity"
	decision.Entitlements[1].Kind = productcatalog.EntitlementCapacityLimit
	decision.Entitlements[1].MeterCode = "organization.member.active"
	decision.Entitlements[1].Period = productcatalog.PeriodNone
	decision.Entitlements[1].Zero = ZeroReject
	delete(decision.IgnoredFields, FieldOrgMemberLimit)
	decision.IgnoredFields[FieldLimitPost] = "post allowance reviewed separately"

	if _, err := BuildDraft(sampleLegacyPlan(), decision); err == nil {
		t.Fatal("BuildDraft() accepted ambiguous zero semantics")
	}
}

func TestBuildDraftAllowsExplicitUnlimitedZero(t *testing.T) {
	decision := sampleDecision()
	decision.Entitlements[1].SourceField = FieldOrgMemberLimit
	decision.Entitlements[1].Code = "organization.member.capacity"
	decision.Entitlements[1].Kind = productcatalog.EntitlementCapacityLimit
	decision.Entitlements[1].MeterCode = "organization.member.active"
	decision.Entitlements[1].Period = productcatalog.PeriodNone
	decision.Entitlements[1].Zero = ZeroUnlimited
	delete(decision.IgnoredFields, FieldOrgMemberLimit)
	decision.IgnoredFields[FieldLimitPost] = "post allowance reviewed separately"

	draft, err := BuildDraft(sampleLegacyPlan(), decision)
	if err != nil {
		t.Fatalf("BuildDraft() error = %v", err)
	}
	entitlement := draft.Catalog.Products[0].Plans[0].Entitlements[1]
	if !entitlement.Unlimited || entitlement.Amount != 0 {
		t.Fatalf("entitlement = %+v", entitlement)
	}
}

func TestBuildDraftRejectsPriceOverflow(t *testing.T) {
	snapshot := sampleLegacyPlan()
	snapshot.Price = math.MaxInt64
	decision := sampleDecision()
	decision.RecurringPrice.Multiplier = 100

	if _, err := BuildDraft(snapshot, decision); err == nil {
		t.Fatal("BuildDraft() accepted overflowing price conversion")
	}
}

func TestSourceChecksumChangesWithLegacySource(t *testing.T) {
	first := sampleLegacyPlan()
	second := sampleLegacyPlan()
	second.LimitPost++

	firstChecksum, err := SourceChecksum(first)
	if err != nil {
		t.Fatal(err)
	}
	secondChecksum, err := SourceChecksum(second)
	if err != nil {
		t.Fatal(err)
	}
	if firstChecksum == secondChecksum {
		t.Fatal("source checksum did not detect legacy row change")
	}
}

func sampleLegacyPlan() LegacyPlanSnapshot {
	return LegacyPlanSnapshot{
		ID:                 7,
		PackageName:        "Gói Pro",
		DurationDays:       365,
		Price:              9900000,
		CreateProduct:      1,
		LimitPost:          5,
		PublicProductLimit: 3,
		CreateProductAI:    0,
		FriendLimit:        50,
		CustomerLimit:      100,
		JoinOrg:            3,
		CreateOrganization: 0,
		OrgMemberLimit:     0,
		OrgProductLimit:    0,
		Color:              "#000000",
		Features:           []string{"Marketing copy"},
		Visibility:         true,
	}
}

func sampleDecision() Decision {
	return Decision{
		CatalogVersion:            "1.0.0",
		ProductCode:               "qhpro",
		ProductDisplayName:        "QHPRO",
		PlanCode:                  "qhpro.legacy_pro",
		PlanVersion:               "1.0.0",
		PlanDisplayName:           "Gói Pro",
		SubjectScope:              productcatalog.SubjectScopeProfile,
		DurationHandling:          "subscription_access_term",
		MarketingFeaturesHandling: "catalog_display_only",
		RecurringPrice: &PriceDecision{
			Code:        "qhpro.legacy_pro.subscription",
			Currency:    "VND",
			Multiplier:  1,
			BillingUnit: "subscription.term",
			Quantity:    1,
		},
		Entitlements: []EntitlementDecision{
			{
				SourceField: FieldCreateProduct,
				ValueMode:   ValueModeFeatureFlag,
				Code:        "product.create.feature",
				Kind:        productcatalog.EntitlementFeatureAccess,
				FeatureCode: "product.create",
				Period:      productcatalog.PeriodNone,
			},
			{
				SourceField: FieldLimitPost,
				ValueMode:   ValueModeAmount,
				Zero:        ZeroReject,
				Code:        "post.publish.allowance",
				Kind:        productcatalog.EntitlementUsageAllowance,
				MeterCode:   "post.publish.accepted",
				Period:      productcatalog.PeriodCalendarMonth,
			},
		},
		IgnoredFields: map[LegacyField]string{
			FieldPublicProductLimit: "business review pending",
			FieldCreateProductAI:    "business review pending",
			FieldFriendLimit:        "capacity semantics pending",
			FieldCustomerLimit:      "capacity semantics pending",
			FieldJoinOrg:            "capacity semantics pending",
			FieldCreateOrganization: "feature semantics pending",
			FieldOrgMemberLimit:     "organization scope pending",
			FieldOrgProductLimit:    "organization scope pending",
		},
	}
}
