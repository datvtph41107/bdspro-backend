package subscriptionpostgres

import (
	"fmt"
	"time"

	"gorm.io/gorm"
	catalogdomain "user/internal/domain/plan"
	"user/internal/usecase/subscription/settlement"
)

type targetRow struct {
	ProductID            uint64
	PlanID               uint64
	PlanVersionID        uint64
	ProductCode          string
	PlanCode             string
	PlanVersion          string
	TierRank             *int32
	SubscriptionTermDays *int32
	TermsChecksum        string
	PublishedAt          *time.Time
	PlanStatus           string
}

func loadTarget(db *gorm.DB, planVersionID uint64) (settlement.Target, error) {
	var rows []targetRow
	err := db.Raw(`SELECT product.id AS product_id, plan.id AS plan_id, pv.id AS plan_version_id,
		product.code AS product_code, plan.code AS plan_code, pv.version AS plan_version,
		plan.tier_rank, pv.subscription_term_days, pv.terms_checksum, pv.published_at, pv.status AS plan_status
		FROM catalog_plan_versions pv
		JOIN catalog_plans plan ON plan.id = pv.plan_id
		JOIN catalog_products product ON product.id = plan.product_id
		WHERE pv.id = ? LIMIT 1`, planVersionID).Scan(&rows).Error
	if err != nil {
		return settlement.Target{}, fmt.Errorf("load settlement target: %w", err)
	}
	if len(rows) != 1 || rows[0].TierRank == nil || rows[0].SubscriptionTermDays == nil {
		return settlement.Target{}, settlement.ErrCatalogTermsConflict
	}
	r := rows[0]
	return settlement.Target{ProductID: r.ProductID, PlanID: r.PlanID, PlanVersionID: r.PlanVersionID,
		ProductCode: r.ProductCode, PlanCode: r.PlanCode, PlanVersion: r.PlanVersion,
		TierRank: *r.TierRank, SubscriptionTermDays: *r.SubscriptionTermDays,
		TermsChecksum: r.TermsChecksum, Published: r.PublishedAt != nil, PlanStatus: catalogdomain.Status(r.PlanStatus)}, nil
}
