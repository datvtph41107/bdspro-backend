package checkoutpostgres

import (
	"context"
	"errors"
	"fmt"

	"gorm.io/gorm"
	"user/internal/models"
	"user/internal/usecase/subscription/checkout"
)

type currentSubscriptionRow struct {
	ID                   uint64
	PlanVersionID        uint64
	PlanCode             string
	TierRank             *int32
	Status               string
	PendingPlanVersionID *uint64
}

func (s *CheckoutStore) CurrentForProduct(ctx context.Context, subject checkout.Subject, productID uint64) (checkout.CurrentSubscription, bool, error) {
	if s == nil || s.db == nil || !subject.IsValid() || productID == 0 {
		return checkout.CurrentSubscription{}, false, checkout.ErrInvalidCommand
	}
	var rows []currentSubscriptionRow
	err := s.db.GetDB(ctx).Table("catalog_subscriptions AS subscription").
		Select(`
			subscription.id,
			subscription.plan_version_id,
			plan.code AS plan_code,
			plan.tier_rank AS tier_rank,
			subscription.status,
			subscription.pending_plan_version_id`).
		Joins("JOIN catalog_plan_versions AS pv ON pv.id = subscription.plan_version_id").
		Joins("JOIN catalog_plans AS plan ON plan.id = pv.plan_id").
		Where("subscription.subject_kind = ? AND subscription.subject_id = ? AND subscription.product_id = ?", string(subject.Kind), subject.ID, productID).
		Where("subscription.status IN ?", []string{string(models.SubscriptionPending), string(models.SubscriptionActive), string(models.SubscriptionPastDue)}).
		Order("subscription.id ASC").
		Limit(2).
		Scan(&rows).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return checkout.CurrentSubscription{}, false, fmt.Errorf("load current subscription: %w", err)
	}
	if len(rows) == 0 {
		return checkout.CurrentSubscription{}, false, nil
	}
	if len(rows) != 1 {
		return checkout.CurrentSubscription{}, false, checkout.ErrAmbiguousCurrentSubscription
	}
	row := rows[0]
	current := checkout.CurrentSubscription{
		ID:            row.ID,
		PlanVersionID: row.PlanVersionID,
		PlanCode:      row.PlanCode,
		Status:        row.Status,
	}
	if row.TierRank != nil {
		current.TierRank = *row.TierRank
	}
	if row.PendingPlanVersionID != nil {
		current.PendingPlanVersionID = *row.PendingPlanVersionID
	}
	return current, true, nil
}
