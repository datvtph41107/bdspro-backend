package usecase

import (
	"context"
	"errors"
	"testing"

	"hub/internal/domain"
	_repo "hub/internal/repo"
)

type versionRepoStub struct {
	_repo.IVersionRepo
	latest       *domain.VersionEntity
	latestErr    error
	createCalled bool
}

func (s *versionRepoStub) GetLatestByAppAndPlatform(context.Context, string, string, string, bool) (*domain.VersionEntity, error) {
	return s.latest, s.latestErr
}

func (s *versionRepoStub) Create(context.Context, *domain.VersionEntity) error {
	s.createCalled = true
	return nil
}

func TestCreateBundleVersionStopsWhenLatestLookupFails(t *testing.T) {
	wantErr := errors.New("database unavailable")
	repo := &versionRepoStub{latestErr: wantErr}
	uc := NewVersionUsecase(repo, nil, nil)
	input := &domain.VersionEntity{
		AppName:     "bdspro",
		Platform:    "android",
		VersionName: "1.2.3",
	}

	got, errDTO := uc.CreateBundleVersion(context.Background(), input)
	if got != nil {
		t.Fatalf("CreateBundleVersion() result = %#v, want nil", got)
	}
	if errDTO == nil {
		t.Fatal("CreateBundleVersion() error = nil, want latest lookup failure")
	}
	if errDTO.Code != 500 {
		t.Fatalf("error code = %d, want 500", errDTO.Code)
	}
	const wantMessage = "Không thể lấy thông tin phiên bản gần nhất: database unavailable"
	if errDTO.Message != wantMessage {
		t.Fatalf("error message = %q, want %q", errDTO.Message, wantMessage)
	}
	if repo.createCalled {
		t.Fatal("Create() called after latest lookup failure")
	}
}
