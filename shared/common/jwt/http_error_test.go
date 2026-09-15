package _jwt

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"common/fault"
	commonhttp "common/httpresponse"

	"github.com/gin-gonic/gin"
)

func TestCallerResolutionFaultMapping(
	t *testing.T,
) {
	tests := []struct {
		name   string
		err    error
		kind   fault.Kind
		code   string
		status int
	}{
		{
			name:   "invalid header",
			err:    ErrAPIKeyHeaderInvalid,
			kind:   fault.KindValidation,
			code:   "auth.api_key_header_invalid",
			status: http.StatusBadRequest,
		},
		{
			name:   "invalid api key",
			err:    ErrAPIKeyInvalid,
			kind:   fault.KindUnauthenticated,
			code:   "auth.api_key_invalid",
			status: http.StatusUnauthorized,
		},
		{
			name:   "verifier unavailable",
			err:    ErrAPIKeyUnavailable,
			kind:   fault.KindUnavailable,
			code:   "auth.api_key_verification_unavailable",
			status: http.StatusServiceUnavailable,
		},
		{
			name:   "deadline exceeded",
			err:    context.DeadlineExceeded,
			kind:   fault.KindUnavailable,
			code:   "auth.api_key_verification_unavailable",
			status: http.StatusServiceUnavailable,
		},
		{
			name:   "authentication required",
			err:    ErrAuthenticationRequired,
			kind:   fault.KindUnauthenticated,
			code:   "auth.authentication_required",
			status: http.StatusUnauthorized,
		},
		{
			name:   "invalid caller",
			err:    ErrCallerClassification,
			kind:   fault.KindUnauthenticated,
			code:   "auth.caller_classification_invalid",
			status: http.StatusUnauthorized,
		},
	}

	for _, test := range tests {
		t.Run(
			test.name,
			func(t *testing.T) {
				err :=
					callerResolutionFault(
						test.err,
					)

				failure, ok :=
					fault.As(err)

				if !ok {
					t.Fatalf(
						"error type = %T",
						err,
					)
				}

				if failure.Kind() !=
					test.kind {
					t.Fatalf(
						"kind = %q, want %q",
						failure.Kind(),
						test.kind,
					)
				}

				if failure.Code() !=
					test.code {
					t.Fatalf(
						"code = %q, want %q",
						failure.Code(),
						test.code,
					)
				}

				problem :=
					commonhttp.ProblemFromError(
						err,
					)

				if problem.Status !=
					test.status {
					t.Fatalf(
						"status = %d, want %d",
						problem.Status,
						test.status,
					)
				}
			},
		)
	}
}

func TestCallerResolutionUnknownFailureFailsClosed(
	t *testing.T,
) {
	dependencyErr := errors.New(
		"postgres password=secret auth lookup failed",
	)

	err := callerResolutionFault(
		dependencyErr,
	)

	failure, ok := fault.As(err)

	if !ok {
		t.Fatalf(
			"error type = %T",
			err,
		)
	}

	if failure.Kind() !=
		fault.KindInternal {
		t.Fatalf(
			"kind = %q",
			failure.Kind(),
		)
	}

	if failure.Code() !=
		"auth.caller_resolution_failed" {
		t.Fatalf(
			"code = %q",
			failure.Code(),
		)
	}

	if !errors.Is(err, dependencyErr) {
		t.Fatal(
			"dependency cause was not preserved",
		)
	}

	problem :=
		commonhttp.ProblemFromError(
			err,
		)

	if problem.Status !=
		http.StatusInternalServerError {
		t.Fatalf(
			"status = %d",
			problem.Status,
		)
	}

	if strings.Contains(
		problem.Detail,
		"secret",
	) ||
		strings.Contains(
			problem.Detail,
			"postgres",
		) {
		t.Fatalf(
			"dependency detail leaked: %q",
			problem.Detail,
		)
	}
}

