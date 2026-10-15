package subscriptionpostgres

import "time"

// settlementReceiptSchema chỉ mô tả schema của durable settlement receipt.
// Logic ghi/đọc vẫn dùng SQL tường minh trong settlement_receipt.go.
type settlementReceiptSchema struct {
	EffectKey            string    `gorm:"column:effect_key;size:128;primaryKey"`
	Fingerprint          string    `gorm:"column:fingerprint;type:char(64);not null"`
	OrderID              uint64    `gorm:"column:order_id;uniqueIndex;not null"`
	SubjectKind          string    `gorm:"column:subject_kind;size:32;not null"`
	SubjectID            string    `gorm:"column:subject_id;size:128;not null"`
	ProductCode          string    `gorm:"column:product_code;size:64;not null"`
	PlanCode             string    `gorm:"column:plan_code;size:64;not null"`
	PlanVersionID        uint64    `gorm:"column:plan_version_id;not null"`
	PlanVersion          string    `gorm:"column:plan_version;size:64;not null"`
	TierRank             int32     `gorm:"column:tier_rank;not null"`
	SubscriptionTermDays int32     `gorm:"column:subscription_term_days;not null"`
	TermsChecksum        string    `gorm:"column:terms_checksum;type:char(64);not null"`
	OccurredAt           time.Time `gorm:"column:occurred_at;not null"`
	Action               string    `gorm:"column:action;size:32;not null"`
	SubscriptionID       *uint64   `gorm:"column:subscription_id"`
	CreatedAt            time.Time `gorm:"column:created_at;autoCreateTime"`
}

func (settlementReceiptSchema) TableName() string {
	return "subscription_settlement_receipts"
}

// AutoMigrateModels đưa private schema về registry canonical của User.
func AutoMigrateModels() []any {
	return []any{&settlementReceiptSchema{}}
}
