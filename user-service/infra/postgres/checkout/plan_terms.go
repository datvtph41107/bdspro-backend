package checkoutpostgres

import (
	"context"
	"fmt"
	"time"

	plan "user/internal/domain/plan"
	"user/internal/models"
	"user/internal/usecase/subscription/checkout"
)

type checkoutTermsRow struct {
	ProductID            uint64
	ProductCode          string
	PlanID               uint64
	PlanCode             string
	PlanVersionID        uint64
	PlanVersion          string
	TierRank             *int32
	SubscriptionTermDays *int32
	TermsChecksum        string
	SubjectScope         plan.SubjectScope
}

func (s *CheckoutStore) ResolveCheckoutTerms(ctx context.Context, planCode string, at time.Time) (checkout.PlanTerms, error) {
	if s == nil || s.db == nil || planCode == "" || at.IsZero() {
		return checkout.PlanTerms{}, checkout.ErrPlanTermsUnavailable
	}
	db := s.db.GetDB(ctx)
	var candidates []checkoutTermsRow
	err := db.Table("catalog_plan_versions AS pv").
		Select(`
			product.id AS product_id,
			product.code AS product_code,
			plan.id AS plan_id,
			plan.code AS plan_code,
			pv.id AS plan_version_id,
			pv.version AS plan_version,
			plan.tier_rank AS tier_rank,
			pv.subscription_term_days AS subscription_term_days,
			pv.terms_checksum AS terms_checksum,
			pv.subject_scope AS subject_scope`).
		Joins("JOIN catalog_plans AS plan ON plan.id = pv.plan_id").
		Joins("JOIN catalog_products AS product ON product.id = plan.product_id").
		Where("plan.code = ?", planCode).
		Where("product.status = ? AND plan.status = ? AND pv.status = ?", plan.StatusActive, plan.StatusActive, plan.StatusActive).
		Where("pv.published_at IS NOT NULL AND pv.effective_from <= ?", at).
		Where("(pv.effective_until IS NULL OR ? < pv.effective_until)", at).
		Order("pv.id ASC").
		Limit(2).
		Scan(&candidates).Error
	if err != nil {
		return checkout.PlanTerms{}, fmt.Errorf("resolve checkout plan version: %w", err)
	}
	if len(candidates) == 0 {
		return checkout.PlanTerms{}, checkout.ErrPlanTermsUnavailable
	}
	if len(candidates) != 1 {
		return checkout.PlanTerms{}, checkout.ErrAmbiguousPlanVersion
	}
	candidate := candidates[0]
	if candidate.TierRank == nil || candidate.SubscriptionTermDays == nil ||
		*candidate.TierRank <= 0 || *candidate.SubscriptionTermDays <= 0 {
		return checkout.PlanTerms{}, checkout.ErrPlanTermsUnavailable
	}

	var prices []models.CatalogPriceItem
	if err := db.Where("plan_version_id = ? AND kind = ?", candidate.PlanVersionID, plan.PriceRecurring).
		Order("id ASC").Limit(2).Find(&prices).Error; err != nil {
		return checkout.PlanTerms{}, fmt.Errorf("resolve checkout recurring price: %w", err)
	}
	if len(prices) != 1 {
		return checkout.PlanTerms{}, checkout.ErrAmbiguousRecurringPrice
	}
	price := prices[0]
	terms := checkout.PlanTerms{
		ProductID:            candidate.ProductID,
		ProductCode:          candidate.ProductCode,
		PlanID:               candidate.PlanID,
		PlanCode:             candidate.PlanCode,
		PlanVersionID:        candidate.PlanVersionID,
		PlanVersion:          candidate.PlanVersion,
		TierRank:             *candidate.TierRank,
		SubscriptionTermDays: *candidate.SubscriptionTermDays,
		TermsChecksum:        candidate.TermsChecksum,
		SubjectScope:         candidate.SubjectScope,
		Price: checkout.Money{
			Currency:    price.Currency,
			AmountMinor: price.AmountMinor,
		},
	}
	if !terms.IsValid() {
		return checkout.PlanTerms{}, checkout.ErrPlanTermsUnavailable
	}
	return terms, nil
}
