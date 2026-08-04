package models

import (
	"encoding/json"
	"time"

	"common/productcatalog"
)

/** CatalogProduct stores one stable product family. */
type CatalogProduct struct {
	ID          uint64                `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Code        string                `gorm:"column:code;type:varchar(128);not null;uniqueIndex" json:"code"`
	DisplayName string                `gorm:"column:display_name;type:varchar(255);not null" json:"displayName"`
	Status      productcatalog.Status `gorm:"column:status;type:varchar(20);not null" json:"status"`
	CreatedBy   *uint64               `gorm:"column:created_by" json:"createdBy,omitempty"`
	UpdatedBy   *uint64               `gorm:"column:updated_by" json:"updatedBy,omitempty"`
	CreatedAt   time.Time             `gorm:"column:created_at;not null;autoCreateTime" json:"createdAt"`
	UpdatedAt   time.Time             `gorm:"column:updated_at;not null;autoUpdateTime" json:"updatedAt"`
}

/** TableName returns the durable catalog product table. */
func (CatalogProduct) TableName() string {
	return "catalog_products"
}

/** CatalogPlan stores one stable commercial plan identity. */
type CatalogPlan struct {
	ID        uint64                `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	ProductID uint64                `gorm:"column:product_id;not null;index" json:"productId"`
	Code      string                `gorm:"column:code;type:varchar(128);not null;uniqueIndex" json:"code"`
	Status    productcatalog.Status `gorm:"column:status;type:varchar(20);not null" json:"status"`
	CreatedBy *uint64               `gorm:"column:created_by" json:"createdBy,omitempty"`
	UpdatedBy *uint64               `gorm:"column:updated_by" json:"updatedBy,omitempty"`
	CreatedAt time.Time             `gorm:"column:created_at;not null;autoCreateTime" json:"createdAt"`
	UpdatedAt time.Time             `gorm:"column:updated_at;not null;autoUpdateTime" json:"updatedAt"`
}

/** TableName returns the durable catalog plan table. */
func (CatalogPlan) TableName() string {
	return "catalog_plans"
}

/** CatalogPlanVersion stores immutable commercial terms. */
type CatalogPlanVersion struct {
	ID             uint64                      `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	PlanID         uint64                      `gorm:"column:plan_id;not null;index" json:"planId"`
	Version        string                      `gorm:"column:version;type:varchar(32);not null" json:"version"`
	DisplayName    string                      `gorm:"column:display_name;type:varchar(255);not null" json:"displayName"`
	Status         productcatalog.Status       `gorm:"column:status;type:varchar(20);not null" json:"status"`
	SubjectScope   productcatalog.SubjectScope `gorm:"column:subject_scope;type:varchar(32);not null" json:"subjectScope"`
	EffectiveFrom  *time.Time                  `gorm:"column:effective_from;type:timestamptz" json:"effectiveFrom,omitempty"`
	EffectiveUntil *time.Time                  `gorm:"column:effective_until;type:timestamptz" json:"effectiveUntil,omitempty"`
	PublishedAt    *time.Time                  `gorm:"column:published_at;type:timestamptz" json:"publishedAt,omitempty"`
	Source         CatalogSource               `gorm:"column:source;type:varchar(32);not null" json:"source"`
	TermsChecksum  string                      `gorm:"column:terms_checksum;type:char(64);not null" json:"termsChecksum"`
	CreatedBy      *uint64                     `gorm:"column:created_by" json:"createdBy,omitempty"`
	UpdatedBy      *uint64                     `gorm:"column:updated_by" json:"updatedBy,omitempty"`
	PublishedBy    *uint64                     `gorm:"column:published_by" json:"publishedBy,omitempty"`
	CreatedAt      time.Time                   `gorm:"column:created_at;not null;autoCreateTime" json:"createdAt"`
	UpdatedAt      time.Time                   `gorm:"column:updated_at;not null;autoUpdateTime" json:"updatedAt"`
}

/** TableName returns the durable plan-version table. */
func (CatalogPlanVersion) TableName() string {
	return "catalog_plan_versions"
}

/** CatalogPlanEntitlement stores one versioned grant. */
type CatalogPlanEntitlement struct {
	ID            uint64                         `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	PlanVersionID uint64                         `gorm:"column:plan_version_id;not null;index" json:"planVersionId"`
	Code          string                         `gorm:"column:code;type:varchar(128);not null" json:"code"`
	Kind          productcatalog.EntitlementKind `gorm:"column:kind;type:varchar(32);not null" json:"kind"`
	FeatureCode   *string                        `gorm:"column:feature_code;type:varchar(128)" json:"featureCode,omitempty"`
	MeterCode     *string                        `gorm:"column:meter_code;type:varchar(128)" json:"meterCode,omitempty"`
	Amount        int64                          `gorm:"column:amount;not null;default:0" json:"amount"`
	Unlimited     bool                           `gorm:"column:unlimited;not null;default:false" json:"unlimited"`
	Period        productcatalog.PeriodKind      `gorm:"column:period;type:varchar(32);not null" json:"period"`
	CreatedBy     *uint64                        `gorm:"column:created_by" json:"createdBy,omitempty"`
	CreatedAt     time.Time                      `gorm:"column:created_at;not null;autoCreateTime" json:"createdAt"`
}

