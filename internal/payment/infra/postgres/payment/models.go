package postgres

import "time"

type orderModel struct {
	ID                   uint64     `gorm:"column:id;primaryKey"`
	SubjectKind          string     `gorm:"column:subject_kind"`
	SubjectID            string     `gorm:"column:subject_id"`
	ProductCode          string     `gorm:"column:product_code"`
	PlanCode             string     `gorm:"column:plan_code"`
	PlanVersionID        uint64     `gorm:"column:plan_version_id"`
	PlanVersion          string     `gorm:"column:plan_version"`
	TierRank             int32      `gorm:"column:tier_rank"`
	SubscriptionTermDays int32      `gorm:"column:subscription_term_days"`
	TermsChecksum        string     `gorm:"column:terms_checksum"`
	Currency             string     `gorm:"column:currency"`
	AmountMinor          int64      `gorm:"column:amount_minor"`
	CommandKey           string     `gorm:"column:command_key"`
	CommandFingerprint   string     `gorm:"column:command_fingerprint"`
	Reference            string     `gorm:"column:reference"`
	Status               string     `gorm:"column:status"`
	ExpiresAt            time.Time  `gorm:"column:expires_at"`
	FundsConfirmedAt     *time.Time `gorm:"column:funds_confirmed_at"`
	CreatedAt            time.Time  `gorm:"column:created_at"`
	UpdatedAt            time.Time  `gorm:"column:updated_at"`
}

func (orderModel) TableName() string { return "commerce_orders" }

type paymentAttemptModel struct {
	ID                 uint64     `gorm:"column:id;primaryKey"`
	OrderID            uint64     `gorm:"column:order_id"`
	CommandKey         string     `gorm:"column:command_key"`
	CommandFingerprint string     `gorm:"column:command_fingerprint"`
	Method             string     `gorm:"column:method"`
	Provider           string     `gorm:"column:provider"`
	ProviderReference  string     `gorm:"column:provider_reference"`
	Status             string     `gorm:"column:status"`
	NextActionKind     string     `gorm:"column:next_action_kind"`
	RedirectURL        string     `gorm:"column:redirect_url"`
	QRPayload          string     `gorm:"column:qr_payload"`
	ActionExpiresAt    *time.Time `gorm:"column:action_expires_at"`
	ExpiresAt          *time.Time `gorm:"column:expires_at"`
	FailureCode        string     `gorm:"column:failure_code"`
	FailureDetail      string     `gorm:"column:failure_detail"`
	CreatedAt          time.Time  `gorm:"column:created_at"`
	UpdatedAt          time.Time  `gorm:"column:updated_at"`
}

func (paymentAttemptModel) TableName() string { return "commerce_payment_attempts" }

type settlementModel struct {
	ID                    uint64    `gorm:"column:id;primaryKey"`
	OrderID               *uint64   `gorm:"column:order_id"`
	Provider              string    `gorm:"column:provider"`
	ProviderTransactionID string    `gorm:"column:provider_transaction_id"`
	Reference             string    `gorm:"column:reference"`
	Currency              string    `gorm:"column:currency"`
	AmountMinor           int64     `gorm:"column:amount_minor"`
	OccurredAt            time.Time `gorm:"column:occurred_at"`
	EvidenceHash          string    `gorm:"column:evidence_hash"`
	Status                string    `gorm:"column:status"`
	ReviewReason          string    `gorm:"column:review_reason"`
	CreatedAt             time.Time `gorm:"column:created_at"`
	UpdatedAt             time.Time `gorm:"column:updated_at"`
}

func (settlementModel) TableName() string { return "commerce_settlements" }

type fulfillmentModel struct {
	ID           uint64     `gorm:"column:id;primaryKey"`
	OrderID      uint64     `gorm:"column:order_id"`
	Status       string     `gorm:"column:status"`
	AttemptCount int        `gorm:"column:attempt_count"`
	AvailableAt  time.Time  `gorm:"column:available_at"`
	LockedBy     string     `gorm:"column:locked_by"`
	ClaimVersion uint64     `gorm:"column:claim_version"`
	LeaseUntil   *time.Time `gorm:"column:lease_until"`
	LastError    string     `gorm:"column:last_error"`
	CompletedAt  *time.Time `gorm:"column:completed_at"`
	CreatedAt    time.Time  `gorm:"column:created_at"`
	UpdatedAt    time.Time  `gorm:"column:updated_at"`
}

func (fulfillmentModel) TableName() string { return "commerce_fulfillments" }

type commandEffectModel struct {
	ID         uint64    `gorm:"column:id;primaryKey"`
	EffectType string    `gorm:"column:effect_type"`
	ScopeID    uint64    `gorm:"column:scope_id"`
	CommandKey string    `gorm:"column:command_key"`
	ActorID    string    `gorm:"column:actor_id"`
	Reason     string    `gorm:"column:reason"`
	Outcome    string    `gorm:"column:outcome"`
	CreatedAt  time.Time `gorm:"column:created_at"`
}

func (commandEffectModel) TableName() string { return "commerce_command_effects" }

type outboxEventModel struct {
	ID            uint64 `gorm:"column:id;primaryKey"`
	EventID       string `gorm:"column:event_id"`
	EventType     string `gorm:"column:event_type"`
	SchemaVersion int    `gorm:"column:schema_version"`
	RoutingKey    string `gorm:"column:routing_key"`
	// Keep JSONB payloads as text at the SQL boundary. A []byte value is bound by
	// PostgreSQL drivers as bytea, which JSONB correctly rejects. The outbox
	// use case still exposes []byte because RabbitMQ publishes an opaque body.
	Payload      string     `gorm:"column:payload;type:jsonb"`
	Status       string     `gorm:"column:status"`
	AttemptCount int        `gorm:"column:attempt_count"`
	AvailableAt  time.Time  `gorm:"column:available_at"`
	LockedBy     string     `gorm:"column:locked_by"`
	ClaimVersion uint64     `gorm:"column:claim_version"`
	LeaseUntil   *time.Time `gorm:"column:lease_until"`
	LastError    string     `gorm:"column:last_error"`
	PublishedAt  *time.Time `gorm:"column:published_at"`
	CreatedAt    time.Time  `gorm:"column:created_at"`
	UpdatedAt    time.Time  `gorm:"column:updated_at"`
}

func (outboxEventModel) TableName() string { return "commerce_outbox_events" }
