package admin

import "common/fault"

var ErrTierRankMustBePositive = fault.Validation(
	"catalog.plan_version.tier_rank_positive",
	"tier rank must be positive",
	fault.FieldViolation{Field: "tier_rank", Description: "must be positive"},
)
