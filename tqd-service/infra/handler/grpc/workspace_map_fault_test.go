package handler_grpc

import (
	"errors"
	"strings"
	"testing"

	_errors "common/errors"
	"tqd/internal"

	"gorm.io/gorm"
)

func TestWorkspaceMapValidationMapsToInvalidArgument(t *testing.T) {
	err := _errors.ReturnError(service.WorkspaceUserIDRequired)
	assertCanonicalStatus(t, workspaceStatusError(err), service.WorkspaceUserIDRequired)
}

func TestWorkspaceMapNotFoundMapsToNotFound(t *testing.T) {
	err := _errors.ReturnError(service.WorkspaceReportNotFound)
	assertCanonicalStatus(t, workspaceStatusError(err), service.WorkspaceReportNotFound)
}

func TestWorkspaceMapPreconditionMapsToFailedPrecondition(t *testing.T) {
	err := _errors.ReturnError(service.WorkspaceReportShareNotReady)
	assertCanonicalStatus(t, workspaceStatusError(err), service.WorkspaceReportShareNotReady)
}

func TestWorkspaceMapRecordNotFoundCompatibility(t *testing.T) {
	assertCanonicalStatus(
		t,
		workspaceStatusError(gorm.ErrRecordNotFound),
		service.WorkspaceRecordNotFound,
	)
}

func TestWorkspaceMapUnknownFailureDoesNotLeak(t *testing.T) {
	st := assertTechnicalStatus(
		t,
		workspaceStatusError(errors.New("postgres password=secret workspace operation failed")),
	)
	if strings.Contains(st.Message(), "secret") || strings.Contains(st.Message(), "postgres") {
		t.Fatalf("dependency detail leaked: %q", st.Message())
	}
}
