package handler

import (
	"testing"

	_fault "common/fault"

	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestAuthTransportErrorMapsCanonicalOTPFaultsToErrorInfo(
	t *testing.T,
) {
	tests := []struct {
		name       string
		errorCode  string
		legacyCode string
		second     string
		message    string
	}{
		{
			name:       "next send limited",
			errorCode:  "user.otp.next_send_limited",
			legacyCode: "1006",
			second:     "30",
			message:    "Hãy thử lại sau 30s",
		},
		{
			name:       "request limited",
			errorCode:  "user.otp.request_limited",
			legacyCode: "1017",
			second:     "45",
			message: "Bạn đã yêu cầu OTP quá nhiều, " +
				"hãy thử lại sau 45s",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := _fault.New(
				_fault.KindResourceExhausted,
				tt.errorCode,
				tt.message,
			).WithMetadata(map[string]string{
				"legacy_code": tt.legacyCode,
				"second":      tt.second,
			})

			grpcErr := authTransportError(err)

			st, ok := status.FromError(grpcErr)
			if !ok {
				t.Fatal("expected gRPC status")
			}

			if st.Code() != codes.ResourceExhausted {
				t.Fatalf(
					"code = %s, want %s",
					st.Code(),
					codes.ResourceExhausted,
				)
			}

			if st.Message() != tt.message {
				t.Fatalf(
					"message = %q, want %q",
					st.Message(),
					tt.message,
				)
			}

			var found *errdetails.ErrorInfo

			for _, detail := range st.Details() {
				info, ok := detail.(*errdetails.ErrorInfo)
				if ok {
					found = info
					break
				}
			}

			if found == nil {
				t.Fatal("canonical ErrorInfo detail not found")
			}

			if found.Metadata["error_code"] != tt.errorCode {
				t.Fatalf(
					"error_code = %q, want %q",
					found.Metadata["error_code"],
					tt.errorCode,
				)
			}

			if found.Metadata["legacy_code"] != tt.legacyCode {
				t.Fatalf(
					"legacy_code = %q, want %q",
					found.Metadata["legacy_code"],
					tt.legacyCode,
				)
			}

			if found.Metadata["second"] != tt.second {
				t.Fatalf(
					"second = %q, want %q",
					found.Metadata["second"],
					tt.second,
				)
			}
		})
	}
}
