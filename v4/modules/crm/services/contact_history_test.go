package services

import (
	"context"
	"errors"
	"testing"

	_dto "common/domain/dto"
	_enum "common/domain/enum"
	domain "mdn4/modules/crm/entity"
)

type stubHistoryWriter struct {
	got []*_dto.HistoryDTO
	err error
}

func (s *stubHistoryWriter) CreateHistory(_ context.Context, h *_dto.HistoryDTO) error {
	s.got = append(s.got, h)
	return s.err
}

func TestContactChangeNotes(t *testing.T) {
	before := &domain.Contact{FullName: "An", Phone: "0901", Address: "Q1", Avatar: "a.png"}
	after := &domain.Contact{FullName: "An Nguyễn", Phone: "0901", Address: "", Avatar: "b.png"}

	got := contactChangeNotes(before, after)
	want := []string{"Họ tên: An → An Nguyễn", "Địa chỉ: Q1 → (trống)", "Cập nhật ảnh đại diện"}
	if len(got) != len(want) {
		t.Fatalf("notes=%v want %v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("notes[%d]=%q want %q", i, got[i], want[i])
		}
	}
	if n := contactChangeNotes(before, before); len(n) != 0 {
		t.Fatalf("no change should produce no notes, got %v", n)
	}
}

func TestRecordContactHistory(t *testing.T) {
	writer := &stubHistoryWriter{err: errors.New("db down")}
	uc := &ContactUsecase{historyWriter: writer}

	uc.recordContactHistory(context.Background(), 7, _enum.HistoryUpdateNote, []string{"Ghi chú: a → b"})

	if len(writer.got) != 1 {
		t.Fatalf("expected one history row, got %d", len(writer.got))
	}
	h := writer.got[0]
	if h.TargetId != 7 || h.TargetType != _enum.TargetHistoryContact || h.ActionType != _enum.HistoryUpdateNote || h.Title != "Cập nhật ghi chú" {
		t.Fatalf("unexpected history: %+v", h)
	}

	(&ContactUsecase{}).recordContactHistory(context.Background(), 7, _enum.HistoryUpdateNote, nil)
}
