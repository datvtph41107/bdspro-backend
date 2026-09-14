package handler_grpc

import (
	"errors"
	"strings"
	"testing"

	"common/fault"

	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestImportRegionGRPCMapsRetryLockConflict(t *testing.T) {
	st, ok := status.FromError(importGRPCError(fault.New(
		fault.KindAborted,
		"tqd.import.retry_lock_conflict",
		"import error 11 could not be locked for retry",
	)))
	if !ok {
		t.Fatal("importGRPCError did not return a gRPC status")
	}
	if st.Code() != codes.Aborted {
		t.Fatalf("gRPC code = %s, want %s", st.Code(), codes.Aborted)
	}
	if importRegionErrorCode(st) != "tqd.import.retry_lock_conflict" {
		t.Fatalf("error_code = %q", importRegionErrorCode(st))
	}
}

func TestImportRegionGRPCMapsInProgressToFailedPrecondition(t *testing.T) {
	st, ok := status.FromError(importGRPCError(fault.New(
		fault.KindPrecondition,
		"tqd.import.in_progress",
		"layer 7 already has an import in progress",
	)))
	if !ok {
		t.Fatal("importGRPCError did not return a gRPC status")
	}
	if st.Code() != codes.FailedPrecondition {
		t.Fatalf("gRPC code = %s, want %s", st.Code(), codes.FailedPrecondition)
	}
}

func TestImportRegionGRPCDoesNotLeakDependencyFailure(t *testing.T) {
	st, ok := status.FromError(importGRPCError(
		errors.New("postgres password=secret connection failed"),
	))
	if !ok {
		t.Fatal("dependency failure did not return a gRPC status")
	}
	if st.Code() != codes.Internal {
		t.Fatalf("gRPC code = %s, want %s", st.Code(), codes.Internal)
	}
	if st.Message() != "import operation failed" {
		t.Fatalf("message = %q, want safe public message", st.Message())
	}
	if strings.Contains(st.Message(), "secret") || strings.Contains(st.Message(), "postgres") {
		t.Fatalf("internal dependency detail leaked: %q", st.Message())
	}
	if importRegionErrorCode(st) != "tqd.import.internal" {
		t.Fatalf("error_code = %q", importRegionErrorCode(st))
	}
}

func importRegionErrorCode(st *status.Status) string {
	for _, detail := range st.Details() {
		if info, ok := detail.(*errdetails.ErrorInfo); ok {
			return info.Metadata["error_code"]
		}
	}
	return ""
}
