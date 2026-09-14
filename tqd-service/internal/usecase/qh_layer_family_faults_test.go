package usecase

import (
	"context"
	"errors"
	"testing"

	"common/fault"
	qh_domain "tqd/internal/domain/qh"
)

type qhLayerFamilyFaultRepoStub struct {
	byID   *qh_domain.QHLayerFamily
	getErr error
}

func (s *qhLayerFamilyFaultRepoStub) Create(context.Context, *qh_domain.QHLayerFamily) error {
	return nil
}
func (s *qhLayerFamilyFaultRepoStub) Update(context.Context, *qh_domain.QHLayerFamily) error {
	return nil
}
func (s *qhLayerFamilyFaultRepoStub) Delete(context.Context, uint64) error { return nil }
func (s *qhLayerFamilyFaultRepoStub) GetByID(context.Context, uint64) (*qh_domain.QHLayerFamily, error) {
	return s.byID, s.getErr
}
func (s *qhLayerFamilyFaultRepoStub) List(context.Context, int, int, string, string) ([]qh_domain.QHLayerFamily, int64, error) {
	return nil, 0, nil
}
func (s *qhLayerFamilyFaultRepoStub) ListClient(context.Context, int, int, string) ([]qh_domain.QHLayerFamily, int64, error) {
	return nil, 0, nil
}

func TestQHLayerFamilyCreateReturnsCanonicalValidationFault(t *testing.T) {
	service := NewQHLayerFamilyUsecase(&qhLayerFamilyFaultRepoStub{}, nil, nil)
	_, err := service.Create(context.Background(), &qh_domain.QHLayerFamily{Name: "   "})
	assertQHLayerFamilyFault(t, err, fault.KindValidation, "tqd.qh_layer_family.name_required")
}

func TestQHLayerFamilyGetReturnsCanonicalNotFoundFault(t *testing.T) {
	service := NewQHLayerFamilyUsecase(&qhLayerFamilyFaultRepoStub{}, nil, nil)
	_, err := service.GetByID(context.Background(), 42)
	assertQHLayerFamilyFault(t, err, fault.KindNotFound, "tqd.qh_layer_family.not_found")
}

func TestQHLayerFamilyBuildPreservesRepositoryFailureAsInternalCause(t *testing.T) {
	dependencyErr := errors.New("postgres unavailable")
	service := NewQHLayerFamilyUsecase(&qhLayerFamilyFaultRepoStub{getErr: dependencyErr}, nil, nil)
	err := service.BuildFamilyPMTiles(context.Background(), 42, 0, 0, "", func(BuildProgress) {})
	if err == nil {
		t.Fatal("BuildFamilyPMTiles returned nil error")
	}
	if !errors.Is(err, dependencyErr) {
		t.Fatalf("error = %v, want wrapped repository failure", err)
	}
	if failure, ok := fault.As(err); ok && failure.Kind() == fault.KindNotFound {
		t.Fatalf("repository failure was incorrectly classified as not found: %+v", failure)
	}
}

func assertQHLayerFamilyFault(t *testing.T, err error, kind fault.Kind, code string) {
	t.Helper()
	failure, ok := fault.As(err)
	if !ok {
		t.Fatalf("error type = %T, want canonical fault", err)
	}
	if failure.Kind() != kind {
		t.Fatalf("kind = %q, want %q", failure.Kind(), kind)
	}
	if failure.Code() != code {
		t.Fatalf("code = %q, want %q", failure.Code(), code)
	}
}
