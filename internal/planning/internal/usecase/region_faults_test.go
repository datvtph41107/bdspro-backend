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

type regionRepoFaultStub struct {
	repo.RegionRepository

	byID       *qh_domain.QHRegion
	byIDErr    error
	syncRef    *repo.RegionSyncReference
	syncRefErr error
}

func (s *regionRepoFaultStub) GetByID(
	context.Context,
	uint64,
) (*qh_domain.QHRegion, error) {
	return s.byID, s.byIDErr
}

func (s *regionRepoFaultStub) GetSyncReferenceByLayerLabel(
	context.Context,
	uint64,
	uint64,
) (*repo.RegionSyncReference, error) {
	return s.syncRef, s.syncRefErr
}

type regionLayerRepoFaultStub struct {
	repo.LayerRepository

	layer *qh_domain.QHLayer
	err   error
}

func (s *regionLayerRepoFaultStub) GetByID(
	context.Context,
	uint64,
) (*qh_domain.QHLayer, error) {
	return s.layer, s.err
}

type regionLabelRepoFaultStub struct {
	repo.QHLabelRepository

	label *qh_domain.QHLabel
	err   error
}

func (s *regionLabelRepoFaultStub) GetByID(
	context.Context,
	uint64,
) (*qh_domain.QHLabel, error) {
	return s.label, s.err
}

func TestRegionCreateLayerNotFoundIsCanonical(t *testing.T) {
	u := &regionUsecaseImpl{
		layerRepo: &regionLayerRepoFaultStub{},
	}

	_, err := u.Create(
		context.Background(),
		7,
		"Region",
		nil,
		0,
	)

	assertRegionFault(
		t,
		err,
		codes.NotFound,
		"tqd.region.layer_not_found",
	)

	if !errors.Is(err, ErrLayerNotFound) {
		t.Fatalf(
			"error = %v, want ErrLayerNotFound cause",
			err,
		)
	}
}

func TestRegionCreateInvalidGeometryIsCanonical(t *testing.T) {
	u := &regionUsecaseImpl{
		layerRepo: &regionLayerRepoFaultStub{
			layer: &qh_domain.QHLayer{},
		},
	}

	_, err := u.Create(
		context.Background(),
		7,
		"Region",
		[]byte("not-json"),
		0,
	)

	assertRegionFault(
		t,
		err,
		codes.InvalidArgument,
		"tqd.region.geometry_invalid",
	)
}

func TestRegionUpdateNotFoundIsCanonical(t *testing.T) {
	name := "Updated"

	u := &regionUsecaseImpl{
		regionRepo: &regionRepoFaultStub{},
	}

	_, err := u.Update(
		context.Background(),
		RegionUpdatePatch{
			RegionID: 88,
			Name:     &name,
		},
	)

	assertRegionFault(
		t,
		err,
		codes.NotFound,
		"tqd.region.not_found",
	)
}

func TestRegionSyncLabelNotFoundIsCanonical(t *testing.T) {
	u := &regionUsecaseImpl{
		layerRepo: &regionLayerRepoFaultStub{
			layer: &qh_domain.QHLayer{},
		},
		labelRepo: &regionLabelRepoFaultStub{},
	}

	_, err := u.SyncRegion(
		context.Background(),
		7,
		11,
	)

	assertRegionFault(
		t,
		err,
		codes.NotFound,
		"tqd.region.label_not_found",
	)
}

func TestRegionSyncReferenceNotFoundIsCanonical(t *testing.T) {
	u := &regionUsecaseImpl{
		layerRepo: &regionLayerRepoFaultStub{
			layer: &qh_domain.QHLayer{},
		},
		labelRepo: &regionLabelRepoFaultStub{
			label: &qh_domain.QHLabel{},
		},
		regionRepo: &regionRepoFaultStub{},
	}

	_, err := u.SyncRegion(
		context.Background(),
		7,
		11,
	)

	assertRegionFault(
		t,
		err,
		codes.NotFound,
		"tqd.region.sync_reference_not_found",
	)
}

func TestRegionDependencyFailurePreservesCause(t *testing.T) {
	dependencyErr := errors.New(
		"postgres password=secret region lookup failed",
	)

	u := &regionUsecaseImpl{
		layerRepo: &regionLayerRepoFaultStub{
			err: dependencyErr,
		},
	}

	_, err := u.Create(
		context.Background(),
		7,
		"Region",
		nil,
		0,
	)

	if _, ok := _errors.As(err); ok {
		t.Fatalf("technical dependency error received application identity: %v", err)
	}

	if !errors.Is(err, dependencyErr) {
		t.Fatalf(
			"error = %v, want preserved dependency cause",
			err,
		)
	}
}

func assertRegionFault(
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
