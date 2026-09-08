package eventing

import (
	"context"
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	domain "notification/internal/domain/delivery"
)

type DeliveryStore struct{ db *gorm.DB }

func NewDeliveryStore(db *gorm.DB) *DeliveryStore { return &DeliveryStore{db: db} }

type deliveryRow struct {
	ID           uint64     `gorm:"column:id"`
	EventID      string     `gorm:"column:event_id"`
	SubjectKind  string     `gorm:"column:subject_kind"`
	SubjectID    string     `gorm:"column:subject_id"`
	Channel      string     `gorm:"column:channel"`
	TemplateCode string     `gorm:"column:template_code"`
	Payload      string     `gorm:"column:payload"`
	Status       string     `gorm:"column:status"`
	AttemptCount int        `gorm:"column:attempt_count"`
	AvailableAt  time.Time  `gorm:"column:available_at"`
	LockedBy     string     `gorm:"column:locked_by"`
	ClaimVersion uint64     `gorm:"column:claim_version"`
	LeaseUntil   *time.Time `gorm:"column:lease_until"`
	LastError    string     `gorm:"column:last_error"`
	SentAt       *time.Time `gorm:"column:sent_at"`
	CreatedAt    time.Time  `gorm:"column:created_at;type:timestamptz;not null;default:now()"`
	UpdatedAt    time.Time  `gorm:"column:updated_at;type:timestamptz;not null;default:now()"`
}

func (deliveryRow) TableName() string { return "notification_delivery_intents" }

func (s *DeliveryStore) ClaimNext(ctx context.Context, workerID string, now time.Time, lease time.Duration) (*domain.Intent, error) {
	if s == nil || s.db == nil {
		return nil, errors.New("delivery store database is nil")
	}
	var claimed deliveryRow
	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var row deliveryRow
		err := tx.Clauses(clause.Locking{Strength: "UPDATE", Options: "SKIP LOCKED"}).
			Where("channel = ? AND status IN ? AND available_at <= ? AND (lease_until IS NULL OR lease_until <= ?)", "push", []string{"pending", "retry", "unknown"}, now, now).
			Order("available_at ASC, id ASC").First(&row).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		}
		if err != nil {
			return err
		}
		newVersion := row.ClaimVersion + 1
		leaseUntil := now.Add(lease)
		result := tx.Model(&deliveryRow{}).Where("id = ? AND claim_version = ?", row.ID, row.ClaimVersion).Updates(map[string]any{
			"status": "running", "locked_by": workerID, "claim_version": newVersion, "lease_until": leaseUntil,
			"attempt_count": gorm.Expr("attempt_count + 1"), "updated_at": now,
		})
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected != 1 {
			return fmt.Errorf("delivery claim lost")
		}
		row.Status = "running"
		row.LockedBy = workerID
		row.ClaimVersion = newVersion
		row.LeaseUntil = &leaseUntil
		row.AttemptCount++
		claimed = row
		return nil
	})
	if err != nil {
		return nil, err
	}
	if claimed.ID == 0 {
		return nil, nil
	}
	return toIntent(claimed), nil
}

func (s *DeliveryStore) MarkSent(ctx context.Context, in domain.Intent, at time.Time) error {
	return s.finish(ctx, in, map[string]any{"status": "sent", "sent_at": at, "last_error": "", "lease_until": nil, "locked_by": "", "updated_at": at})
}
func (s *DeliveryStore) MarkRetry(ctx context.Context, in domain.Intent, at time.Time, cause string) error {
	return s.finish(ctx, in, map[string]any{"status": "retry", "available_at": at, "last_error": cause, "lease_until": nil, "locked_by": "", "updated_at": time.Now().UTC()})
}
func (s *DeliveryStore) MarkUnknown(ctx context.Context, in domain.Intent, at time.Time, cause string) error {
	return s.finish(ctx, in, map[string]any{"status": "unknown", "available_at": at, "last_error": cause, "lease_until": nil, "locked_by": "", "updated_at": time.Now().UTC()})
}
func (s *DeliveryStore) MarkFailed(ctx context.Context, in domain.Intent, cause string) error {
	return s.finish(ctx, in, map[string]any{"status": "failed", "last_error": cause, "lease_until": nil, "locked_by": "", "updated_at": time.Now().UTC()})
}
func (s *DeliveryStore) finish(ctx context.Context, in domain.Intent, changes map[string]any) error {
	result := s.db.WithContext(ctx).Model(&deliveryRow{}).Where("id = ? AND status = 'running' AND locked_by = ? AND claim_version = ?", in.ID, in.LockedBy, in.ClaimVersion).Updates(changes)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected != 1 {
		return fmt.Errorf("delivery claim lost")
	}
	return nil
}
func toIntent(r deliveryRow) *domain.Intent {
	return &domain.Intent{ID: r.ID, EventID: r.EventID, SubjectKind: r.SubjectKind, SubjectID: r.SubjectID, Channel: r.Channel, TemplateCode: r.TemplateCode, Payload: r.Payload, Status: domain.Status(r.Status), AttemptCount: r.AttemptCount, AvailableAt: r.AvailableAt, LockedBy: r.LockedBy, ClaimVersion: r.ClaimVersion, LeaseUntil: r.LeaseUntil, LastError: r.LastError, SentAt: r.SentAt}
}
