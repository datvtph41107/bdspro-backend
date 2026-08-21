package application

import (
	"context"
	"time"

	"payment/internal/commerce/domain"
)

type Clock func() time.Time

type ReferenceGenerator interface {
	NewOrderReference() (string, error)
}

type OrderStore interface {
	FindOrderByID(ctx context.Context, orderID uint64) (domain.Order, bool, error)
	FindOrderByCommand(ctx context.Context, subject domain.Subject, productCode, commandKey string) (domain.Order, bool, error)
	FindOrderByReference(ctx context.Context, reference string) (domain.Order, bool, error)
	CreateOrder(ctx context.Context, order domain.Order) (domain.Order, bool, error)
}

type SettlementDecision struct {
	Settlement        domain.Settlement
	ConfirmOrderID    uint64
	FundsConfirmedAt  *time.Time
	CreateFulfillment bool
}

type SettlementResult struct {
	Settlement  domain.Settlement
	Fulfillment *domain.Fulfillment
	Replay      bool
}

type SettlementStore interface {
	FindSettlementByProviderTransaction(ctx context.Context, provider, transactionID string) (domain.Settlement, bool, error)
	AcceptSettlement(ctx context.Context, decision SettlementDecision) (SettlementResult, error)
	MarkSettlementRequiresReview(ctx context.Context, settlementID uint64, reason string) error
}

type FulfillmentStore interface {
	ClaimFulfillment(ctx context.Context, workerID string, now time.Time, lease time.Duration) (domain.FulfillmentClaim, bool, error)
	CompleteFulfillment(ctx context.Context, fulfillmentID uint64, workerID string, claimVersion uint64, now time.Time) error
	RetryFulfillment(ctx context.Context, fulfillmentID uint64, workerID string, claimVersion uint64, availableAt time.Time, reason string) error
	RequireReview(ctx context.Context, fulfillmentID uint64, workerID string, claimVersion uint64, now time.Time, reason string) error
}

type RecoveryStore interface {
	RedriveFulfillment(ctx context.Context, command domain.RedriveCommand, now time.Time) (changed bool, replay bool, err error)
}

type SubscriptionPort interface {
	ApplySettlement(ctx context.Context, effect domain.SettlementEffect) error
}

type PaymentAttemptStore interface {
	FindAttemptByCommand(ctx context.Context, orderID uint64, commandKey string) (domain.PaymentAttempt, bool, error)
	CreateAttempt(ctx context.Context, attempt domain.PaymentAttempt) (domain.PaymentAttempt, bool, error)
	UpdateAttemptFromProvider(ctx context.Context, attemptID uint64, provider domain.ProviderAttempt, now time.Time) (domain.PaymentAttempt, error)
}

type AttemptProvider interface {
	Code() string
	Supports(method domain.PaymentMethod) bool
	CreateAttempt(ctx context.Context, order domain.Order, method domain.PaymentMethod, commandKey string) (domain.ProviderAttempt, error)
}

type OutboxMessage struct {
	ID            uint64
	EventKey      string
	EventType     string
	SchemaVersion int
	RoutingKey    string
	Payload       []byte
	AttemptCount  int
	ClaimVersion  uint64
	LockedBy      string
}

type OutboxStore interface {
	ClaimOutbox(ctx context.Context, workerID string, now time.Time, lease time.Duration) (OutboxMessage, bool, error)
	MarkOutboxPublished(ctx context.Context, id uint64, workerID string, claimVersion uint64, now time.Time) error
	RetryOutbox(ctx context.Context, id uint64, workerID string, claimVersion uint64, availableAt time.Time, reason string) error
}

type EventPublisher interface {
	Publish(ctx context.Context, message OutboxMessage) error
}
