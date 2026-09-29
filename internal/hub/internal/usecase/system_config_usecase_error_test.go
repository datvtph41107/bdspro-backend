package usecase

import (
	"context"
	"errors"
	"testing"

	"hub/internal/domain"
	"hub/internal/dto"
	_repo "hub/internal/repo"
	sharepb "pb/types/shared"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type systemConfigRepoStub struct {
	_repo.ISystemConfigRepo
	result *domain.SystemConfigEntity
	err    error
}

func (s *systemConfigRepoStub) GetByKey(context.Context, string) (*domain.SystemConfigEntity, error) {
	return s.result, s.err
}

func TestSystemConfigGetByKeyReturnsEntity(t *testing.T) {
	want := &domain.SystemConfigEntity{}
	uc := NewSystemConfigUsecase(dto.SystemConfigPersist{}, &systemConfigRepoStub{result: want})

	got, err := uc.GetSystemConfigByKey(context.Background(), "config-key")
	if err != nil {
		t.Fatalf("GetSystemConfigByKey() error = %v, want nil", err)
	}
	if got != want {
		t.Fatalf("GetSystemConfigByKey() entity = %p, want %p", got, want)
	}
}

func TestSystemConfigGetByKeyPreservesNotFoundContract(t *testing.T) {
	const key = "missing-key"
	const wantMessage = "Không tìm thấy config với key: " + key
	uc := NewSystemConfigUsecase(dto.SystemConfigPersist{}, &systemConfigRepoStub{})

	got, err := uc.GetSystemConfigByKey(context.Background(), key)
	if got != nil {
		t.Fatalf("GetSystemConfigByKey() entity = %#v, want nil", got)
	}
	if err == nil {
		t.Fatal("GetSystemConfigByKey() error = nil, want legacy not-found error")
	}

	st, ok := status.FromError(err)
	if !ok {
		t.Fatalf("GetSystemConfigByKey() error = %v, want gRPC status", err)
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
}

func TestSystemConfigGetByKeyPassesThroughRepositoryFailure(t *testing.T) {
	wantErr := errors.New("database unavailable")
	uc := NewSystemConfigUsecase(dto.SystemConfigPersist{}, &systemConfigRepoStub{err: wantErr})

	got, err := uc.GetSystemConfigByKey(context.Background(), "config-key")
	if got != nil {
		t.Fatalf("GetSystemConfigByKey() entity = %#v, want nil", got)
	}
	if !errors.Is(err, wantErr) {
		t.Fatalf("GetSystemConfigByKey() error = %v, want repository error %v", err, wantErr)
	}
}
