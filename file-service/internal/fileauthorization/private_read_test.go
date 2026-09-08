package fileauthorization

import (
	"context"
	"errors"
	"testing"
)

type fakeAuthorizer struct {
	err   error
	calls int
}

func (f *fakeAuthorizer) Require(context.Context, Capability) error {
	f.calls++
	return f.err
}

func TestAuthorizePrivateReadDoesNotCallGlobalWhenOrdinaryAccessAllows(t *testing.T) {
	global := &fakeAuthorizer{}
	err := AuthorizePrivateRead(context.Background(), func() (bool, error) { return true, nil }, global)
	if err != nil {
		t.Fatalf("AuthorizePrivateRead() error = %v", err)
	}
	if global.calls != 0 {
		t.Fatalf("global authority calls = %d, want 0", global.calls)
	}
}

func TestAuthorizePrivateReadAllowsExplicitGlobalCapability(t *testing.T) {
	global := &fakeAuthorizer{}
	err := AuthorizePrivateRead(context.Background(), func() (bool, error) { return false, nil }, global)
	if err != nil {
		t.Fatalf("AuthorizePrivateRead() error = %v", err)
	}
	if global.calls != 1 {
		t.Fatalf("global authority calls = %d, want 1", global.calls)
	}
}

func TestAuthorizePrivateReadDeniesWhenBothAuthoritiesDeny(t *testing.T) {
	global := &fakeAuthorizer{err: ErrDenied}
	err := AuthorizePrivateRead(context.Background(), func() (bool, error) { return false, nil }, global)
	if !errors.Is(err, ErrDenied) {
		t.Fatalf("AuthorizePrivateRead() error = %v, want ErrDenied", err)
	}
}

func TestAuthorizePrivateReadFailsClosedWhenGlobalAuthorityUnavailable(t *testing.T) {
	global := &fakeAuthorizer{err: ErrUnavailable}
	err := AuthorizePrivateRead(context.Background(), func() (bool, error) { return false, nil }, global)
	if !errors.Is(err, ErrUnavailable) {
		t.Fatalf("AuthorizePrivateRead() error = %v, want ErrUnavailable", err)
	}
}

func TestAuthorizePrivateReadPreservesOrdinaryRepositoryErrorWhenGlobalDenies(t *testing.T) {
	ordinaryErr := errors.New("postgres down")
	global := &fakeAuthorizer{err: ErrDenied}
	err := AuthorizePrivateRead(context.Background(), func() (bool, error) { return false, ordinaryErr }, global)
	if !errors.Is(err, ordinaryErr) {
		t.Fatalf("AuthorizePrivateRead() error = %v, want ordinary error", err)
	}
}

func TestAuthorizePrivateReadGlobalCapabilityCanBypassMissingOrdinaryEvidence(t *testing.T) {
	ordinaryErr := errors.New("file access row unavailable")
	global := &fakeAuthorizer{}
	err := AuthorizePrivateRead(context.Background(), func() (bool, error) { return false, ordinaryErr }, global)
	if err != nil {
		t.Fatalf("AuthorizePrivateRead() error = %v", err)
	}
}
