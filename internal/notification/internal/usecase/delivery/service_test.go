package delivery

import (
	"context"
	"errors"
	"testing"
	"time"

	domain "notification/internal/domain/delivery"
)

type storeProbe struct {
	intent *domain.Intent
	mark   string
}

func (p *storeProbe) ClaimNext(context.Context, string, time.Time, time.Duration) (*domain.Intent, error) {
	return p.intent, nil
}
func (p *storeProbe) MarkSent(context.Context, domain.Intent, time.Time) error {
	p.mark = "sent"
	return nil
}
func (p *storeProbe) MarkRetry(context.Context, domain.Intent, time.Time, string) error {
	p.mark = "retry"
	return nil
}
func (p *storeProbe) MarkUnknown(context.Context, domain.Intent, time.Time, string) error {
	p.mark = "unknown"
	return nil
}
func (p *storeProbe) MarkFailed(context.Context, domain.Intent, string) error {
	p.mark = "failed"
	return nil
}

type tokensProbe struct {
	tokens []string
	err    error
}

func (p tokensProbe) GetPushTokensByProfileId(context.Context, uint64) ([]string, error) {
	return p.tokens, p.err
}

type pushProbe struct {
	err   error
	calls int
}

func (p *pushProbe) SendPushNotification(context.Context, []string, string, string, map[string]string) error {
	p.calls++
	return p.err
}

func intent() *domain.Intent {
	return &domain.Intent{ID: 1, EventID: "event-1", SubjectKind: "profile", SubjectID: "42", Channel: "push", TemplateCode: "payment_completed", Payload: "pro|VND|100|42", Status: domain.StatusRunning, LockedBy: "w1", ClaimVersion: 2, AttemptCount: 1}
}

func TestProcessOneMarksProviderResponseFailureUnknown(t *testing.T) {
	store := &storeProbe{intent: intent()}
	push := &pushProbe{err: errors.New("response lost")}
	service := NewService(store, tokensProbe{tokens: []string{"opaque-token"}}, push, time.Second)
	processed, err := service.ProcessOne(context.Background(), "w1")
	if err != nil || !processed {
		t.Fatalf("ProcessOne() = %v, %v", processed, err)
	}
	if store.mark != "unknown" {
		t.Fatalf("mark = %q, want unknown", store.mark)
	}
}

func TestProcessOneRetriesTokenAuthorityFailure(t *testing.T) {
	store := &storeProbe{intent: intent()}
	service := NewService(store, tokensProbe{err: errors.New("user unavailable")}, &pushProbe{}, time.Second)
	processed, err := service.ProcessOne(context.Background(), "w1")
	if err != nil || !processed {
		t.Fatalf("ProcessOne() = %v, %v", processed, err)
	}
	if store.mark != "retry" {
		t.Fatalf("mark = %q, want retry", store.mark)
	}
}

func TestProcessOneFencesUnsupportedSubjectBeforeProvider(t *testing.T) {
	in := intent()
	in.SubjectKind = "organization"
	store := &storeProbe{intent: in}
	push := &pushProbe{}
	service := NewService(store, tokensProbe{tokens: []string{"opaque-token"}}, push, time.Second)
	processed, err := service.ProcessOne(context.Background(), "w1")
	if err != nil || !processed {
		t.Fatalf("ProcessOne() = %v, %v", processed, err)
	}
	if store.mark != "failed" {
		t.Fatalf("mark = %q, want failed", store.mark)
	}
	if push.calls != 0 {
		t.Fatalf("push calls = %d, want 0", push.calls)
	}
}
