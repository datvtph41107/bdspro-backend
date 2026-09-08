package fulfillment

import (
	"context"
	"errors"
	"time"

	payment "payment/internal/domain/payment"
)

type Clock func() time.Time

type Store interface {
	ClaimFulfillment(ctx context.Context, workerID string, now time.Time, lease time.Duration) (payment.FulfillmentClaim, bool, error)
	CompleteFulfillment(ctx context.Context, fulfillmentID uint64, workerID string, claimVersion uint64, now time.Time) error
	RetryFulfillment(ctx context.Context, fulfillmentID uint64, workerID string, claimVersion uint64, availableAt time.Time, reason string) error
	RequireReview(ctx context.Context, fulfillmentID uint64, workerID string, claimVersion uint64, now time.Time, reason string) error
}

type SubscriptionPort interface {
	ApplySettlement(ctx context.Context, effect payment.SettlementEffect) error
}

type Service struct {
	store        Store
	subscription SubscriptionPort
	now          Clock
	retryDelay   time.Duration
}

func NewService(store Store, subscription SubscriptionPort, now Clock, retryDelay time.Duration) *Service {
	return &Service{store: store, subscription: subscription, now: now, retryDelay: retryDelay}
}

func (s *Service) ProcessClaim(ctx context.Context, claim payment.FulfillmentClaim, workerID string) error {
	if s == nil || s.store == nil || s.subscription == nil || s.now == nil ||
		claim.Fulfillment.ID == 0 || claim.Order.ID == 0 || workerID == "" {
		return payment.ErrInvalidCommand
	}
	if err := ctx.Err(); err != nil {
		return err
	}
	if claim.Order.FundsConfirmedAt == nil || claim.Order.FundsConfirmedAt.IsZero() {
		return s.store.RequireReview(ctx, claim.Fulfillment.ID, workerID, claim.Fulfillment.ClaimVersion, s.now().UTC(), "order has no durable funds-confirmed time")
	}

	effect := payment.SettlementEffect{
		EffectKey:            claim.Order.SettlementEffectKey(),
		OrderID:              claim.Order.ID,
		Subject:              claim.Order.Subject,
		ProductCode:          claim.Order.Terms.ProductCode,
		PlanCode:             claim.Order.Terms.PlanCode,
		PlanVersionID:        claim.Order.Terms.PlanVersionID,
		PlanVersion:          claim.Order.Terms.PlanVersion,
		TierRank:             claim.Order.Terms.TierRank,
		SubscriptionTermDays: claim.Order.Terms.SubscriptionTermDays,
		TermsChecksum:        claim.Order.Terms.TermsChecksum,
		OccurredAt:           claim.Order.FundsConfirmedAt.UTC(),
	}

	err := s.subscription.ApplySettlement(ctx, effect)
	now := s.now().UTC()
	switch {
	case err == nil:
		return s.store.CompleteFulfillment(ctx, claim.Fulfillment.ID, workerID, claim.Fulfillment.ClaimVersion, now)
	case errors.Is(err, payment.ErrSubscriptionTransient):
		return s.store.RetryFulfillment(ctx, claim.Fulfillment.ID, workerID, claim.Fulfillment.ClaimVersion, now.Add(s.retryDelay), err.Error())
	case errors.Is(err, payment.ErrSubscriptionPermanent), errors.Is(err, payment.ErrSubscriptionEffectConflict):
		return s.store.RequireReview(ctx, claim.Fulfillment.ID, workerID, claim.Fulfillment.ClaimVersion, now, err.Error())
	default:
		return s.store.RetryFulfillment(ctx, claim.Fulfillment.ID, workerID, claim.Fulfillment.ClaimVersion, now.Add(s.retryDelay), err.Error())
	}
}
