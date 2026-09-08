package checkout

import "errors"

var (
	ErrInvalidCommand               = errors.New("invalid checkout command")
	ErrCheckoutCommandConflict      = errors.New("checkout command conflict")
	ErrPlanTermsUnavailable         = errors.New("plan terms unavailable")
	ErrAmbiguousPlanVersion         = errors.New("multiple effective plan versions")
	ErrAmbiguousRecurringPrice      = errors.New("recurring price is missing or ambiguous")
	ErrAmbiguousCurrentSubscription = errors.New("multiple current subscriptions for subject/product")
	ErrSubjectScope                 = errors.New("plan does not allow checkout subject")
	ErrPendingChange                = errors.New("subscription already has a pending plan change")
	ErrSamePlan                     = errors.New("target plan is not an upgrade")
	ErrDowngradeRequiresSchedule    = errors.New("downgrade must be scheduled by subscription owner")
)
