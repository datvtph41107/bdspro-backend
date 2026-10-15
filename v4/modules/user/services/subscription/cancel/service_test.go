package cancel

import (
	"context"
	"errors"
	"testing"
	"time"
)

type storeStub struct {
	command Command
	result  Result
	err     error
}

func (s *storeStub) CancelCurrent(_ context.Context, command Command) (Result, error) {
	s.command = command
	return s.result, s.err
}

func TestCancelValidatesAndDelegatesCanonicalCommand(t *testing.T) {
	now := time.Date(2026, 9, 8, 8, 0, 0, 0, time.FixedZone("ICT", 7*60*60))
	store := &storeStub{result: Result{Changed: true}}
	service := NewService(store, func() time.Time { return now })

	result, err := service.Cancel(context.Background(), 29, " qhpro ", " Người dùng yêu cầu hủy ")
	if err != nil {
		t.Fatalf("Cancel() error = %v", err)
	}
	if !result.Changed || store.command.ProfileID != 29 || store.command.ProductCode != "qhpro" ||
		store.command.Reason != "Người dùng yêu cầu hủy" || !store.command.OccurredAt.Equal(now.UTC()) {
		t.Fatalf("unexpected cancellation command/result: command=%+v result=%+v", store.command, result)
	}
}

func TestCancelRejectsIncompleteCommandWithoutCallingStore(t *testing.T) {
	store := &storeStub{err: errors.New("must not be called")}
	service := NewService(store, time.Now)
	for _, testCase := range []struct {
		profileID uint64
		product   string
		reason    string
	}{
		{0, "qhpro", "requested"},
		{29, "", "requested"},
		{29, "qhpro", ""},
	} {
		if _, err := service.Cancel(context.Background(), testCase.profileID, testCase.product, testCase.reason); !errors.Is(err, ErrInvalid) {
			t.Fatalf("Cancel(%+v) error = %v, want ErrInvalid", testCase, err)
		}
	}
}
