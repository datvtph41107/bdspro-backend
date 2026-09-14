package fault

import (
	"testing"

	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestToGRPCMapsValidationIdentityAndFieldViolation(t *testing.T) {
	err := Validation(
		"catalog.plan_version.tier_rank_positive",
		"tier rank must be positive",
		FieldViolation{Field: "tier_rank", Description: "must be positive"},
	)

	grpcErr := ToGRPC(err)
	st, ok := status.FromError(grpcErr)
	if !ok {
		t.Fatal("ToGRPC did not return a gRPC status")
	}
	if st.Code() != codes.InvalidArgument {
		t.Fatalf("gRPC code = %s, want %s", st.Code(), codes.InvalidArgument)
	}
	if st.Message() != "tier rank must be positive" {
		t.Fatalf("message = %q", st.Message())
	}

	var foundInfo, foundViolation bool
	for _, detail := range st.Details() {
		switch value := detail.(type) {
		case *errdetails.ErrorInfo:
			foundInfo = value.Reason == "VALIDATION_ERROR" &&
				value.Domain == errorDomain &&
				value.Metadata["error_code"] == "catalog.plan_version.tier_rank_positive"
		case *errdetails.BadRequest:
			foundViolation = len(value.FieldViolations) == 1 &&
				value.FieldViolations[0].Field == "tier_rank" &&
				value.FieldViolations[0].Description == "must be positive"
		}
	}
	if !foundInfo || !foundViolation {
		t.Fatalf("missing canonical details: info=%v violation=%v details=%v", foundInfo, foundViolation, st.Details())
	}
}

func TestToGRPCMapsAbortedIdentity(t *testing.T) {
	err := New(
		KindAborted,
		"tqd.import.retry_lock_conflict",
		"import retry could not acquire lock",
	)

	st, ok := status.FromError(ToGRPC(err))
	if !ok {
		t.Fatal("ToGRPC did not return a gRPC status")
	}
	if st.Code() != codes.Aborted {
		t.Fatalf("gRPC code = %s, want %s", st.Code(), codes.Aborted)
	}
	if st.Message() != "import retry could not acquire lock" {
		t.Fatalf("message = %q", st.Message())
	}

	var foundInfo bool
	for _, detail := range st.Details() {
		if info, ok := detail.(*errdetails.ErrorInfo); ok {
			foundInfo = info.Reason == "ABORTED" &&
				info.Domain == errorDomain &&
				info.Metadata["error_code"] == "tqd.import.retry_lock_conflict"
		}
	}
	if !foundInfo {
		t.Fatalf("missing canonical aborted ErrorInfo: details=%v", st.Details())
	}
}

func TestToGRPCHidesUnknownInternalFailure(t *testing.T) {
	st, _ := status.FromError(ToGRPC(assertionError("database password leaked here")))
	if st.Code() != codes.Internal || st.Message() != "internal server error" {
		t.Fatalf("unexpected internal mapping: code=%s message=%q", st.Code(), st.Message())
	}
}

type assertionError string

func (e assertionError) Error() string { return string(e) }
