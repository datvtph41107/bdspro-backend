package admin

import (
	"fmt"

	"common/fault"
)

var (
	ErrTierRankMustBePositive = fault.Validation(
		"catalog.plan_version.tier_rank_positive",
		"tier rank must be positive",
		fault.FieldViolation{Field: "tier_rank", Description: "must be positive"},
	)
	ErrPlanVersionIDRequired = fault.Validation(
		"catalog.plan_version.id_required",
		"plan version id is required",
		fault.FieldViolation{Field: "plan_version_id", Description: "is required"},
	)
	ErrActorIDRequired = fault.Validation(
		"catalog.plan_version.actor_id_required",
		"actor id is required",
		fault.FieldViolation{Field: "actor_id", Description: "is required"},
	)
	ErrProductDisplayNameRequired = fault.Validation(
		"catalog.plan_version.product_display_name_required",
		"product display name is required",
		fault.FieldViolation{Field: "product_display_name", Description: "is required"},
	)
	ErrPageSizeOutOfRange = fault.Validation(
		"catalog.plan_version.page_size_invalid",
		"page size must be between 1 and 100",
		fault.FieldViolation{Field: "page_size", Description: "must be between 1 and 100"},
	)
)

func invalidPlanStatusFault(status string) error {
	return fault.Validation(
		"catalog.plan_version.status_invalid",
		fmt.Sprintf("invalid plan status %q", status),
		fault.FieldViolation{Field: "status", Description: "must be draft, active, or retired"},
	)
}

func invalidPlanTermsFault(cause error) error {
	if cause == nil {
		return fault.Validation("catalog.plan_version.terms_invalid", "plan terms are invalid")
	}
	return fault.Wrap(
		cause,
		fault.KindValidation,
		"catalog.plan_version.terms_invalid",
		cause.Error(),
	)
}

func subscriptionTermRequiredFault(planCode string) error {
	return fault.Validation(
		"catalog.plan_version.subscription_term_required",
		fmt.Sprintf("plan %q requires subscription_term_days", planCode),
		fault.FieldViolation{Field: "subscription_term_days", Description: "must be positive"},
	)
}
