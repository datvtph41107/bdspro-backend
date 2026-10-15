package plan

import "time"

type Status string

const (
	StatusDraft   Status = "draft"
	StatusActive  Status = "active"
	StatusRetired Status = "retired"
)

type SubjectScope string

const (
	SubjectScopeProfile      SubjectScope = "profile"
	SubjectScopeOrganization SubjectScope = "organization"
	SubjectScopeAny          SubjectScope = "any"
)

type EntitlementKind string

const (
	EntitlementFeatureAccess  EntitlementKind = "feature_access"
	EntitlementUsageAllowance EntitlementKind = "usage_allowance"
	EntitlementCapacityLimit  EntitlementKind = "capacity_limit"
)

type PeriodKind string

const (
	PeriodNone              PeriodKind = "none"
	PeriodDay               PeriodKind = "day"
	PeriodCalendarMonth     PeriodKind = "calendar_month"
	PeriodSubscriptionCycle PeriodKind = "subscription_cycle"
	PeriodLifetime          PeriodKind = "lifetime"
)

type PriceKind string

const (
	PriceRecurring PriceKind = "recurring"
	PriceAddOn     PriceKind = "add_on"
	PriceOverage   PriceKind = "overage"
)

/** Snapshot is an immutable product contract snapshot. */
type Snapshot struct {
	Version  string
	Products []Product
}

/** Product groups stable plan families. */
type Product struct {
	Code        string
	DisplayName string
	Plans       []PlanVersion
}

/** PlanVersion freezes commercial terms for one effective period. */
type PlanVersion struct {
	ProductCode          string
	PlanCode             string
	Version              string
	DisplayName          string
	Status               Status
	SubjectScope         SubjectScope
	SubscriptionTermDays int32
	EffectiveFrom        time.Time
	EffectiveUntil       *time.Time
	Entitlements         []Entitlement
	Operations           []OperationBinding
	Prices               []PriceItem
}

/** Entitlement grants access, allowance, or capacity. */
type Entitlement struct {
	Code        string
	Kind        EntitlementKind
	FeatureCode string
	MeterCode   string
	Amount      int64
	Unlimited   bool
	Period      PeriodKind
}

/** OperationBinding binds one business operation to this version's commercial terms. */
type OperationBinding struct {
	Code           string
	FeatureCode    string
	MeterCode      string
	UnitsPerAction int64
}

/** PriceItem describes a recurring, add-on, or overage price. */
type PriceItem struct {
	Code        string
	Kind        PriceKind
	Currency    string
	AmountMinor int64
	BillingUnit string
	MeterCode   string
	Quantity    int64
}
