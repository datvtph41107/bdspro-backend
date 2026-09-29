package plan

import (
	"time"
)

/** PlanVersionAggregate is one durable catalog aggregate. */
type PlanVersionAggregate struct {
	ID                   uint64
	ProductID            uint64
	PlanID               uint64
	ProductCode          string
	ProductDisplayName   string
	ProductStatus        Status
	PlanCode             string
	TierRank             int32
	PlanStatus           Status
	Version              string
	DisplayName          string
	Status               Status
	SubjectScope         SubjectScope
	SubscriptionTermDays int32
	EffectiveFrom        *time.Time
	EffectiveUntil       *time.Time
	PublishedAt          *time.Time
	TermsChecksum        string
	Entitlements         []Entitlement
	Operations           []OperationBinding
	Prices               []PriceItem
}

/** Contract returns a defensive commercial contract. */
func (a PlanVersionAggregate) Contract(
	status Status,
	effectiveFrom time.Time,
	effectiveUntil *time.Time,
) PlanVersion {
	return PlanVersion{
		ProductCode:          a.ProductCode,
		PlanCode:             a.PlanCode,
		Version:              a.Version,
		DisplayName:          a.DisplayName,
		Status:               status,
		SubjectScope:         a.SubjectScope,
		SubscriptionTermDays: a.SubscriptionTermDays,
		EffectiveFrom:        effectiveFrom,
		EffectiveUntil:       copyTime(effectiveUntil),
		Entitlements:         append([]Entitlement(nil), a.Entitlements...),
		Operations:           append([]OperationBinding(nil), a.Operations...),
		Prices:               append([]PriceItem(nil), a.Prices...),
	}
}

/** CurrentContract returns the aggregate's persisted contract. */
func (a PlanVersionAggregate) CurrentContract() PlanVersion {
	effectiveFrom := time.Time{}
	if a.EffectiveFrom != nil {
		effectiveFrom = *a.EffectiveFrom
	}
	return a.Contract(a.Status, effectiveFrom, a.EffectiveUntil)
}

/** Copy returns a defensive aggregate snapshot. */
func (a PlanVersionAggregate) Copy() PlanVersionAggregate {
	a.EffectiveFrom = copyTime(a.EffectiveFrom)
	a.EffectiveUntil = copyTime(a.EffectiveUntil)
	a.PublishedAt = copyTime(a.PublishedAt)
	a.Entitlements = append([]Entitlement(nil), a.Entitlements...)
	a.Operations = append([]OperationBinding(nil), a.Operations...)
	a.Prices = append([]PriceItem(nil), a.Prices...)
	return a
}

func copyTime(value *time.Time) *time.Time {
	if value == nil {
		return nil
	}
	copied := *value
	return &copied
}
