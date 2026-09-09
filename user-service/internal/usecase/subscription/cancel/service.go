package cancel

import (
	"context"
	"errors"
	"strings"
	"time"

	commonrequest "common/request"
	subscriptiondomain "user/internal/domain/subscription"
)

var (
	ErrUnavailable = errors.New("subscription cancellation is unavailable")
	ErrNotFound    = errors.New("active paid subscription was not found")
	ErrAmbiguous   = errors.New("multiple active paid subscriptions found")
	ErrInvalid     = errors.New("invalid subscription cancellation")
)

type Command struct {
	ProfileID   uint64
	ProductCode string
	Reason      string
	RequestID   string
	OperationID string
	OccurredAt  time.Time
}

type Result struct {
	Subscription subscriptiondomain.AdminProjection
	Changed      bool
}

// Store owns the single transaction that locks the current subscription,
// revokes access and appends lifecycle evidence.
type Store interface {
	CancelCurrent(context.Context, Command) (Result, error)
}

type Service struct {
	store Store
	now   func() time.Time
}

func NewService(store Store, now func() time.Time) *Service {
	if now == nil {
		now = time.Now
	}
	return &Service{store: store, now: now}
}

func (s *Service) Cancel(ctx context.Context, profileID uint64, productCode, reason string) (Result, error) {
	if s == nil || s.store == nil || s.now == nil {
		return Result{}, ErrUnavailable
	}
	if err := ctx.Err(); err != nil {
		return Result{}, err
	}
	productCode = strings.TrimSpace(productCode)
	reason = strings.TrimSpace(reason)
	if profileID == 0 || productCode == "" || reason == "" || len(reason) > 500 {
		return Result{}, ErrInvalid
	}
	requestID, _ := commonrequest.RequestIDFromContext(ctx)
	operationID, _ := commonrequest.OperationIDFromContext(ctx)
	return s.store.CancelCurrent(ctx, Command{
		ProfileID: profileID, ProductCode: productCode, Reason: reason,
		RequestID: requestID, OperationID: operationID, OccurredAt: s.now().UTC(),
	})
}
