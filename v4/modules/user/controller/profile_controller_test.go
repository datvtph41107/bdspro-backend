package controller_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	sharepb "pb/types/shared"
	"userv4/modules/user/controller"
)

type stubProfileService struct {
	resp       *sharepb.UserV3Proto
	v3Resp     *sharepb.TotalResponseProto
	err        error
	lastID     uint64
	lastAccess string
	lastTS     int64
}

func (s *stubProfileService) GetInfoMe(ctx context.Context) (*sharepb.UserV3Proto, error) {
	return s.resp, s.err
}

func (s *stubProfileService) GetProfileInfoV3(ctx context.Context, id uint64, access string, timestamp int64) (*sharepb.TotalResponseProto, error) {
	s.lastID = id
	s.lastAccess = access
	s.lastTS = timestamp
	return s.v3Resp, s.err
}

type stubDemoService struct {
	enable *bool
	err    error
}

func (s *stubDemoService) SetDemoMode(ctx context.Context, enable bool) error {
	s.enable = &enable
	return s.err
}

func TestProfileController_RegisterHTTP_InfoMe(t *testing.T) {
	ctrl := controller.NewProfileController(&stubProfileService{
		resp: &sharepb.UserV3Proto{
			FullName:  "Alice",
			Timestamp: 1727160000000,
		},
	}, &stubDemoService{})

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

func TestProfileController_RegisterHTTP_ProfileInfoV3(t *testing.T) {
	svc := &stubProfileService{
		v3Resp: &sharepb.TotalResponseProto{
			User: &sharepb.UserV3Proto{FullName: "Bob", ProfileId: 267},
		},
	}
	ctrl := controller.NewProfileController(svc, &stubDemoService{})

	mux := http.NewServeMux()
	controller.NewHandler(ctrl).RegisterHTTP(mux)

	req := httptest.NewRequest(http.MethodGet, "/v3/user/profile/p/267?timestamp=100", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status=%d body=%s", rec.Code, rec.Body.String())
	}
	if svc.lastID != 267 || svc.lastAccess != "p" || svc.lastTS != 100 {
		t.Fatalf("args id=%d access=%q ts=%d", svc.lastID, svc.lastAccess, svc.lastTS)
	}

	var got map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
		t.Fatal(err)
	}
	if got["code"].(float64) != 0 {
		t.Fatalf("code=%v body=%s", got["code"], rec.Body.String())
	}
	data, ok := got["data"].(map[string]any)
	if !ok {
		t.Fatalf("data=%v body=%s", got["data"], rec.Body.String())
	}
	user, ok := data["user"].(map[string]any)
	if !ok {
		t.Fatalf("data.user=%v body=%s", data["user"], rec.Body.String())
	}
	if user["fullName"] != "Bob" {
		t.Fatalf("user.fullName=%v body=%s", user["fullName"], rec.Body.String())
	}
}

func TestProfileController_RegisterHTTP_DemoMode(t *testing.T) {
	demo := &stubDemoService{}
	ctrl := controller.NewProfileController(&stubProfileService{}, demo)

	mux := http.NewServeMux()
	controller.NewHandler(ctrl).RegisterHTTP(mux)

	body, _ := json.Marshal(map[string]any{"enable": true})
	req := httptest.NewRequest(http.MethodPost, "/v2/user/demo-mode", bytes.NewReader(body))
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status=%d body=%s", rec.Code, rec.Body.String())
	}
	if demo.enable == nil || !*demo.enable {
		t.Fatalf("expected enable=true, got %+v", demo.enable)
	}
	var got map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
		t.Fatal(err)
	}
	if got["code"].(float64) != 0 {
		t.Fatalf("code=%v body=%s", got["code"], rec.Body.String())
	}
}