/** TableName returns the durable entitlement table. */
func (CatalogPlanEntitlement) TableName() string {
	return "catalog_plan_entitlements"
}

/** CatalogPriceItem stores one versioned commercial price. */
type CatalogPriceItem struct {
	ID            uint64                   `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	PlanVersionID uint64                   `gorm:"column:plan_version_id;not null;index" json:"planVersionId"`
	Code          string                   `gorm:"column:code;type:varchar(128);not null" json:"code"`
	Kind          productcatalog.PriceKind `gorm:"column:kind;type:varchar(20);not null" json:"kind"`
	Currency      string                   `gorm:"column:currency;type:char(3);not null" json:"currency"`
	AmountMinor   int64                    `gorm:"column:amount_minor;not null" json:"amountMinor"`
	BillingUnit   string                   `gorm:"column:billing_unit;type:varchar(128);not null" json:"billingUnit"`
	MeterCode     *string                  `gorm:"column:meter_code;type:varchar(128)" json:"meterCode,omitempty"`
	Quantity      int64                    `gorm:"column:quantity;not null" json:"quantity"`
	CreatedBy     *uint64                  `gorm:"column:created_by" json:"createdBy,omitempty"`
	CreatedAt     time.Time                `gorm:"column:created_at;not null;autoCreateTime" json:"createdAt"`
}

/** TableName returns the durable price-item table. */
func (CatalogPriceItem) TableName() string {
	return "catalog_price_items"
}

type CatalogSource string

const (
	CatalogSourceNative         CatalogSource = "native"
	CatalogSourceLegacyBackfill CatalogSource = "legacy_backfill"
	CatalogSourceImport         CatalogSource = "import"
)

type SubscriptionSubjectKind string

const (
	SubscriptionSubjectProfile      SubscriptionSubjectKind = "profile"
	SubscriptionSubjectOrganization SubscriptionSubjectKind = "organization"
)

type SubscriptionStatus string

const (
	SubscriptionPending  SubscriptionStatus = "pending"
	SubscriptionActive   SubscriptionStatus = "active"
	SubscriptionPastDue  SubscriptionStatus = "past_due"
	SubscriptionCanceled SubscriptionStatus = "canceled"
	SubscriptionExpired  SubscriptionStatus = "expired"
)

/** CatalogSubscription binds one subject to one plan version. */
type CatalogSubscription struct {
	ID                 uint64                  `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	SubscriptionKey    string                  `gorm:"column:subscription_key;type:varchar(128);not null;uniqueIndex" json:"subscriptionKey"`
	SubjectKind        SubscriptionSubjectKind `gorm:"column:subject_kind;type:varchar(32);not null;index:idx_catalog_subscription_subject" json:"subjectKind"`
	SubjectID          string                  `gorm:"column:subject_id;type:varchar(128);not null;index:idx_catalog_subscription_subject" json:"subjectId"`
	ProductID          uint64                  `gorm:"column:product_id;not null;index" json:"productId"`
	PlanVersionID      uint64                  `gorm:"column:plan_version_id;not null;index" json:"planVersionId"`
	Status             SubscriptionStatus      `gorm:"column:status;type:varchar(20);not null" json:"status"`
	StartedAt          time.Time               `gorm:"column:started_at;type:timestamptz;not null" json:"startedAt"`
	CurrentPeriodStart time.Time               `gorm:"column:current_period_start;type:timestamptz;not null" json:"currentPeriodStart"`
	CurrentPeriodEnd   time.Time               `gorm:"column:current_period_end;type:timestamptz;not null" json:"currentPeriodEnd"`
	AccessUntil        *time.Time              `gorm:"column:access_until;type:timestamptz" json:"accessUntil,omitempty"`
	AutoRenew          bool                    `gorm:"column:auto_renew;not null;default:false" json:"autoRenew"`
	OrderReference     string                  `gorm:"column:order_reference;type:varchar(128)" json:"orderReference,omitempty"`
	CanceledAt         *time.Time              `gorm:"column:canceled_at;type:timestamptz" json:"canceledAt,omitempty"`
	EndedAt            *time.Time              `gorm:"column:ended_at;type:timestamptz" json:"endedAt,omitempty"`
	CreatedBy          *uint64                 `gorm:"column:created_by" json:"createdBy,omitempty"`
	UpdatedBy          *uint64                 `gorm:"column:updated_by" json:"updatedBy,omitempty"`
	CreatedAt          time.Time               `gorm:"column:created_at;not null;autoCreateTime" json:"createdAt"`
	UpdatedAt          time.Time               `gorm:"column:updated_at;not null;autoUpdateTime" json:"updatedAt"`
}

