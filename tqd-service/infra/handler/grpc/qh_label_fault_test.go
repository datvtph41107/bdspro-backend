package handler_grpc

import (
	"errors"
	"testing"

	"common/fault"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestQHLabelErrorPreservesCanonicalNotFoundFault(t *testing.T) {
	st, ok := status.FromError(qhLabelError(fault.New(
		fault.KindNotFound,
		"tqd.qh_label.source_not_found",
		"source label 42 was not found",
	)))
	if !ok {
		t.Fatal("qhLabelError did not return a gRPC status")
	}
	if st.Code() != codes.NotFound {
		t.Fatalf("gRPC code = %s, want %s", st.Code(), codes.NotFound)
	}
	if errorCodeFromQHLabelStatus(st) != "tqd.qh_label.source_not_found" {
		t.Fatalf("error_code = %q", errorCodeFromQHLabelStatus(st))
	}
}

func TestQHLabelErrorMapsNameConflictToAlreadyExists(t *testing.T) {
	st, ok := status.FromError(qhLabelError(fault.New(
		fault.KindConflict,
		"tqd.qh_label.name_conflict",
		"label name already exists in layer",
	)))
	if !ok {
		t.Fatal("qhLabelError did not return a gRPC status")
	}
	if st.Code() != codes.AlreadyExists {
		t.Fatalf("gRPC code = %s, want %s", st.Code(), codes.AlreadyExists)
	}
	if errorCodeFromQHLabelStatus(st) != "tqd.qh_label.name_conflict" {
		t.Fatalf("error_code = %q", errorCodeFromQHLabelStatus(st))
	}
}

func TestQHLabelErrorCarriesLayerMismatchFieldViolation(t *testing.T) {
	st, ok := status.FromError(qhLabelError(fault.Validation(
		"tqd.qh_label.source_layer_mismatch",
		"source label belongs to another layer",
		fault.FieldViolation{Field: "source_label_ids", Description: "contains a label from another layer"},
	)))
	if !ok {
		t.Fatal("qhLabelError did not return a gRPC status")
	}
	if st.Code() != codes.InvalidArgument {
		t.Fatalf("gRPC code = %s, want %s", st.Code(), codes.InvalidArgument)
	}
	if errorCodeFromQHLabelStatus(st) != "tqd.qh_label.source_layer_mismatch" {
		t.Fatalf("error_code = %q", errorCodeFromQHLabelStatus(st))
	}

	var field string
	for _, detail := range st.Details() {
		if badRequest, ok := detail.(*errdetails.BadRequest); ok && len(badRequest.FieldViolations) > 0 {
			field = badRequest.FieldViolations[0].Field
		}
	}
	if field != "source_label_ids" {
		t.Fatalf("field violation = %q, want source_label_ids", field)
	}
}

func TestQHLabelErrorDoesNotLeakDependencyFailure(t *testing.T) {
	st, ok := status.FromError(qhLabelError(errors.New("postgres password=secret connection failed")))
	if !ok {
		t.Fatal("dependency failure did not return a gRPC status")
	}
	if st.Code() != codes.Internal {
		t.Fatalf("gRPC code = %s, want %s", st.Code(), codes.Internal)
	}
	if st.Message() != "label operation failed" {
		t.Fatalf("message = %q, want safe public message", st.Message())
	}
	if errorCodeFromQHLabelStatus(st) != "tqd.qh_label.internal" {
		t.Fatalf("error_code = %q", errorCodeFromQHLabelStatus(st))
	}
}

func errorCodeFromQHLabelStatus(st *status.Status) string {
	for _, detail := range st.Details() {
		if info, ok := detail.(*errdetails.ErrorInfo); ok {
			return info.Metadata["error_code"]
		}
	}
	return ""
}
