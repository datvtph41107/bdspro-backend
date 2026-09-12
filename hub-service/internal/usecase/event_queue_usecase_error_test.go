package usecase

import (
	"context"
	"errors"
	"testing"

	"hub/internal/domain"
	"hub/internal/dto"
	"hub/internal/repo"
)

type eventQueueRepoLookupStub struct {
	repo.IEventQueueRepo
	event        *domain.EventQueueEntity
	err          error
	updateCalled bool
	deleteCalled bool
}

func (s *eventQueueRepoLookupStub) GetByID(context.Context, uint64) (*domain.EventQueueEntity, error) {
	return s.event, s.err
}

func (s *eventQueueRepoLookupStub) Update(context.Context, uint64, *domain.EventQueueEntity) error {
	s.updateCalled = true
	return nil
}

func (s *eventQueueRepoLookupStub) Delete(context.Context, uint64) error {
	s.deleteCalled = true
	return nil
}

func TestEventQueueDetailDistinguishesAbsenceFromTechnicalFailure(t *testing.T) {
	t.Run("missing remains legacy 404", func(t *testing.T) {
		uc := NewEventQueueUsecase(&eventQueueRepoLookupStub{})

		got, errDTO := uc.Detail(context.Background(), 42)
		if got != nil {
			t.Fatalf("Detail() result = %#v, want nil", got)
		}
		if errDTO == nil {
			t.Fatal("Detail() error = nil, want not-found error")
		}
		if errDTO.Code != 404 {
			t.Fatalf("error code = %d, want 404", errDTO.Code)
		}
		if errDTO.Message != "Không tìm thấy event queue" {
			t.Fatalf("error message = %q, want %q", errDTO.Message, "Không tìm thấy event queue")
		}
	})

	t.Run("database failure becomes 500", func(t *testing.T) {
		uc := NewEventQueueUsecase(&eventQueueRepoLookupStub{err: errors.New("database unavailable")})

		got, errDTO := uc.Detail(context.Background(), 42)
		if got != nil {
			t.Fatalf("Detail() result = %#v, want nil", got)
		}
		if errDTO == nil {
			t.Fatal("Detail() error = nil, want technical error")
		}
		if errDTO.Code != 500 {
			t.Fatalf("error code = %d, want 500", errDTO.Code)
		}
		const wantMessage = "Lỗi khi lấy event queue: database unavailable"
		if errDTO.Message != wantMessage {
			t.Fatalf("error message = %q, want %q", errDTO.Message, wantMessage)
		}
	})

	t.Run("success returns entity", func(t *testing.T) {
		want := &domain.EventQueueEntity{}
		uc := NewEventQueueUsecase(&eventQueueRepoLookupStub{event: want})

		got, errDTO := uc.Detail(context.Background(), 42)
		if errDTO != nil {
			t.Fatalf("Detail() error = %#v, want nil", errDTO)
		}
		if got != want {
			t.Fatalf("Detail() result = %p, want %p", got, want)
		}
	})
}

func TestEventQueueUpdateStopsOnLookupFailure(t *testing.T) {
	repoStub := &eventQueueRepoLookupStub{err: errors.New("database unavailable")}
	uc := NewEventQueueUsecase(repoStub)

	got, errDTO := uc.Update(context.Background(), 42, dto.EventQueueSaveDTO{})
	if got != nil {
		t.Fatalf("Update() result = %#v, want nil", got)
	}
	if errDTO == nil || errDTO.Code != 500 {
		t.Fatalf("Update() error = %#v, want code 500", errDTO)
	}
	if repoStub.updateCalled {
		t.Fatal("Update() called repository Update after lookup failure")
	}
}

func TestEventQueueDeleteStopsOnLookupFailure(t *testing.T) {
	repoStub := &eventQueueRepoLookupStub{err: errors.New("database unavailable")}
	uc := NewEventQueueUsecase(repoStub)

	errDTO := uc.Delete(context.Background(), 42)
	if errDTO == nil || errDTO.Code != 500 {
		t.Fatalf("Delete() error = %#v, want code 500", errDTO)
	}
	if repoStub.deleteCalled {
		t.Fatal("Delete() called repository Delete after lookup failure")
	}
}
