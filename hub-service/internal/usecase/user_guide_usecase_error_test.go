package usecase

import (
	"context"
	"errors"
	"testing"

	"hub/internal/domain"
	_repo "hub/internal/repo"
	sharepb "pb/types/shared"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type userGuideRepoStub struct {
	_repo.IUserGuideRepo
	result *domain.UserGuideEntity
	err    error
}

func (s *userGuideRepoStub) GetByKey(context.Context, string) (*domain.UserGuideEntity, error) {
	return s.result, s.err
}

func TestUserGuideUsecaseGetByKeyReturnsEntity(t *testing.T) {
	want := &domain.UserGuideEntity{}
	uc := NewUserGuideUsecase(&userGuideRepoStub{result: want}, nil)

	got, err := uc.GetByKey(context.Background(), "guide-key")
	if err != nil {
		t.Fatalf("GetByKey() error = %v, want nil", err)
	}
	if got != want {
		t.Fatalf("GetByKey() entity = %p, want %p", got, want)
	}
}

func TestUserGuideUsecaseGetByKeyPreservesNotFoundContract(t *testing.T) {
	const key = "missing-key"
	const wantMessage = "Không tìm thấy user guide với key: " + key

	uc := NewUserGuideUsecase(&userGuideRepoStub{}, nil)
	got, err := uc.GetByKey(context.Background(), key)
	if got != nil {
		t.Fatalf("GetByKey() entity = %#v, want nil", got)
	}
	if err == nil {
		t.Fatal("GetByKey() error = nil, want legacy not-found error")
	}

	st, ok := status.FromError(err)
	if !ok {
		t.Fatalf("GetByKey() error = %v, want gRPC status", err)
	}
	if st.Code() != codes.Internal {
		t.Fatalf("status code = %v, want %v", st.Code(), codes.Internal)
	}
	if st.Message() != wantMessage {
		t.Fatalf("status message = %q, want %q", st.Message(), wantMessage)
	}

	details := st.Details()
	if len(details) != 1 {
		t.Fatalf("details len = %d, want 1", len(details))
	}
	detail, ok := details[0].(*sharepb.ErrorResponse)
	if !ok {
		t.Fatalf("detail type = %T, want *shared.ErrorResponse", details[0])
	}
	if detail.Code != 404 {
		t.Fatalf("detail code = %d, want 404", detail.Code)
	}
	if detail.Message != wantMessage {
		t.Fatalf("detail message = %q, want %q", detail.Message, wantMessage)
	}
	if detail.Second != nil {
		t.Fatalf("detail second = %v, want nil", detail.Second)
	}
}

func TestUserGuideUsecaseGetByKeyPassesThroughRepositoryFailure(t *testing.T) {
	wantErr := errors.New("database unavailable")
	uc := NewUserGuideUsecase(&userGuideRepoStub{err: wantErr}, nil)

	got, err := uc.GetByKey(context.Background(), "guide-key")
	if got != nil {
		t.Fatalf("GetByKey() entity = %#v, want nil", got)
	}
	if !errors.Is(err, wantErr) {
		t.Fatalf("GetByKey() error = %v, want repository error %v", err, wantErr)
	}
}
