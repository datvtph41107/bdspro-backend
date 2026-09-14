package httpresponse

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	_request "common/request"
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
	if len(problem.Errors) != 1 || problem.Errors[0].Field != "tier_rank" {
		t.Fatalf("errors = %+v", problem.Errors)
	}
}

func TestWriteProblemFailsClosedOnInvalidStatus(t *testing.T) {
	recorder := httptest.NewRecorder()
	WriteProblem(context.Background(), recorder, Problem{Detail: "do not trust caller status"})
	if recorder.Code != http.StatusInternalServerError {
		t.Fatalf("status = %d, want 500", recorder.Code)
	}
}
