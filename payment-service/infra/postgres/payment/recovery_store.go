package postgres

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	domain "payment/internal/domain/payment"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const redriveEffectType = "payment.fulfillment.redrive"

func (s *Store) RedriveFulfillment(ctx context.Context, command domain.RedriveCommand, now time.Time) (changed bool, replay bool, err error) {
	if command.FulfillmentID == 0 || strings.TrimSpace(command.CommandKey) == "" || strings.TrimSpace(command.ActorID) == "" || now.IsZero() {
		return false, false, domain.ErrInvalidCommand
	}
	if s == nil || s.db == nil {
		return false, false, errors.New("commerce database is not configured")
	}
	err = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var fulfillment fulfillmentModel
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&fulfillment, command.FulfillmentID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return domain.ErrFulfillmentNotFound
			}
			return fmt.Errorf("lock fulfillment for redrive: %w", err)
		}

		effect := commandEffectModel{
			EffectType: redriveEffectType,
			ScopeID:    command.FulfillmentID,
			CommandKey: command.CommandKey,
			ActorID:    command.ActorID,
			Reason:     command.Reason,
			Outcome:    "accepted",
			CreatedAt:  now,
		}
		insert := tx.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "effect_type"}, {Name: "scope_id"}, {Name: "command_key"}},
			DoNothing: true,
		}).Create(&effect)
		if insert.Error != nil {
			return fmt.Errorf("consume fulfillment redrive command: %w", insert.Error)
		}
		if insert.RowsAffected == 0 {
			replay = true
			return nil
		}

		if fulfillment.Status != string(domain.FulfillmentRequiresReview) {
			if err := tx.Model(&commandEffectModel{}).Where("id = ?", effect.ID).Update("outcome", "noop").Error; err != nil {
				return fmt.Errorf("record fulfillment redrive no-op: %w", err)
			}
			return nil
		}

		result := tx.Model(&fulfillmentModel{}).Where("id = ? AND status = ?", fulfillment.ID, domain.FulfillmentRequiresReview).Updates(map[string]any{
			"status":       domain.FulfillmentPending,
			"available_at": now,
			"locked_by":    "",
			"lease_until":  nil,
			"last_error":   "",
			"updated_at":   now,
		})
		if result.Error != nil {
			return fmt.Errorf("redrive commerce fulfillment: %w", result.Error)
		}
		if result.RowsAffected != 1 {
			return domain.ErrRedriveNotAllowed
		}
		if err := tx.Model(&commandEffectModel{}).Where("id = ?", effect.ID).Update("outcome", "scheduled").Error; err != nil {
			return fmt.Errorf("record fulfillment redrive outcome: %w", err)
		}
		changed = true
		return nil
	})
	return changed, replay, err
}
