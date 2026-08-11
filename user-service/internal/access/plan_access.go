package access

import commonaccess "common/access"

/**
 * PlanAccess là quyền của một operation trong một plan version.
 */
type PlanAccess struct {
	Allowed   bool
	Unlimited bool
	Limit     int64
	Period    commonaccess.Period
}
