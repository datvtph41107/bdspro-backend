package application

import (
	"context"
	"testing"
	"time"

	"payment/internal/commerce/domain"
)

type fakeAttemptProvider struct {
	code string
}

func (p fakeAttemptProvider) Code() string { return p.code }
func (p fakeAttemptProvider) Supports(method domain.PaymentMethod) bool {
	return method == domain.PaymentMethodQR || method == domain.PaymentMethodBankTransfer
}
func (p fakeAttemptProvider) CreateAttempt(_ context.Context, order domain.Order, method domain.PaymentMethod, _ string) (domain.ProviderAttempt, error) {
	expires := order.CreatedAt.Add(5 * time.Minute)
	return domain.ProviderAttempt{
		ProviderReference: "provider-" + order.Reference,
		Status:            domain.AttemptPendingAction,
		NextAction: domain.NextAction{
			Kind:      domain.NextActionDisplayQR,
			QRPayload: "PAY:" + order.Reference,
			ExpiresAt: &expires,
		},
		ExpiresAt: &expires,
	}, nil
}

func TestAttemptCommandReplayAndConflict(t *testing.T) {
	now := time.Date(2026, 8, 20, 1, 0, 0, 0, time.UTC)
	store := newMemoryStore()
	orderSvc := NewOrderService(store, staticReference{value: "QHP-100"}, func() time.Time { return now })
	order, _, err := orderSvc.CreateOrder(context.Background(), CreateOrderCommand{
		Subject: domain.Subject{Kind: domain.SubjectProfile, ID: "u-1"}, Terms: terms(), CommandKey: "checkout-1", ExpiresAt: now.Add(time.Hour),
	})
	if err != nil {
		t.Fatal(err)
	}

	svc := NewAttemptService(store, store, []AttemptProvider{fakeAttemptProvider{code: "sepay_qr"}}, func() time.Time { return now.Add(time.Minute) })
	command := CreateAttemptCommand{OrderID: order.ID, Method: domain.PaymentMethodQR, CommandKey: "attempt-1"}
	first, created, err := svc.CreateAttempt(context.Background(), command)
	if err != nil || !created {
		t.Fatalf("first attempt created=%v err=%v", created, err)
	}
	second, created, err := svc.CreateAttempt(context.Background(), command)
	if err != nil || created || first.ID != second.ID {
		t.Fatalf("replay created=%v err=%v first=%d second=%d", created, err, first.ID, second.ID)
	}

	command.Method = domain.PaymentMethodBankTransfer
	if _, _, err := svc.CreateAttempt(context.Background(), command); err != domain.ErrAttemptCommandConflict {
		t.Fatalf("changed payload err=%v", err)
	}
}

func TestAttemptCannotStartAfterFundsConfirmed(t *testing.T) {
	now := time.Date(2026, 8, 20, 1, 0, 0, 0, time.UTC)
	store := newMemoryStore()
	orderSvc := NewOrderService(store, staticReference{value: "QHP-101"}, func() time.Time { return now })
	order, _, _ := orderSvc.CreateOrder(context.Background(), CreateOrderCommand{
		Subject: domain.Subject{Kind: domain.SubjectProfile, ID: "u-1"}, Terms: terms(), CommandKey: "checkout-2", ExpiresAt: now.Add(time.Hour),
	})
	order.Status = domain.OrderFundsConfirmed
	store.orders[order.ID] = order

	svc := NewAttemptService(store, store, []AttemptProvider{fakeAttemptProvider{code: "sepay_qr"}}, func() time.Time { return now.Add(time.Minute) })
	_, _, err := svc.CreateAttempt(context.Background(), CreateAttemptCommand{OrderID: order.ID, Method: domain.PaymentMethodQR, CommandKey: "attempt-2"})
	if err != domain.ErrAttemptNotAllowed {
		t.Fatalf("err=%v", err)
	}
}
