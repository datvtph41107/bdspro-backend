package postgres

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"common/events/paymentcompleted"
	domain "payment/internal/domain/payment"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func (s *Store) ClaimFulfillment(ctx context.Context, workerID string, now time.Time, lease time.Duration) (domain.FulfillmentClaim, bool, error) {
	if strings.TrimSpace(workerID) == "" || now.IsZero() || lease <= 0 {
		return domain.FulfillmentClaim{}, false, domain.ErrInvalidCommand
	}
	if s == nil || s.db == nil {
		return domain.FulfillmentClaim{}, false, errors.New("commerce database is not configured")
	}
	var claim domain.FulfillmentClaim
	found := false
	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var row fulfillmentModel
		err := tx.Clauses(clause.Locking{Strength: "UPDATE", Options: "SKIP LOCKED"}).
			Where("((status IN ? AND available_at <= ?) OR (status = ? AND lease_until <= ?))",
				[]string{string(domain.FulfillmentPending), string(domain.FulfillmentRetry)}, now,
				string(domain.FulfillmentRunning), now).
			Order("available_at ASC, id ASC").
			First(&row).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		}
		if err != nil {
			return fmt.Errorf("claim commerce fulfillment: %w", err)
		}

		leaseUntil := now.Add(lease)
		if err := tx.Model(&fulfillmentModel{}).Where("id = ?", row.ID).Updates(map[string]any{
			"status":        domain.FulfillmentRunning,
			"locked_by":     workerID,
			"claim_version": gorm.Expr("claim_version + 1"),
			"lease_until":   leaseUntil,
			"attempt_count": gorm.Expr("attempt_count + 1"),
			"updated_at":    now,
		}).Error; err != nil {
			return fmt.Errorf("mark commerce fulfillment claimed: %w", err)
		}
		if err := tx.First(&row, row.ID).Error; err != nil {
			return fmt.Errorf("reload commerce fulfillment claim: %w", err)
		}
		var order orderModel
		if err := tx.First(&order, row.OrderID).Error; err != nil {
			return fmt.Errorf("load commerce order for fulfillment: %w", err)
		}
		claim = domain.FulfillmentClaim{Fulfillment: fulfillmentFromModel(row), Order: orderFromModel(order)}
		found = true
		return nil
	})
	if err != nil {
		return domain.FulfillmentClaim{}, false, err
	}
	return claim, found, nil
}

func (s *Store) CompleteFulfillment(ctx context.Context, fulfillmentID uint64, workerID string, claimVersion uint64, now time.Time) error {
	if fulfillmentID == 0 || strings.TrimSpace(workerID) == "" || claimVersion == 0 || now.IsZero() {
		return domain.ErrInvalidCommand
	}
	if s == nil || s.db == nil {
		return errors.New("commerce database is not configured")
	}
	completedAt := now.UTC()
	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var fulfillment fulfillmentModel
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("id = ? AND status = ? AND locked_by = ? AND claim_version = ?", fulfillmentID, domain.FulfillmentRunning, workerID, claimVersion).
			First(&fulfillment).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return domain.ErrClaimLost
			}
			return fmt.Errorf("lock commerce fulfillment completion: %w", err)
		}
		var order orderModel
		if err := tx.First(&order, fulfillment.OrderID).Error; err != nil {
			return fmt.Errorf("load commerce order for completion event: %w", err)
		}
		if order.FundsConfirmedAt == nil || order.FundsConfirmedAt.IsZero() {
			return domain.ErrInvalidCommand
		}
		result := tx.Model(&fulfillmentModel{}).
			Where("id = ? AND status = ? AND locked_by = ? AND claim_version = ?", fulfillmentID, domain.FulfillmentRunning, workerID, claimVersion).
			Updates(map[string]any{
				"status": domain.FulfillmentCompleted, "locked_by": "", "lease_until": nil,
				"last_error": "", "completed_at": &completedAt, "updated_at": completedAt,
			})
		if result.Error != nil {
			return fmt.Errorf("complete commerce fulfillment: %w", result.Error)
		}
		if result.RowsAffected != 1 {
			return domain.ErrClaimLost
		}

		domainOrder := orderFromModel(order)
		event := paymentcompleted.V1{
			EventID: domainOrder.CompletedEventID(), OrderID: order.ID,
			SubjectKind: order.SubjectKind, SubjectID: order.SubjectID,
			ProductCode: order.ProductCode, PlanCode: order.PlanCode,
			PlanVersionID: order.PlanVersionID, PlanVersion: order.PlanVersion,
			Currency: order.Currency, AmountMinor: order.AmountMinor,
			FundsConfirmedAt: order.FundsConfirmedAt.UTC(), CompletedAt: completedAt,
		}
		payload, err := json.Marshal(event)
		if err != nil {
			return fmt.Errorf("marshal payment completed event: %w", err)
		}
		outbox := outboxEventModel{
			EventID: event.EventID, EventType: paymentcompleted.EventTypeV1,
			SchemaVersion: paymentcompleted.SchemaVersion, RoutingKey: paymentcompleted.EventTypeV1,
			Payload: string(payload), Status: outboxPending, AvailableAt: completedAt,
			CreatedAt: completedAt, UpdatedAt: completedAt,
		}
		if err := tx.Clauses(clause.OnConflict{DoNothing: true}).Create(&outbox).Error; err != nil {
			return fmt.Errorf("append payment completed outbox: %w", err)
		}
		return nil
	})
}

func (s *Store) RetryFulfillment(ctx context.Context, fulfillmentID uint64, workerID string, claimVersion uint64, availableAt time.Time, reason string) error {
	return s.finalizeClaim(ctx, fulfillmentID, workerID, claimVersion, map[string]any{
		"status":       domain.FulfillmentRetry,
		"locked_by":    "",
		"lease_until":  nil,
		"available_at": availableAt,
		"last_error":   reason,
		"updated_at":   availableAt,
	})
}

func (s *Store) RequireReview(ctx context.Context, fulfillmentID uint64, workerID string, claimVersion uint64, now time.Time, reason string) error {
	return s.finalizeClaim(ctx, fulfillmentID, workerID, claimVersion, map[string]any{
		"status":      domain.FulfillmentRequiresReview,
		"locked_by":   "",
		"lease_until": nil,
		"last_error":  reason,
		"updated_at":  now,
	})
}

func (s *Store) finalizeClaim(ctx context.Context, fulfillmentID uint64, workerID string, claimVersion uint64, values map[string]any) error {
	if fulfillmentID == 0 || strings.TrimSpace(workerID) == "" || claimVersion == 0 {
		return domain.ErrInvalidCommand
	}
	if s == nil || s.db == nil {
		return errors.New("commerce database is not configured")
	}
	result := s.db.WithContext(ctx).Model(&fulfillmentModel{}).
		Where("id = ? AND status = ? AND locked_by = ? AND claim_version = ?", fulfillmentID, domain.FulfillmentRunning, workerID, claimVersion).
		Updates(values)
	if result.Error != nil {
		return fmt.Errorf("finalize commerce fulfillment claim: %w", result.Error)
	}
	if result.RowsAffected != 1 {
		return domain.ErrClaimLost
	}
	return nil
}
