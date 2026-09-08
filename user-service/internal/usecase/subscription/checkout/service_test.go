package checkout

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	plan "user/internal/domain/plan"
)

type noopAttemptPort struct{}

func (noopAttemptPort) CreateAttempt(context.Context, CreatePaymentAttemptCommand) (PaymentAttempt, error) {
	return PaymentAttempt{}, nil
}

type noopContactStore struct{}

func (noopContactStore) UpdateCheckoutContact(_ context.Context, _ uint64, contact Contact) (Contact, error) {
	contact.Phone = "0900000000"
	return contact, nil
}

type fakePlanResolver struct {
	terms PlanTerms
	err   error
	calls int
}

func (f *fakePlanResolver) ResolveCheckoutTerms(context.Context, string, time.Time) (PlanTerms, error) {
	f.calls++
	return f.terms, f.err
}

type fakeSubscriptions struct {
	current CurrentSubscription
	found   bool
	err     error
	calls   int
}

func (f *fakeSubscriptions) CurrentForProduct(context.Context, Subject, uint64) (CurrentSubscription, bool, error) {
	f.calls++
	return f.current, f.found, f.err
}

type fakeCommands struct {
	byKey       map[string]CheckoutCommand
	findErr     error
	createErr   error
	findCalls   int
	createCalls int
}

func newFakeCommands() *fakeCommands {
	return &fakeCommands{byKey: make(map[string]CheckoutCommand)}
}

func checkoutCommandKey(subject Subject, commandKey string) string {
	return fmt.Sprintf("%s:%s:%s", subject.Kind, subject.ID, commandKey)
}

func (f *fakeCommands) FindByCommand(_ context.Context, subject Subject, commandKey string) (CheckoutCommand, bool, error) {
	f.findCalls++
	if f.findErr != nil {
		return CheckoutCommand{}, false, f.findErr
	}
	command, found := f.byKey[checkoutCommandKey(subject, commandKey)]
	return command, found, nil
}

func (f *fakeCommands) CreateCommand(_ context.Context, command CheckoutCommand) (CheckoutCommand, bool, error) {
	f.createCalls++
	if f.createErr != nil {
		return CheckoutCommand{}, false, f.createErr
	}
	key := checkoutCommandKey(command.Subject, command.CommandKey)
	if existing, found := f.byKey[key]; found {
		return existing, false, nil
	}
	f.byKey[key] = command
	return command, true, nil
}

type fakePayment struct {
	order   PaymentOrder
	created bool
	err     error
	calls   int
	last    CreatePaymentOrderCommand
}

func (f *fakePayment) CreateOrder(_ context.Context, command CreatePaymentOrderCommand) (PaymentOrder, bool, error) {
	f.calls++
	f.last = command
	return f.order, f.created, f.err
}

func validTerms() PlanTerms {
	return PlanTerms{
		ProductID: 1, ProductCode: "qhpro", PlanID: 2, PlanCode: "qhpro.pro",
		PlanVersionID: 3, PlanVersion: "2.0.0", TierRank: 20, SubscriptionTermDays: 30,
		TermsChecksum: "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
		SubjectScope:  plan.SubjectScopeProfile,
		Price:         Money{Currency: "VND", AmountMinor: 19900000},
	}
}

func newCheckoutFixture() (*Service, *fakePlanResolver, *fakeSubscriptions, *fakeCommands, *fakePayment) {
	plans := &fakePlanResolver{terms: validTerms()}
	subs := &fakeSubscriptions{}
	commands := newFakeCommands()
	payment := &fakePayment{order: PaymentOrder{ID: 9, Reference: "QHP9", Status: "pending_funds"}, created: true}
	now := func() time.Time { return time.Date(2026, 8, 19, 10, 30, 0, 0, time.UTC) }
	return NewService(plans, subs, commands, payment, noopAttemptPort{}, noopContactStore{}, now), plans, subs, commands, payment
}

