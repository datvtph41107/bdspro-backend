package controller_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"userv4/modules/auth/controller"
	"userv4/modules/auth/dto"
)

type stubPermissionService struct {
	resp *dto.PermissionKeysResponse
	err  error
}

func (s *stubPermissionService) GetMyPermissionKeys(ctx context.Context) (*dto.PermissionKeysResponse, error) {
	return s.resp, s.err
}

func TestPermissionController_RegisterHTTP_PermissionMe(t *testing.T) {
	ctrl := controller.NewPermissionController(&stubPermissionService{
		resp: &dto.PermissionKeysResponse{
			Keys:      []string{"product.read", "deal.write"},
			Timestamp: 1727160000000,
		},
	})

	mux := http.NewServeMux()
	ctrl.RegisterHTTP(mux)

	req := httptest.NewRequest(http.MethodGet, "/v2/auth/permission/me", nil)
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
	keys, ok := data["keys"].([]any)
	if !ok || len(keys) != 2 || keys[0] != "product.read" {
		t.Fatalf("data.keys=%v body=%s", data["keys"], rec.Body.String())
	}
}
