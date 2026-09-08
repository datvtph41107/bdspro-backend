package settlement

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"sync"
	"testing"
	"time"

	catalogdomain "user/internal/domain/plan"
	subscriptiondomain "user/internal/domain/subscription"
	"user/internal/models"
)

func testEffect() Effect {
	return Effect{
		EffectKey: "payment.order.91", OrderID: 91,
		SubjectKind: models.SubscriptionSubjectProfile, SubjectID: "42",
		ProductCode: "qhpro", PlanCode: "qhpro.pro", PlanVersionID: 7,
		PlanVersion: "2.0.0", TierRank: 20, SubscriptionTermDays: 30,
		TermsChecksum: "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
		OccurredAt:    time.Date(2026, 8, 19, 10, 0, 0, 0, time.UTC),
	}
}

func testTarget(e Effect) Target {
	return Target{ProductID: 1, PlanID: 2, PlanVersionID: e.PlanVersionID, ProductCode: e.ProductCode,
		PlanCode: e.PlanCode, PlanVersion: e.PlanVersion, TierRank: e.TierRank,
		SubscriptionTermDays: e.SubscriptionTermDays, TermsChecksum: e.TermsChecksum, Published: true, PlanStatus: catalogdomain.StatusActive}
}

func activeCurrent(tier int32) Current {
	start := time.Date(2026, 8, 1, 0, 0, 0, 0, time.UTC)
	end := start.AddDate(0, 1, 0)
	return Current{Found: true, TierRank: tier, Aggregate: subscriptiondomain.Aggregate{
		ID: 5, SubscriptionKey: "sub_5", SubjectKind: models.SubscriptionSubjectProfile, SubjectID: "42",
		ProductID: 1, ProductCode: "qhpro", PlanVersionID: 6, PlanCode: "qhpro.basic", PlanVersion: "1.0.0",
		PlanTermsChecksum: "bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb",
		Status:            models.SubscriptionActive, StartedAt: start, CurrentPeriodStart: start, CurrentPeriodEnd: end,
		AccessUntil: &end, CreatedAt: start, UpdatedAt: start,
	}}
}

func TestDecideActivationUsesOccurredAtAndTerm(t *testing.T) {
	e := testEffect()
	now := e.OccurredAt.Add(time.Second)
	d, err := Decide(e, testTarget(e), Current{}, now)
	if err != nil {
		t.Fatal(err)
	}
	if d.Action != ActionActivated || !d.Subscription.StartedAt.Equal(e.OccurredAt) ||
		!d.Subscription.CurrentPeriodEnd.Equal(e.OccurredAt.AddDate(0, 0, 30)) {
		t.Fatalf("decision = %+v", d)
	}
	sum := sha256.Sum256([]byte(e.EffectKey))
	wantKey := "sub_payment_" + hex.EncodeToString(sum[:])
	if d.Subscription.SubscriptionKey != wantKey || d.Subscription.OrderReference != e.EffectKey {
		t.Fatalf("subscription identity = %q, order reference = %q", d.Subscription.SubscriptionKey, d.Subscription.OrderReference)
	}
}

func TestSubscriptionKeyDoesNotCollideWhenPaymentDatabaseReusesOrderID(t *testing.T) {
	first := testEffect()
	second := first
	first.EffectKey = "payment.order.91.QHP-FIRST"
	second.EffectKey = "payment.order.91.QHP-SECOND"

	if first.SubscriptionKey() == second.SubscriptionKey() {
		t.Fatalf("different durable Payment identities produced the same subscription key: %q", first.SubscriptionKey())
	}
	if len(first.SubscriptionKey()) > 128 {
		t.Fatalf("subscription key exceeds catalog_subscriptions.subscription_key: %d", len(first.SubscriptionKey()))
	}
}

func TestDecideUpgradePreservesCurrentPeriod(t *testing.T) {
	e := testEffect()
	current := activeCurrent(10)
	d, err := Decide(e, testTarget(e), current, e.OccurredAt.Add(time.Second))
	if err != nil {
		t.Fatal(err)
	}
	if d.Action != ActionUpgraded || d.Subscription.PlanVersionID != e.PlanVersionID ||
		!d.Subscription.CurrentPeriodStart.Equal(current.Aggregate.CurrentPeriodStart) ||
		!d.Subscription.CurrentPeriodEnd.Equal(current.Aggregate.CurrentPeriodEnd) {
		t.Fatalf("decision = %+v", d)
	}
}

func TestDecideSameAndDowngradeAreExplicitTerminalOutcomes(t *testing.T) {
	e := testEffect()
	for name, tc := range map[string]struct {
		tier   int32
		action Action
	}{
		"same": {20, ActionRejectedSameTier}, "down": {30, ActionRejectedDowngrade},
	} {
		t.Run(name, func(t *testing.T) {
			current := activeCurrent(tc.tier)
			d, err := Decide(e, testTarget(e), current, e.OccurredAt.Add(time.Second))
			if err != nil || d.Action != tc.action {
				t.Fatalf("decision=%+v err=%v", d, err)
			}
		})
	}
}

func TestFingerprintChangesWithOccurredAt(t *testing.T) {
	e := testEffect()
	a, _ := e.Fingerprint()
	e.OccurredAt = e.OccurredAt.Add(time.Nanosecond)
	b, _ := e.Fingerprint()
	if a == b {
		t.Fatal("fingerprint ignored OccurredAt")
	}
}

type memoryStore struct {
	mu      sync.Mutex
	records map[string]struct {
		fingerprint string
		result      Result
	}
	calls     int
	loseFirst bool
}

func (m *memoryStore) ApplySettlement(_ context.Context, input Acceptance) (Result, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.records == nil {
		m.records = map[string]struct {
			fingerprint string
			result      Result
		}{}
	}
	if old, ok := m.records[input.Effect.EffectKey]; ok {
		if old.fingerprint != input.Fingerprint {
			return Result{}, ErrEffectConflict
		}
		res := old.result
		res.Replayed = true
		return res, nil
	}
	m.calls++
	res := Result{Action: ActionActivated, SubscriptionID: 11}
	m.records[input.Effect.EffectKey] = struct {
		fingerprint string
		result      Result
	}{input.Fingerprint, res}
	if m.loseFirst {
		m.loseFirst = false
		return Result{}, errors.New("response lost")
	}
	return res, nil
}

func TestServiceResponseLossRetryDoesNotDuplicateEffect(t *testing.T) {
	store := &memoryStore{loseFirst: true}
	svc := NewService(store)
	e := testEffect()
	if _, err := svc.ApplySettlement(context.Background(), e); err == nil {
		t.Fatal("first call should lose response")
	}
	res, err := svc.ApplySettlement(context.Background(), e)
	if err != nil {
		t.Fatal(err)
	}
	if !res.Replayed || store.calls != 1 {
		t.Fatalf("result=%+v calls=%d", res, store.calls)
	}
}

func TestServiceSameEffectKeyChangedPayloadConflicts(t *testing.T) {
	store := &memoryStore{}
	svc := NewService(store)
	e := testEffect()
	if _, err := svc.ApplySettlement(context.Background(), e); err != nil {
		t.Fatal(err)
	}
	e.OccurredAt = e.OccurredAt.Add(time.Second)
	if _, err := svc.ApplySettlement(context.Background(), e); !errors.Is(err, ErrEffectConflict) {
		t.Fatalf("err=%v", err)
	}
}
