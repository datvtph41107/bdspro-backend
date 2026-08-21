package application

import (
	"context"
	"errors"
	"fmt"
	"time"

	"payment/internal/commerce/domain"
)

type FulfillmentService struct {
	store        FulfillmentStore
	subscription SubscriptionPort
	now          Clock
	retryDelay   time.Duration
}

func NewFulfillmentService(store FulfillmentStore, subscription SubscriptionPort, now Clock, retryDelay time.Duration) *FulfillmentService {
	return &FulfillmentService{store: store, subscription: subscription, now: now, retryDelay: retryDelay}
}

func (s *FulfillmentService) ProcessClaim(ctx context.Context, claim domain.FulfillmentClaim, workerID string) error {
	if s == nil || s.store == nil || s.subscription == nil || s.now == nil ||
		claim.Fulfillment.ID == 0 || claim.Order.ID == 0 || workerID == "" {
		return domain.ErrInvalidCommand
	}
	if err := ctx.Err(); err != nil {
		return err
	}
	if claim.Order.FundsConfirmedAt == nil || claim.Order.FundsConfirmedAt.IsZero() {
		return s.store.RequireReview(ctx, claim.Fulfillment.ID, workerID, claim.Fulfillment.ClaimVersion, s.now().UTC(), "order has no durable funds-confirmed time")
	}

	effect := domain.SettlementEffect{
		EffectKey:            fmt.Sprintf("payment.order.%d", claim.Order.ID),
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
	case errors.Is(err, domain.ErrSubscriptionTransient):
		return s.store.RetryFulfillment(ctx, claim.Fulfillment.ID, workerID, claim.Fulfillment.ClaimVersion, now.Add(s.retryDelay), err.Error())
	case errors.Is(err, domain.ErrSubscriptionPermanent), errors.Is(err, domain.ErrSubscriptionEffectConflict):
		return s.store.RequireReview(ctx, claim.Fulfillment.ID, workerID, claim.Fulfillment.ClaimVersion, now, err.Error())
	default:
		return s.store.RetryFulfillment(ctx, claim.Fulfillment.ID, workerID, claim.Fulfillment.ClaimVersion, now.Add(s.retryDelay), err.Error())
	}
}
