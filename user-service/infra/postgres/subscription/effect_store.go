package subscriptionpostgres

import (
	"fmt"

	"gorm.io/gorm"
	"user/internal/modules/subscription/settlement"
)

type effectRow struct {
	EffectKey      string
	Fingerprint    string
	Action         string
	SubscriptionID *uint64
}

func loadEffect(db *gorm.DB, effectKey string) (effectRow, bool, error) {
	var rows []effectRow
	err := db.Raw(`SELECT effect_key, fingerprint, action, subscription_id
		FROM commercial_subscription_effects WHERE effect_key = ? LIMIT 1`, effectKey).Scan(&rows).Error
	if err != nil {
		return effectRow{}, false, fmt.Errorf("load subscription effect: %w", err)
	}
	if len(rows) == 0 {
		return effectRow{}, false, nil
	}
	return rows[0], true, nil
}

func loadEffectByOrder(db *gorm.DB, orderID uint64) (effectRow, bool, error) {
	var rows []effectRow
	err := db.Raw(`SELECT effect_key, fingerprint, action, subscription_id
		FROM commercial_subscription_effects WHERE order_id = ? LIMIT 1`, orderID).Scan(&rows).Error
	if err != nil {
		return effectRow{}, false, fmt.Errorf("load subscription effect by order: %w", err)
	}
	if len(rows) == 0 {
		return effectRow{}, false, nil
	}
	return rows[0], true, nil
}

func insertEffect(db *gorm.DB, input settlement.Acceptance, action settlement.Action, subscriptionID uint64) error {
	e := input.Effect
	if err := db.Exec(`INSERT INTO commercial_subscription_effects
		(effect_key, fingerprint, order_id, subject_kind, subject_id, product_code, plan_code,
		 plan_version_id, plan_version, tier_rank, subscription_term_days, terms_checksum,
		 occurred_at, action, subscription_id)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		e.EffectKey, input.Fingerprint, e.OrderID, string(e.SubjectKind), e.SubjectID, e.ProductCode, e.PlanCode,
		e.PlanVersionID, e.PlanVersion, e.TierRank, e.SubscriptionTermDays, e.TermsChecksum,
		e.OccurredAt, string(action), subscriptionID).Error; err != nil {
		return fmt.Errorf("save durable subscription settlement effect: %w", err)
	}
	return nil
}
