package application

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"payment/internal/commerce/domain"
)

type staticReference struct {
	value string
	err   error
}

func (r staticReference) NewOrderReference() (string, error) { return r.value, r.err }

type memoryStore struct {
	mu sync.Mutex

	nextOrder       uint64
	nextSettlement  uint64
	nextFulfillment uint64
	nextAttempt     uint64

	orders             map[uint64]domain.Order
	orderByCommand     map[string]uint64
	orderByReference   map[string]uint64
	settlements        map[uint64]domain.Settlement
	settlementByTx     map[string]uint64
	fulfillments       map[uint64]domain.Fulfillment
	fulfillmentByOrder map[uint64]uint64
	redriveEffects     map[string]struct{}
	attempts           map[uint64]domain.PaymentAttempt
	attemptByCommand   map[string]uint64
}

func newMemoryStore() *memoryStore {
	return &memoryStore{
		nextOrder:          1,
		nextSettlement:     1,
		nextFulfillment:    1,
		nextAttempt:        1,
		orders:             map[uint64]domain.Order{},
		orderByCommand:     map[string]uint64{},
		orderByReference:   map[string]uint64{},
		settlements:        map[uint64]domain.Settlement{},
		settlementByTx:     map[string]uint64{},
		fulfillments:       map[uint64]domain.Fulfillment{},
		fulfillmentByOrder: map[uint64]uint64{},
		redriveEffects:     map[string]struct{}{},
		attempts:           map[uint64]domain.PaymentAttempt{},
		attemptByCommand:   map[string]uint64{},
	}
}

func commandScope(subject domain.Subject, product, command string) string {
	return fmt.Sprintf("%s|%s|%s|%s", subject.Kind, subject.ID, product, command)
}

func providerTxScope(provider, tx string) string { return provider + "|" + tx }

func (s *memoryStore) FindOrderByID(_ context.Context, orderID uint64) (domain.Order, bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	order, ok := s.orders[orderID]
	return order, ok, nil
}

func (s *memoryStore) FindOrderByCommand(_ context.Context, subject domain.Subject, productCode, commandKey string) (domain.Order, bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	id, ok := s.orderByCommand[commandScope(subject, productCode, commandKey)]
	if !ok {
		return domain.Order{}, false, nil
	}
	return s.orders[id], true, nil
}

func (s *memoryStore) FindOrderByReference(_ context.Context, reference string) (domain.Order, bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	id, ok := s.orderByReference[reference]
	if !ok {
		return domain.Order{}, false, nil
	}
	return s.orders[id], true, nil
}

func (s *memoryStore) CreateOrder(_ context.Context, order domain.Order) (domain.Order, bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	scope := commandScope(order.Subject, order.Terms.ProductCode, order.CommandKey)
	if id, ok := s.orderByCommand[scope]; ok {
		existing := s.orders[id]
		if existing.CommandFingerprint != order.CommandFingerprint {
			return domain.Order{}, false, domain.ErrOrderCommandConflict
		}
		return existing, false, nil
	}
	order.ID = s.nextOrder
	s.nextOrder++
	s.orders[order.ID] = order
	s.orderByCommand[scope] = order.ID
	s.orderByReference[order.Reference] = order.ID
	return order, true, nil
}

func (s *memoryStore) FindSettlementByProviderTransaction(_ context.Context, provider, transactionID string) (domain.Settlement, bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	id, ok := s.settlementByTx[providerTxScope(provider, transactionID)]
	if !ok {
		return domain.Settlement{}, false, nil
	}
	return s.settlements[id], true, nil
}

func (s *memoryStore) AcceptSettlement(_ context.Context, decision SettlementDecision) (SettlementResult, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	key := providerTxScope(decision.Settlement.Provider, decision.Settlement.ProviderTransactionID)
	if id, ok := s.settlementByTx[key]; ok {
		existing := s.settlements[id]
		if existing.EvidenceHash != decision.Settlement.EvidenceHash {
			existing.Status = domain.SettlementRequiresReview
			existing.ReviewReason = "provider transaction evidence changed"
			s.settlements[id] = existing
			return SettlementResult{}, domain.ErrProviderEvidenceConflict
		}
		return SettlementResult{Settlement: existing, Replay: true}, nil
	}

	settlement := decision.Settlement
	settlement.ID = s.nextSettlement
	s.nextSettlement++
	s.settlements[settlement.ID] = settlement
	s.settlementByTx[key] = settlement.ID
	result := SettlementResult{Settlement: settlement}

	if decision.ConfirmOrderID != 0 {
		order, ok := s.orders[decision.ConfirmOrderID]
		if !ok {
			return SettlementResult{}, domain.ErrOrderNotFound
		}
		order.Status = domain.OrderFundsConfirmed
		if decision.FundsConfirmedAt != nil {
			at := decision.FundsConfirmedAt.UTC()
			order.FundsConfirmedAt = &at
		}
		s.orders[order.ID] = order
	}
	if decision.CreateFulfillment {
		if existingID, ok := s.fulfillmentByOrder[decision.ConfirmOrderID]; ok {
			existing := s.fulfillments[existingID]
			result.Fulfillment = &existing
			return result, nil
		}
		now := settlement.CreatedAt
		fulfillment := domain.Fulfillment{
			ID:          s.nextFulfillment,
			OrderID:     decision.ConfirmOrderID,
			Status:      domain.FulfillmentPending,
			AvailableAt: now,
			CreatedAt:   now,
			UpdatedAt:   now,
		}
		s.nextFulfillment++
		s.fulfillments[fulfillment.ID] = fulfillment
		s.fulfillmentByOrder[fulfillment.OrderID] = fulfillment.ID
		result.Fulfillment = &fulfillment
	}
	return result, nil
}

