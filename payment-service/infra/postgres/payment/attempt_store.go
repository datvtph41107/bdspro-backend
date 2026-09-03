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

func (s *Store) FindOrderByID(ctx context.Context, orderID uint64) (domain.Order, bool, error) {
	if orderID == 0 || s == nil || s.db == nil {
		return domain.Order{}, false, domain.ErrInvalidCommand
	}
	var row orderModel
	err := s.db.WithContext(ctx).First(&row, orderID).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return domain.Order{}, false, nil
	}
	if err != nil {
		return domain.Order{}, false, fmt.Errorf("find commerce order by id: %w", err)
	}
	return orderFromModel(row), true, nil
}

func (s *Store) FindAttemptByCommand(ctx context.Context, orderID uint64, commandKey string) (domain.PaymentAttempt, bool, error) {
	if orderID == 0 || strings.TrimSpace(commandKey) == "" || s == nil || s.db == nil {
		return domain.PaymentAttempt{}, false, domain.ErrInvalidCommand
	}
	var row paymentAttemptModel
	err := s.db.WithContext(ctx).
		Where("order_id = ? AND command_key = ?", orderID, strings.TrimSpace(commandKey)).
		First(&row).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return domain.PaymentAttempt{}, false, nil
	}
	if err != nil {
		return domain.PaymentAttempt{}, false, fmt.Errorf("find payment attempt by command: %w", err)
	}
	return paymentAttemptFromModel(row), true, nil
}

func (s *Store) CreateAttempt(ctx context.Context, attempt domain.PaymentAttempt) (domain.PaymentAttempt, bool, error) {
	if !attempt.IsValid() || s == nil || s.db == nil {
		return domain.PaymentAttempt{}, false, domain.ErrInvalidCommand
	}
	row := paymentAttemptToModel(attempt)
	created := false
	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var order orderModel
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&order, attempt.OrderID).Error; err != nil {
			return fmt.Errorf("lock commerce order for attempt: %w", err)
		}
		if order.Status != string(domain.OrderPendingFunds) {
			return domain.ErrAttemptNotAllowed
		}
		result := tx.Clauses(clause.OnConflict{DoNothing: true}).Create(&row)
		if result.Error != nil {
			return fmt.Errorf("create payment attempt: %w", result.Error)
		}
		created = result.RowsAffected == 1
		if !created {
			if err := tx.Where("order_id = ? AND command_key = ?", attempt.OrderID, attempt.CommandKey).First(&row).Error; err != nil {
				return fmt.Errorf("reload payment attempt replay: %w", err)
			}
			if row.CommandFingerprint != attempt.CommandFingerprint {
				return domain.ErrAttemptCommandConflict
			}
		}
		return nil
	})
	if err != nil {
		return domain.PaymentAttempt{}, false, err
	}
	return paymentAttemptFromModel(row), created, nil
}

func (s *Store) UpdateAttemptFromProvider(ctx context.Context, attemptID uint64, provider domain.ProviderAttempt, now time.Time) (domain.PaymentAttempt, error) {
	if attemptID == 0 || !provider.IsValid() || now.IsZero() || s == nil || s.db == nil {
		return domain.PaymentAttempt{}, domain.ErrInvalidCommand
	}
	values := map[string]any{
		"provider_reference": provider.ProviderReference,
		"status":             provider.Status,
		"next_action_kind":   provider.NextAction.Kind,
		"redirect_url":       provider.NextAction.RedirectURL,
		"qr_payload":         provider.NextAction.QRPayload,
		"action_expires_at":  provider.NextAction.ExpiresAt,
		"expires_at":         provider.ExpiresAt,
		"updated_at":         now.UTC(),
	}
	result := s.db.WithContext(ctx).Model(&paymentAttemptModel{}).Where("id = ?", attemptID).Updates(values)
	if result.Error != nil {
		return domain.PaymentAttempt{}, fmt.Errorf("update payment attempt provider state: %w", result.Error)
	}
	if result.RowsAffected != 1 {
		return domain.PaymentAttempt{}, domain.ErrAttemptNotFound
	}
	var row paymentAttemptModel
	if err := s.db.WithContext(ctx).First(&row, attemptID).Error; err != nil {
		return domain.PaymentAttempt{}, fmt.Errorf("reload payment attempt: %w", err)
	}
	return paymentAttemptFromModel(row), nil
}
