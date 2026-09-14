package handler_grpc

import (
	"errors"
	"testing"

	"common/fault"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestQHLayerFamilyErrorPreservesCanonicalFault(t *testing.T) {
	st, ok := status.FromError(qhLayerFamilyError(fault.New(
		fault.KindNotFound,
		"tqd.qh_layer_family.not_found",
		"layer family 42 was not found",
	)))
	if !ok {
		t.Fatal("qhLayerFamilyError did not return a gRPC status")
	}
	if st.Code() != codes.NotFound {
		t.Fatalf("gRPC code = %s, want %s", st.Code(), codes.NotFound)
	}
	if errorCodeFromQHLayerFamilyStatus(st) != "tqd.qh_layer_family.not_found" {
		t.Fatalf("error_code = %q", errorCodeFromQHLayerFamilyStatus(st))
	}
}

func TestQHLayerFamilyValidationCarriesFieldViolation(t *testing.T) {
	st, ok := status.FromError(qhLayerFamilyValidation(
		"tqd.qh_layer_family.id_required",
		"id is required",
		"id",
	))
	if !ok {
		t.Fatal("validation did not return a gRPC status")
	}
	if st.Code() != codes.InvalidArgument {
		t.Fatalf("gRPC code = %s, want %s", st.Code(), codes.InvalidArgument)
	}
	if errorCodeFromQHLayerFamilyStatus(st) != "tqd.qh_layer_family.id_required" {
		t.Fatalf("error_code = %q", errorCodeFromQHLayerFamilyStatus(st))
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

func TestQHLayerFamilyErrorDoesNotLeakDependencyFailure(t *testing.T) {
	st, ok := status.FromError(qhLayerFamilyError(errors.New("postgres password=secret connection failed")))
	if !ok {
		t.Fatal("dependency failure did not return a gRPC status")
	}
	if st.Code() != codes.Internal {
		t.Fatalf("gRPC code = %s, want %s", st.Code(), codes.Internal)
	}
	if st.Message() != "layer family operation failed" {
		t.Fatalf("message = %q", st.Message())
	}
	if errorCodeFromQHLayerFamilyStatus(st) != "tqd.qh_layer_family.internal" {
		t.Fatalf("error_code = %q", errorCodeFromQHLayerFamilyStatus(st))
	}
}

func errorCodeFromQHLayerFamilyStatus(st *status.Status) string {
	for _, detail := range st.Details() {
		if info, ok := detail.(*errdetails.ErrorInfo); ok {
			return info.Metadata["error_code"]
		}
	}
	return ""
}
