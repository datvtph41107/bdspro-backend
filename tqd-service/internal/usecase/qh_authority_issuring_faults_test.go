package usecase

import (
	"context"
	"testing"

	"common/fault"
	qh_domain "tqd/internal/domain/qh"
)

type qhAuthorityIssuringRepoStub struct {
	byID   *qh_domain.QHAuthorityIssuring
	byCode *qh_domain.QHAuthorityIssuring
}

func (s *qhAuthorityIssuringRepoStub) Create(context.Context, *qh_domain.QHAuthorityIssuring) error {
	return nil
}
func (s *qhAuthorityIssuringRepoStub) Update(context.Context, *qh_domain.QHAuthorityIssuring) error {
	return nil
}
func (s *qhAuthorityIssuringRepoStub) Delete(context.Context, uint64) error { return nil }
func (s *qhAuthorityIssuringRepoStub) GetByID(context.Context, uint64) (*qh_domain.QHAuthorityIssuring, error) {
	return s.byID, nil
}
func (s *qhAuthorityIssuringRepoStub) GetByCode(context.Context, string) (*qh_domain.QHAuthorityIssuring, error) {
	return s.byCode, nil
}
func (s *qhAuthorityIssuringRepoStub) List(context.Context, int, int) ([]qh_domain.QHAuthorityIssuring, int64, error) {
	return nil, 0, nil
}

func TestQHAuthorityIssuringCreateReturnsCanonicalValidationFault(t *testing.T) {
	service := NewQHAuthorityIssuringUsecase(&qhAuthorityIssuringRepoStub{})
	_, err := service.Create(context.Background(), &qh_domain.QHAuthorityIssuring{Name: "", Code: "code"})
	assertQHAuthorityIssuringFault(t, err, fault.KindValidation, "tqd.qh_authority_issuring.name_required")
}

func TestQHAuthorityIssuringCreateReturnsCanonicalConflictFault(t *testing.T) {
	service := NewQHAuthorityIssuringUsecase(&qhAuthorityIssuringRepoStub{
		byCode: &qh_domain.QHAuthorityIssuring{Code: "dup"},
	})
	_, err := service.Create(context.Background(), &qh_domain.QHAuthorityIssuring{Name: "Authority", Code: "dup"})
	assertQHAuthorityIssuringFault(t, err, fault.KindConflict, "tqd.qh_authority_issuring.code_conflict")
}

func TestQHAuthorityIssuringGetReturnsCanonicalNotFoundFault(t *testing.T) {
	service := NewQHAuthorityIssuringUsecase(&qhAuthorityIssuringRepoStub{})
	_, err := service.GetByID(context.Background(), 42)
	assertQHAuthorityIssuringFault(t, err, fault.KindNotFound, "tqd.qh_authority_issuring.not_found")
}

func assertQHAuthorityIssuringFault(t *testing.T, err error, kind fault.Kind, code string) {
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
