package usecase

import (
	"context"
	"errors"
	"testing"

	"hub/internal/domain"
	"hub/internal/dto"
	_repo "hub/internal/repo"
)

type eventQueueRepoStub struct {
	_repo.IEventQueueRepo
	result      *domain.EventQueueEntity
	err         error
	updateCalls int
	deleteCalls int
}

func (s *eventQueueRepoStub) GetByID(context.Context, uint64) (*domain.EventQueueEntity, error) {
	return s.result, s.err
}

func (s *eventQueueRepoStub) Update(context.Context, uint64, *domain.EventQueueEntity) error {
	s.updateCalls++
	return nil
}

func (s *eventQueueRepoStub) Delete(context.Context, uint64) error {
	s.deleteCalls++
	return nil
}

func TestEventQueueDetailPreservesNotFoundContract(t *testing.T) {
	uc := NewEventQueueUsecase(&eventQueueRepoStub{})

	got, errDTO := uc.Detail(context.Background(), 42)
	if got != nil {
		t.Fatalf("Detail() entity = %#v, want nil", got)
	}
	if errDTO == nil {
		t.Fatal("Detail() error = nil, want not-found error")
	}
	if errDTO.Code != 404 {
		t.Fatalf("Detail() code = %d, want 404", errDTO.Code)
	}
	if errDTO.Message != "Không tìm thấy event queue" {
		t.Fatalf("Detail() message = %q, want not-found message", errDTO.Message)
	}
}

func TestEventQueueDetailPreservesRepositoryFailure(t *testing.T) {
	wantErr := errors.New("database unavailable")
	uc := NewEventQueueUsecase(&eventQueueRepoStub{err: wantErr})

	got, errDTO := uc.Detail(context.Background(), 42)
	if got != nil {
		t.Fatalf("Detail() entity = %#v, want nil", got)
	}
	if errDTO == nil {
		t.Fatal("Detail() error = nil, want technical error")
	}
	if errDTO.Code != 500 {
		t.Fatalf("Detail() code = %d, want 500", errDTO.Code)
	}
	const wantMessage = "Lỗi khi lấy chi tiết event queue: database unavailable"
	if errDTO.Message != wantMessage {
		t.Fatalf("Detail() message = %q, want %q", errDTO.Message, wantMessage)
	}
}

func TestEventQueueUpdatePreservesNotFoundContract(t *testing.T) {
	repoStub := &eventQueueRepoStub{}
	uc := NewEventQueueUsecase(repoStub)

	got, errDTO := uc.Update(context.Background(), 42, dto.EventQueueSaveDTO{})
	if got != nil {
		t.Fatalf("Update() entity = %#v, want nil", got)
	}
	if errDTO == nil || errDTO.Code != 404 {
		t.Fatalf("Update() error = %#v, want code 404", errDTO)
	}
	if repoStub.updateCalls != 0 {
		t.Fatalf("Update() repo update calls = %d, want 0", repoStub.updateCalls)
	}
}

func TestEventQueueUpdateStopsOnLookupFailure(t *testing.T) {
	repoStub := &eventQueueRepoStub{err: errors.New("database unavailable")}
	uc := NewEventQueueUsecase(repoStub)

	got, errDTO := uc.Update(context.Background(), 42, dto.EventQueueSaveDTO{})
	if got != nil {
		t.Fatalf("Update() entity = %#v, want nil", got)
	}
	if errDTO == nil || errDTO.Code != 500 {
		t.Fatalf("Update() error = %#v, want code 500", errDTO)
	}
	const wantMessage = "Lỗi khi lấy event queue trước khi cập nhật: database unavailable"
	if errDTO.Message != wantMessage {
		t.Fatalf("Update() message = %q, want %q", errDTO.Message, wantMessage)
	}
	if repoStub.updateCalls != 0 {
		t.Fatalf("Update() repo update calls = %d, want 0", repoStub.updateCalls)
	}
}

func TestEventQueueDeletePreservesNotFoundContract(t *testing.T) {
	repoStub := &eventQueueRepoStub{}
	uc := NewEventQueueUsecase(repoStub)

	errDTO := uc.Delete(context.Background(), 42)
	if errDTO == nil || errDTO.Code != 404 {
		t.Fatalf("Delete() error = %#v, want code 404", errDTO)
	}
	if repoStub.deleteCalls != 0 {
		t.Fatalf("Delete() repo delete calls = %d, want 0", repoStub.deleteCalls)
	}
}

func TestEventQueueDeleteStopsOnLookupFailure(t *testing.T) {
	repoStub := &eventQueueRepoStub{err: errors.New("database unavailable")}
	uc := NewEventQueueUsecase(repoStub)

	errDTO := uc.Delete(context.Background(), 42)
	if errDTO == nil || errDTO.Code != 500 {
		t.Fatalf("Delete() error = %#v, want code 500", errDTO)
	}
	const wantMessage = "Lỗi khi lấy event queue trước khi xóa: database unavailable"
	if errDTO.Message != wantMessage {
		t.Fatalf("Delete() message = %q, want %q", errDTO.Message, wantMessage)
	}
	if repoStub.deleteCalls != 0 {
		t.Fatalf("Delete() repo delete calls = %d, want 0", repoStub.deleteCalls)
	}
}
