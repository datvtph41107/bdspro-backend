package _jwt

import (
	"context"
	"errors"
	"fmt"

	_errors "common/errors"
	commonhttp "common/httpresponse"

	"github.com/gin-gonic/gin"
)

// AuthorizationHeaderMissingFault preserves the historical direct-HTTP
// representation while canonical identity lives in common/errors.
func AuthorizationHeaderMissingFault() error {
	return _errors.ReturnError(_errors.AuthorizationHeaderMissing)
}

func TokenMissingFault() error {
	return _errors.ReturnError(_errors.TokenMissing)
}

func InvalidOrExpiredTokenFault(cause error) error {
	return _errors.ReturnError(
		_errors.TokenInvalidOrExpired,
		_errors.WithCause(cause),
	)
}

func abortWithFault(c *gin.Context, err error) {
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
	return InvalidOrExpiredTokenFault(cause)
}

func actorContextConflictFault(cause error) error {
	return fmt.Errorf("actor context conflict: %w", cause)
}

func tempTokenRouteForbiddenFault() error {
	return _errors.ReturnError(_errors.TempTokenRouteForbidden)
}

func callerContextConflictFault(cause error) error {
	return fmt.Errorf("caller context conflict: %w", cause)
}

func callerResolutionFault(err error) error {
	switch {
	case errors.Is(err, ErrAPIKeyHeaderInvalid):
		return _errors.ReturnError(
			_errors.APIKeyHeaderInvalid,
			_errors.WithCause(err),
		)

	case errors.Is(err, ErrAPIKeyInvalid):
		return _errors.ReturnError(
			_errors.APIKeyInvalid,
			_errors.WithCause(err),
		)

	case errors.Is(err, ErrAPIKeyUnavailable),
		errors.Is(err, context.DeadlineExceeded):
		return _errors.ReturnError(
			_errors.APIKeyVerificationUnavailable,
			_errors.WithCause(err),
		)

	case errors.Is(err, ErrAuthenticationRequired):
		return _errors.ReturnError(
			_errors.AuthenticationRequired,
			_errors.WithCause(err),
		)

	case errors.Is(err, ErrCallerClassification):
		return _errors.ReturnError(
			_errors.CallerClassificationInvalid,
			_errors.WithCause(err),
		)

	default:
		return fmt.Errorf("resolve caller identity: %w", err)
	}
}
