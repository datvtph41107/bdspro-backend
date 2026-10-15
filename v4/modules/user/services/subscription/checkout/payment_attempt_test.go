package checkout

import (
	"context"
	"errors"
	"testing"
)

type fakeAttemptPort struct {
	calls   int
	command CreatePaymentAttemptCommand
	result  PaymentAttempt
	err     error
}

func (f *fakeAttemptPort) CreateAttempt(_ context.Context, command CreatePaymentAttemptCommand) (PaymentAttempt, error) {
	f.calls++
	f.command = command
	return f.result, f.err
}

func TestCreatePaymentAttemptDelegatesToPort(t *testing.T) {
	attempts := &fakeAttemptPort{result: PaymentAttempt{ID: 11, OrderID: 7, Method: "qr"}}
	service := &Service{attempts: attempts}
	subject := Subject{Kind: SubjectProfile, ID: "42"}

	got, err := service.CreatePaymentAttempt(context.Background(), CreatePaymentAttemptCommand{
		Subject: subject, OrderID: 7, Method: " qr ", CommandKey: " attempt-1 ",
	})
	if err != nil {
		t.Fatalf("CreatePaymentAttempt() error = %v", err)
	}
	if attempts.calls != 1 {
		t.Fatalf("attempt calls = %d, want 1", attempts.calls)
	}
	if attempts.command.Method != "qr" || attempts.command.CommandKey != "attempt-1" {
		t.Fatalf("command not normalized: %+v", attempts.command)
	}
	if got.ID != 11 {
		t.Fatalf("attempt id = %d, want 11", got.ID)
	}
}

func TestCreatePaymentAttemptRejectsInvalidBeforePort(t *testing.T) {
	attempts := &fakeAttemptPort{}
	service := &Service{attempts: attempts}
	_, err := service.CreatePaymentAttempt(context.Background(), CreatePaymentAttemptCommand{})
	if !errors.Is(err, ErrInvalidCommand) {
		t.Fatalf("error = %v, want ErrInvalidCommand", err)
	}
	if attempts.calls != 0 {
		t.Fatalf("attempt calls = %d, want 0", attempts.calls)
	}
}

func TestCreatePaymentAttemptPropagatesContextCancellation(t *testing.T) {
	attempts := &fakeAttemptPort{}
	service := &Service{attempts: attempts}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	_, err := service.CreatePaymentAttempt(ctx, CreatePaymentAttemptCommand{
		Subject: Subject{Kind: SubjectProfile, ID: "42"}, OrderID: 7, Method: "qr", CommandKey: "attempt-1",
	})
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("error = %v, want context.Canceled", err)
	}
	if attempts.calls != 0 {
		t.Fatalf("attempt calls = %d, want 0", attempts.calls)
	}
}
