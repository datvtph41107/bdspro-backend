package usecase

import (
	"context"
	"errors"
	"testing"

	"common/fault"

	qh_domain "tqd/internal/domain/qh"
	"tqd/internal/enums"
	"tqd/internal/interface/repo"
)

type importLayerRepoStub struct {
	repo.LayerRepository
	layer *qh_domain.QHLayer
	err   error
}

func (s *importLayerRepoStub) GetByID(context.Context, uint64) (*qh_domain.QHLayer, error) {
	return s.layer, s.err
}

type importRegionRepoStub struct {
	repo.RegionRepository
	row     *qh_domain.QHRegionImportErrorLog
	getErr  error
	locked  bool
	lockErr error
}

func (s *importRegionRepoStub) GetImportErrorLogByID(context.Context, uint64) (*qh_domain.QHRegionImportErrorLog, error) {
	return s.row, s.getErr
}

func (s *importRegionRepoStub) TryBeginImportErrorRetry(context.Context, uint64) (bool, error) {
	return s.locked, s.lockErr
}

func TestImportRegionLayerNotFoundIsCanonical(t *testing.T) {
	u := &importUsecaseImpl{
		layerRepo: &importLayerRepoStub{},
	}
	_, err := u.validateLayerForImport(context.Background(), 7)
	assertImportRegionFault(t, err, fault.KindNotFound, "tqd.import.layer_not_found")
}

func TestImportRegionAlreadyProcessingIsCanonical(t *testing.T) {
	u := &importUsecaseImpl{
		layerRepo: &importLayerRepoStub{
			layer: &qh_domain.QHLayer{
				ID:           7,
				ImportStatus: enums.LayerImportStatusProcessing,
			},
		},
	}
	_, err := u.validateLayerForImport(context.Background(), 7)
	assertImportRegionFault(t, err, fault.KindPrecondition, "tqd.import.in_progress")
}

func TestImportRegionPreservesLayerRepositoryFailure(t *testing.T) {
	dependencyErr := errors.New("postgres password=secret connection failed")
	u := &importUsecaseImpl{
		layerRepo: &importLayerRepoStub{err: dependencyErr},
	}
	_, err := u.validateLayerForImport(context.Background(), 7)

	if !errors.Is(err, dependencyErr) {
		t.Fatalf("error = %v, want wrapped/original repository failure", err)
	}
	if _, ok := fault.As(err); ok {
		t.Fatalf("repository failure was incorrectly classified as business fault: %v", err)
	}
}

func TestImportRegionRetryNotFoundIsCanonical(t *testing.T) {
	u := &importUsecaseImpl{
		regionRepo: &importRegionRepoStub{},
	}
	_, err := u.RetryImportError(context.Background(), 11)
	assertImportRegionFault(t, err, fault.KindNotFound, "tqd.import.error_not_found")
}

func TestImportRegionRetryInProgressIsCanonical(t *testing.T) {
	u := &importUsecaseImpl{
		regionRepo: &importRegionRepoStub{
			row: &qh_domain.QHRegionImportErrorLog{
				ID:     11,
				Status: qh_domain.ImportErrorStatusProcessing,
			},
		},
	}
	_, err := u.RetryImportError(context.Background(), 11)
	assertImportRegionFault(t, err, fault.KindPrecondition, "tqd.import.retry_in_progress")
}

func TestImportRegionRetryNotAllowedIsCanonical(t *testing.T) {
	u := &importUsecaseImpl{
		regionRepo: &importRegionRepoStub{
			row: &qh_domain.QHRegionImportErrorLog{
				ID:     11,
				Status: qh_domain.ImportErrorStatusResolved,
			},
		},
	}
	_, err := u.RetryImportError(context.Background(), 11)
	assertImportRegionFault(t, err, fault.KindPrecondition, "tqd.import.retry_not_allowed")
}

func TestImportRegionRetryLockConflictIsCanonical(t *testing.T) {
	u := &importUsecaseImpl{
		regionRepo: &importRegionRepoStub{
			row: &qh_domain.QHRegionImportErrorLog{
				ID:     11,
				Status: qh_domain.ImportErrorStatusPending,
			},
			locked: false,
		},
	}
	_, err := u.RetryImportError(context.Background(), 11)
	assertImportRegionFault(t, err, fault.KindAborted, "tqd.import.retry_lock_conflict")
}

func assertImportRegionFault(t *testing.T, err error, kind fault.Kind, code string) {
	t.Helper()
	failure, ok := fault.As(err)
	if !ok {
		t.Fatalf("error type = %T, want canonical fault: %v", err, err)
	}
	if failure.Kind() != kind {
		t.Fatalf("kind = %q, want %q", failure.Kind(), kind)
	}
	if failure.Code() != code {
		t.Fatalf("code = %q, want %q", failure.Code(), code)
	}
}
