package delivery

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	domain "notification/internal/domain/delivery"
)

type Store interface {
	ClaimNext(ctx context.Context, workerID string, now time.Time, lease time.Duration) (*domain.Intent, error)
	MarkSent(ctx context.Context, intent domain.Intent, sentAt time.Time) error
	MarkRetry(ctx context.Context, intent domain.Intent, availableAt time.Time, cause string) error
	MarkUnknown(ctx context.Context, intent domain.Intent, availableAt time.Time, cause string) error
	MarkFailed(ctx context.Context, intent domain.Intent, cause string) error
}

type TokenResolver interface {
	GetPushTokensByProfileId(ctx context.Context, profileID uint64) ([]string, error)
}

type PushSender interface {
	SendPushNotification(ctx context.Context, tokens []string, title string, message string, data map[string]string) error
}

type Service struct {
	store  Store
	tokens TokenResolver
	push   PushSender
	now    func() time.Time
	lease  time.Duration
}

func NewService(store Store, tokens TokenResolver, push PushSender, lease time.Duration) *Service {
	if lease <= 0 {
		lease = 30 * time.Second
	}
	return &Service{store: store, tokens: tokens, push: push, now: time.Now, lease: lease}
}

// ProcessOne executes one durable delivery responsibility. External push
// failures are classified UNKNOWN because the current Firebase adapter cannot
// prove whether the provider accepted a request before the response was lost.
// Retrying UNKNOWN therefore intentionally provides at-least-once semantics.
func (s *Service) ProcessOne(ctx context.Context, workerID string) (bool, error) {
	if s == nil || s.store == nil || s.tokens == nil || s.push == nil {
		return false, fmt.Errorf("delivery service missing dependency")
	}
	now := s.now().UTC()
	intent, err := s.store.ClaimNext(ctx, workerID, now, s.lease)
	if err != nil || intent == nil {
		return false, err
	}
	if !strings.EqualFold(intent.SubjectKind, "profile") {
		return true, s.store.MarkFailed(ctx, *intent, "push delivery supports profile subject only")
	}
	profileID, err := strconv.ParseUint(strings.TrimSpace(intent.SubjectID), 10, 64)
	if err != nil || profileID == 0 {
		return true, s.store.MarkFailed(ctx, *intent, "invalid profile subject id")
	}
	tokens, err := s.tokens.GetPushTokensByProfileId(ctx, profileID)
	if err != nil {
		return true, s.store.MarkRetry(ctx, *intent, now.Add(backoff(intent.AttemptCount)), err.Error())
	}
	if len(tokens) == 0 {
		return true, s.store.MarkFailed(ctx, *intent, "no push token")
	}

	title, body := render(intent.TemplateCode, intent.Payload)
	data := map[string]string{"event_id": intent.EventID, "template_code": intent.TemplateCode}
	if err := s.push.SendPushNotification(ctx, tokens, title, body, data); err != nil {
		var effectErr *domain.EffectError
		if errors.As(err, &effectErr) {
			switch effectErr.Kind {
			case domain.EffectFailurePermanent:
				return true, s.store.MarkFailed(ctx, *intent, effectErr.Error())
			case domain.EffectFailureRetryable:
				return true, s.store.MarkRetry(ctx, *intent, now.Add(backoff(intent.AttemptCount)), effectErr.Error())
			case domain.EffectFailureUnknown:
				return true, s.store.MarkUnknown(ctx, *intent, now.Add(backoff(intent.AttemptCount)), effectErr.Error())
			}
		}
		// An unclassified provider error is conservatively UNKNOWN because the
		// adapter cannot prove the external effect did not happen.
		return true, s.store.MarkUnknown(ctx, *intent, now.Add(backoff(intent.AttemptCount)), err.Error())
	}
	return true, s.store.MarkSent(ctx, *intent, s.now().UTC())
}

func render(template, payload string) (string, string) {
	switch template {
	case "payment_completed":
		return "Thanh toán thành công", "Giao dịch của bạn đã được xác nhận: " + payload
	default:
		return "QHPRO", payload
	}
}

func backoff(attempt int) time.Duration {
	if attempt < 1 {
		attempt = 1
	}
	if attempt > 6 {
		attempt = 6
	}
	return time.Duration(1<<uint(attempt-1)) * time.Second
}
