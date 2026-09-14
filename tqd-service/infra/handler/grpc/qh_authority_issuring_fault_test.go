package handler_grpc

import (
	"errors"
	"testing"

	"common/fault"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestQHAuthorityIssuringErrorPreservesCanonicalFault(t *testing.T) {
	err := qhAuthorityIssuringError(fault.New(
		fault.KindNotFound,
		"tqd.qh_authority_issuring.not_found",
		"authority issuring 42 was not found",
	))
	st, ok := status.FromError(err)
	if !ok {
		t.Fatal("qhAuthorityIssuringError did not return a gRPC status")
	}
	if st.Code() != codes.NotFound {
		t.Fatalf("gRPC code = %s, want %s", st.Code(), codes.NotFound)
	}
	if errorCodeFromQHAuthorityStatus(st) != "tqd.qh_authority_issuring.not_found" {
		t.Fatalf("error_code = %q", errorCodeFromQHAuthorityStatus(st))
	}
}

func TestQHAuthorityIssuringValidationCarriesFieldViolation(t *testing.T) {
	st, ok := status.FromError(qhAuthorityIssuringValidation(
		"tqd.qh_authority_issuring.id_required",
		"id is required",
		"id",
	))
	if !ok {
		t.Fatal("validation did not return a gRPC status")
	}
	if st.Code() != codes.InvalidArgument {
		t.Fatalf("gRPC code = %s, want %s", st.Code(), codes.InvalidArgument)
	}
	if errorCodeFromQHAuthorityStatus(st) != "tqd.qh_authority_issuring.id_required" {
		t.Fatalf("error_code = %q", errorCodeFromQHAuthorityStatus(st))
	}

	var field string
	for _, detail := range st.Details() {
		if badRequest, ok := detail.(*errdetails.BadRequest); ok && len(badRequest.FieldViolations) > 0 {
			field = badRequest.FieldViolations[0].Field
		}
	}
	if field != "id" {
		t.Fatalf("field violation = %q, want id", field)
	}
}

func TestQHAuthorityIssuringErrorDoesNotLeakDependencyFailure(t *testing.T) {
	st, ok := status.FromError(qhAuthorityIssuringError(errors.New("postgres password=secret connection failed")))
	if !ok {
		t.Fatal("dependency failure did not return a gRPC status")
	}
	if st.Code() != codes.Internal {
		t.Fatalf("gRPC code = %s, want %s", st.Code(), codes.Internal)
	}
	if st.Message() != "authority issuring operation failed" {
		t.Fatalf("message = %q", st.Message())
	}
	if errorCodeFromQHAuthorityStatus(st) != "tqd.qh_authority_issuring.internal" {
		t.Fatalf("error_code = %q", errorCodeFromQHAuthorityStatus(st))
	}
}

func errorCodeFromQHAuthorityStatus(st *status.Status) string {
	for _, detail := range st.Details() {
		if info, ok := detail.(*errdetails.ErrorInfo); ok {
			return info.Metadata["error_code"]
		}
	}
	return ""
}
