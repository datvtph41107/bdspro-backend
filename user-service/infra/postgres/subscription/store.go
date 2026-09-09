package subscriptionpostgres

import (
	"context"
	"fmt"
	"strings"
	"time"

	_db "common/db"
	"gorm.io/gorm"
	"user/internal/usecase/subscription/settlement"
)

type SettlementStore struct {
	db  *_db.TransactionRepo
	now func() time.Time
}

func NewSettlementStore(db *_db.TransactionRepo) *SettlementStore {
	return &SettlementStore{db: db, now: func() time.Time { return time.Now().UTC() }}
}
func (s *SettlementStore) ApplySettlement(ctx context.Context, input settlement.Acceptance) (settlement.Result, error) {
	if s == nil || s.db == nil || strings.TrimSpace(input.Fingerprint) == "" {
		return settlement.Result{}, settlement.ErrInvalidEffect
	}
	if err := input.Effect.Validate(); err != nil {
		return settlement.Result{}, err
	}
	var result settlement.Result
	err := s.db.WithTransaction(ctx, func(txCtx context.Context) error {
		db := s.db.GetDB(txCtx)
		if err := advisoryLock(db, "subscription-effect:"+input.Effect.EffectKey); err != nil {
			return err
		}
		if existing, found, err := loadSettlementReceipt(db, input.Effect.EffectKey); err != nil {
			return err
		} else if found {
			if existing.Fingerprint != input.Fingerprint {
				return settlement.ErrEffectConflict
			}
			result = settlement.Result{Action: settlement.Action(existing.Action), Replayed: true}
			if existing.SubscriptionID != nil {
				result.SubscriptionID = *existing.SubscriptionID
			}
			return nil
		}
		target, err := loadTarget(db, input.Effect.PlanVersionID)
		if err != nil {
			return err
		}
		if !target.Matches(input.Effect) {
			return settlement.ErrCatalogTermsConflict
		}
		lockKey := fmt.Sprintf("subscription:%s:%s:%d", input.Effect.SubjectKind, input.Effect.SubjectID, target.ProductID)
		if err := advisoryLock(db, lockKey); err != nil {
			return err
		}

		current, err := loadCurrent(db, input.Effect, target.ProductID)
		if err != nil {
			return err
		}
		decision, err := settlement.Decide(input.Effect, target, current, s.now())
		if err != nil {
			return err
		}

		subscriptionID := uint64(0)
		switch decision.Action {
		case settlement.ActionActivated:
			subscriptionID, err = insertSubscription(db, decision.Subscription)
		case settlement.ActionUpgraded:
			subscriptionID = decision.Subscription.ID
			err = updateSubscription(db, decision.Subscription)
		case settlement.ActionRejectedSameTier, settlement.ActionRejectedDowngrade:
			subscriptionID = decision.Subscription.ID
		default:
			return settlement.ErrCurrentState
		}
		if err != nil {
			return err
		}

		if decision.Action == settlement.ActionActivated || decision.Action == settlement.ActionUpgraded {
			if err := appendEvent(db, input.Effect, decision, subscriptionID); err != nil {
				return err
			}
		}
		if err := insertSettlementReceipt(db, input, decision.Action, subscriptionID); err != nil {
			return err
		}
		result = settlement.Result{Action: decision.Action, SubscriptionID: subscriptionID}
		return nil
	})
	return result, err
}

func advisoryLock(db *gorm.DB, key string) error {
	if err := db.Exec("SELECT pg_advisory_xact_lock(hashtextextended(?, 0))", key).Error; err != nil {
		return fmt.Errorf("lock subscription settlement %q: %w", key, err)
	}
	return nil
}
