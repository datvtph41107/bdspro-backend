package plan

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"sort"
	"testing"
	"time"
)

func TestTermsChecksumIgnoresLifecycleState(t *testing.T) {
	start := time.Date(2026, 8, 1, 0, 0, 0, 0, time.UTC)
	plan := sampleCatalog(start).Products[0].Plans[0]
	first, err := TermsChecksum(plan)
	if err != nil {
		t.Fatal(err)
	}

	end := start.AddDate(0, 1, 0)
	plan.Status = StatusRetired
	plan.EffectiveFrom = start.AddDate(0, 0, 1)
	plan.EffectiveUntil = &end
	second, err := TermsChecksum(plan)
	if err != nil {
		t.Fatal(err)
	}
	if first != second {
		t.Fatal("lifecycle state changed the commercial checksum")
	}
}

func TestTermsChecksumDetectsCommercialChange(t *testing.T) {
	start := time.Date(2026, 8, 1, 0, 0, 0, 0, time.UTC)
	plan := sampleCatalog(start).Products[0].Plans[0]
	first, err := TermsChecksum(plan)
	if err != nil {
		t.Fatal(err)
	}

	plan.Entitlements[1].Amount++
	second, err := TermsChecksum(plan)
	if err != nil {
		t.Fatal(err)
	}
	if first == second {
		t.Fatal("commercial change did not change the checksum")
	}
}

func TestTermsChecksumDetectsOperationPolicyChange(t *testing.T) {
	start := time.Date(2026, 8, 1, 0, 0, 0, 0, time.UTC)
	plan := sampleCatalog(start).Products[0].Plans[0]
	first, err := TermsChecksum(plan)
	if err != nil {
		t.Fatal(err)
	}

	plan.Operations[0].UnitsPerAction = 2
	second, err := TermsChecksum(plan)
	if err != nil {
		t.Fatal(err)
	}
	if first == second {
		t.Fatal("operation policy change did not change the checksum")
	}
}

func TestTermsChecksumIsOrderIndependent(t *testing.T) {
	start := time.Date(2026, 8, 1, 0, 0, 0, 0, time.UTC)
	plan := sampleCatalog(start).Products[0].Plans[0]
	plan.Prices = append(plan.Prices, PriceItem{
		Code:        "workspace.report.overage",
		Kind:        PriceOverage,
		Currency:    "VND",
		AmountMinor: 100000,
		BillingUnit: "accepted.report",
		MeterCode:   "workspace.report_generation.accepted",
		Quantity:    1,
	})
	plan.Operations = append(plan.Operations, OperationBinding{
		Code:        "workspace.report.download",
		FeatureCode: "workspace.report.generate",
	})
	first, err := TermsChecksum(plan)
	if err != nil {
		t.Fatal(err)
	}

	plan.Entitlements[0], plan.Entitlements[1] =
		plan.Entitlements[1], plan.Entitlements[0]
	plan.Operations[0], plan.Operations[1] = plan.Operations[1], plan.Operations[0]
	plan.Prices[0], plan.Prices[1] = plan.Prices[1], plan.Prices[0]
	second, err := TermsChecksum(plan)
	if err != nil {
		t.Fatal(err)
	}
	if first != second {
		t.Fatal("row ordering changed the commercial checksum")
	}
}

func TestTermsChecksumDetectsSubscriptionTermChange(t *testing.T) {
	start := time.Date(2026, 8, 1, 0, 0, 0, 0, time.UTC)
	plan := sampleCatalog(start).Products[0].Plans[0]
	first, err := TermsChecksum(plan)
	if err != nil {
		t.Fatal(err)
	}

	plan.SubscriptionTermDays++
	second, err := TermsChecksum(plan)
	if err != nil {
		t.Fatal(err)
	}
	if first == second {
		t.Fatal("subscription term change did not change the commercial checksum")
	}
}

func TestTermsChecksumPreservesLegacyPayloadWhenSubscriptionTermIsMissing(t *testing.T) {
	start := time.Date(2026, 8, 1, 0, 0, 0, 0, time.UTC)
	plan := sampleCatalog(start).Products[0].Plans[0]
	plan.SubscriptionTermDays = 0

	got, err := TermsChecksum(plan)
	if err != nil {
		t.Fatalf("TermsChecksum legacy plan error = %v", err)
	}

	entitlements := append([]Entitlement(nil), plan.Entitlements...)
	operations := append([]OperationBinding(nil), plan.Operations...)
	prices := append([]PriceItem(nil), plan.Prices...)
	sort.Slice(entitlements, func(i, j int) bool { return entitlements[i].Code < entitlements[j].Code })
	sort.Slice(operations, func(i, j int) bool { return operations[i].Code < operations[j].Code })
	sort.Slice(prices, func(i, j int) bool { return prices[i].Code < prices[j].Code })

	legacyPayload, err := json.Marshal(struct {
		ProductCode  string
		PlanCode     string
		Version      string
		DisplayName  string
		SubjectScope SubjectScope
		Entitlements []Entitlement
		Operations   []OperationBinding
		Prices       []PriceItem
	}{
		ProductCode:  plan.ProductCode,
		PlanCode:     plan.PlanCode,
		Version:      plan.Version,
		DisplayName:  plan.DisplayName,
		SubjectScope: plan.SubjectScope,
		Entitlements: entitlements,
		Operations:   operations,
		Prices:       prices,
	})
	if err != nil {
		t.Fatal(err)
	}
	legacySum := sha256.Sum256(legacyPayload)
	want := hex.EncodeToString(legacySum[:])
	if got != want {
		t.Fatalf("legacy checksum changed: got %s want %s", got, want)
	}

	catalog := sampleCatalog(start)
	catalog.Products[0].Plans[0] = plan
	if _, err := NewRegistry(catalog); err == nil {
		t.Fatal("NewRegistry() accepted an active legacy plan without explicit subscription term")
	}
}
