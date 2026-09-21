package worker

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"sync/atomic"
	"testing"
	"time"

	"payment/internal/usecase/outbox"
)

type idleStore struct{}

func (idleStore) ClaimOutbox(ctx context.Context, _ string, _ time.Time, _ time.Duration) (outbox.Message, bool, error) {
	if err := ctx.Err(); err != nil {
		return outbox.Message{}, false, err
	}
	return outbox.Message{}, false, nil
}
func (idleStore) MarkOutboxPublished(context.Context, uint64, string, uint64, time.Time) error {
	return nil
}
func (idleStore) RetryOutbox(context.Context, uint64, string, uint64, time.Time, string) error {
	return nil
}

type idlePublisher struct{}

func (idlePublisher) Publish(context.Context, outbox.Message) error { return nil }

func TestReconnectPolicyIsDeterministicBoundedAndAttemptAware(t *testing.T) {
	policy := ReconnectPolicy{Base: 8 * time.Millisecond, Max: 30 * time.Millisecond}

	first, firstCapped := policy.Delay("payment-outbox-host-1", 1)
	if firstCapped || first < 4*time.Millisecond || first > 8*time.Millisecond {
		t.Fatalf("first delay=%s capped=%v", first, firstCapped)
	}
	if got, capped := policy.Delay("payment-outbox-host-1", 1); got != first || capped {
		t.Fatalf("deterministic first delay=%s capped=%v, want %s,false", got, capped, first)
	}

	third, thirdCapped := policy.Delay("payment-outbox-host-1", 3)
	if !thirdCapped || third < 15*time.Millisecond || third > 30*time.Millisecond {
		t.Fatalf("third delay=%s capped=%v, want capped within [15ms,30ms]", third, thirdCapped)
	}
}

func TestOutboxSupervisorSurvivesInitialBrokerOutage(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var calls atomic.Int32
	factory := func() (*outbox.Service, func(), error) {
		call := calls.Add(1)
		if call < 3 {
			return nil, nil, errors.New("rabbit unavailable")
		}
		cancel()
		service := outbox.NewService(
			idleStore{},
			idlePublisher{},
			time.Now,
			outbox.RetryPolicy{Base: time.Millisecond, Max: 2 * time.Millisecond},
		)
		return service, func() {}, nil
	}

	supervisor := NewOutboxSupervisor(
		factory,
		"payment-outbox-test",
		time.Second,
		time.Millisecond,
		ReconnectPolicy{Base: time.Millisecond, Max: 2 * time.Millisecond},
	)
	if err := supervisor.Run(ctx); err != nil {
		t.Fatalf("Run() = %v, want nil after context cancellation", err)
	}
	if got := calls.Load(); got != 3 {
		t.Fatalf("factory calls = %d, want 3", got)
	}
}

func TestCappedReconnectEscalatesOnce(t *testing.T) {
	var buffer bytes.Buffer
	logger := slog.New(slog.NewJSONHandler(&buffer, nil))
	escalated := false

	logOutboxReconnectFailure(
		logger,
		errors.New("down"),
		30*time.Second,
		3,
		time.Minute,
		true,
		&escalated,
	)
	logOutboxReconnectFailure(
		logger,
		errors.New("still down"),
		30*time.Second,
		4,
		2*time.Minute,
		true,
		&escalated,
	)

	lines := bytes.Split(bytes.TrimSpace(buffer.Bytes()), []byte("\n"))
	if len(lines) != 2 {
		t.Fatalf("log lines = %d, want 2", len(lines))
	}
	var first, second map[string]any
	if err := json.Unmarshal(lines[0], &first); err != nil {
		t.Fatal(err)
	}
	if err := json.Unmarshal(lines[1], &second); err != nil {
		t.Fatal(err)
	}
	if first["level"] != "ERROR" || second["level"] != "WARN" {
		t.Fatalf("levels = %v, %v; want ERROR then WARN", first["level"], second["level"])
	}
	if first["failure_count"] != float64(3) || second["failure_count"] != float64(4) {
		t.Fatalf("failure counts = %v, %v", first["failure_count"], second["failure_count"])
	}
}
