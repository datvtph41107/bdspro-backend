package httperror

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
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

func TestWriteGRPCMapsCanonicalConflictToHTTPProblem(t *testing.T) {
	const (
		wantMessage = "invitation already exists for this user"
		wantCode    = "bdspro.deal_invitation.already_exists"
	)

	grpcStatus, detailsErr := status.New(
		codes.AlreadyExists,
		wantMessage,
	).WithDetails(
		&errdetails.ErrorInfo{
			Reason: "CONFLICT",
			Domain: "qhpro.backend",
			Metadata: map[string]string{
				"error_code": wantCode,
			},
		},
	)
	if detailsErr != nil {
		t.Fatal(detailsErr)
	}

	recorder := httptest.NewRecorder()
	WriteGRPC(recorder, grpcStatus.Err())

	if recorder.Code != http.StatusConflict {
		t.Fatalf(
			"HTTP status = %d, want %d",
			recorder.Code,
			http.StatusConflict,
		)
	}

	var response struct {
		Type     string            `json:"type"`
		Status   int               `json:"status"`
		Detail   string            `json:"detail"`
		Code     string            `json:"code"`
		Reason   string            `json:"reason"`
		Domain   string            `json:"domain"`
		Metadata map[string]string `json:"metadata"`
	}

	if err := json.Unmarshal(
		recorder.Body.Bytes(),
		&response,
	); err != nil {
		t.Fatal(err)
	}

	if response.Type !=
		"urn:qhpro:error:bdspro.deal_invitation.already_exists" {
		t.Fatalf(
			"type = %q",
			response.Type,
		)
	}

	if response.Status != http.StatusConflict {
		t.Fatalf(
			"body status = %d",
			response.Status,
		)
	}

	if response.Detail != wantMessage {
		t.Fatalf(
			"detail = %q",
			response.Detail,
		)
	}

	if response.Code != wantCode {
		t.Fatalf(
			"code = %q",
			response.Code,
		)
	}

	if response.Reason != "CONFLICT" {
		t.Fatalf(
			"reason = %q",
			response.Reason,
		)
	}

	if response.Domain != "qhpro.backend" {
		t.Fatalf(
			"domain = %q",
			response.Domain,
		)
	}

	if response.Metadata["error_code"] != wantCode {
		t.Fatalf(
			"metadata error_code = %q",
			response.Metadata["error_code"],
		)
	}
}

func TestWriteGRPCPreservesCanonicalOTPCompatibilityEnvelope(t *testing.T) {
	tests := []struct {
		name       string
		errorCode  string
		legacyCode string
		message    string
		second     string
	}{
		{
			name:       "next send limited",
			errorCode:  "user.otp.next_send_limited",
			legacyCode: "1006",
			message:    "Hãy thử lại sau 30s",
			second:     "30",
		},
		{
			name:       "request limited",
			errorCode:  "user.otp.request_limited",
			legacyCode: "1017",
			message:    "Bạn đã yêu cầu OTP quá nhiều, hãy thử lại sau 45s",
			second:     "45",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			st, err := status.New(
				codes.ResourceExhausted,
				tt.message,
			).WithDetails(
				&errdetails.ErrorInfo{
					Reason: "RESOURCE_EXHAUSTED",
					Domain: "qhpro.backend",
					Metadata: map[string]string{
						"error_code":  tt.errorCode,
						"legacy_code": tt.legacyCode,
						"second":      tt.second,
					},
				},
			)
			if err != nil {
				t.Fatal(err)
			}

			recorder := httptest.NewRecorder()

			WriteGRPC(recorder, st.Err())

			if recorder.Code != http.StatusOK {
				t.Fatalf(
					"status = %d, want %d",
					recorder.Code,
					http.StatusOK,
				)
			}

			var response struct {
				Code    int32  `json:"code"`
				Message string `json:"message"`
				Second  *int32 `json:"second"`
			}

			if err := json.Unmarshal(
				recorder.Body.Bytes(),
				&response,
			); err != nil {
				t.Fatal(err)
			}

			legacyCode, err := strconv.Atoi(tt.legacyCode)
			if err != nil {
				t.Fatal(err)
			}

			second, err := strconv.Atoi(tt.second)
			if err != nil {
				t.Fatal(err)
			}

			if response.Code != int32(legacyCode) {
				t.Fatalf(
					"code = %d, want %d",
					response.Code,
					legacyCode,
				)
			}

			if response.Message != tt.message {
				t.Fatalf(
					"message = %q, want %q",
					response.Message,
					tt.message,
				)
			}

			if response.Second == nil ||
				*response.Second != int32(second) {
				t.Fatalf(
					"second = %v, want %d",
					response.Second,
					second,
				)
			}
		})
	}
}

