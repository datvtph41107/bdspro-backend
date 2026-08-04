package productcatalog

import (
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
	first, err := TermsChecksum(plan)
	if err != nil {
		t.Fatal(err)
	}

	plan.Entitlements[0], plan.Entitlements[1] =
		plan.Entitlements[1], plan.Entitlements[0]
	plan.Prices[0], plan.Prices[1] = plan.Prices[1], plan.Prices[0]
	second, err := TermsChecksum(plan)
	if err != nil {
		t.Fatal(err)
	}
	if first != second {
		t.Fatal("row ordering changed the commercial checksum")
	}
}
