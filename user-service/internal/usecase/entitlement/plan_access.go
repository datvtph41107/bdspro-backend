package evaluate

import "user/internal/domain/entitlement"

/**
 * PlanAccess là quyền của một operation trong một plan version.
 */
type PlanAccess struct {
	Allowed        bool
	FeatureCode    string
	MeterCode      string
	UnitsPerAction int64
	Unlimited      bool
	Limit          int64
	Period         access.Period
}
