package accesspostgres

import (
	"common/operation"
	"context"
	"errors"
	"time"
	"user/internal/domain/entitlement"
	"user/internal/models"
	useraccess "user/internal/usecase/entitlement"

	"gorm.io/gorm"
)

type subscriptionRow struct {
	ID                 uint64
	SubjectKind        string
	SubjectID          string
	PlanVersionID      uint64
	Status             string
	StartedAt          time.Time
	CurrentPeriodStart time.Time
	CurrentPeriodEnd   time.Time
	AccessUntil        *time.Time
	PlanCode           string
	PlanVersion        string
}

/**
 * SubscriptionStore đọc catalog_subscriptions và plan version từ PostgreSQL.
 */
type SubscriptionStore struct {
	db *gorm.DB
}

func NewSubscriptionStore(db *gorm.DB) *SubscriptionStore {
	return &SubscriptionStore{db: db}
}

func (s *SubscriptionStore) FindSubscriptionForOperation(
	ctx context.Context,
	subjectType string,
	subjectID string,
	operationCode operation.Code,
	now time.Time,
) (useraccess.Subscription, bool, error) {
	if s == nil || s.db == nil {
		return useraccess.Subscription{}, false, errors.New("subscription database is not configured")
	}

	// Operation is the consumer-visible business question. Product ownership stays
	// inside the catalog: a subject may have multiple current products, so picking
	// the newest subscription before considering the operation is not authoritative.
	//
	// Limit(2) is deliberate. Zero means no current subscription grants the
	// operation, one is the canonical authority, and more than one is ambiguous
	// until the business defines an explicit composition/precedence rule.
	var rows []subscriptionRow
	err := s.db.WithContext(ctx).
		Table("catalog_subscriptions AS s").
		Select(`
			s.id,
			s.subject_kind,
			s.subject_id,
			s.plan_version_id,
			s.status,
			s.started_at,
			s.current_period_start,
			s.current_period_end,
			s.access_until,
			p.code AS plan_code,
			pv.version AS plan_version
		`).
		Joins("JOIN catalog_plan_versions pv ON pv.id = s.plan_version_id").
		Joins("JOIN catalog_plans p ON p.id = pv.plan_id").
		Joins(
			"JOIN catalog_plan_operation_policies op ON op.plan_version_id = s.plan_version_id AND op.operation_code = ?",
			string(operationCode),
		).
		Where("s.subject_kind = ? AND s.subject_id = ?", subjectType, subjectID).
		Where("s.started_at <= ?", now).
		Where(
			`(s.status = ? AND COALESCE(s.access_until, s.current_period_end) > ?)
			 OR (s.status IN ? AND s.access_until IS NOT NULL AND s.access_until > ?)`,
			models.SubscriptionActive,
			now,
			[]models.SubscriptionStatus{models.SubscriptionPastDue, models.SubscriptionCanceled},
			now,
		).
		Order("s.id ASC").
		Limit(2).
		Find(&rows).Error
	if err != nil {
		return useraccess.Subscription{}, false, err
	}
	if len(rows) == 0 {
		return useraccess.Subscription{}, false, nil
	}
	if len(rows) > 1 {
		return useraccess.Subscription{}, false, useraccess.ErrAmbiguousSubscription
	}

	row := rows[0]
	subject := access.Subject{
		Type: access.SubjectType(row.SubjectKind),
		ID:   row.SubjectID,
	}

	return useraccess.Subscription{
		ID:            row.ID,
		Subject:       subject,
		PlanVersionID: row.PlanVersionID,
		PlanCode:      row.PlanCode,
		PlanVersion:   row.PlanVersion,
		Status:        row.Status,
		StartedAt:     row.StartedAt,
		PeriodStart:   row.CurrentPeriodStart,
		PeriodEnd:     row.CurrentPeriodEnd,
		AccessUntil:   row.AccessUntil,
	}, true, nil
}
