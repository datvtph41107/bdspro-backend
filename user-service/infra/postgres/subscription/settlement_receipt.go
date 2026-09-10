package subscriptionpostgres

import (
	"fmt"

	"gorm.io/gorm"
	"user/internal/usecase/subscription/settlement"
)

type settlementReceiptRow struct {
	EffectKey      string
	Fingerprint    string
	Action         string
	SubscriptionID *uint64
}

func loadSettlementReceipt(db *gorm.DB, effectKey string) (settlementReceiptRow, bool, error) {
	var rows []settlementReceiptRow
	err := db.Raw(`SELECT effect_key, fingerprint, action, subscription_id
		FROM subscription_settlement_receipts WHERE effect_key = ? LIMIT 1`, effectKey).Scan(&rows).Error
	if err != nil {
		return settlementReceiptRow{}, false, fmt.Errorf("load subscription settlement receipt: %w", err)
	}
	if len(rows) == 0 {
		return settlementReceiptRow{}, false, nil
	}
	return rows[0], true, nil
}

func insertSettlementReceipt(db *gorm.DB, input settlement.Acceptance, action settlement.Action, subscriptionID uint64) error {
	e := input.Effect
	if err := db.Exec(`INSERT INTO subscription_settlement_receipts
		(effect_key, fingerprint, order_id, subject_kind, subject_id, product_code, plan_code,
		 plan_version_id, plan_version, tier_rank, subscription_term_days, terms_checksum,
		 occurred_at, action, subscription_id)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		e.EffectKey, input.Fingerprint, e.OrderID, string(e.SubjectKind), e.SubjectID, e.ProductCode, e.PlanCode,
		e.PlanVersionID, e.PlanVersion, e.TierRank, e.SubscriptionTermDays, e.TermsChecksum,
		e.OccurredAt, string(action), subscriptionID).Error; err != nil {
		return fmt.Errorf("save durable subscription settlement receipt: %w", err)
	}
	return nil
}
