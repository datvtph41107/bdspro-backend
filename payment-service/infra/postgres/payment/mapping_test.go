package postgres

import (
	"testing"
	"time"

	domain "payment/internal/domain/payment"
)

func TestOrderMappingPreservesFrozenCommercialTerms(t *testing.T) {
	now := time.Date(2026, 8, 19, 9, 0, 0, 0, time.UTC)
	confirmed := now.Add(time.Minute)
	order := domain.Order{
		ID:      7,
		Subject: domain.Subject{Kind: domain.SubjectProfile, ID: "42"},
		Terms: domain.CommercialTerms{
			ProductCode: "qhpro", PlanCode: "qhpro.pro", PlanVersionID: 11, PlanVersion: "2.0.0",
			TierRank: 20, SubscriptionTermDays: 30,
			TermsChecksum: "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
			Price:         domain.Money{Currency: domain.CurrencyVND, AmountMinor: 19900000},
		},
		CommandKey: "checkout-1", CommandFingerprint: "bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb",
		Reference: "QHPORDER", Status: domain.OrderFundsConfirmed,
		ExpiresAt: now.Add(time.Hour), FundsConfirmedAt: &confirmed, CreatedAt: now,
	}
	got := orderFromModel(orderToModel(order))
	if got.ID != order.ID || got.Subject != order.Subject || got.Terms != order.Terms ||
		got.CommandKey != order.CommandKey || got.CommandFingerprint != order.CommandFingerprint ||
		got.Reference != order.Reference || got.Status != order.Status ||
		got.FundsConfirmedAt == nil || !got.FundsConfirmedAt.Equal(confirmed) {
		t.Fatalf("order mapping lost durable terms: got=%+v want=%+v", got, order)
	}
}

func TestSettlementAndFulfillmentMappingPreserveEvidenceAndFence(t *testing.T) {
	now := time.Date(2026, 8, 19, 9, 0, 0, 0, time.UTC)
	settlement := domain.Settlement{
		ID: 3, OrderID: 7, Provider: "sepay", ProviderTransactionID: "tx-1", Reference: "QHPORDER",
		Amount: domain.Money{Currency: domain.CurrencyVND, AmountMinor: 19900000}, OccurredAt: now,
		EvidenceHash: "cccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccc",
		Status:       domain.SettlementFundsConfirmed, CreatedAt: now.Add(time.Second),
	}
	mappedSettlement := settlementFromModel(settlementToModel(settlement))
	if mappedSettlement.OrderID != settlement.OrderID || mappedSettlement.Amount != settlement.Amount ||
		!mappedSettlement.OccurredAt.Equal(settlement.OccurredAt) || mappedSettlement.EvidenceHash != settlement.EvidenceHash {
		t.Fatalf("settlement mapping lost money evidence: got=%+v want=%+v", mappedSettlement, settlement)
	}

	lease := now.Add(time.Minute)
	completed := now.Add(2 * time.Minute)
	row := fulfillmentModel{
		ID: 5, OrderID: 7, Status: string(domain.FulfillmentRunning), AttemptCount: 2,
		AvailableAt: now, LockedBy: "worker-b", ClaimVersion: 9, LeaseUntil: &lease,
		LastError: "previous", CompletedAt: &completed, CreatedAt: now, UpdatedAt: now,
	}
	mappedFulfillment := fulfillmentFromModel(row)
	if mappedFulfillment.ClaimVersion != 9 || mappedFulfillment.LockedBy != "worker-b" ||
		!mappedFulfillment.LeaseUntil.Equal(lease) || mappedFulfillment.AttemptCount != 2 {
		t.Fatalf("fulfillment mapping lost fencing data: %+v", mappedFulfillment)
	}
}
