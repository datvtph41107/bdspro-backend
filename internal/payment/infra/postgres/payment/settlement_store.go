package postgres

import (
	"context"
	"errors"
	"fmt"

	domain "payment/internal/domain/payment"
	"payment/internal/usecase/settlement"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func (s *Store) FindSettlementByProviderTransaction(ctx context.Context, provider, transactionID string) (domain.Settlement, bool, error) {
	if s == nil || s.db == nil {
		return domain.Settlement{}, false, errors.New("commerce database is not configured")
	}
	var row settlementModel
	err := s.db.WithContext(ctx).
		Where("provider = ? AND provider_transaction_id = ?", provider, transactionID).
		First(&row).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return domain.Settlement{}, false, nil
	}
	if err != nil {
		return domain.Settlement{}, false, fmt.Errorf("find commerce settlement: %w", err)
	}
	return settlementFromModel(row), true, nil
}

func (s *Store) AcceptSettlement(ctx context.Context, decision settlement.Decision) (settlement.Result, error) {
	if s == nil || s.db == nil {
		return settlement.Result{}, errors.New("commerce database is not configured")
	}

	var accepted settlement.Result
	var outcomeErr error
	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if decision.CommandEffect != nil {
			replay, err := consumeCommandEffect(tx, *decision.CommandEffect)
			if err != nil {
				return err
			}
			if replay {
				accepted.Replay = true
			}
		}
		row := settlementToModel(decision.Settlement)
		create := tx.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "provider"}, {Name: "provider_transaction_id"}},
			DoNothing: true,
		}).Create(&row)
		if create.Error != nil {
			return fmt.Errorf("insert commerce settlement: %w", create.Error)
		}
		if create.RowsAffected == 0 {
			var existing settlementModel
			if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
				Where("provider = ? AND provider_transaction_id = ?", decision.Settlement.Provider, decision.Settlement.ProviderTransactionID).
				First(&existing).Error; err != nil {
				return fmt.Errorf("load settlement conflict winner: %w", err)
			}
			if existing.EvidenceHash != decision.Settlement.EvidenceHash {
				if err := tx.Model(&settlementModel{}).Where("id = ?", existing.ID).Updates(map[string]any{
					"status":        domain.SettlementRequiresReview,
					"review_reason": "provider transaction evidence changed",
					"updated_at":    gorm.Expr("NOW()"),
				}).Error; err != nil {
					return fmt.Errorf("mark settlement evidence conflict: %w", err)
				}
				existing.Status = string(domain.SettlementRequiresReview)
				existing.ReviewReason = "provider transaction evidence changed"
				accepted.Settlement = settlementFromModel(existing)
				outcomeErr = domain.ErrProviderEvidenceConflict
				return nil
			}
			accepted.Settlement = settlementFromModel(existing)
			accepted.Replay = true
			return nil
		}

		accepted.Settlement = settlementFromModel(row)
		if decision.ConfirmOrderID == 0 || !decision.CreateFulfillment {
			return nil
		}

		var order orderModel
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&order, decision.ConfirmOrderID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return domain.ErrOrderNotFound
			}
			return fmt.Errorf("lock commerce order for settlement: %w", err)
		}
		if order.Status != string(domain.OrderPendingFunds) {
			reason := "additional settlement for order already classified"
			if err := tx.Model(&settlementModel{}).Where("id = ?", row.ID).Updates(map[string]any{
				"status":        domain.SettlementRequiresReview,
				"review_reason": reason,
				"updated_at":    gorm.Expr("NOW()"),
			}).Error; err != nil {
				return fmt.Errorf("mark additional settlement for review: %w", err)
			}
			row.Status = string(domain.SettlementRequiresReview)
			row.ReviewReason = reason
			accepted.Settlement = settlementFromModel(row)
			return nil
		}

		if decision.FundsConfirmedAt == nil || decision.FundsConfirmedAt.IsZero() {
			return errors.New("funds confirmed settlement requires occurred_at")
		}
		confirmedAt := decision.FundsConfirmedAt.UTC()
		if err := tx.Model(&orderModel{}).Where("id = ? AND status = ?", order.ID, domain.OrderPendingFunds).Updates(map[string]any{
			"status":             domain.OrderFundsConfirmed,
			"funds_confirmed_at": confirmedAt,
			"updated_at":         gorm.Expr("NOW()"),
		}).Error; err != nil {
			return fmt.Errorf("confirm commerce order funds: %w", err)
		}

		fulfillment := fulfillmentModel{
			OrderID:     order.ID,
			Status:      string(domain.FulfillmentPending),
			AvailableAt: row.CreatedAt,
			CreatedAt:   row.CreatedAt,
			UpdatedAt:   row.CreatedAt,
		}
		insertFulfillment := tx.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "order_id"}},
			DoNothing: true,
		}).Create(&fulfillment)
		if insertFulfillment.Error != nil {
			return fmt.Errorf("create commerce fulfillment: %w", insertFulfillment.Error)
		}
		if insertFulfillment.RowsAffected == 0 {
			if err := tx.Where("order_id = ?", order.ID).First(&fulfillment).Error; err != nil {
				return fmt.Errorf("load commerce fulfillment conflict winner: %w", err)
			}
		}
		mapped := fulfillmentFromModel(fulfillment)
		accepted.Fulfillment = &mapped
		return nil
	})
	if err != nil {
		return settlement.Result{}, err
	}
	if outcomeErr != nil {
		return accepted, outcomeErr
	}
	return accepted, nil
}

func consumeCommandEffect(tx *gorm.DB, effect domain.CommandEffect) (bool, error) {
	if tx == nil || !effect.IsValid() {
		return false, domain.ErrInvalidCommand
	}
	row := commandEffectModel{
		EffectType: effect.EffectType, ScopeID: effect.ScopeID, CommandKey: effect.CommandKey,
		ActorID: effect.ActorID, Reason: effect.Reason, Outcome: effect.Outcome, CreatedAt: effect.CreatedAt,
	}
	insert := tx.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "effect_type"}, {Name: "scope_id"}, {Name: "command_key"}},
		DoNothing: true,
	}).Create(&row)
	if insert.Error != nil {
		return false, fmt.Errorf("consume settlement operator command: %w", insert.Error)
	}
	if insert.RowsAffected != 0 {
		return false, nil
	}
	var existing commandEffectModel
	if err := tx.Where("effect_type = ? AND scope_id = ? AND command_key = ?", effect.EffectType, effect.ScopeID, effect.CommandKey).First(&existing).Error; err != nil {
		return false, fmt.Errorf("load settlement operator command winner: %w", err)
	}
	if existing.ActorID != effect.ActorID || existing.Reason != effect.Reason || existing.Outcome != effect.Outcome {
		return false, domain.ErrCommandEvidenceConflict
	}
	return true, nil
}

func (s *Store) MarkSettlementRequiresReview(ctx context.Context, settlementID uint64, reason string) error {
	if s == nil || s.db == nil {
		return errors.New("commerce database is not configured")
	}
	result := s.db.WithContext(ctx).Model(&settlementModel{}).Where("id = ?", settlementID).Updates(map[string]any{
		"status":        domain.SettlementRequiresReview,
		"review_reason": reason,
		"updated_at":    gorm.Expr("NOW()"),
	})
	if result.Error != nil {
		return fmt.Errorf("mark commerce settlement review: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return errors.New("commerce settlement not found")
	}
	return nil
}
