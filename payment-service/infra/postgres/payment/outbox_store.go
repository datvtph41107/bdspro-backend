package postgres

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	domain "payment/internal/domain/payment"
	"payment/internal/usecase/outbox"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const (
	outboxPending   = "pending"
	outboxRunning   = "running"
	outboxPublished = "published"
	outboxRetry     = "retry"
)

func (s *Store) ClaimOutbox(ctx context.Context, workerID string, now time.Time, lease time.Duration) (outbox.Message, bool, error) {
	if strings.TrimSpace(workerID) == "" || now.IsZero() || lease <= 0 || s == nil || s.db == nil {
		return outbox.Message{}, false, domain.ErrInvalidCommand
	}
	var message outbox.Message
	found := false
	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var row outboxEventModel
		err := tx.Clauses(clause.Locking{Strength: "UPDATE", Options: "SKIP LOCKED"}).
			Where("((status IN ? AND available_at <= ?) OR (status = ? AND lease_until <= ?))",
				[]string{outboxPending, outboxRetry}, now, outboxRunning, now).
			Order("available_at ASC, id ASC").First(&row).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		}
		if err != nil {
			return fmt.Errorf("claim commerce outbox: %w", err)
		}
		leaseUntil := now.Add(lease)
		if err := tx.Model(&outboxEventModel{}).Where("id = ?", row.ID).Updates(map[string]any{
			"status": outboxRunning, "locked_by": workerID, "lease_until": leaseUntil,
			"claim_version": gorm.Expr("claim_version + 1"), "attempt_count": gorm.Expr("attempt_count + 1"), "updated_at": now,
		}).Error; err != nil {
			return fmt.Errorf("mark commerce outbox claimed: %w", err)
		}
		if err := tx.First(&row, row.ID).Error; err != nil {
			return fmt.Errorf("reload commerce outbox claim: %w", err)
		}
		message = outbox.Message{
			ID: row.ID, EventID: row.EventID, EventType: row.EventType, SchemaVersion: row.SchemaVersion,
			RoutingKey: row.RoutingKey, Payload: []byte(row.Payload), AttemptCount: row.AttemptCount,
			ClaimVersion: row.ClaimVersion, LockedBy: row.LockedBy,
		}
		found = true
		return nil
	})
	return message, found, err
}

func (s *Store) MarkOutboxPublished(ctx context.Context, id uint64, workerID string, claimVersion uint64, now time.Time) error {
	publishedAt := now.UTC()
	return s.finalizeOutbox(ctx, id, workerID, claimVersion, map[string]any{
		"status": outboxPublished, "locked_by": "", "lease_until": nil, "last_error": "",
		"published_at": &publishedAt, "updated_at": publishedAt,
	})
}

func (s *Store) RetryOutbox(ctx context.Context, id uint64, workerID string, claimVersion uint64, availableAt time.Time, reason string) error {
	return s.finalizeOutbox(ctx, id, workerID, claimVersion, map[string]any{
		"status": outboxRetry, "locked_by": "", "lease_until": nil, "available_at": availableAt.UTC(),
		"last_error": reason, "updated_at": availableAt.UTC(),
	})
}

func (s *Store) finalizeOutbox(ctx context.Context, id uint64, workerID string, claimVersion uint64, values map[string]any) error {
	if id == 0 || strings.TrimSpace(workerID) == "" || claimVersion == 0 || s == nil || s.db == nil {
		return domain.ErrInvalidCommand
	}
	result := s.db.WithContext(ctx).Model(&outboxEventModel{}).
		Where("id = ? AND status = ? AND locked_by = ? AND claim_version = ?", id, outboxRunning, workerID, claimVersion).
		Updates(values)
	if result.Error != nil {
		return fmt.Errorf("finalize commerce outbox: %w", result.Error)
	}
	if result.RowsAffected != 1 {
		return domain.ErrClaimLost
	}
	return nil
}
