package wallethandler

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"testing"
	"time"

	domain "payment/internal/domain/payment"
	settlementuc "payment/internal/usecase/settlement"
	paymentpb "pb/types/payment"

	"google.golang.org/grpc/metadata"
)

func TestExtractCommercialOrderReferenceCanonicalizesCase(t *testing.T) {
	got, ok := ExtractCommercialOrderReference(
		"bank transfer qhp-d9ee180c81828aa0",
	)
	if !ok || got != "QHP-D9EE180C81828AA0" {
		t.Fatalf("reference = %q, found = %v", got, ok)
	}
}

func TestParseVNDMinor(t *testing.T) {
	t.Parallel()
	for _, test := range []struct {
		input string
		want  int64
		ok    bool
	}{
		{input: "99000", want: 99000, ok: true},
		{input: "50000.00", want: 50000, ok: true},
		{input: "50000.01", ok: false},
		{input: "0", ok: false},
		{input: "NaN", ok: false},
	} {
		got, err := parseVNDMinor(test.input)
		if test.ok && (err != nil || got != test.want) {
			t.Fatalf("parseVNDMinor(%q) = %d, %v", test.input, got, err)
		}
		if !test.ok && err == nil {
			t.Fatalf("parseVNDMinor(%q) unexpectedly succeeded", test.input)
		}
	}
}

func TestSepayAPIKeyCredentialUsesOfficialAuthorizationWireFormat(t *testing.T) {
	t.Parallel()
	secret := "sepay-webhook-secret"
	digest := sha256.Sum256([]byte("Apikey " + secret))
	ctx := metadata.NewIncomingContext(context.Background(), metadata.Pairs(
		"x-qhpro-provider-credential-sha256", hex.EncodeToString(digest[:]),
	))

	if !validSepayAPIKeyCredential(ctx, secret) {
		t.Fatal("official SePay API-key credential was rejected")
	}
	if validSepayAPIKeyCredential(ctx, "") {
		t.Fatal("empty SePay API key was accepted")
	}
	if validProviderCredential(ctx, secret) {
		t.Fatal("raw key without SePay's Apikey scheme was accepted")
	}
}

func TestSepayGatewayIPNAcceptsOnlyAuthenticatedCapturedPayment(t *testing.T) {
	t.Parallel()
	now := time.Date(2026, 9, 8, 8, 0, 0, 0, time.UTC)
	store := &gatewaySettlementStore{order: domain.Order{
		ID: 41, Reference: "QHP-0123456789ABCDEF",
		Status: domain.OrderPendingFunds, ExpiresAt: now.Add(time.Hour),
		Terms: domain.CommercialTerms{Price: domain.Money{Currency: domain.CurrencyVND, AmountMinor: 99000}},
	}}
	settlement := settlementuc.NewService(store, store, func() time.Time { return now })
	handler := NewPaymentHandler(nil, nil, settlement, "", "ipn-secret")
	digest := sha256.Sum256([]byte("ipn-secret"))
	ctx := metadata.NewIncomingContext(context.Background(), metadata.Pairs(
		"x-qhpro-provider-credential-sha256", hex.EncodeToString(digest[:]),
	))

	_, err := handler.SepayGatewayIPN(ctx, &paymentpb.SepayGatewayIPNRequest{
		NotificationType: "ORDER_PAID",
		Order: &paymentpb.SepayGatewayOrder{
			OrderStatus: "CAPTURED", OrderCurrency: "VND", OrderAmount: "99000.00",
			OrderInvoiceNumber: "QHP-0123456789ABCDEF",
		},
		Transaction: &paymentpb.SepayGatewayTransaction{
			Id: "transaction-1", TransactionId: "provider-transaction-1",
			TransactionType: "PAYMENT", TransactionStatus: "APPROVED",
			TransactionAmount: "99000", TransactionCurrency: "VND",
			TransactionDate: "2026-09-08 15:00:00",
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	if store.accepted.Settlement.OrderID != 41 || !store.accepted.CreateFulfillment {
		t.Fatalf("settlement decision = %#v", store.accepted)
	}
	if got := store.accepted.Settlement.ProviderTransactionID; got != "provider-transaction-1" {
		t.Fatalf("provider transaction = %q", got)
	}
}

type gatewaySettlementStore struct {
	order    domain.Order
	accepted settlementuc.Decision
}

func (s *gatewaySettlementStore) FindOrderByReference(_ context.Context, reference string) (domain.Order, bool, error) {
	return s.order, reference == s.order.Reference, nil
}

func (s *gatewaySettlementStore) FindSettlementByProviderTransaction(context.Context, string, string) (domain.Settlement, bool, error) {
	return domain.Settlement{}, false, nil
}

func (s *gatewaySettlementStore) AcceptSettlement(_ context.Context, decision settlementuc.Decision) (settlementuc.Result, error) {
	s.accepted = decision
	return settlementuc.Result{Settlement: decision.Settlement}, nil
}

func (s *gatewaySettlementStore) MarkSettlementRequiresReview(context.Context, uint64, string) error {
	return nil
}
