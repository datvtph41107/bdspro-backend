package controller_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"mdn4/modules/user/controller"
	"mdn4/modules/user/services"
	sharepb "pb/types/shared"
	userpb "pb/types/user"
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

func (s *stubProfileService) GetProfileInfo(ctx context.Context, id uint64) (*userpb.ProfilePublicInfo, error) {
	return &userpb.ProfilePublicInfo{FullName: "Public"}, s.err
}

func (s *stubProfileService) ListAccountMe(ctx context.Context) (*userpb.AccountInfoResponse, error) {
	return &userpb.AccountInfoResponse{Total: 0}, s.err
}

func (s *stubProfileService) SetInfo(ctx context.Context, req *userpb.SetInfoRequest) error {
	return s.err
}

func (s *stubProfileService) SetPerson(ctx context.Context, req *sharepb.UserV3Proto) error {
	return s.err
}

func (s *stubProfileService) UpdatePrivacySetting(ctx context.Context, req *userpb.UpdatePrivacySettingRequest) error {
	return s.err
}

func (s *stubProfileService) QuickSetup(ctx context.Context, req *sharepb.UserV3Proto) error {
	return s.err
}

func (s *stubProfileService) SearchProfilePublic(ctx context.Context, text string, page, size uint32) (*userpb.SearchProfilePublicResponse, error) {
	return &userpb.SearchProfilePublicResponse{}, s.err
}

func (s *stubProfileService) FindByPhone(ctx context.Context, phone string) (*userpb.FindByPhoneResponse, error) {
	return &userpb.FindByPhoneResponse{Phone: phone}, s.err
}

func (s *stubProfileService) GetProfileCompletion(ctx context.Context) (*userpb.ProfileCompletionResponse, error) {
	return &userpb.ProfileCompletionResponse{Percentage: 50}, s.err
}

func (s *stubProfileService) GetStatusOnline(ctx context.Context, id uint64) (*userpb.StatusOnlineResponse, error) {
	return &userpb.StatusOnlineResponse{ProfileId: id, Online: true}, s.err
}

type stubDemoService struct {
	enable *bool
	err    error
}

func (s *stubDemoService) SetDemoMode(ctx context.Context, enable bool) error {
	s.enable = &enable
	return s.err
}

type stubTileSessionService struct {
	result *services.TileSessionResult
	err    error
}

func (s *stubTileSessionService) GenTileSession(ctx context.Context) (*services.TileSessionResult, error) {
	return s.result, s.err
}

func testHandler(profile services.ProfileService, demo services.DemoService, session services.TileSessionService) *controller.Handler {
	return controller.NewHandler(
		controller.NewProfileController(profile, demo, session),
		nil, nil, nil, nil, nil,
	)
}

func TestProfileController_RegisterHTTP_InfoMe(t *testing.T) {
	mux := http.NewServeMux()
	testHandler(&stubProfileService{
		resp: &sharepb.UserV3Proto{
			FullName:  "Alice",
			Timestamp: 1727160000000,
		},
	}, &stubDemoService{}, &stubTileSessionService{}).RegisterHTTP(mux)

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
	mux := http.NewServeMux()
	testHandler(svc, &stubDemoService{}, &stubTileSessionService{}).RegisterHTTP(mux)

	req := httptest.NewRequest(http.MethodGet, "/v3/user/profile/p/267?timestamp=100", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status=%d body=%s", rec.Code, rec.Body.String())
	}
	if svc.lastID != 267 || svc.lastAccess != "p" || svc.lastTS != 100 {
		t.Fatalf("args id=%d access=%q ts=%d", svc.lastID, svc.lastAccess, svc.lastTS)
	}
}

func TestProfileController_RegisterHTTP_DemoMode(t *testing.T) {
	demo := &stubDemoService{}
	mux := http.NewServeMux()
	testHandler(&stubProfileService{}, demo, &stubTileSessionService{}).RegisterHTTP(mux)

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
}

func TestProfileController_RegisterHTTP_Session(t *testing.T) {
	session := &stubTileSessionService{
		result: &services.TileSessionResult{
			SessionK:          "12345",
			SessionEncryptKey: "abc=",
			ExpiresIn:         3600,
		},
	}
	mux := http.NewServeMux()
	testHandler(&stubProfileService{}, &stubDemoService{}, session).RegisterHTTP(mux)

	req := httptest.NewRequest(http.MethodPost, "/v2/user/profile/session", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status=%d body=%s", rec.Code, rec.Body.String())
	}
	var got map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
		t.Fatal(err)
	}
	data, ok := got["data"].(map[string]any)
	if !ok {
		t.Fatalf("data=%v body=%s", got["data"], rec.Body.String())
	}
	if data["sessionK"] != "12345" {
		t.Fatalf("sessionK=%v", data["sessionK"])
	}
	if data["sessionEncryptKey"] != "abc=" {
		t.Fatalf("sessionEncryptKey=%v", data["sessionEncryptKey"])
	}
}
