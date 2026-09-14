package usecase

import (
	"context"
	"errors"
	"testing"

	"common/fault"
	qh_domain "tqd/internal/domain/qh"
	"tqd/internal/interface/repo"
)

type qhLabelMergeRepoStub struct {
	repo.QHLabelRepository
	byID   map[uint64]*qh_domain.QHLabel
	byName *qh_domain.QHLabel
	getErr error
}

func (s *qhLabelMergeRepoStub) GetByID(context.Context, uint64) (*qh_domain.QHLabel, error) {
	if s.getErr != nil {
		return nil, s.getErr
	}
	return nil, nil
}

func (s *qhLabelMergeRepoStub) GetByNameAndLayer(context.Context, string, uint64) (*qh_domain.QHLabel, error) {
	return s.byName, nil
}

func (s *qhLabelMergeRepoStub) label(id uint64) *qh_domain.QHLabel {
	if s.byID == nil {
		return nil
	}
	return s.byID[id]
}

type qhLabelMergeRepoWithLookup struct {
	*qhLabelMergeRepoStub
}

func (s *qhLabelMergeRepoWithLookup) GetByID(_ context.Context, id uint64) (*qh_domain.QHLabel, error) {
	if s.getErr != nil {
		return nil, s.getErr
	}
	return s.label(id), nil
}

func TestQHLabelMergeReturnsCanonicalSourceNotFoundFault(t *testing.T) {
	repository := &qhLabelMergeRepoWithLookup{qhLabelMergeRepoStub: &qhLabelMergeRepoStub{}}
	service := NewQHLabelUsecase(repository, nil, nil)

	_, err := service.Merge(context.Background(), 7, []uint64{42}, &qh_domain.QHLabel{Name: "merged"})
	assertQHLabelFault(t, err, fault.KindNotFound, "tqd.qh_label.source_not_found")
}

func TestQHLabelMergeReturnsCanonicalLayerMismatchFault(t *testing.T) {
	repository := &qhLabelMergeRepoWithLookup{qhLabelMergeRepoStub: &qhLabelMergeRepoStub{
		byID: map[uint64]*qh_domain.QHLabel{
			1: {ID: 1, LayerID: 9, Name: "source"},
		},
	}}
	service := NewQHLabelUsecase(repository, nil, nil)

	_, err := service.Merge(context.Background(), 7, []uint64{1}, &qh_domain.QHLabel{Name: "merged"})
	assertQHLabelFault(t, err, fault.KindValidation, "tqd.qh_label.source_layer_mismatch")
}

func TestQHLabelMergeReturnsCanonicalNameConflictFault(t *testing.T) {
	repository := &qhLabelMergeRepoWithLookup{qhLabelMergeRepoStub: &qhLabelMergeRepoStub{
		byID: map[uint64]*qh_domain.QHLabel{
			1: {ID: 1, LayerID: 7, Name: "source"},
		},
		byName: &qh_domain.QHLabel{ID: 2, LayerID: 7, Name: "merged"},
	}}
	service := NewQHLabelUsecase(repository, nil, nil)

	_, err := service.Merge(context.Background(), 7, []uint64{1}, &qh_domain.QHLabel{Name: "merged"})
	assertQHLabelFault(t, err, fault.KindConflict, "tqd.qh_label.name_conflict")
}

func TestQHLabelMergePreservesRepositoryFailureAsInternalCause(t *testing.T) {
	dependencyErr := errors.New("postgres password=secret connection failed")
	repository := &qhLabelMergeRepoWithLookup{qhLabelMergeRepoStub: &qhLabelMergeRepoStub{getErr: dependencyErr}}
	service := NewQHLabelUsecase(repository, nil, nil)

	_, err := service.Merge(context.Background(), 7, []uint64{1}, &qh_domain.QHLabel{Name: "merged"})
	if err == nil {
		t.Fatal("Merge returned nil error")
	}
	if !errors.Is(err, dependencyErr) {
		t.Fatalf("error = %v, want wrapped repository failure", err)
	}
	if _, ok := fault.As(err); ok {
		t.Fatalf("repository failure was incorrectly classified inside usecase: %v", err)
	}
}

func assertQHLabelFault(t *testing.T, err error, kind fault.Kind, code string) {
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
