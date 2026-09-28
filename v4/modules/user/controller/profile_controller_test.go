package controller_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	sharepb "pb/types/shared"
	"userv4/modules/user/controller"
)

type stubProfileService struct {
	resp *sharepb.UserV3Proto
	err  error
}

func (s *stubProfileService) GetInfoMe(ctx context.Context) (*sharepb.UserV3Proto, error) {
	return s.resp, s.err
}

func TestProfileController_RegisterHTTP_InfoMe(t *testing.T) {
	ctrl := controller.NewProfileController(&stubProfileService{
		resp: &sharepb.UserV3Proto{
			FullName:  "Alice",
			Timestamp: 1727160000000,
		},
	})

	mux := http.NewServeMux()
	controller.NewHandler(ctrl).RegisterHTTP(mux)

	req := httptest.NewRequest(http.MethodGet, "/v2/user/profile/info/me", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status=%d body=%s", rec.Code, rec.Body.String())
	}

	var got map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
		t.Fatal(err)
	}
	if got["code"].(float64) != 0 {
		t.Fatalf("code=%v", got["code"])
	}
	if got["timestamp"].(float64) != 1727160000000 {
		t.Fatalf("timestamp=%v body=%s", got["timestamp"], rec.Body.String())
	}
	data, ok := got["data"].(map[string]any)
	if !ok {
		t.Fatalf("data=%v body=%s", got["data"], rec.Body.String())
	}
	if data["fullName"] != "Alice" {
		t.Fatalf("data.fullName=%v body=%s", data["fullName"], rec.Body.String())
	}
	if _, hasTS := data["timestamp"]; hasTS {
		t.Fatalf("timestamp should be lifted out of data; body=%s", rec.Body.String())
	}
}
