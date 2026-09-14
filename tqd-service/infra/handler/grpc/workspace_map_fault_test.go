package handler_grpc

import (
	"errors"
	"strings"
	"testing"

	"common/fault"

	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

func TestWorkspaceMapValidationMapsToInvalidArgument(
	t *testing.T,
) {
	st, ok := status.FromError(workspaceStatusError(
		fault.Validation(
			"tqd.workspace.user_id_required",
			"user_id is required",
		),
	))
	if !ok {
		t.Fatal("workspaceStatusError did not return gRPC status")
	}

	if st.Code() != codes.InvalidArgument {
		t.Fatalf(
			"code = %s, want %s",
			st.Code(),
			codes.InvalidArgument,
		)
	}

	if workspaceMapErrorCode(st) !=
		"tqd.workspace.user_id_required" {
		t.Fatalf(
			"error_code = %q",
			workspaceMapErrorCode(st),
		)
	}
}

func TestWorkspaceMapNotFoundMapsToNotFound(
	t *testing.T,
) {
	st, _ := status.FromError(workspaceStatusError(
		fault.New(
			fault.KindNotFound,
			"tqd.workspace.report_not_found",
			"report not found",
		),
	))

	if st.Code() != codes.NotFound {
		t.Fatalf(
			"code = %s, want %s",
			st.Code(),
			codes.NotFound,
		)
	}
}

func TestWorkspaceMapPreconditionMapsToFailedPrecondition(
	t *testing.T,
) {
	st, _ := status.FromError(workspaceStatusError(
		fault.New(
			fault.KindPrecondition,
			"tqd.workspace.report_share_not_ready",
			"report is not ready to share",
		),
	))

	if st.Code() != codes.FailedPrecondition {
		t.Fatalf(
			"code = %s, want %s",
			st.Code(),
			codes.FailedPrecondition,
		)
	}
}

func TestWorkspaceMapRecordNotFoundCompatibility(
	t *testing.T,
) {
	st, _ := status.FromError(
		workspaceStatusError(gorm.ErrRecordNotFound),
	)

	if st.Code() != codes.NotFound {
		t.Fatalf(
			"code = %s, want %s",
			st.Code(),
			codes.NotFound,
		)
	}

	if workspaceMapErrorCode(st) !=
		"tqd.workspace.record_not_found" {
		t.Fatalf(
			"error_code = %q",
			workspaceMapErrorCode(st),
		)
	}
}

func TestWorkspaceMapUnknownFailureDoesNotLeak(
	t *testing.T,
) {
	dependencyErr := errors.New(
		"postgres password=secret workspace operation failed",
	)

	st, ok := status.FromError(
		workspaceStatusError(dependencyErr),
	)
	if !ok {
		t.Fatal("workspaceStatusError did not return gRPC status")
	}

	if st.Code() != codes.Internal {
		t.Fatalf(
			"code = %s, want %s",
			st.Code(),
			codes.Internal,
		)
	}

	if st.Message() != "workspace operation failed" {
		t.Fatalf(
			"message = %q",
			st.Message(),
		)
	}

	if strings.Contains(st.Message(), "secret") ||
		strings.Contains(st.Message(), "postgres") {
		t.Fatalf(
			"dependency detail leaked: %q",
			st.Message(),
		)
	}

	if workspaceMapErrorCode(st) != "tqd.workspace.internal" {
		t.Fatalf(
			"error_code = %q",
			workspaceMapErrorCode(st),
		)
	}
}

func workspaceMapErrorCode(st *status.Status) string {
	for _, detail := range st.Details() {
		info, ok := detail.(*errdetails.ErrorInfo)
		if ok {
			return info.Metadata["error_code"]
		}
	}

	return ""
}
