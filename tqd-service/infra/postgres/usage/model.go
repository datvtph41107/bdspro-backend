package postgresusage

import "time"

/**
 * model maps trực tiếp tới quota_usage_events.
 */
type model struct {
	ID uint64 `gorm:"column:id;primaryKey;autoIncrement"`

	UsageKey string `gorm:"column:usage_key;type:char(64);not null;uniqueIndex"`

	SubjectType string `gorm:"column:subject_type;type:varchar(32);not null;index:idx_quota_usage_subject_meter_period"`
	SubjectID   string `gorm:"column:subject_id;type:varchar(128);not null;index:idx_quota_usage_subject_meter_period"`
	Operation   string `gorm:"column:operation;type:varchar(128);not null"`
	MeterCode   string `gorm:"column:meter_code;type:varchar(128);not null;index:idx_quota_usage_subject_meter_period"`

	OperationID    string  `gorm:"column:operation_id;type:varchar(128);not null"`
	IdempotencyKey string  `gorm:"column:idempotency_key;type:varchar(128);not null"`
	CommandKey     string  `gorm:"column:command_key;type:varchar(128);not null"`
	ReservationID  *string `gorm:"column:reservation_id;type:varchar(128)"`

	Amount int64 `gorm:"column:amount;not null"`

	PeriodStart time.Time `gorm:"column:period_start;type:timestamptz;not null;index:idx_quota_usage_subject_meter_period"`
	PeriodEnd   time.Time `gorm:"column:period_end;type:timestamptz;not null;index:idx_quota_usage_subject_meter_period"`
	CreatedAt   time.Time `gorm:"column:created_at;type:timestamptz;not null"`

	SubscriptionID *uint64 `gorm:"column:subscription_id"`
	PlanCode       *string `gorm:"column:plan_code;type:varchar(128)"`
	PlanVersion    *string `gorm:"column:plan_version;type:varchar(64)"`
	PolicyVersion  *string `gorm:"column:policy_version;type:varchar(64)"`
	LimitSnapshot  *int64  `gorm:"column:limit_snapshot"`
}

func (model) TableName() string {
	return "quota_usage_events"
}