func TestWriteGRPCDoesNotApplyOTPCompatibilityToOtherCanonicalQuotaFaults(
	t *testing.T,
) {
	st, err := status.New(
		codes.ResourceExhausted,
		"quota exhausted",
	).WithDetails(
		&errdetails.ErrorInfo{
			Reason: "RESOURCE_EXHAUSTED",
			Domain: "qhpro.backend",
			Metadata: map[string]string{
				"error_code": "quota.exhausted",
				"second":     "30",
			},
		},
	)
	if err != nil {
		t.Fatal(err)
	}

	recorder := httptest.NewRecorder()

	WriteGRPC(recorder, st.Err())

	if recorder.Code != http.StatusTooManyRequests {
		t.Fatalf(
			"status = %d, want %d",
			recorder.Code,
			http.StatusTooManyRequests,
		)
	}
}

func TestWriteGRPCRejectsMalformedCanonicalOTPCompatibilityMetadata(
	t *testing.T,
) {
	tests := []struct {
		name       string
		errorCode  string
		legacyCode string
		second     string
	}{
		{
			name:       "wrong legacy code",
			errorCode:  "user.otp.next_send_limited",
			legacyCode: "1017",
			second:     "30",
		},
		{
			name:       "invalid second",
			errorCode:  "user.otp.request_limited",
			legacyCode: "1017",
			second:     "not-a-number",
		},
		{
			name:       "negative second",
			errorCode:  "user.otp.request_limited",
			legacyCode: "1017",
			second:     "-1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			st, err := status.New(
				codes.ResourceExhausted,
				"otp cooldown",
			).WithDetails(
				&errdetails.ErrorInfo{
					Reason: "RESOURCE_EXHAUSTED",
					Domain: "qhpro.backend",
					Metadata: map[string]string{
						"error_code":  tt.errorCode,
						"legacy_code": tt.legacyCode,
						"second":      tt.second,
					},
				},
			)
			if err != nil {
				t.Fatal(err)
			}

			recorder := httptest.NewRecorder()

			WriteGRPC(recorder, st.Err())

			if recorder.Code != http.StatusTooManyRequests {
				t.Fatalf(
					"status = %d, want %d",
					recorder.Code,
					http.StatusTooManyRequests,
				)
			}
		})
	}
}

func TestWriteGRPCRequiresCanonicalOTPErrorInfoProvenance(t *testing.T) {
	tests := []struct {
		name       string
		grpcCode   codes.Code
		reason     string
		domain     string
		wantStatus int
	}{
		{
			name:       "wrong grpc code",
			grpcCode:   codes.Internal,
			reason:     "RESOURCE_EXHAUSTED",
			domain:     "qhpro.backend",
			wantStatus: http.StatusInternalServerError,
		},
		{
			name:       "wrong reason",
			grpcCode:   codes.ResourceExhausted,
			reason:     "INTERNAL",
			domain:     "qhpro.backend",
			wantStatus: http.StatusTooManyRequests,
		},
		{
			name:       "wrong domain",
			grpcCode:   codes.ResourceExhausted,
			reason:     "RESOURCE_EXHAUSTED",
			domain:     "other.domain",
			wantStatus: http.StatusTooManyRequests,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			st, err := status.New(
				tt.grpcCode,
				"otp cooldown",
			).WithDetails(
				&errdetails.ErrorInfo{
					Reason: tt.reason,
					Domain: tt.domain,
					Metadata: map[string]string{
						"error_code":  "user.otp.next_send_limited",
						"legacy_code": "1006",
						"second":      "30",
					},
				},
			)
			if err != nil {
				t.Fatal(err)
			}

			recorder := httptest.NewRecorder()
			WriteGRPC(recorder, st.Err())

			if recorder.Code != tt.wantStatus {
				t.Fatalf(
					"status = %d, want %d",
					recorder.Code,
					tt.wantStatus,
				)
			}
		})
	}
}
