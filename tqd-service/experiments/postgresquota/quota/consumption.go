package quota

import (
	commonaccess "common/access"
	commonmetering "common/metering"
	"errors"
	"time"
)

var (
	ErrExceeded       = errors.New("quota exceeded")
	ErrPeriodConflict = errors.New("quota period conflicts with existing pool")
)

// Consumption is the resolved allowance mutation for one accepted command.
// Plan/version is intentionally absent: those are evidence, not pool identity.
type Consumption struct {
	Subject commonaccess.Subject
	Meter   commonmetering.Code
	Amount  int64
	Limit   int64

	PeriodStart time.Time
	PeriodEnd   time.Time
}

func (c Consumption) IsValid() bool {
	return c.Subject.IsValid() &&
		c.Meter.IsValid() &&
		c.Amount > 0 &&
		c.Limit >= 0 &&
		c.Amount <= c.Limit &&
		!c.PeriodStart.IsZero() &&
		c.PeriodEnd.After(c.PeriodStart)
}
