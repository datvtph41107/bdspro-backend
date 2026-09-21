package httperror

import (
	commonerrors "common/errors"
	"common/request"
	"common/validation"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
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

func TestValidationResponseHasHTTP400AndCorrelation(t *testing.T) {
	ctx := request.WithRequestID(context.Background(), "req_9770ca1e2a97738baee1447ecd49df9f")
	err, _ := commonerrors.ValidationStatus(validation.New("draft.tier_rank", "PLAN_TIER_RANK_INVALID", "tier rank must be positive"), "user.catalog")
	recorder := httptest.NewRecorder()
	WriteGRPCContext(ctx, recorder, err)
	var result map[string]any
	if err := json.Unmarshal(recorder.Body.Bytes(), &result); err != nil {
		t.Fatal(err)
	}
	if recorder.Code != 400 || result["code"] != float64(3) || result["reason"] != "PLAN_TIER_RANK_INVALID" || result["request_id"] == nil || result["errors"] == nil {
		t.Fatalf("response=%s", recorder.Body.String())
	}
}

func TestLegacyBusinessErrorReturnsHTTP200(t *testing.T) {
	for _, code := range []int32{400, 401, 403, 404, 409, 429, 500, 501, 503, 0, 1004, 1016} {
		recorder := httptest.NewRecorder()
		WriteGRPC(recorder, commonerrors.ReturnError(code, "diagnostic"))
		if recorder.Code != http.StatusOK {
			t.Fatalf("legacy %d returned HTTP %d, want 200", code, recorder.Code)
		}
		var response struct {
			Status     int   `json:"status"`
			LegacyCode int32 `json:"legacy_code"`
		}
		if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
			t.Fatal(err)
		}
		if response.Status != http.StatusOK {
			t.Fatalf("legacy %d body status = %d, want 200", code, response.Status)
		}
		if response.LegacyCode != code {
			t.Fatalf("legacy %d legacy_code = %d", code, response.LegacyCode)
		}
	}
}

func TestInternalErrorDetailsAreNotPublic(t *testing.T) {
	st, _ := status.New(codes.Internal, "sql secret").WithDetails(&errdetails.ErrorInfo{Reason: "DATABASE_FAILURE", Metadata: map[string]string{"dsn": "private"}})
	recorder := httptest.NewRecorder()
	WriteGRPC(recorder, st.Err())
	if strings.Contains(recorder.Body.String(), "private") || strings.Contains(recorder.Body.String(), "sql secret") {
		t.Fatal(recorder.Body.String())
	}
}

func TestAuthBusinessErrorMessageIsPublic(t *testing.T) {
	recorder := httptest.NewRecorder()
	WriteGRPC(recorder, commonerrors.ReturnError(1004, "Mã không chính xác. Bạn còn 2 lần."))
	if recorder.Code != http.StatusOK {
		t.Fatalf("HTTP status = %d, want %d", recorder.Code, http.StatusOK)
	}
	var response struct {
		Message    string `json:"message"`
		LegacyCode int32  `json:"legacy_code"`
		Status     int    `json:"status"`
	}
	if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
		t.Fatal(err)
	}
	if response.LegacyCode != 1004 {
		t.Fatalf("legacy_code = %d", response.LegacyCode)
	}
	if response.Status != http.StatusOK {
		t.Fatalf("status = %d", response.Status)
	}
	if response.Message != "Mã không chính xác. Bạn còn 2 lần." {
		t.Fatalf("message = %q", response.Message)
	}
}

func TestOwnerInternalBusinessMessageIsPublic(t *testing.T) {
	recorder := httptest.NewRecorder()
	WriteGRPC(recorder, commonerrors.ReturnError(500, "Không tạo được hồ sơ người dùng. Vui lòng thử lại"))
	if recorder.Code != http.StatusOK {
		t.Fatalf("HTTP status = %d, want %d", recorder.Code, http.StatusOK)
	}
	var response struct {
		Message string `json:"message"`
	}
	if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
		t.Fatal(err)
	}
	if response.Message != "Không tạo được hồ sơ người dùng. Vui lòng thử lại" {
		t.Fatalf("message = %q", response.Message)
	}
}

func TestOTPActivatedBusinessErrorReturnsHTTP200(t *testing.T) {
	recorder := httptest.NewRecorder()
	WriteGRPC(recorder, commonerrors.ReturnError(1016, "OTP đã kích hoạt"))
	if recorder.Code != http.StatusOK {
		t.Fatalf("HTTP status = %d, want %d", recorder.Code, http.StatusOK)
	}
	var response struct {
		Code       int32  `json:"code"`
		LegacyCode int32  `json:"legacy_code"`
		Message    string `json:"message"`
		Status     int    `json:"status"`
	}
	if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
		t.Fatal(err)
	}
	if response.LegacyCode != 1016 || response.Status != http.StatusOK {
		t.Fatalf("response = %+v", response)
	}
	if response.Message != "OTP đã kích hoạt" {
		t.Fatalf("message = %q", response.Message)
	}
}
