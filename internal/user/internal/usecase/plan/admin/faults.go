package admin

import (
	"fmt"

	_errors "common/errors"
	"user/internal"
)

var (
	ErrTierRankMustBePositive = _errors.ReturnError(
		service.PlanTierRankMustBePositive,
		_errors.WithViolations(_errors.FieldViolation{Field: "tier_rank", Description: "must be positive"}),
	)
	ErrPlanVersionIDRequired = _errors.ReturnError(
		service.PlanVersionIDRequired,
		_errors.WithViolations(_errors.FieldViolation{Field: "plan_version_id", Description: "is required"}),
	)
	ErrActorIDRequired = _errors.ReturnError(
		service.PlanActorIDRequired,
		_errors.WithViolations(_errors.FieldViolation{Field: "actor_id", Description: "is required"}),
	)
	ErrProductDisplayNameRequired = _errors.ReturnError(
		service.PlanProductDisplayNameRequired,
		_errors.WithViolations(_errors.FieldViolation{Field: "product_display_name", Description: "is required"}),
	)
	ErrPageSizeOutOfRange = _errors.ReturnError(
		service.PlanPageSizeOutOfRange,
		_errors.WithViolations(_errors.FieldViolation{Field: "page_size", Description: "must be between 1 and 100"}),
	)
)

func invalidPlanStatusFault(status string) error {
	return _errors.ReturnError(
		service.PlanStatusInvalid,
		_errors.WithPublicMessage(fmt.Sprintf("invalid plan status %q", status)),
		_errors.WithViolations(_errors.FieldViolation{Field: "status", Description: "must be draft, active, or retired"}),
	)
}

func invalidPlanTermsFault(cause error) error {
	if cause == nil {
		return _errors.ReturnError(service.PlanTermsInvalid)
	}
	return _errors.ReturnError(
		service.PlanTermsInvalid,
		_errors.WithCause(cause),
		_errors.WithPublicMessage(cause.Error()),
	)
}

func subscriptionTermRequiredFault(planCode string) error {
	return _errors.ReturnError(
		service.PlanSubscriptionTermRequired,
		_errors.WithPublicMessage(fmt.Sprintf("plan %q requires subscription_term_days", planCode)),
		_errors.WithViolations(_errors.FieldViolation{Field: "subscription_term_days", Description: "must be positive"}),
	)
}
