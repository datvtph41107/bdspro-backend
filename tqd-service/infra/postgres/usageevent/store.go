package usageeventpostgres

import (
	commonaccess "common/access"
	"common/businessoperation"
	"context"
	"errors"
	"time"
	"tqd/internal/usage"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

/**
 * Store lưu usage events trong PostgreSQL.
 */
type Store struct {
	db *gorm.DB
}

func NewStore(db *gorm.DB) *Store {
	return &Store{db: db}
}

func (s *Store) SaveUsageEvent(ctx context.Context, event usage.Event) (bool, error) {
	if s == nil || s.db == nil {
		return false, errors.New("usage database is not configured")
	}
	if !event.IsValid() {
		return false, errors.New("usage event is invalid")
	}

	row := model{
		UsageKey:       event.UsageKey,
		SubjectType:    string(event.Subject.Type),
		SubjectID:      event.Subject.ID,
		Operation:      string(event.Operation),
		OperationID:    event.OperationID,
		IdempotencyKey: event.IdempotencyKey,
		CommandKey:     event.CommandKey,
		ReservationID:  event.ReservationID,
		Amount:         event.Amount,
		PeriodStart:    event.PeriodStart,
		PeriodEnd:      event.PeriodEnd,
		CreatedAt:      event.CreatedAt,
	}

	result := s.db.WithContext(ctx).
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "usage_key"}},
			DoNothing: true,
		}).
		Create(&row)
	if result.Error != nil {
		return false, result.Error
	}

	return result.RowsAffected == 1, nil
}

func (s *Store) FindUsageEventByKey(
	ctx context.Context,
	usageKey string,
) (usage.Event, bool, error) {
	if s == nil || s.db == nil {
		return usage.Event{}, false, errors.New("usage database is not configured")
	}
	if usageKey == "" {
		return usage.Event{}, false, nil
	}

	var row model
	err := s.db.WithContext(ctx).Where("usage_key = ?", usageKey).First(&row).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return usage.Event{}, false, nil
	}
	if err != nil {
		return usage.Event{}, false, err
	}
	return usage.Event{
		UsageKey:       row.UsageKey,
		Subject:        commonaccess.Subject{Type: commonaccess.SubjectType(row.SubjectType), ID: row.SubjectID},
		Operation:      businessoperation.Code(row.Operation),
		OperationID:    row.OperationID,
		IdempotencyKey: row.IdempotencyKey,
		CommandKey:     row.CommandKey,
		ReservationID:  row.ReservationID,
		Amount:         row.Amount,
		PeriodStart:    row.PeriodStart,
		PeriodEnd:      row.PeriodEnd,
		CreatedAt:      row.CreatedAt,
	}, true, nil
}

func (s *Store) HasUsageForReservation(
	ctx context.Context,
	reservationID string,
) (bool, error) {
	if s == nil || s.db == nil {
		return false, errors.New("usage database is not configured")
	}
	if reservationID == "" {
		return false, nil
	}

	var count int64
	err := s.db.WithContext(ctx).
		Model(&model{}).
		Where("reservation_id = ?", reservationID).
		Limit(1).
		Count(&count).Error
	return count > 0, err
}

func (s *Store) SumUsage(
	ctx context.Context,
	subject commonaccess.Subject,
	operation businessoperation.Code,
	periodStart time.Time,
	periodEnd time.Time,
) (int64, error) {
	if s == nil || s.db == nil {
		return 0, errors.New("usage database is not configured")
	}

	var total int64
	err := s.db.WithContext(ctx).
		Model(&model{}).
		Select("COALESCE(SUM(amount), 0)").
		Where("subject_type = ?", subject.Type).
		Where("subject_id = ?", subject.ID).
		Where("operation = ?", operation).
		Where("period_start = ?", periodStart).
		Where("period_end = ?", periodEnd).
		Scan(&total).Error
	return total, err
}
