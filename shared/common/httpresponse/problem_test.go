package httpresponse

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"common/fault"
	_request "common/request"
)

func TestWriteProblemUsesRFC9457MediaTypeAndCorrelation(
	t *testing.T,
) {
	ctx := _request.WithRequestID(
		context.Background(),
		"req-123",
	)

	ctx, err := _request.BindOperationID(
		ctx,
		"op-456",
	)
	if err != nil {
		t.Fatal(err)
	}

	recorder := httptest.NewRecorder()

	WriteProblem(
		ctx,
		recorder,
		Problem{
			Type: "urn:qhpro:error:" +
				"catalog.plan_version.tier_rank_positive",
			Title:  "Invalid request",
			Status: http.StatusBadRequest,
			Detail: "tier rank must be positive",
			Code: "catalog.plan_version." +
				"tier_rank_positive",
			Errors: []FieldProblem{
				{
					Field:  "tier_rank",
					Detail: "must be positive",
				},
			},
		},
	)

	if recorder.Code != http.StatusBadRequest {
		t.Fatalf(
			"status = %d, want %d",
			recorder.Code,
			http.StatusBadRequest,
		)
	}

	if got := recorder.Header().Get("Content-Type"); got != ProblemMediaType {
		t.Fatalf(
			"content type = %q, want %q",
			got,
			ProblemMediaType,
		)
	}

	var problem Problem

	if err := json.Unmarshal(
		recorder.Body.Bytes(),
		&problem,
	); err != nil {
		t.Fatal(err)
	}

	if problem.Code !=
		"catalog.plan_version.tier_rank_positive" {
		t.Fatalf(
			"code = %q",
			problem.Code,
		)
	}

	if problem.RequestID != "req-123" ||
		problem.OperationID != "op-456" {
		t.Fatalf(
			"correlation = request:%q operation:%q",
			problem.RequestID,
			problem.OperationID,
		)
	}

	if problem.Instance !=
		"urn:qhpro:request:req-123" {
		t.Fatalf(
			"instance = %q",
			problem.Instance,
		)
	}
}

func TestStatusForKindOwnsHTTPMapping(t *testing.T) {
	tests := []struct {
		kind   fault.Kind
		status int
	}{
		{fault.KindValidation, http.StatusBadRequest},
		{fault.KindUnauthenticated, http.StatusUnauthorized},
		{fault.KindPermissionDenied, http.StatusForbidden},
		{fault.KindNotFound, http.StatusNotFound},
		{fault.KindConflict, http.StatusConflict},
		{fault.KindPrecondition, http.StatusPreconditionFailed},
		{fault.KindAborted, http.StatusConflict},
		{fault.KindResourceExhausted, http.StatusTooManyRequests},
		{fault.KindUnavailable, http.StatusServiceUnavailable},
		{fault.KindInternal, http.StatusInternalServerError},
	}

	for _, test := range tests {
		if got := StatusForKind(test.kind); got != test.status {
			t.Fatalf(
				"kind %q status = %d, want %d",
				test.kind,
				got,
				test.status,
			)
		}
	}
}

func TestProblemFromErrorPreservesStableIdentity(
	t *testing.T,
) {
	err := fault.Validation(
		"catalog.plan_version.tier_rank_positive",
		"tier rank must be positive",
		fault.FieldViolation{
			Field:       "tier_rank",
			Description: "must be positive",
		},
	).WithMetadata(
		map[string]string{
			"plan_id": "7",
		},
	)

	problem := ProblemFromError(err)

	if problem.Status != http.StatusBadRequest {
		t.Fatalf(
			"status = %d",
			problem.Status,
		)
	}

	if problem.Code !=
		"catalog.plan_version.tier_rank_positive" {
		t.Fatalf(
			"code = %q",
			problem.Code,
		)
	}

	if len(problem.Errors) != 1 ||
		problem.Errors[0].Field != "tier_rank" {
		t.Fatalf(
			"errors = %+v",
			problem.Errors,
		)
	}

	if problem.Metadata["plan_id"] != "7" {
		t.Fatalf(
			"metadata = %+v",
			problem.Metadata,
		)
	}
}

func TestProblemFromUnknownErrorFailsClosed(
	t *testing.T,
) {
	problem := ProblemFromError(
		errors.New(
			"postgres password=secret authentication failed",
		),
	)

	if problem.Status !=
		http.StatusInternalServerError {
		t.Fatalf(
			"status = %d",
			problem.Status,
		)
	}

	if problem.Code != "common.internal" {
		t.Fatalf(
			"code = %q",
			problem.Code,
		)
	}

	if problem.Detail != "internal server error" {
		t.Fatalf(
			"detail = %q",
			problem.Detail,
		)
	}

	if strings.Contains(problem.Detail, "secret") ||
		strings.Contains(problem.Detail, "postgres") {
		t.Fatalf(
			"internal cause leaked: %q",
			problem.Detail,
		)
	}
}

func TestWriteProblemLegacyEnvelopePreservesNumericCodeAndAddsIdentity(
	t *testing.T,
) {
	ctx := _request.WithRequestID(
		context.Background(),
		"req-legacy",
	)

	recorder := httptest.NewRecorder()

	WriteProblem(
		ctx,
		recorder,
		NewProblem(
			http.StatusUnauthorized,
			"auth.token_invalid_or_expired",
			"Invalid or expired token",
		),
		WithLegacyJSONEnvelope(),
	)

	if recorder.Code != http.StatusUnauthorized {
		t.Fatalf(
			"status = %d, want %d",
			recorder.Code,
			http.StatusUnauthorized,
		)
	}

	if got := recorder.Header().Get("Content-Type"); got != LegacyJSONMediaType {
		t.Fatalf(
			"content type = %q, want %q",
			got,
			LegacyJSONMediaType,
		)
	}

	var response struct {
		Code      int    `json:"code"`
		Message   string `json:"message"`
		ErrorCode string `json:"error_code"`
		Status    int    `json:"status"`
		RequestID string `json:"request_id"`
	}

	if err := json.Unmarshal(
		recorder.Body.Bytes(),
		&response,
	); err != nil {
		t.Fatal(err)
	}

	if response.Code != http.StatusUnauthorized {
		t.Fatalf(
			"legacy code = %d",
			response.Code,
		)
	}

	if response.Message != "Invalid or expired token" {
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

	if response.Status != http.StatusUnauthorized {
		t.Fatalf(
			"status field = %d",
			response.Status,
		)
	}

	if response.RequestID != "req-legacy" {
		t.Fatalf(
			"request_id = %q",
			response.RequestID,
		)
	}
}
