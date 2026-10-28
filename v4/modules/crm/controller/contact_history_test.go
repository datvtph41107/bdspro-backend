package controller

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	_dto "common/domain/dto"
	notificationpb "pb/types/notification"
)

type stubContactHistoryReader struct {
	gotID   uint64
	gotPage _dto.Pagable
}

func (s *stubContactHistoryReader) ContactHistory(_ context.Context, contactID uint64, page _dto.Pagable) ([]*notificationpb.HistoryDTO, int64, error) {
	s.gotID, s.gotPage = contactID, page
	return []*notificationpb.HistoryDTO{{TargetId: contactID, Title: "Chuyển giai đoạn"}}, 1, nil
}

func TestContactHistoryRoute(t *testing.T) {
	reader := &stubContactHistoryReader{}
	mux := http.NewServeMux()
	NewContactController(nil, nil, reader).RegisterHTTP(mux)

	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/v2/crm/contact/7/history?page=1&size=5", nil))
	if rec.Code != http.StatusOK {
		t.Fatalf("status=%d body=%s", rec.Code, rec.Body.String())
	}
	if reader.gotID != 7 || reader.gotPage.Page != 1 || reader.gotPage.Size != 5 {
		t.Fatalf("reader got id=%d page=%+v", reader.gotID, reader.gotPage)
	}
	var body struct {
		Data       []map[string]any `json:"data"`
		Pagination struct {
			Total int64 `json:"total"`
		} `json:"pagination"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatal(err)
	}
	if len(body.Data) != 1 || body.Data[0]["title"] != "Chuyển giai đoạn" || body.Pagination.Total != 1 {
		t.Fatalf("unexpected body: %s", rec.Body.String())
	}
}
