package access

import (
	"common/businessoperation"
	"time"
)

/**
 * Period cho biết quota được tính lại theo khoảng thời gian nào.
 */
type Period string

const (
	PeriodNone              Period = "none"
	PeriodDay               Period = "day"
	PeriodCalendarMonth     Period = "calendar_month"
	PeriodSubscriptionCycle Period = "subscription_cycle"
	PeriodLifetime          Period = "lifetime"
)

/**
 * Result là kết quả access mà service nghiệp vụ cần dùng.
 *
 * G4 dừng ở đây. Result không chứa used/reserved vì phần đó thuộc quota runtime.
 */
type Result struct {
	Subject   Subject
	Operation businessoperation.Code
	Allowed   bool

	Unlimited bool
	Limit     int64
	Period    Period

	PeriodStart time.Time
	PeriodEnd   time.Time

	SubscriptionID uint64
	PlanCode       string
	PlanVersion    string
}

/**
 * UsesQuota cho biết operation này có quota cần kiểm tra runtime hay không.
 */
func (r Result) UsesQuota() bool {
	if !r.Allowed || r.Unlimited {
		return false
	}
	return r.Limit > 0 && r.Period != PeriodNone
}

/**
 * IsValid kiểm tra các field cần thiết trước khi Result đi vào quota service.
 */
func (r Result) IsValid() bool {
	if !r.Subject.IsValid() || !r.Operation.IsValid() {
		return false
	}

	if !r.Allowed {
		return true
	}

	if r.Unlimited {
		return r.Limit == 0
	}

	if r.Limit < 0 {
		return false
	}

	if r.Limit == 0 {
		return r.Period == PeriodNone
	}

	if r.Period == PeriodNone {
		return false
	}

	if r.Period != PeriodLifetime {
		if r.PeriodStart.IsZero() || r.PeriodEnd.IsZero() || !r.PeriodEnd.After(r.PeriodStart) {
			return false
		}
	}

	return true
}
