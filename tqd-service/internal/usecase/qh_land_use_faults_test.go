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
		codes.InvalidArgument,
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
		codes.NotFound,
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
		codes.NotFound,
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
		codes.AlreadyExists,
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

	if _, ok := _errors.As(err); ok {
		t.Fatalf(
			"technical dependency error was incorrectly classified: %v",
			err,
		)
	}
}

func assertQHLandUseFault(
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
