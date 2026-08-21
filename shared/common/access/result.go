package access

import (
	"common/codeformat"
	commonmetering "common/metering"
	"common/operation"
	"strings"
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
	Operation operation.Code
	Allowed   bool
	Metering  Metering

	Unlimited bool
	Limit     int64
	Period    Period

	PeriodStart time.Time
	PeriodEnd   time.Time

	SubscriptionID uint64
	PlanCode       string
	PlanVersion    string
}

/** Description: Commercial evidence attached to an entitlement decision. */
type Metering struct {
	FeatureCode    string
	MeterCode      commonmetering.Code
	UnitsPerAction int64
	PolicyVersion  string
}

/** Description: Reports whether this decision consumes a metered allowance. */
func (m Metering) IsMetered() bool {
	return m.MeterCode != ""
}

/** Description: Validates the policy evidence required by an allowed action. */
func (m Metering) IsValid() bool {
	if !codeformat.IsValid(strings.TrimSpace(m.FeatureCode)) {
		return false
	}
	if strings.TrimSpace(m.PolicyVersion) == "" {
		return false
	}
	if !m.IsMetered() {
		return m.UnitsPerAction == 0
	}
	if !m.MeterCode.IsValid() {
		return false
	}
	return m.UnitsPerAction > 0
}

/**
 * UsesQuota cho biết operation này có quota cần kiểm tra runtime hay không.
 */
func (r Result) UsesQuota() bool {
	if !r.Allowed || r.Unlimited || !r.Metering.IsMetered() {
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
	if !r.Metering.IsValid() {
		return false
	}
	if !r.Metering.IsMetered() {
		return !r.Unlimited && r.Limit == 0 && r.Period == PeriodNone
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