func (s *memoryStore) MarkSettlementRequiresReview(_ context.Context, settlementID uint64, reason string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	settlement, ok := s.settlements[settlementID]
	if !ok {
		return errors.New("settlement missing")
	}
	settlement.Status = domain.SettlementRequiresReview
	settlement.ReviewReason = reason
	s.settlements[settlementID] = settlement
	return nil
}

func (s *memoryStore) ClaimFulfillment(_ context.Context, workerID string, now time.Time, lease time.Duration) (domain.FulfillmentClaim, bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for id, fulfillment := range s.fulfillments {
		claimable := (fulfillment.Status == domain.FulfillmentPending || fulfillment.Status == domain.FulfillmentRetry) && !fulfillment.AvailableAt.After(now)
		reclaimable := fulfillment.Status == domain.FulfillmentRunning && !fulfillment.LeaseUntil.After(now)
		if !claimable && !reclaimable {
			continue
		}
		fulfillment.Status = domain.FulfillmentRunning
		fulfillment.LockedBy = workerID
		fulfillment.ClaimVersion++
		fulfillment.LeaseUntil = now.Add(lease)
		fulfillment.AttemptCount++
		fulfillment.UpdatedAt = now
		s.fulfillments[id] = fulfillment
		return domain.FulfillmentClaim{Fulfillment: fulfillment, Order: s.orders[fulfillment.OrderID]}, true, nil
	}
	return domain.FulfillmentClaim{}, false, nil
}

func (s *memoryStore) verifyClaim(fulfillmentID uint64, workerID string, claimVersion uint64) (domain.Fulfillment, error) {
	fulfillment, ok := s.fulfillments[fulfillmentID]
	if !ok {
		return domain.Fulfillment{}, domain.ErrFulfillmentNotFound
	}
	if fulfillment.Status != domain.FulfillmentRunning || fulfillment.LockedBy != workerID || fulfillment.ClaimVersion != claimVersion {
		return domain.Fulfillment{}, domain.ErrClaimLost
	}
	return fulfillment, nil
}

func (s *memoryStore) CompleteFulfillment(_ context.Context, fulfillmentID uint64, workerID string, claimVersion uint64, now time.Time) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	fulfillment, err := s.verifyClaim(fulfillmentID, workerID, claimVersion)
	if err != nil {
		return err
	}
	fulfillment.Status = domain.FulfillmentCompleted
	fulfillment.LockedBy = ""
	fulfillment.LeaseUntil = time.Time{}
	fulfillment.UpdatedAt = now
	completed := now
	fulfillment.CompletedAt = &completed
	s.fulfillments[fulfillment.ID] = fulfillment
	return nil
}

func (s *memoryStore) RetryFulfillment(_ context.Context, fulfillmentID uint64, workerID string, claimVersion uint64, availableAt time.Time, reason string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	fulfillment, err := s.verifyClaim(fulfillmentID, workerID, claimVersion)
	if err != nil {
		return err
	}
	fulfillment.Status = domain.FulfillmentRetry
	fulfillment.LockedBy = ""
	fulfillment.LeaseUntil = time.Time{}
	fulfillment.AvailableAt = availableAt
	fulfillment.LastError = reason
	fulfillment.UpdatedAt = availableAt
	s.fulfillments[fulfillment.ID] = fulfillment
	return nil
}