func TestUpdateCheckoutContactNormalizesAndPersists(t *testing.T) {
	service, _, _, _, _ := newCheckoutFixture()
	contact, err := service.UpdateContact(context.Background(), 42, Contact{FullName: " Nguyễn Văn A ", Email: " USER@Example.com "})
	if err != nil {
		t.Fatalf("UpdateContact() error = %v", err)
	}
	if contact.FullName != "Nguyễn Văn A" || contact.Email != "user@example.com" || contact.Phone != "0900000000" {
		t.Fatalf("UpdateContact() = %+v", contact)
	}
}

func TestCheckoutInitialPurchaseFreezesServerTermsBeforePayment(t *testing.T) {
	service, _, _, commands, payment := newCheckoutFixture()
	subject := Subject{Kind: SubjectProfile, ID: "42"}

	result, err := service.Checkout(context.Background(), subject, "qhpro.pro", "checkout-1")
	if err != nil {
		t.Fatal(err)
	}
	if commands.createCalls != 1 || payment.calls != 1 || !result.Created || result.Order.ID != 9 {
		t.Fatalf("calls/result = command:%d payment:%d result:%+v", commands.createCalls, payment.calls, result)
	}
	stored, found := commands.byKey[checkoutCommandKey(subject, "checkout-1")]
	if !found || stored.Terms != validTerms() {
		t.Fatalf("durable checkout command = %+v found=%v", stored, found)
	}
	if stored.RequestHash != RequestHashFor(subject, "qhpro.pro") {
		t.Fatalf("request hash = %q, want canonical hash", stored.RequestHash)
	}
	if payment.last.CommandKey != "checkout-1" || payment.last.Subject != subject || payment.last.Terms != stored.Terms {
		t.Fatalf("payment command did not use frozen server terms: %+v", payment.last)
	}
}

func TestCheckoutRequestHashRejectsTamperedDurableCommand(t *testing.T) {
	service, _, _, commands, payment := newCheckoutFixture()
	subject := Subject{Kind: SubjectProfile, ID: "42"}
	commands.byKey[checkoutCommandKey(subject, "tampered")] = CheckoutCommand{
		Subject:     subject,
		PlanCode:    "qhpro.pro",
		RequestHash: "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
		Terms:       validTerms(),
		CommandKey:  "tampered",
		CreatedAt:   time.Date(2026, 8, 19, 10, 30, 0, 0, time.UTC),
	}
	_, err := service.Checkout(context.Background(), subject, "qhpro.pro", "tampered")
	if !errors.Is(err, ErrCheckoutCommandConflict) {
		t.Fatalf("tampered command error = %v, want %v", err, ErrCheckoutCommandConflict)
	}
	if payment.calls != 0 {
		t.Fatalf("tampered durable command reached payment: %d calls", payment.calls)
	}
}

func TestCheckoutReplayUsesFrozenTermsBeforeCurrentState(t *testing.T) {
	service, plans, subs, commands, payment := newCheckoutFixture()
	subject := Subject{Kind: SubjectProfile, ID: "42"}

	first, err := service.Checkout(context.Background(), subject, "qhpro.pro", "checkout-replay")
	if err != nil {
		t.Fatal(err)
	}
	frozen := first.Terms

	// The world changes after the first durable Checkout command: Plan could
	// advance to a successor and fulfillment could already activate the target.
	plans.err = errors.New("catalog must not be read on durable replay")
	subs.found = true
	subs.current = CurrentSubscription{ID: 7, PlanVersionID: 3, PlanCode: "qhpro.pro", TierRank: 20, Status: "active"}
	payment.created = false

	replay, err := service.Checkout(context.Background(), subject, "qhpro.pro", "checkout-replay")
	if err != nil {
		t.Fatalf("replay error = %v", err)
	}
	if plans.calls != 1 || subs.calls != 1 || commands.createCalls != 1 {
		t.Fatalf("replay re-read mutable state: plans=%d subs=%d creates=%d", plans.calls, subs.calls, commands.createCalls)
	}
	if payment.calls != 2 || replay.Created || replay.Terms != frozen || payment.last.Terms != frozen {
		t.Fatalf("replay did not preserve frozen command: payment=%d result=%+v last=%+v", payment.calls, replay, payment.last)
	}
}

