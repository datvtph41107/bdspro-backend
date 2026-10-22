package subscriptionpostgres

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"time"

	subscriptiondomain "mdn4/modules/user/domain/subscription"
	cancelapp "mdn4/modules/user/services/subscription/cancel"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// @bind: mdn4/modules/user/services/subscription/cancel.Store
type CancellationStore struct{ db *gorm.DB }

func NewCancellationStore(db *gorm.DB) *CancellationStore { return &CancellationStore{db: db} }

func (s *CancellationStore) CancelCurrent(ctx context.Context, command cancelapp.Command) (cancelapp.Result, error) {
	if s == nil || s.db == nil {
		return cancelapp.Result{}, cancelapp.ErrUnavailable
	}
	var output cancelapp.Result
	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var candidates []struct{ ID uint64 }
		if err := tx.Table("catalog_subscriptions AS s").
			Select("s.id").
			Joins("JOIN catalog_products product ON product.id = s.product_id").
			Where("s.subject_kind = ? AND s.subject_id = ?", "profile", strconv.FormatUint(command.ProfileID, 10)).
			Where("product.code = ?", command.ProductCode).
			Where("s.status IN ?", []string{"active", "past_due"}).
			Order("s.id ASC").Limit(2).
			Clauses(clause.Locking{Strength: "UPDATE", Table: clause.Table{Name: "s"}}).
			Scan(&candidates).Error; err != nil {
			return fmt.Errorf("lock current subscription: %w", err)
		}
		if len(candidates) > 1 {
			return cancelapp.ErrAmbiguous
		}
		if len(candidates) == 0 {
			// A response may be lost after commit. Returning the latest canceled
			// aggregate makes an immediate retry observable and side-effect free.
			var latest []subscriptionProjectionRow
			query := projectionBase(tx.Table("catalog_subscriptions AS s")).
				Where("s.subject_kind = ? AND s.subject_id = ?", "profile", strconv.FormatUint(command.ProfileID, 10)).
				Where("product.code = ? AND s.status = ? AND s.canceled_at IS NOT NULL", command.ProductCode, "canceled")
			if err := projectionSelect(query).Order("s.canceled_at DESC, s.id DESC").Limit(1).Scan(&latest).Error; err != nil {
				return fmt.Errorf("load canceled subscription: %w", err)
			}
			if len(latest) == 0 {
				return cancelapp.ErrNotFound
			}
			output = cancelapp.Result{Subscription: latest[0].domain(), Changed: false}
			return nil
		}

		id := candidates[0].ID
		result := tx.Exec(`UPDATE catalog_subscriptions
			SET status = 'canceled', auto_renew = FALSE, access_until = ?, canceled_at = ?, ended_at = ?,
				pending_plan_version_id = NULL, pending_effective_at = NULL, updated_by = ?, updated_at = ?
			WHERE id = ? AND status IN ('active','past_due')`,
			command.OccurredAt, command.OccurredAt, command.OccurredAt, command.ProfileID, command.OccurredAt, id)
		if result.Error != nil {
			return fmt.Errorf("cancel subscription: %w", result.Error)
		}
		if result.RowsAffected != 1 {
			return cancelapp.ErrNotFound
		}

		metadata, _ := json.Marshal(map[string]any{"mode": "immediate", "productCode": command.ProductCode})
		if err := tx.Exec(`INSERT INTO catalog_subscription_events
			(subscription_id, event_type, from_plan_version_id, to_plan_version_id, effective_at, reason,
			 actor_kind, actor_id, request_id, operation_id, metadata, created_at)
			SELECT id, 'subscription.canceled', plan_version_id, NULL, ?, ?, 'profile', ?, ?, ?, ?::jsonb, ?
			FROM catalog_subscriptions WHERE id = ?`,
			command.OccurredAt, command.Reason, strconv.FormatUint(command.ProfileID, 10),
			command.RequestID, command.OperationID, string(metadata), command.OccurredAt, id).Error; err != nil {
			return fmt.Errorf("append cancellation event: %w", err)
		}

		var rows []subscriptionProjectionRow
		if err := projectionSelect(projectionBase(tx.Table("catalog_subscriptions AS s")).Where("s.id = ?", id)).Scan(&rows).Error; err != nil {
			return fmt.Errorf("load canceled projection: %w", err)
		}
		if len(rows) != 1 {
			return errors.New("canceled subscription projection was not found")
		}
		output = cancelapp.Result{Subscription: rows[0].domain(), Changed: true}
		return nil
	})
	return output, err
}

// Named scan shape keeps the store transaction independent from GORM model
// mutation while reusing the canonical Admin projection columns.
type subscriptionProjectionRow struct {
	SubscriptionID       uint64
	SubjectKind          string
	SubjectID            string
	ProductCode          string
	ProductDisplayName   string
	PlanCode             string
	PlanDisplayName      string
	PlanVersionID        uint64
	PlanVersion          string
	Status               string
	StartedAt            time.Time
	CurrentPeriodStart   time.Time
	CurrentPeriodEnd     time.Time
	AccessUntil          *time.Time
	PendingPlanVersionID *uint64
	PendingEffectiveAt   *time.Time
	OrderReference       string
	CanceledAt           *time.Time
	EndedAt              *time.Time
	AutoRenew            bool
	UpdatedAt            time.Time
}

func (row subscriptionProjectionRow) domain() subscriptiondomain.AdminProjection {
	return subscriptiondomain.AdminProjection{
		SubscriptionID: row.SubscriptionID, SubjectKind: row.SubjectKind, SubjectID: row.SubjectID,
		ProductCode: row.ProductCode, ProductDisplayName: row.ProductDisplayName,
		PlanCode: row.PlanCode, PlanDisplayName: row.PlanDisplayName,
		PlanVersionID: row.PlanVersionID, PlanVersion: row.PlanVersion, Status: row.Status,
		StartedAt: row.StartedAt, CurrentPeriodStart: row.CurrentPeriodStart, CurrentPeriodEnd: row.CurrentPeriodEnd,
		AccessUntil: row.AccessUntil, PendingPlanVersionID: row.PendingPlanVersionID,
		PendingEffectiveAt: row.PendingEffectiveAt, OrderReference: row.OrderReference,
		CanceledAt: row.CanceledAt, EndedAt: row.EndedAt, AutoRenew: row.AutoRenew, UpdatedAt: row.UpdatedAt,
	}
}

var _ cancelapp.Store = (*CancellationStore)(nil)
