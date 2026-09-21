package httpresponse

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	_errors "common/errors"
	_request "common/request"

	"google.golang.org/grpc/codes"
)

func TestWriteProblemUsesRFC9457MediaTypeAndCorrelation(t *testing.T) {
	ctx := _request.WithRequestID(context.Background(), "req-123")
	ctx, err := _request.BindOperationID(ctx, "op-456")
	if err != nil {
		t.Fatal(err)
	}

	recorder := httptest.NewRecorder()
	WriteProblem(ctx, recorder, Problem{
		Type:   "urn:qhpro:error:catalog.plan_version.tier_rank_positive",
		Title:  "Invalid request",
		Status: http.StatusBadRequest,
		Detail: "tier rank must be positive",
		Code:   "catalog.plan_version.tier_rank_positive",
		Errors: []FieldProblem{{Field: "tier_rank", Detail: "must be positive"}},
	})

	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d", recorder.Code, http.StatusBadRequest)
	}
	if got := recorder.Header().Get("Content-Type"); got != ProblemMediaType {
		t.Fatalf("content type = %q, want %q", got, ProblemMediaType)
	}

	var problem Problem
	if err := json.Unmarshal(recorder.Body.Bytes(), &problem); err != nil {
		t.Fatal(err)
	}
	if problem.Code != "catalog.plan_version.tier_rank_positive" {
		t.Fatalf("code = %q", problem.Code)
	}
	if problem.RequestID != "req-123" || problem.OperationID != "op-456" {
		t.Fatalf("correlation = request:%q operation:%q", problem.RequestID, problem.OperationID)
	}
	if problem.Instance != "urn:qhpro:request:req-123" {
		t.Fatalf("instance = %q", problem.Instance)
	}
}

func TestStatusForRPCOwnsHTTPMapping(t *testing.T) {
	tests := []struct {
		rpc    codes.Code
		status int
	}{
		{codes.InvalidArgument, http.StatusBadRequest},
		{codes.Unauthenticated, http.StatusUnauthorized},
		{codes.PermissionDenied, http.StatusForbidden},
		{codes.NotFound, http.StatusNotFound},
		{codes.AlreadyExists, http.StatusConflict},
		{codes.FailedPrecondition, http.StatusPreconditionFailed},
		{codes.Aborted, http.StatusConflict},
		{codes.ResourceExhausted, http.StatusTooManyRequests},
		{codes.Unavailable, http.StatusServiceUnavailable},
		{codes.Internal, http.StatusInternalServerError},
	}
	for _, test := range tests {
		if got := StatusForRPC(test.rpc); got != test.status {
			t.Fatalf("rpc %q status = %d, want %d", test.rpc, got, test.status)
		}
	}
}

func TestProblemFromErrorPreservesCanonicalAndLegacyIdentity(t *testing.T) {
	spec := _errors.MustSpec(
		100900,
		"TEST_TIER_RANK_POSITIVE",
		"tier rank must be positive",
		codes.InvalidArgument,
		_errors.LegacyProblemCode("catalog.plan_version.tier_rank_positive"),
	)
	err := _errors.ReturnError(
		spec,
		_errors.WithViolations(_errors.FieldViolation{
			Field:       "tier_rank",
			Description: "must be positive",
		}),
		_errors.WithMetadata(map[string]string{"plan_id": "7"}),
	)

	problem := ProblemFromError(err)
	if problem.Status != http.StatusBadRequest {
		t.Fatalf("status = %d", problem.Status)
	}
	if problem.Code != "catalog.plan_version.tier_rank_positive" {
		t.Fatalf("code = %q", problem.Code)
	}
	if problem.Reason != "TEST_TIER_RANK_POSITIVE" || problem.Domain != _errors.ErrorDomain {
		t.Fatalf("identity = reason:%q domain:%q", problem.Reason, problem.Domain)
	}
	if len(problem.Errors) != 1 || problem.Errors[0].Field != "tier_rank" {
		t.Fatalf("errors = %+v", problem.Errors)
	}
	if problem.Metadata["plan_id"] != "7" {
		t.Fatalf("metadata = %+v", problem.Metadata)
	}
}

func TestProblemFromUnknownErrorFailsClosedWithoutFakeApplicationIdentity(t *testing.T) {
	problem := ProblemFromError(errors.New("postgres password=secret authentication failed"))
	if problem.Status != http.StatusInternalServerError {
		t.Fatalf("status = %d", problem.Status)
	}
	if problem.Code != "" {
		t.Fatalf("code = %q, want empty", problem.Code)
	}
	if problem.Type != "about:blank" {
		t.Fatalf("type = %q", problem.Type)
	}
	if problem.Detail != "internal server error" {
		t.Fatalf("detail = %q", problem.Detail)
	}
	if strings.Contains(problem.Detail, "secret") || strings.Contains(problem.Detail, "postgres") {
		t.Fatalf("internal cause leaked: %q", problem.Detail)
	}
}

func TestWriteProblemLegacyEnvelopePreservesNumericCodeAndAddsIdentity(t *testing.T) {
	ctx := _request.WithRequestID(context.Background(), "req-legacy")
	recorder := httptest.NewRecorder()
	WriteProblem(
		ctx,
		recorder,
		NewProblem(http.StatusUnauthorized, "auth.token_invalid_or_expired", "Invalid or expired token"),
		WithLegacyJSONEnvelope(),
	)

	if recorder.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d", recorder.Code)
	}
	var response struct {
		Code      int    `json:"code"`
		Message   string `json:"message"`
		ErrorCode string `json:"error_code"`
		Status    int    `json:"status"`
		RequestID string `json:"request_id"`
	}
	if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
		t.Fatal(err)
	}
	if response.Code != http.StatusUnauthorized ||
		response.ErrorCode != "auth.token_invalid_or_expired" ||
		response.Status != http.StatusUnauthorized ||
		response.RequestID != "req-legacy" {
		t.Fatalf("response = %+v", response)
	}
}

func TestWriteProblemLegacyDataEnvelopePreservesExactHistoricalShape(t *testing.T) {
	ctx := _request.WithRequestID(context.Background(), "req-historical")
	ctx, err := _request.BindOperationID(ctx, "op-historical")
	if err != nil {
		t.Fatal(err)
	}
	recorder := httptest.NewRecorder()
	WriteProblem(
		ctx,
		recorder,
		NewProblem(http.StatusUnauthorized, "auth.token_invalid_or_expired", "Invalid or expired token"),
		WithLegacyDataJSONEnvelope(),
	)

	want := `{"code":401,"message":"Invalid or expired token","data":{}}` + "\n"
	if got := recorder.Body.String(); got != want {
		t.Fatalf("body = %q, want %q", got, want)
	}
	for _, forbidden := range []string{"error_code", "status", "request_id", "operation_id"} {
		if strings.Contains(recorder.Body.String(), forbidden) {
			t.Fatalf("historical body leaked %q: %q", forbidden, recorder.Body.String())
		}
	}
}
