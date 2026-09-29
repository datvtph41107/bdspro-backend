package usecase

import (
	"context"
	"errors"
	"google.golang.org/grpc/codes"
	"testing"

	_errors "common/errors"

	qh_domain "tqd/internal/domain/qh"
	"tqd/internal/dto"
	"tqd/internal/enums"
	"tqd/internal/interface/repo"
)

type workspaceMapRepoFaultStub struct {
	repo.IMapWorkspaceRepo

	parcel    *dto.ParcelWorkspacePreviewRow
	parcelErr error

	report    *qh_domain.QHUserReported
	reportErr error
}

func (s *workspaceMapRepoFaultStub) GetParcelWorkspacePreview(
	context.Context,
	uint64,
) (*dto.ParcelWorkspacePreviewRow, error) {
	return s.parcel, s.parcelErr
}

func (s *workspaceMapRepoFaultStub) GetGeneratedReport(
	context.Context,
	uint64,
	uint64,
) (*qh_domain.QHUserReported, error) {
	return s.report, s.reportErr
}

func TestWorkspaceMapFollowParcelValidationIsCanonical(
	t *testing.T,
) {
	u := &mapWorkspaceUsecase{}

	_, err := u.FollowParcel(
		context.Background(),
		dto.FollowParcelRequestDTO{},
	)

	assertWorkspaceMapFault(
		t,
		err,
		codes.InvalidArgument,
		"tqd.workspace.user_id_required",
	)
}

func TestWorkspaceMapParcelNotFoundIsCanonical(
	t *testing.T,
) {
	u := &mapWorkspaceUsecase{
		repo: &workspaceMapRepoFaultStub{},
	}

	_, err := u.FollowParcel(
		context.Background(),
		dto.FollowParcelRequestDTO{
			UserID:   7,
			ParcelID: 11,
		},
	)

	assertWorkspaceMapFault(
		t,
		err,
		codes.NotFound,
		"tqd.workspace.parcel_not_found",
	)
}

func TestWorkspaceMapReportNotFoundIsCanonical(
	t *testing.T,
) {
	u := &mapWorkspaceUsecase{
		repo: &workspaceMapRepoFaultStub{},
	}

	_, err := u.GetGeneratedReport(
		context.Background(),
		7,
		99,
	)

	assertWorkspaceMapFault(
		t,
		err,
		codes.NotFound,
		"tqd.workspace.report_not_found",
	)
}

func TestWorkspaceMapReportRegenerationPreconditionIsCanonical(
	t *testing.T,
) {
	report := &qh_domain.QHUserReported{
		ID:     99,
		UserID: 7,
		Status: enums.GeneratedReportStatusProcessing,
	}

	u := &mapWorkspaceUsecase{
		repo: &workspaceMapRepoFaultStub{
			report: report,
		},
	}

	err := u.RegenerateReport(
		context.Background(),
		7,
		99,
	)

	assertWorkspaceMapFault(
		t,
		err,
		codes.FailedPrecondition,
		"tqd.workspace.report_regeneration_not_allowed",
	)
}

func TestWorkspaceMapReportSharePreconditionIsCanonical(
	t *testing.T,
) {
	report := &qh_domain.QHUserReported{
		ID:     99,
		UserID: 7,
		Status: enums.GeneratedReportStatusProcessing,
	}

	u := &mapWorkspaceUsecase{
		repo: &workspaceMapRepoFaultStub{
			report: report,
		},
	}

	_, err := u.ShareReport(
		context.Background(),
		7,
		99,
	)

	assertWorkspaceMapFault(
		t,
		err,
		codes.FailedPrecondition,
		"tqd.workspace.report_share_not_ready",
	)
}

func TestWorkspaceMapDependencyFailurePreservesCause(
	t *testing.T,
) {
	dependencyErr := errors.New(
		"postgres password=secret workspace lookup failed",
	)

	u := &mapWorkspaceUsecase{
		repo: &workspaceMapRepoFaultStub{
			parcelErr: dependencyErr,
		},
	}

	_, err := u.FollowParcel(
		context.Background(),
		dto.FollowParcelRequestDTO{
			UserID:   7,
			ParcelID: 11,
		},
	)

	if !errors.Is(err, dependencyErr) {
		t.Fatalf(
			"error = %v, want preserved dependency cause",
			err,
		)
	}

	if _, ok := _errors.As(err); ok {
		t.Fatalf(
			"raw repository failure was incorrectly classified: %v",
			err,
		)
	}
}

func assertWorkspaceMapFault(
	t *testing.T,
	err error,
	rpcCode codes.Code,
	code string,
) {
	t.Helper()

	application, ok := _errors.As(err)
	if !ok {
		t.Fatalf(
			"error type = %T, want canonical application error: %v",
			err,
			err,
		)
	}

	if application.RPCCode() != rpcCode {
		t.Fatalf(
			"kind = %q, want %q",
			application.RPCCode(),
			rpcCode,
		)
	}

	if application.Spec().LegacyProblemCode() != code {
		t.Fatalf(
			"code = %q, want %q",
			application.Spec().LegacyProblemCode(),
			code,
		)
	}
}
