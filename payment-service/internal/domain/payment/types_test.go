package domain

import (
	"strings"
	"testing"
	"time"
)

func TestCommercialTermsRejectMalformedMoneyAndChecksum(t *testing.T) {
	terms := CommercialTerms{
		ProductCode: "qhpro", PlanCode: "qhpro.pro", PlanVersionID: 1, PlanVersion: "1.0.0",
		TierRank: 1, SubscriptionTermDays: 30,
		TermsChecksum: "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
		Price:         Money{Currency: CurrencyVND, AmountMinor: 100},
	}
	if !terms.IsValid() {
		t.Fatal("valid commercial terms rejected")
	}
	terms.Price.AmountMinor = 0
	if terms.IsValid() {
		t.Fatal("zero amount accepted")
	}
	terms.Price.AmountMinor = 100
	terms.TermsChecksum = "not-a-checksum"
	if terms.IsValid() {
		t.Fatal("malformed checksum accepted")
	}
}

func TestOrderFingerprintChangesForAnyImmutableCommercialInput(t *testing.T) {
	subject := Subject{Kind: SubjectProfile, ID: "42"}
	terms := CommercialTerms{
		ProductCode: "qhpro", PlanCode: "qhpro.pro", PlanVersionID: 1, PlanVersion: "1.0.0",
		TierRank: 1, SubscriptionTermDays: 30,
		TermsChecksum: "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
		Price:         Money{Currency: CurrencyVND, AmountMinor: 100},
	}
	first := OrderFingerprint(subject, terms)
	if len(first) != 64 {
		t.Fatalf("fingerprint length = %d", len(first))
	}
	terms.Price.AmountMinor++
	if second := OrderFingerprint(subject, terms); second == first {
		t.Fatal("money mutation did not change order fingerprint")
	}
	terms.Price.AmountMinor = 100
	if second := OrderFingerprint(Subject{Kind: SubjectOrganization, ID: "42"}, terms); second == first {
		t.Fatal("subject mutation did not change order fingerprint")
	}
}

func TestCrossOwnerIdentitiesIncludeImmutableOrderReference(t *testing.T) {
	first := Order{ID: 1, Reference: "QHP-AAAAAAAAAAAAAAAA"}
	second := Order{ID: 1, Reference: "QHP-BBBBBBBBBBBBBBBB"}

	if first.SettlementEffectKey() == second.SettlementEffectKey() {
		t.Fatal("settlement effect identity collides after Payment database recreation")
	}
	if first.CompletedEventID() == second.CompletedEventID() {
		t.Fatal("completed event identity collides after Payment database recreation")
	}
	if got := first.SettlementEffectKey(); got != "payment.order.1.QHP-AAAAAAAAAAAAAAAA" {
		t.Fatalf("settlement effect key = %q", got)
	}
	if got := first.CompletedEventID(); got != "payment.completed.order.1.QHP-AAAAAAAAAAAAAAAA" {
		t.Fatalf("completed event id = %q", got)
	}
	if got := (Order{ID: 1}).SettlementEffectKey(); got != "" {
		t.Fatalf("incomplete order produced effect key %q", got)
	}
}

func TestProviderEvidenceHashUsesExplicitEvidenceWhenPresent(t *testing.T) {
	event := ProviderEvent{
		Provider: "sepay", TransactionID: "tx-1", Reference: "QHP-1",
		Amount:     Money{Currency: CurrencyVND, AmountMinor: 100},
		OccurredAt: time.Date(2026, 8, 28, 0, 0, 0, 0, time.UTC),
		Type:       ProviderFundsConfirmed, EvidenceHash: strings.Repeat("A", 64),
	}
	if got := ProviderEvidenceHash(event); got != strings.Repeat("a", 64) {
		t.Fatalf("explicit evidence hash = %q", got)
	}
	event.EvidenceHash = "raw-evidence-from-provider"
	if got := ProviderEvidenceHash(event); len(got) != 64 || got == event.EvidenceHash {
		t.Fatalf("raw evidence was not hashed = %q", got)
	}
	event.EvidenceHash = ""
	if got := ProviderEvidenceHash(event); len(got) != 64 {
		t.Fatalf("derived evidence hash = %q", got)
	}
}