func TestCheckoutSameCommandDifferentPlanConflictsBeforeMutableReads(t *testing.T) {
	service, plans, subs, _, payment := newCheckoutFixture()
	subject := Subject{Kind: SubjectProfile, ID: "42"}
	if _, err := service.Checkout(context.Background(), subject, "qhpro.pro", "checkout-conflict"); err != nil {
		t.Fatal(err)
	}
	planCalls, subscriptionCalls, paymentCalls := plans.calls, subs.calls, payment.calls

	_, err := service.Checkout(context.Background(), subject, "qhpro.enterprise", "checkout-conflict")
	if !errors.Is(err, ErrCheckoutCommandConflict) {
		t.Fatalf("conflicting replay error = %v", err)
	}
	if plans.calls != planCalls || subs.calls != subscriptionCalls || payment.calls != paymentCalls {
		t.Fatalf("conflicting replay reached mutable dependencies")
	}
}

func TestCheckoutUpgradeUsesTierRankAndKeepsServerSnapshot(t *testing.T) {
	service, _, subs, _, payment := newCheckoutFixture()
	subs.found = true
	subs.current = CurrentSubscription{ID: 7, PlanVersionID: 1, PlanCode: "qhpro.basic", TierRank: 10, Status: "active"}

	_, err := service.Checkout(context.Background(), Subject{Kind: SubjectProfile, ID: "42"}, "qhpro.pro", "checkout-upgrade")
	if err != nil {
		t.Fatal(err)
	}
	if payment.calls != 1 || payment.last.Terms.TierRank != 20 || payment.last.Terms.Price.AmountMinor != 19900000 {
		t.Fatalf("upgrade payment command = %+v", payment.last)
	}
}

func TestCheckoutSameTierAndDowngradeNeverCreatePaymentOrder(t *testing.T) {
	cases := []struct {
		name        string
		currentTier int32
		want        error
	}{
		{name: "same", currentTier: 20, want: ErrSamePlan},
		{name: "downgrade", currentTier: 30, want: ErrDowngradeRequiresSchedule},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			service, _, subs, commands, payment := newCheckoutFixture()
			subs.found = true
			subs.current = CurrentSubscription{ID: 7, PlanVersionID: 1, PlanCode: "current", TierRank: tc.currentTier, Status: "active"}
			_, err := service.Checkout(context.Background(), Subject{Kind: SubjectProfile, ID: "42"}, "qhpro.pro", "checkout-x")
			if !errors.Is(err, tc.want) {
				t.Fatalf("Checkout error = %v, want %v", err, tc.want)
			}
			if commands.createCalls != 0 || payment.calls != 0 {
				t.Fatalf("command/payment created for %s: commands=%d payment=%d", tc.name, commands.createCalls, payment.calls)
			}
		})
	}
}

func TestCheckoutPendingChangeBlocksPayment(t *testing.T) {
	service, _, subs, commands, payment := newCheckoutFixture()
	subs.found = true
	subs.current = CurrentSubscription{ID: 7, PlanVersionID: 1, PlanCode: "qhpro.basic", TierRank: 10, Status: "active", PendingPlanVersionID: 99}
	_, err := service.Checkout(context.Background(), Subject{Kind: SubjectProfile, ID: "42"}, "qhpro.pro", "checkout-x")
	if !errors.Is(err, ErrPendingChange) || commands.createCalls != 0 || payment.calls != 0 {
		t.Fatalf("pending checkout err=%v command calls=%d payment calls=%d", err, commands.createCalls, payment.calls)
	}
}