func TestAbortWithFaultPreservesLegacyEnvelopeAndAddsCanonicalIdentity(
	t *testing.T,
) {
	gin.SetMode(gin.TestMode)

	recorder := httptest.NewRecorder()

	ctx, _ :=
		gin.CreateTestContext(
			recorder,
		)

	ctx.Request =
		httptest.NewRequest(
			http.MethodGet,
			"/private",
			nil,
		)

	abortWithFault(
		ctx,
		invalidTokenFault(
			errors.New(
				"jwt signature key=secret",
			),
		),
	)

	if !ctx.IsAborted() {
		t.Fatal(
			"gin context was not aborted",
		)
	}

	if recorder.Code !=
		http.StatusUnauthorized {
		t.Fatalf(
			"status = %d",
			recorder.Code,
		)
	}

	if got := recorder.Header().Get(
		"Content-Type",
	); got != commonhttp.LegacyJSONMediaType {
		t.Fatalf(
			"content type = %q",
			got,
		)
	}

	var response struct {
		Code      int    `json:"code"`
		Message   string `json:"message"`
		ErrorCode string `json:"error_code"`
		Status    int    `json:"status"`
	}

	if err := json.Unmarshal(
		recorder.Body.Bytes(),
		&response,
	); err != nil {
		t.Fatal(err)
	}

	if response.Code !=
		http.StatusUnauthorized {
		t.Fatalf(
			"legacy code = %d",
			response.Code,
		)
	}

	if response.Message !=
		"Invalid or expired token" {
		t.Fatalf(
			"message = %q",
			response.Message,
		)
	}

	if response.ErrorCode !=
		"auth.token_invalid_or_expired" {
		t.Fatalf(
			"error_code = %q",
			response.ErrorCode,
		)
	}

	if response.Status !=
		http.StatusUnauthorized {
		t.Fatalf(
			"status field = %d",
			response.Status,
		)
	}

	if strings.Contains(
		response.Message,
		"secret",
	) {
		t.Fatalf(
			"token cause leaked: %q",
			response.Message,
		)
	}
}

func TestTempTokenRouteFaultMapsToForbidden(
	t *testing.T,
) {
	problem :=
		commonhttp.ProblemFromError(
			tempTokenRouteForbiddenFault(),
		)

	if problem.Status !=
		http.StatusForbidden {
		t.Fatalf(
			"status = %d",
			problem.Status,
		)
	}

	if problem.Code !=
		"auth.temp_token_route_forbidden" {
		t.Fatalf(
			"code = %q",
			problem.Code,
		)
	}
}

func TestPublicJWTCompatibilityFaultsOwnStableIdentity(
	t *testing.T,
) {
	cause := errors.New(
		"jwt parser technical failure",
	)

	tests := []struct {
		name    string
		err     error
		code    string
		message string
		cause   error
	}{
		{
			name:    "authorization header missing",
			err:     AuthorizationHeaderMissingFault(),
			code:    "auth.authorization_header_missing",
			message: "Authorization header is missing",
		},
		{
			name:    "token missing",
			err:     TokenMissingFault(),
			code:    "auth.token_missing",
			message: "Missing token",
		},
		{
			name: "invalid or expired token",
			err: InvalidOrExpiredTokenFault(
				cause,
			),
			code:    "auth.token_invalid_or_expired",
			message: "Invalid or expired token",
			cause:   cause,
		},
	}

	for _, test := range tests {
		t.Run(
			test.name,
			func(t *testing.T) {
				failure, ok := fault.As(
					test.err,
				)

				if !ok {
					t.Fatalf(
						"error type = %T",
						test.err,
					)
				}

				if failure.Kind() !=
					fault.KindUnauthenticated {
					t.Fatalf(
						"kind = %q",
						failure.Kind(),
					)
				}

				if failure.Code() !=
					test.code {
					t.Fatalf(
						"code = %q, want %q",
						failure.Code(),
						test.code,
					)
				}

				if failure.PublicMessage() !=
					test.message {
					t.Fatalf(
						"message = %q, want %q",
						failure.PublicMessage(),
						test.message,
					)
				}

				problem :=
					commonhttp.ProblemFromError(
						test.err,
					)

				if problem.Status !=
					http.StatusUnauthorized {
					t.Fatalf(
						"status = %d, want %d",
						problem.Status,
						http.StatusUnauthorized,
					)
				}

				if problem.Detail !=
					test.message {
					t.Fatalf(
						"detail = %q, want %q",
						problem.Detail,
						test.message,
					)
				}

				if test.cause != nil &&
					!errors.Is(
						test.err,
						test.cause,
					) {
					t.Fatal(
						"technical cause was not preserved",
					)
				}
			},
		)
	}
}
