package catalogdomain

import (
	"time"

	"common/productcatalog"
)

/** PlanVersionAggregate is one durable catalog aggregate. */
type PlanVersionAggregate struct {
	ID                 uint64
	ProductID          uint64
	PlanID             uint64
	ProductCode        string
	ProductDisplayName string
	ProductStatus      productcatalog.Status
	PlanCode           string
	PlanStatus         productcatalog.Status
	Version            string
	DisplayName        string
	Status             productcatalog.Status
	SubjectScope       productcatalog.SubjectScope
	EffectiveFrom      *time.Time
	EffectiveUntil     *time.Time
	PublishedAt        *time.Time
	TermsChecksum      string
	Entitlements       []productcatalog.Entitlement
	Prices             []productcatalog.PriceItem
}

/** Contract returns a defensive commercial contract. */
func (a PlanVersionAggregate) Contract(
	status productcatalog.Status,
	effectiveFrom time.Time,
	effectiveUntil *time.Time,
) productcatalog.PlanVersion {
	return productcatalog.PlanVersion{
		ProductCode:    a.ProductCode,
		PlanCode:       a.PlanCode,
		Version:        a.Version,
		DisplayName:    a.DisplayName,
		Status:         status,
		SubjectScope:   a.SubjectScope,
		EffectiveFrom:  effectiveFrom,
		EffectiveUntil: copyTime(effectiveUntil),
		Entitlements:   append([]productcatalog.Entitlement(nil), a.Entitlements...),
		Prices:         append([]productcatalog.PriceItem(nil), a.Prices...),
	}
}

/** CurrentContract returns the aggregate's persisted contract. */
func (a PlanVersionAggregate) CurrentContract() productcatalog.PlanVersion {
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
	a.Entitlements = append([]productcatalog.Entitlement(nil), a.Entitlements...)
	a.Prices = append([]productcatalog.PriceItem(nil), a.Prices...)
	return a
}

func copyTime(value *time.Time) *time.Time {
	if value == nil {
		return nil
	}
	copied := *value
	return &copied
}