func TestCheckoutFailsClosedWhenCommercialTermsAreMissing(t *testing.T) {
	mutations := []struct {
		name string
		fn   func(*PlanTerms)
	}{
		{name: "tier", fn: func(t *PlanTerms) { t.TierRank = 0 }},
		{name: "term", fn: func(t *PlanTerms) { t.SubscriptionTermDays = 0 }},
		{name: "price", fn: func(t *PlanTerms) { t.Price.AmountMinor = 0 }},
		{name: "checksum", fn: func(t *PlanTerms) { t.TermsChecksum = "" }},
	}
	for _, tc := range mutations {
		t.Run(tc.name, func(t *testing.T) {
			service, plans, _, commands, payment := newCheckoutFixture()
			terms := plans.terms
			tc.fn(&terms)
			plans.terms = terms
			_, err := service.Checkout(context.Background(), Subject{Kind: SubjectProfile, ID: "42"}, "qhpro.pro", "checkout-x")
			if !errors.Is(err, ErrPlanTermsUnavailable) || commands.createCalls != 0 || payment.calls != 0 {
				t.Fatalf("missing %s err=%v command calls=%d payment calls=%d", tc.name, err, commands.createCalls, payment.calls)
			}
		})
	}
}

func TestCheckoutEnforcesSubjectScope(t *testing.T) {
	service, plans, _, commands, payment := newCheckoutFixture()
	plans.terms.SubjectScope = plan.SubjectScopeOrganization
	_, err := service.Checkout(context.Background(), Subject{Kind: SubjectProfile, ID: "42"}, "qhpro.pro", "checkout-x")
	if !errors.Is(err, ErrSubjectScope) || commands.createCalls != 0 || payment.calls != 0 {
		t.Fatalf("scope err=%v command calls=%d payment calls=%d", err, commands.createCalls, payment.calls)
	}
}

func TestCheckoutPersistenceFailureStopsPayment(t *testing.T) {
	service, _, _, commands, payment := newCheckoutFixture()
	commands.createErr = errors.New("database unavailable")
	_, err := service.Checkout(context.Background(), Subject{Kind: SubjectProfile, ID: "42"}, "qhpro.pro", "checkout-db-fail")
	if err == nil || payment.calls != 0 {
		t.Fatalf("persistence failure err=%v payment calls=%d", err, payment.calls)
	}
}

func TestCheckoutRejectsInvalidCommandBeforeDependencies(t *testing.T) {
	service, plans, subs, commands, payment := newCheckoutFixture()
	_, err := service.Checkout(context.Background(), Subject{}, "", "")
	if !errors.Is(err, ErrInvalidCommand) {
		t.Fatalf("invalid command error = %v", err)
	}
	if plans.calls != 0 || subs.calls != 0 || commands.findCalls != 0 || payment.calls != 0 {
		t.Fatalf("invalid command reached dependencies: plans=%d subs=%d commands=%d payment=%d", plans.calls, subs.calls, commands.findCalls, payment.calls)
	}
}

func TestClassifyPlanChange(t *testing.T) {
	cases := []struct {
		current int32
		target  int32
		want    PlanChange
	}{
		{current: 10, target: 10, want: PlanChangeSame},
		{current: 10, target: 20, want: PlanChangeUpgrade},
		{current: 20, target: 10, want: PlanChangeDowngrade},
	}
	for _, tc := range cases {
		got, err := ClassifyPlanChange(tc.current, tc.target)
		if err != nil || got != tc.want {
			t.Fatalf("ClassifyPlanChange(%d,%d) = %q err=%v want=%q", tc.current, tc.target, got, err, tc.want)
		}
	}
	if _, err := ClassifyPlanChange(0, 10); !errors.Is(err, ErrPlanTermsUnavailable) {
		t.Fatalf("invalid tier error = %v", err)
	}
}
