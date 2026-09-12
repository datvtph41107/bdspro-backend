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

func TestWriteGRPCDoesNotExposeUnavailableTransportDetails(t *testing.T) {
	recorder := httptest.NewRecorder()
	WriteGRPC(recorder, status.Error(
		codes.Unavailable,
		`connection error: desc = "transport: Error while dialing: dial tcp 127.0.0.1:8219: connect: connection refused"`,
	))

	if recorder.Code != http.StatusServiceUnavailable {
		t.Fatalf("HTTP status = %d, want %d", recorder.Code, http.StatusServiceUnavailable)
	}

	var response struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	}
	if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
		t.Fatal(err)
	}
	if response.Code != int(codes.Unavailable) {
		t.Fatalf("response code = %d, want %d", response.Code, codes.Unavailable)
	}
	if response.Message != "máy chủ đang bận vui lòng thử lại sau (8219)" {
		t.Fatalf("response message = %q", response.Message)
	}
}

func TestWriteGRPCUnavailableWithoutTransportDetails(t *testing.T) {
	recorder := httptest.NewRecorder()
	WriteGRPC(recorder, status.Error(codes.Unavailable, "admin usage projection unavailable"))

	if recorder.Code != http.StatusServiceUnavailable {
		t.Fatalf("HTTP status = %d, want %d", recorder.Code, http.StatusServiceUnavailable)
	}

	var response struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	}
	if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
		t.Fatal(err)
	}
	if response.Code != int(codes.Unavailable) {
		t.Fatalf("response code = %d, want %d", response.Code, codes.Unavailable)
	}
	if response.Message != "admin usage projection unavailable" {
		t.Fatalf("response message = %q", response.Message)
	}
}
