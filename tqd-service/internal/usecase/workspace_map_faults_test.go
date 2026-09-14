package usecase

import (
	"context"
	"errors"
	"testing"

	"common/fault"

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
		fault.KindValidation,
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
		fault.KindNotFound,
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
		fault.KindNotFound,
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
		fault.KindPrecondition,
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
		fault.KindPrecondition,
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

	if _, ok := fault.As(err); ok {
		t.Fatalf(
			"raw repository failure was incorrectly classified: %v",
			err,
		)
	}
}

func assertWorkspaceMapFault(
	t *testing.T,
	err error,
	kind fault.Kind,
	code string,
) {
	t.Helper()

	failure, ok := fault.As(err)
	if !ok {
		t.Fatalf(
			"error type = %T, want canonical fault: %v",
			err,
			err,
		)
	}

	if failure.Kind() != kind {
		t.Fatalf(
			"kind = %q, want %q",
			failure.Kind(),
			kind,
		)
	}

	if failure.Code() != code {
		t.Fatalf(
			"code = %q, want %q",
			failure.Code(),
			code,
		)
	}
}
