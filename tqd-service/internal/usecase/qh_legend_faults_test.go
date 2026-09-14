package usecase

import (
	"context"
	"errors"
	"testing"

	"common/fault"

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
		fault.KindValidation,
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
		fault.KindNotFound,
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
		fault.KindNotFound,
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
		fault.KindConflict,
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

	if _, ok := fault.As(err); ok {
		t.Fatalf(
			"technical dependency error was incorrectly classified: %v",
			err,
		)
	}
}

func assertQHLayerLegendFault(
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
