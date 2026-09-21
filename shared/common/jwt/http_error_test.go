package _jwt

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	_errors "common/errors"
	commonhttp "common/httpresponse"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/codes"
)

func TestCallerResolutionCanonicalMapping(t *testing.T) {
	tests := []struct {
		name   string
		err    error
		rpc    codes.Code
		code   string
		status int
	}{
		{"invalid header", ErrAPIKeyHeaderInvalid, codes.InvalidArgument, "auth.api_key_header_invalid", http.StatusBadRequest},
		{"invalid api key", ErrAPIKeyInvalid, codes.Unauthenticated, "auth.api_key_invalid", http.StatusUnauthorized},
		{"verifier unavailable", ErrAPIKeyUnavailable, codes.Unavailable, "auth.api_key_verification_unavailable", http.StatusServiceUnavailable},
		{"deadline exceeded", context.DeadlineExceeded, codes.Unavailable, "auth.api_key_verification_unavailable", http.StatusServiceUnavailable},
		{"authentication required", ErrAuthenticationRequired, codes.Unauthenticated, "auth.authentication_required", http.StatusUnauthorized},
		{"invalid caller", ErrCallerClassification, codes.Unauthenticated, "auth.caller_classification_invalid", http.StatusUnauthorized},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := callerResolutionFault(test.err)
			application, ok := _errors.As(err)
			if !ok {
				t.Fatalf("error type = %T", err)
			}
			if application.RPCCode() != test.rpc {
				t.Fatalf("rpc = %q, want %q", application.RPCCode(), test.rpc)
			}
			if application.Spec().LegacyProblemCode() != test.code {
				t.Fatalf("legacy problem code = %q, want %q", application.Spec().LegacyProblemCode(), test.code)
			}
			problem := commonhttp.ProblemFromError(err)
			if problem.Status != test.status || problem.Code != test.code {
				t.Fatalf("problem = %+v", problem)
			}
		})
	}
}

func TestCallerResolutionUnknownFailureFailsClosed(t *testing.T) {
	dependencyErr := errors.New("postgres password=secret auth lookup failed")
	err := callerResolutionFault(dependencyErr)

	if _, ok := _errors.As(err); ok {
		t.Fatalf("unknown technical failure became an application error: %v", err)
	}
	if !errors.Is(err, dependencyErr) {
		t.Fatal("dependency cause was not preserved")
	}

	problem := commonhttp.ProblemFromError(err)
	if problem.Status != http.StatusInternalServerError || problem.Code != "" {
		t.Fatalf("problem = %+v", problem)
	}
	if strings.Contains(problem.Detail, "secret") || strings.Contains(problem.Detail, "postgres") {
		t.Fatalf("dependency detail leaked: %q", problem.Detail)
	}
}

func TestAbortWithFaultPreservesLegacyEnvelopeAndCanonicalIdentity(t *testing.T) {
	gin.SetMode(gin.TestMode)
	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)
	ctx.Request = httptest.NewRequest(http.MethodGet, "/private", nil)

	abortWithFault(ctx, invalidTokenFault(errors.New("jwt signature key=secret")))

	if !ctx.IsAborted() {
		t.Fatal("gin context was not aborted")
	}
	var response struct {
		Code      int    `json:"code"`
		Message   string `json:"message"`
		ErrorCode string `json:"error_code"`
		Status    int    `json:"status"`
	}
	if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
		t.Fatal(err)
	}
	if recorder.Code != http.StatusUnauthorized ||
		response.Code != http.StatusUnauthorized ||
		response.ErrorCode != "auth.token_invalid_or_expired" ||
		response.Status != http.StatusUnauthorized ||
		response.Message != "Invalid or expired token" {
		t.Fatalf("response = %+v status=%d", response, recorder.Code)
	}
	if strings.Contains(response.Message, "secret") {
		t.Fatalf("token cause leaked: %q", response.Message)
	}
}

func TestTempTokenRouteFaultMapsToForbidden(t *testing.T) {
	problem := commonhttp.ProblemFromError(tempTokenRouteForbiddenFault())
	if problem.Status != http.StatusForbidden ||
		problem.Code != "auth.temp_token_route_forbidden" {
		t.Fatalf("problem = %+v", problem)
	}
}

func TestPublicJWTCompatibilityFaultsOwnStableCanonicalIdentity(t *testing.T) {
	cause := errors.New("jwt parser technical failure")
	tests := []struct {
		name    string
		err     error
		code    string
		message string
		cause   error
	}{
		{"authorization header missing", AuthorizationHeaderMissingFault(), "auth.authorization_header_missing", "Authorization header is missing", nil},
		{"token missing", TokenMissingFault(), "auth.token_missing", "Missing token", nil},
		{"invalid or expired token", InvalidOrExpiredTokenFault(cause), "auth.token_invalid_or_expired", "Invalid or expired token", cause},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			application, ok := _errors.As(test.err)
			if !ok {
				t.Fatalf("error type = %T", test.err)
			}
			if application.RPCCode() != codes.Unauthenticated ||
				application.Spec().LegacyProblemCode() != test.code ||
				application.PublicMessage() != test.message {
				t.Fatalf("application = key:%q rpc:%q code:%q message:%q",
					application.Key(), application.RPCCode(), application.Spec().LegacyProblemCode(), application.PublicMessage())
			}
			problem := commonhttp.ProblemFromError(test.err)
			if problem.Status != http.StatusUnauthorized ||
				problem.Code != test.code ||
				problem.Detail != test.message {
				t.Fatalf("problem = %+v", problem)
			}
			if test.cause != nil && !errors.Is(test.err, test.cause) {
				t.Fatal("technical cause was not preserved")
			}
		})
	}
}
