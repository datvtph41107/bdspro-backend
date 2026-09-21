package usecase

import (
	"context"
	"google.golang.org/grpc/codes"
	"testing"

	_errors "common/errors"
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
	assertQHAuthorityIssuringFault(t, err, codes.InvalidArgument, "tqd.qh_authority_issuring.name_required")
}

func TestQHAuthorityIssuringCreateReturnsCanonicalConflictFault(t *testing.T) {
	service := NewQHAuthorityIssuringUsecase(&qhAuthorityIssuringRepoStub{
		byCode: &qh_domain.QHAuthorityIssuring{Code: "dup"},
	})
	_, err := service.Create(context.Background(), &qh_domain.QHAuthorityIssuring{Name: "Authority", Code: "dup"})
	assertQHAuthorityIssuringFault(t, err, codes.AlreadyExists, "tqd.qh_authority_issuring.code_conflict")
}

func TestQHAuthorityIssuringGetReturnsCanonicalNotFoundFault(t *testing.T) {
	service := NewQHAuthorityIssuringUsecase(&qhAuthorityIssuringRepoStub{})
	_, err := service.GetByID(context.Background(), 42)
	assertQHAuthorityIssuringFault(t, err, codes.NotFound, "tqd.qh_authority_issuring.not_found")
}

func assertQHAuthorityIssuringFault(t *testing.T, err error, rpcCode codes.Code, code string) {
	t.Helper()
	application, ok := _errors.As(err)
	if !ok {
		t.Fatalf("error type = %T, want canonical application error", err)
	}
	if application.RPCCode() != rpcCode {
		t.Fatalf("kind = %q, want %q", application.RPCCode(), rpcCode)
	}
	if application.Spec().LegacyProblemCode() != code {
		t.Fatalf("code = %q, want %q", application.Spec().LegacyProblemCode(), code)
	}
}
