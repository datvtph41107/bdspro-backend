package handler_grpc

import (
	"errors"
	"strings"
	"testing"

	_errors "common/errors"
	"tqd/internal"
)

func TestImportRegionGRPCMapsRetryLockConflict(t *testing.T) {
	err := _errors.ReturnError(
		service.ImportRetryLockConflict,
		_errors.WithPublicMessage("import error 11 could not be locked for retry"),
	)
	assertCanonicalStatus(t, importGRPCError(err), service.ImportRetryLockConflict)
}

func TestImportRegionGRPCMapsInProgressToFailedPrecondition(t *testing.T) {
	err := _errors.ReturnError(
		service.ImportInProgress,
		_errors.WithPublicMessage("layer 7 already has an import in progress"),
	)
	assertCanonicalStatus(t, importGRPCError(err), service.ImportInProgress)
}

func TestImportRegionGRPCDoesNotLeakDependencyFailure(t *testing.T) {
	st := assertTechnicalStatus(
		t,
		importGRPCError(errors.New("postgres password=secret connection failed")),
	)
	if strings.Contains(st.Message(), "secret") || strings.Contains(st.Message(), "postgres") {
		t.Fatalf("internal dependency detail leaked: %q", st.Message())
	}
}
