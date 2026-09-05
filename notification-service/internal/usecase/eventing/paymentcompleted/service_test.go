package paymentcompleted

import (
	"context"
	"testing"
	"time"

	sharedevent "common/events/paymentcompleted"
	"notification/internal/domain/eventing"
)

type reactionStoreProbe struct {
	facts  []eventing.PaymentCompleted
	replay bool
	err    error
}

func (p *reactionStoreProbe) AcceptOnce(_ context.Context, fact eventing.PaymentCompleted) (bool, error) {
	p.facts = append(p.facts, fact)
	return p.replay, p.err
}

func validPaymentCompleted() sharedevent.V1 {
	now := time.Date(2026, 8, 28, 8, 0, 0, 0, time.UTC)
	return sharedevent.V1{
		EventID: "payment.completed.order.42", OrderID: 42,
		SubjectKind: "profile", SubjectID: "7", ProductCode: "qhpro", PlanCode: "pro",
		PlanVersionID: 3, PlanVersion: "3", Currency: "VND", AmountMinor: 199000,
		FundsConfirmedAt: now, CompletedAt: now.Add(time.Minute),
	}
}

func TestHandleNormalizesValidatedIntegrationFact(t *testing.T) {
	store := &reactionStoreProbe{}
	service := NewService(store)
	gotReplay, err := service.Handle(context.Background(), validPaymentCompleted())
	if err != nil || gotReplay {
		t.Fatalf("Handle() = replay=%v err=%v", gotReplay, err)
	}
	if len(store.facts) != 1 {
		t.Fatalf("store calls = %d, want 1", len(store.facts))
	}
	if store.facts[0].EventID != "payment.completed.order.42" || store.facts[0].AmountMinor != 199000 {
		t.Fatalf("normalized fact = %+v", store.facts[0])
	}
}

func TestHandleRejectsMalformedEventBeforeDurableStore(t *testing.T) {
	store := &reactionStoreProbe{}
	service := NewService(store)
	e := validPaymentCompleted()
	e.EventID = ""
	if _, err := service.Handle(context.Background(), e); err == nil {
		t.Fatal("Handle() error = nil, want invalid event")
	}
	if len(store.facts) != 0 {
		t.Fatal("malformed event reached durable store")
	}
}

func TestPushIntentOnlyForProfileSubject(t *testing.T) {
	fact := eventing.PaymentCompleted{SubjectKind: "profile"}
	if !ShouldCreatePushIntent(fact) {
		t.Fatal("profile payment completion should create push intent")
	}
	fact.SubjectKind = "organization"
	if ShouldCreatePushIntent(fact) {
		t.Fatal("organization payment completion should not assume profile push tokens")
	}
}
