package commercegrpc

import (
	"testing"

	adminusecase "payment/internal/usecase/admin"

	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestAdminPaymentErrorPreservesCanonicalPageSizeFault(t *testing.T) {
	st, ok := status.FromError(adminPaymentError(adminusecase.ErrPageSizeOutOfRange))
	if !ok {
		t.Fatal("adminPaymentError did not return a gRPC status")
	}
	if st.Code() != codes.InvalidArgument {
		t.Fatalf("gRPC code = %s, want %s", st.Code(), codes.InvalidArgument)
	}
	if st.Message() != "page size must be between 1 and 100" {
		t.Fatalf("message = %q", st.Message())
	}

	var errorCode string
	var field string
	for _, detail := range st.Details() {
		switch typed := detail.(type) {
		case *errdetails.ErrorInfo:
			errorCode = typed.Metadata["error_code"]
		case *errdetails.BadRequest:
			if len(typed.FieldViolations) > 0 {
				field = typed.FieldViolations[0].Field
			}
		}
	}
	if errorCode != "payment.admin.page_size_out_of_range" {
		t.Fatalf("error_code = %q", errorCode)
	}
	if field != "page_size" {
		t.Fatalf("field violation = %q", field)
	}
}
