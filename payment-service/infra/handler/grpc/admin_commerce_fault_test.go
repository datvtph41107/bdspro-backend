package commercegrpc

import (
	"testing"

	_errors "common/errors"
	"payment/internal"

	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestAdminPaymentErrorPreservesCanonicalPageSizeError(t *testing.T) {
	err := _errors.ReturnError(
		service.AdminPageSizeOutOfRange,
		_errors.WithViolations(_errors.FieldViolation{
			Field:       "page_size",
			Description: "must be between 1 and 100",
		}),
	)
	st, ok := status.FromError(adminPaymentError(err))
	if !ok {
		t.Fatal("adminPaymentError did not return a gRPC status")
	}
	if st.Code() != codes.InvalidArgument {
		t.Fatalf("gRPC code = %s, want %s", st.Code(), codes.InvalidArgument)
	}

	var errorCode, applicationCode, field string
	for _, detail := range st.Details() {
		switch typed := detail.(type) {
		case *errdetails.ErrorInfo:
			errorCode = typed.Metadata["error_code"]
			applicationCode = typed.Metadata["application_code"]
		case *errdetails.BadRequest:
			if len(typed.FieldViolations) > 0 {
				field = typed.FieldViolations[0].Field
			}
		}
	}
	if errorCode != "payment.admin.page_size_out_of_range" {
		t.Fatalf("error_code = %q", errorCode)
	}
	if applicationCode != "510001" {
		t.Fatalf("application_code = %q", applicationCode)
	}
	if field != "page_size" {
		t.Fatalf("field violation = %q", field)
	}
}
