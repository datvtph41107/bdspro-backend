package usecase

import (
	"context"
	"errors"
	"testing"

	"common/fault"

	qh_domain "tqd/internal/domain/qh"
	"tqd/internal/interface/repo"
)

type qhLandUseRepoFaultStub struct {
	repo.QHLandUseRepository

	byID              *qh_domain.QHLandUse
	byIDErr           error
	byLayerAndLandUse *qh_domain.QHLandUse
	byLayerErr        error
}

func (s *qhLandUseRepoFaultStub) GetByID(
	context.Context,
	uint64,
) (*qh_domain.QHLandUse, error) {
	return s.byID, s.byIDErr
}

func (s *qhLandUseRepoFaultStub) GetByLayerAndLandUse(
	context.Context,
	uint64,
	uint64,
) (*qh_domain.QHLandUse, error) {
	return s.byLayerAndLandUse, s.byLayerErr
}

type qhLandUseLayerRepoFaultStub struct {
	repo.LayerRepository
	layer *qh_domain.QHLayer
	err   error
}

func (s *qhLandUseLayerRepoFaultStub) GetByID(
	context.Context,
	uint64,
) (*qh_domain.QHLayer, error) {
	return s.layer, s.err
}

func TestQHLandUseNameRequiredIsCanonical(t *testing.T) {
	u := &qhLayerLandUseGroupUsecase{}

	_, err := u.resolveLandUse(
		context.Background(),
		&qh_domain.QHLandUse{},
	)

	assertQHLandUseFault(
		t,
		err,
		fault.KindValidation,
		"tqd.land_use.name_required",
	)
}

func TestQHLandUseRecordNotFoundIsCanonical(t *testing.T) {
	u := &qhLayerLandUseGroupUsecase{
		repo: &qhLandUseRepoFaultStub{},
	}

	_, err := u.GetByID(context.Background(), 88)

	assertQHLandUseFault(
		t,
		err,
		fault.KindNotFound,
		"tqd.land_use.record_not_found",
	)
}

func TestQHLandUseLayerNotFoundIsCanonical(t *testing.T) {
	u := &qhLayerLandUseGroupUsecase{
		layerRepo: &qhLandUseLayerRepoFaultStub{},
	}

	_, err := u.CreateLayerLandUse(
		context.Background(),
		7,
		11,
	)

	assertQHLandUseFault(
		t,
		err,
		fault.KindNotFound,
		"tqd.land_use.layer_not_found",
	)
}

func TestQHLandUseDuplicateIsCanonical(t *testing.T) {
	layer := &qh_domain.QHLayer{}
	layer.ID = 7

	landUse := &qh_domain.QHLandUse{}
	landUse.ID = 11

	duplicate := &qh_domain.QHLandUse{}
	duplicate.ID = 11

	u := &qhLayerLandUseGroupUsecase{
		layerRepo: &qhLandUseLayerRepoFaultStub{
			layer: layer,
		},
		repo: &qhLandUseRepoFaultStub{
			byID:              landUse,
			byLayerAndLandUse: duplicate,
		},
	}

	_, err := u.CreateLayerLandUse(
		context.Background(),
		7,
		11,
	)

	assertQHLandUseFault(
		t,
		err,
		fault.KindConflict,
		"tqd.land_use.duplicate",
	)
}

func TestQHLandUsePreservesDependencyCause(t *testing.T) {
	dependencyErr := errors.New(
		"postgres password=secret land-use lookup failed",
	)

	u := &qhLayerLandUseGroupUsecase{
		repo: &qhLandUseRepoFaultStub{
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

func assertQHLandUseFault(
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
