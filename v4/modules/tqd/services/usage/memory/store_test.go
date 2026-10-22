package memoryusage

import (
	"context"
	"errors"
	"mdn4/modules/tqd/services/usage"
	"testing"
)

func TestInvalidUsageIsAnErrorRatherThanSuccessfulReplay(t *testing.T) {
	created, err := NewStore().SaveUsageEvent(context.Background(), usage.Event{})
	if err == nil || created {
		t.Fatalf("invalid event: created=%v err=%v", created, err)
	}
}

func TestCanceledContextDoesNotPretendUsageWasAccepted(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	created, err := NewStore().SaveUsageEvent(ctx, usage.Event{})
	if !errors.Is(err, context.Canceled) || created {
		t.Fatalf("canceled write: created=%v err=%v", created, err)
	}
}