/** TableName returns the durable subscription table. */
func (CatalogSubscription) TableName() string {
	return "catalog_subscriptions"
}

/** CatalogSubscriptionEvent records append-only lifecycle history. */
type CatalogSubscriptionEvent struct {
	ID                uint64          `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	SubscriptionID    uint64          `gorm:"column:subscription_id;not null;index" json:"subscriptionId"`
	EventType         string          `gorm:"column:event_type;type:varchar(64);not null" json:"eventType"`
	FromPlanVersionID *uint64         `gorm:"column:from_plan_version_id" json:"fromPlanVersionId,omitempty"`
	ToPlanVersionID   *uint64         `gorm:"column:to_plan_version_id" json:"toPlanVersionId,omitempty"`
	EffectiveAt       time.Time       `gorm:"column:effective_at;type:timestamptz;not null" json:"effectiveAt"`
	Reason            string          `gorm:"column:reason;type:text;not null" json:"reason"`
	ActorKind         string          `gorm:"column:actor_kind;type:varchar(32);not null" json:"actorKind"`
	ActorID           string          `gorm:"column:actor_id;type:varchar(128);not null" json:"actorId"`
	RequestID         string          `gorm:"column:request_id;type:varchar(128);not null" json:"requestId"`
	OperationID       string          `gorm:"column:operation_id;type:varchar(128);not null" json:"operationId"`
	Metadata          json.RawMessage `gorm:"column:metadata;type:jsonb;not null" json:"metadata"`
	CreatedAt         time.Time       `gorm:"column:created_at;not null;autoCreateTime" json:"createdAt"`
}

/** TableName returns the append-only subscription event table. */
func (CatalogSubscriptionEvent) TableName() string {
	return "catalog_subscription_events"
}

type CatalogMigrationStatus string

const (
	CatalogMigrationPendingReview CatalogMigrationStatus = "pending_review"
	CatalogMigrationReviewed      CatalogMigrationStatus = "reviewed"
	CatalogMigrationDraftCreated  CatalogMigrationStatus = "draft_created"
	CatalogMigrationShadowRead    CatalogMigrationStatus = "shadow_read"
	CatalogMigrationCutover       CatalogMigrationStatus = "cutover"
	CatalogMigrationRolledBack    CatalogMigrationStatus = "rolled_back"
	CatalogMigrationBlocked       CatalogMigrationStatus = "blocked"
)

/** LegacyPlanCatalogMapping tracks plan backfill and rollback. */
type LegacyPlanCatalogMapping struct {
	LegacyPlanID   uint64                 `gorm:"column:legacy_plan_id;primaryKey" json:"legacyPlanId"`
	PlanVersionID  *uint64                `gorm:"column:plan_version_id;uniqueIndex" json:"planVersionId,omitempty"`
	SourceChecksum string                 `gorm:"column:source_checksum;type:char(64);not null" json:"sourceChecksum"`
	Status         CatalogMigrationStatus `gorm:"column:status;type:varchar(32);not null" json:"status"`
	Decision       json.RawMessage        `gorm:"column:decision;type:jsonb;not null" json:"decision"`
	ReviewedBy     *uint64                `gorm:"column:reviewed_by" json:"reviewedBy,omitempty"`
	ReviewedAt     *time.Time             `gorm:"column:reviewed_at;type:timestamptz" json:"reviewedAt,omitempty"`
	MigratedAt     *time.Time             `gorm:"column:migrated_at;type:timestamptz" json:"migratedAt,omitempty"`
	CreatedAt      time.Time              `gorm:"column:created_at;not null;autoCreateTime" json:"createdAt"`
	UpdatedAt      time.Time              `gorm:"column:updated_at;not null;autoUpdateTime" json:"updatedAt"`
}

/** TableName returns the legacy plan mapping table. */
func (LegacyPlanCatalogMapping) TableName() string {
	return "legacy_plan_catalog_mappings"
}

/** LegacyMembershipSubscriptionMapping tracks membership cutover. */
type LegacyMembershipSubscriptionMapping struct {
	LegacyMembershipID       uint64                 `gorm:"column:legacy_membership_id;primaryKey" json:"legacyMembershipId"`
	SubscriptionID           *uint64                `gorm:"column:subscription_id;uniqueIndex" json:"subscriptionId,omitempty"`
	SourceChecksum           string                 `gorm:"column:source_checksum;type:char(64);not null" json:"sourceChecksum"`
	Status                   CatalogMigrationStatus `gorm:"column:status;type:varchar(32);not null" json:"status"`
	Decision                 json.RawMessage        `gorm:"column:decision;type:jsonb;not null" json:"decision"`
	ReviewedBy               *uint64                `gorm:"column:reviewed_by" json:"reviewedBy,omitempty"`
	ReviewedAt               *time.Time             `gorm:"column:reviewed_at;type:timestamptz" json:"reviewedAt,omitempty"`
	MigratedAt               *time.Time             `gorm:"column:migrated_at;type:timestamptz" json:"migratedAt,omitempty"`
	BlockedReason            string                 `gorm:"column:blocked_reason;type:text;not null" json:"blockedReason,omitempty"`
	LastRunID                string                 `gorm:"column:last_run_id;type:varchar(128);not null" json:"lastRunId,omitempty"`
	LastError                string                 `gorm:"column:last_error;type:text;not null" json:"lastError,omitempty"`
	CutoverEvidenceChecksum  string                 `gorm:"column:cutover_evidence_checksum;type:char(64);not null" json:"cutoverEvidenceChecksum,omitempty"`
	CutoverAt                *time.Time             `gorm:"column:cutover_at;type:timestamptz" json:"cutoverAt,omitempty"`
	CutoverBy                string                 `gorm:"column:cutover_by;type:varchar(128);not null" json:"cutoverBy,omitempty"`
	RollbackEvidenceChecksum string                 `gorm:"column:rollback_evidence_checksum;type:char(64);not null" json:"rollbackEvidenceChecksum,omitempty"`
	RolledBackAt             *time.Time             `gorm:"column:rolled_back_at;type:timestamptz" json:"rolledBackAt,omitempty"`
	RolledBackBy             string                 `gorm:"column:rolled_back_by;type:varchar(128);not null" json:"rolledBackBy,omitempty"`
	CreatedAt                time.Time              `gorm:"column:created_at;not null;autoCreateTime" json:"createdAt"`
	UpdatedAt                time.Time              `gorm:"column:updated_at;not null;autoUpdateTime" json:"updatedAt"`
}

/** TableName returns the legacy membership mapping table. */
func (LegacyMembershipSubscriptionMapping) TableName() string {
	return "legacy_membership_subscription_mappings"
}
