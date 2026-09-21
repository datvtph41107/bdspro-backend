package usecase

import (
	"context"
	"errors"
	"google.golang.org/grpc/codes"
	"testing"

	_errors "common/errors"

	qh_domain "tqd/internal/domain/qh"
	"tqd/internal/interface/repo"
)

type qhLayerLegendRepoFaultStub struct {
	repo.QHLayerLegendRepository

	byID            *qh_domain.QHLayerLegend
	byIDErr         error
	byLayerAndLabel *qh_domain.QHLayerLegend
	byLayerLabelErr error
}

func (s *qhLayerLegendRepoFaultStub) GetByID(
	context.Context,
	uint64,
) (*qh_domain.QHLayerLegend, error) {
	return s.byID, s.byIDErr
}

func (s *qhLayerLegendRepoFaultStub) GetByLayerAndLabel(
	context.Context,
	uint64,
	uint64,
) (*qh_domain.QHLayerLegend, error) {
	return s.byLayerAndLabel, s.byLayerLabelErr
}

type qhLayerLegendLayerRepoFaultStub struct {
	repo.LayerRepository
	layer *qh_domain.QHLayer
	err   error
}

func (s *qhLayerLegendLayerRepoFaultStub) GetByID(
	context.Context,
	uint64,
) (*qh_domain.QHLayer, error) {
	return s.layer, s.err
}

func TestQHLayerLegendInvalidTypeIsCanonical(t *testing.T) {
	u := &qhLayerLegendUsecase{}

	err := u.validateLegendTypes("unsupported", "")

	assertQHLayerLegendFault(
		t,
		err,
		codes.InvalidArgument,
		"tqd.legend.legend_type_invalid",
	)
}

func TestQHLayerLegendRecordNotFoundIsCanonical(t *testing.T) {
	u := &qhLayerLegendUsecase{
		repo: &qhLayerLegendRepoFaultStub{},
	}

	_, err := u.GetByID(context.Background(), 88)

	assertQHLayerLegendFault(
		t,
		err,
		codes.NotFound,
		"tqd.legend.record_not_found",
	)
}

func TestQHLayerLegendLayerNotFoundIsCanonical(t *testing.T) {
	u := &qhLayerLegendUsecase{
		layerRepo: &qhLayerLegendLayerRepoFaultStub{},
	}

	_, _, err := u.BatchCreate(
		context.Background(),
		7,
		nil,
	)

	assertQHLayerLegendFault(
		t,
		err,
		codes.NotFound,
		"tqd.legend.layer_not_found",
	)
}

func TestQHLayerLegendDuplicateIsCanonical(t *testing.T) {
	existing := &qh_domain.QHLayerLegend{}
	existing.ID = 99

	u := &qhLayerLegendUsecase{
		repo: &qhLayerLegendRepoFaultStub{
			byLayerAndLabel: existing,
		},
	}

	err := u.checkDuplicate(
		context.Background(),
		7,
		11,
		0,
	)

	assertQHLayerLegendFault(
		t,
		err,
		codes.AlreadyExists,
		"tqd.legend.duplicate",
	)
}

func TestQHLayerLegendPreservesDependencyCause(t *testing.T) {
	dependencyErr := errors.New(
		"postgres password=secret legend lookup failed",
	)

	u := &qhLayerLegendUsecase{
		repo: &qhLayerLegendRepoFaultStub{
			byIDErr: dependencyErr,
		},
	}

	_, err := u.GetByID(context.Background(), 88)

	if !errors.Is(err, dependencyErr) {
		t.Fatalf(
			"error = %v, want preserved dependency cause",
			err,
		)
	}

	if _, ok := _errors.As(err); ok {
		t.Fatalf(
			"technical dependency error was incorrectly classified: %v",
			err,
		)
	}
}

func assertQHLayerLegendFault(
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
