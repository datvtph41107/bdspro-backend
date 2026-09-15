package _jwt

import (
	"context"
	"errors"

	"common/fault"
	commonhttp "common/httpresponse"

	"github.com/gin-gonic/gin"
)

func abortWithFault(
	c *gin.Context,
	err error,
) {
	problem := commonhttp.ProblemFromError(err)

	commonhttp.WriteProblem(
		c.Request.Context(),
		c.Writer,
		problem,
		commonhttp.WithLegacyJSONEnvelope(),
	)

	c.Abort()
}

func invalidTokenFault(cause error) error {
	return fault.Wrap(
		cause,
		fault.KindUnauthenticated,
		"auth.token_invalid_or_expired",
		"Invalid or expired token",
	)
}

func actorContextConflictFault(
	cause error,
) error {
	return fault.Wrap(
		cause,
		fault.KindInternal,
		"auth.actor_context_conflict",
		"actor context conflict",
	)
}

func tempTokenRouteForbiddenFault() error {
	return fault.New(
		fault.KindPermissionDenied,
		"auth.temp_token_route_forbidden",
		"Token không đúng",
	)
}

func callerContextConflictFault(
	cause error,
) error {
	return fault.Wrap(
		cause,
		fault.KindInternal,
		"auth.caller_context_conflict",
		"caller context conflict",
	)
}

func callerResolutionFault(err error) error {
	switch {
	case errors.Is(
		err,
		ErrAPIKeyHeaderInvalid,
	):
		return fault.Wrap(
			err,
			fault.KindValidation,
			"auth.api_key_header_invalid",
			"Invalid API key header",
		)

	case errors.Is(
		err,
		ErrAPIKeyInvalid,
	):
		return fault.Wrap(
			err,
			fault.KindUnauthenticated,
			"auth.api_key_invalid",
			"API key verification failed",
		)

	case errors.Is(
		err,
		ErrAPIKeyUnavailable,
	),
		errors.Is(
			err,
			context.DeadlineExceeded,
		):
		return fault.Wrap(
			err,
			fault.KindUnavailable,
			"auth.api_key_verification_unavailable",
			"API key verification unavailable",
		)

	case errors.Is(
		err,
		ErrAuthenticationRequired,
	):
		return fault.Wrap(
			err,
			fault.KindUnauthenticated,
			"auth.authentication_required",
			"Authorization header missing",
		)

	case errors.Is(
		err,
		ErrCallerClassification,
	):
		return fault.Wrap(
			err,
			fault.KindUnauthenticated,
			"auth.caller_classification_invalid",
			"Invalid caller classification",
		)

	default:
		return fault.Wrap(
			err,
			fault.KindInternal,
			"auth.caller_resolution_failed",
			"authentication resolution failed",
		)
	}
}