func (s *memoryStore) RequireReview(_ context.Context, fulfillmentID uint64, workerID string, claimVersion uint64, now time.Time, reason string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	fulfillment, err := s.verifyClaim(fulfillmentID, workerID, claimVersion)
	if err != nil {
		return err
	}
	fulfillment.Status = domain.FulfillmentRequiresReview
	fulfillment.LockedBy = ""
	fulfillment.LeaseUntil = time.Time{}
	fulfillment.LastError = reason
	fulfillment.UpdatedAt = now
	s.fulfillments[fulfillment.ID] = fulfillment
	return nil
}

func (s *memoryStore) RedriveFulfillment(_ context.Context, command domain.RedriveCommand, now time.Time) (changed bool, replay bool, err error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	fulfillment, ok := s.fulfillments[command.FulfillmentID]
	if !ok {
		return false, false, domain.ErrFulfillmentNotFound
	}
	effectKey := fmt.Sprintf("fulfillment:%d:redrive:%s", command.FulfillmentID, command.CommandKey)
	if _, ok := s.redriveEffects[effectKey]; ok {
		return false, true, nil
	}
	// Consume identity before state classification so a no-op cannot come alive later.
	s.redriveEffects[effectKey] = struct{}{}
	if fulfillment.Status != domain.FulfillmentRequiresReview {
		return false, false, nil
	}
	fulfillment.Status = domain.FulfillmentPending
	fulfillment.AvailableAt = now
	fulfillment.LastError = ""
	fulfillment.UpdatedAt = now
	s.fulfillments[fulfillment.ID] = fulfillment
	return true, false, nil
}

type idempotentSubscription struct {
	mu                 sync.Mutex
	effects            map[string]domain.SettlementEffect
	businessApplyCount int
	loseFirstResponse  bool
}

func newIdempotentSubscription(loseFirstResponse bool) *idempotentSubscription {
	return &idempotentSubscription{effects: map[string]domain.SettlementEffect{}, loseFirstResponse: loseFirstResponse}
}

func (s *idempotentSubscription) ApplySettlement(_ context.Context, effect domain.SettlementEffect) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if existing, ok := s.effects[effect.EffectKey]; ok {
		if existing != effect {
			return domain.ErrSubscriptionEffectConflict
		}
		return nil
	}
	s.effects[effect.EffectKey] = effect
	s.businessApplyCount++
	if s.loseFirstResponse {
		s.loseFirstResponse = false
		return domain.ErrSubscriptionTransient
	}
	return nil
}

func terms() domain.CommercialTerms {
	return domain.CommercialTerms{
		ProductCode:          "qhpro",
		PlanCode:             "qhpro.pro",
		PlanVersionID:        11,
		PlanVersion:          "2.0.0",
		TierRank:             20,
		SubscriptionTermDays: 30,
		TermsChecksum:        "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
		Price:                domain.Money{Currency: domain.CurrencyVND, AmountMinor: 19900000},
	}
}

func subject() domain.Subject { return domain.Subject{Kind: domain.SubjectProfile, ID: "42"} }

func attemptCommandScope(orderID uint64, commandKey string) string {
	return fmt.Sprintf("%d|%s", orderID, commandKey)
}

func (s *memoryStore) FindAttemptByCommand(_ context.Context, orderID uint64, commandKey string) (domain.PaymentAttempt, bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	id, ok := s.attemptByCommand[attemptCommandScope(orderID, commandKey)]
	if !ok {
		return domain.PaymentAttempt{}, false, nil
	}
	return s.attempts[id], true, nil
}

func (s *memoryStore) CreateAttempt(_ context.Context, attempt domain.PaymentAttempt) (domain.PaymentAttempt, bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	key := attemptCommandScope(attempt.OrderID, attempt.CommandKey)
	if id, ok := s.attemptByCommand[key]; ok {
		existing := s.attempts[id]
		if existing.CommandFingerprint != attempt.CommandFingerprint {
			return domain.PaymentAttempt{}, false, domain.ErrAttemptCommandConflict
		}
		return existing, false, nil
	}
	attempt.ID = s.nextAttempt
	s.nextAttempt++
	s.attempts[attempt.ID] = attempt
	s.attemptByCommand[key] = attempt.ID
	return attempt, true, nil
}

func (s *memoryStore) UpdateAttemptFromProvider(_ context.Context, attemptID uint64, provider domain.ProviderAttempt, now time.Time) (domain.PaymentAttempt, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	attempt, ok := s.attempts[attemptID]
	if !ok {
		return domain.PaymentAttempt{}, domain.ErrAttemptNotFound
	}
	attempt.ProviderReference = provider.ProviderReference
	attempt.Status = provider.Status
	attempt.NextAction = provider.NextAction
	attempt.ExpiresAt = provider.ExpiresAt
	attempt.UpdatedAt = now
	s.attempts[attemptID] = attempt
	return attempt, nil
}
