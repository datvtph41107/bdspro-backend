package subscriptionpostgres

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"
	catalogdomain "user/internal/domain/plan"
	subscriptiondomain "user/internal/domain/subscription"
	"user/internal/models"
	"user/internal/usecase/subscription/settlement"
)

type currentRow struct {
	ID                   uint64
	SubscriptionKey      string
	SubjectKind          string
	SubjectID            string
	ProductID            uint64
	PlanVersionID        uint64
	Status               string
	StartedAt            time.Time
	CurrentPeriodStart   time.Time
	CurrentPeriodEnd     time.Time
	AccessUntil          *time.Time
	AutoRenew            bool
	OrderReference       string
	PendingPlanVersionID *uint64
	PendingEffectiveAt   *time.Time
	CanceledAt           *time.Time
	EndedAt              *time.Time
	CreatedAt            time.Time
	UpdatedAt            time.Time
	ProductCode          string
	PlanCode             string
	PlanVersion          string
	PlanTermsChecksum    string
	PlanStatus           string
	TierRank             *int32
}

func (r currentRow) aggregate() subscriptiondomain.Aggregate {
	return subscriptiondomain.Aggregate{
		ID: r.ID, SubscriptionKey: r.SubscriptionKey,
		SubjectKind: models.SubscriptionSubjectKind(r.SubjectKind), SubjectID: r.SubjectID,
		ProductID: r.ProductID, ProductCode: r.ProductCode, PlanVersionID: r.PlanVersionID,
		PlanCode: r.PlanCode, PlanVersion: r.PlanVersion, PlanTermsChecksum: r.PlanTermsChecksum,
		PlanStatus: catalogdomain.Status(r.PlanStatus), Status: models.SubscriptionStatus(r.Status),
		StartedAt: r.StartedAt, CurrentPeriodStart: r.CurrentPeriodStart, CurrentPeriodEnd: r.CurrentPeriodEnd,
		AccessUntil: r.AccessUntil, AutoRenew: r.AutoRenew, OrderReference: r.OrderReference,
		CanceledAt: r.CanceledAt, EndedAt: r.EndedAt, CreatedAt: r.CreatedAt, UpdatedAt: r.UpdatedAt,
	}
}

func loadCurrent(db *gorm.DB, effect settlement.Effect, productID uint64) (settlement.Current, error) {
	var rows []currentRow
	err := db.Raw(`SELECT s.id, s.subscription_key, s.subject_kind, s.subject_id, s.product_id,
		s.plan_version_id, s.status, s.started_at, s.current_period_start, s.current_period_end,
		s.access_until, s.auto_renew, s.order_reference, s.pending_plan_version_id, s.pending_effective_at,
		s.canceled_at, s.ended_at, s.created_at, s.updated_at, product.code AS product_code,
		plan.code AS plan_code, pv.version AS plan_version, pv.terms_checksum AS plan_terms_checksum,
		pv.status AS plan_status, plan.tier_rank
		FROM catalog_subscriptions s
		JOIN catalog_plan_versions pv ON pv.id = s.plan_version_id
		JOIN catalog_plans plan ON plan.id = pv.plan_id
		JOIN catalog_products product ON product.id = s.product_id
		WHERE s.subject_kind = ? AND s.subject_id = ? AND s.product_id = ?
		  AND s.status IN ('pending','active','past_due')
		ORDER BY s.id ASC LIMIT 2 FOR UPDATE OF s`, string(effect.SubjectKind), effect.SubjectID, productID).Scan(&rows).Error
	if err != nil {
		return settlement.Current{}, fmt.Errorf("load current subscription for settlement: %w", err)
	}
	if len(rows) == 0 {
		return settlement.Current{}, nil
	}
	if len(rows) != 1 || rows[0].TierRank == nil {
		return settlement.Current{}, settlement.ErrCurrentState
	}
	r := rows[0]
	agg := r.aggregate()
	return settlement.Current{Aggregate: agg, TierRank: *r.TierRank, Found: true,
		HasPendingChange: r.PendingPlanVersionID != nil || r.PendingEffectiveAt != nil}, nil
}

func insertSubscription(db *gorm.DB, a subscriptiondomain.Aggregate) (uint64, error) {
	var rows []struct{ ID uint64 }
	err := db.Raw(`INSERT INTO catalog_subscriptions
		(subscription_key, subject_kind, subject_id, product_id, plan_version_id, status, started_at,
		 current_period_start, current_period_end, access_until, auto_renew, order_reference, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?) RETURNING id`,
		a.SubscriptionKey, string(a.SubjectKind), a.SubjectID, a.ProductID, a.PlanVersionID, string(a.Status),
		a.StartedAt, a.CurrentPeriodStart, a.CurrentPeriodEnd, a.AccessUntil, a.AutoRenew, a.OrderReference,
		a.CreatedAt, a.UpdatedAt).Scan(&rows).Error
	if err != nil || len(rows) != 1 {
		if err == nil {
			err = errors.New("subscription insert returned no id")
		}
		return 0, fmt.Errorf("activate subscription: %w", err)
	}
	return rows[0].ID, nil
}

func updateSubscription(db *gorm.DB, a subscriptiondomain.Aggregate) error {
	result := db.Exec(`UPDATE catalog_subscriptions SET plan_version_id = ?, order_reference = ?, updated_at = ?
		WHERE id = ? AND status = 'active'`, a.PlanVersionID, a.OrderReference, a.UpdatedAt, a.ID)
	if result.Error != nil {
		return fmt.Errorf("upgrade subscription: %w", result.Error)
	}
	if result.RowsAffected != 1 {
		return settlement.ErrCurrentState
	}
	return nil
}

func appendEvent(db *gorm.DB, effect settlement.Effect, decision settlement.Decision, subscriptionID uint64) error {
	metadata, _ := json.Marshal(map[string]any{"effectKey": effect.EffectKey, "orderId": effect.OrderID})
	eventType := "subscription.activated"
	if decision.Action == settlement.ActionUpgraded {
		eventType = "subscription.upgraded"
	}
	if err := db.Exec(`INSERT INTO catalog_subscription_events
		(subscription_id, event_type, from_plan_version_id, to_plan_version_id, effective_at, reason,
		 actor_kind, actor_id, request_id, operation_id, metadata, created_at)
		VALUES (?, ?, ?, ?, ?, ?, 'internal_service', 'payment-service', '', ?, ?::jsonb, ?)`,
		subscriptionID, eventType, decision.FromPlanVersionID, decision.ToPlanVersionID, effect.OccurredAt,
		"payment settlement applied", effect.EffectKey, string(metadata), time.Now().UTC()).Error; err != nil {
		return fmt.Errorf("append subscription settlement event: %w", err)
	}
	return nil
}
