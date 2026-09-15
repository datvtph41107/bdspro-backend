package httperror

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestStatusFromCodeMapsQuotaExhaustion(t *testing.T) {
	if got := statusFromCode(codes.ResourceExhausted); got != http.StatusTooManyRequests {
		t.Fatalf("statusFromCode(ResourceExhausted) = %d, want %d", got, http.StatusTooManyRequests)
	}
}

func TestStatusFromCodeMapsBusinessConflict(t *testing.T) {
	if got := statusFromCode(codes.Aborted); got != http.StatusConflict {
		t.Fatalf("statusFromCode(Aborted) = %d, want %d", got, http.StatusConflict)
	}
	if got := statusFromCode(codes.FailedPrecondition); got != http.StatusPreconditionFailed {
		t.Fatalf("statusFromCode(FailedPrecondition) = %d, want %d", got, http.StatusPreconditionFailed)
	}
}

func TestWriteGRPCPreservesSemanticErrorInfo(t *testing.T) {
	grpcStatus, err := status.New(codes.ResourceExhausted, "quota exhausted").WithDetails(
		&errdetails.ErrorInfo{
			Reason: "QUOTA_EXHAUSTED",
			Domain: "tqd.qhpro",
			Metadata: map[string]string{
				"meter":     "workspace.report_generation.accepted",
				"limit":     "3",
				"used":      "3",
				"remaining": "0",
			},
		},
	)
	if err != nil {
		t.Fatal(err)
	}
	recorder := httptest.NewRecorder()
	WriteGRPC(recorder, grpcStatus.Err())
	if recorder.Code != http.StatusTooManyRequests {
		t.Fatalf("HTTP status = %d, want %d", recorder.Code, http.StatusTooManyRequests)
	}
	var response struct {
		Reason   string            `json:"reason"`
		Domain   string            `json:"domain"`
		Metadata map[string]string `json:"metadata"`
	}
	if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
		t.Fatal(err)
	}
	if response.Reason != "QUOTA_EXHAUSTED" || response.Domain != "tqd.qhpro" ||
		response.Metadata["meter"] != "workspace.report_generation.accepted" ||
		response.Metadata["limit"] != "3" || response.Metadata["used"] != "3" {
		t.Fatalf("response = %+v", response)
	}
}

func TestWriteGRPCMapsCanonicalNotFoundToHTTPProblem(t *testing.T) {
	grpcStatus, detailsErr := status.New(
		codes.NotFound,
		"role group not found",
	).WithDetails(
		&errdetails.ErrorInfo{
			Reason: "NOT_FOUND",
			Domain: "qhpro.backend",
			Metadata: map[string]string{
				"error_code": "iam.role_group.not_found",
			},
		},
	)
	if detailsErr != nil {
		t.Fatal(detailsErr)
	}

	recorder := httptest.NewRecorder()
	WriteGRPC(recorder, grpcStatus.Err())

	if recorder.Code != http.StatusNotFound {
		t.Fatalf(
			"HTTP status = %d, want %d",
			recorder.Code,
			http.StatusNotFound,
		)
	}

	var response struct {
		Type   string `json:"type"`
		Status int    `json:"status"`
		Detail string `json:"detail"`
		Code   string `json:"code"`
		Reason string `json:"reason"`
		Domain string `json:"domain"`
	}
	if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
		t.Fatal(err)
	}

	if response.Type != "urn:qhpro:error:iam.role_group.not_found" {
		t.Fatalf("type = %q", response.Type)
	}
	if response.Status != http.StatusNotFound {
		t.Fatalf("body status = %d", response.Status)
	}
	if response.Detail != "role group not found" {
		t.Fatalf("detail = %q", response.Detail)
	}
	if response.Code != "iam.role_group.not_found" {
		t.Fatalf("code = %q", response.Code)
	}
	if response.Reason != "NOT_FOUND" {
		t.Fatalf("reason = %q", response.Reason)
	}
	if response.Domain != "qhpro.backend" {
		t.Fatalf("domain = %q", response.Domain)
	}
}
