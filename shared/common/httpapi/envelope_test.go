package httpapi_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	_routes "common/routes"
	"userv4/internal/httpapi"
)

func captureDefaultLog(t *testing.T) *bytes.Buffer {
	t.Helper()
	var buf bytes.Buffer
	prev := slog.Default()
	slog.SetDefault(slog.New(slog.NewTextHandler(&buf, nil)))
	t.Cleanup(func() { slog.SetDefault(prev) })
	return &buf
}

func TestWriteServiceError_LegacyExceptKeepsStatus(t *testing.T) {
	logs := captureDefaultLog(t)
	w := httptest.NewRecorder()
	httpapi.WriteServiceError(w, &_routes.Except{Code: 401, Message: "invalid token"})

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("status=%d body=%s", w.Code, w.Body.String())
	}
	if !strings.Contains(w.Body.String(), `"message":"invalid token"`) {
		t.Fatalf("body=%s", w.Body.String())
	}
	if logs.Len() != 0 {
		t.Fatalf("4xx should not be logged as error: %s", logs.String())
	}
}

func TestBindJSON_LogsCauseOfServerError(t *testing.T) {
	logs := captureDefaultLog(t)
	h := httpapi.BindJSON(func(_ context.Context, _ httpapi.Empty) (*httpapi.Result, error) {
		return nil, errors.New("redis: connection refused")
	})
	w := httptest.NewRecorder()
	h(w, httptest.NewRequest(http.MethodPost, "/v2/auth/token/refresh", nil))

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("status=%d", w.Code)
	}
	got := logs.String()
	for _, want := range []string{"redis: connection refused", "/v2/auth/token/refresh", "http_status_code=500"} {
		if !strings.Contains(got, want) {
			t.Fatalf("log missing %q: %s", want, got)
		}
	}
}

func TestWriteOK_OmitsEmptyOptionalFields(t *testing.T) {
	w := httptest.NewRecorder()
	httpapi.WriteOK(w, map[string]int64{"product_me": 1})

	if w.Code != http.StatusOK {
		t.Fatalf("status=%d", w.Code)
	}

	var body map[string]any
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatal(err)
	}
	if body["code"] != float64(0) {
		t.Fatalf("code=%v", body["code"])
	}
	if body["message"] != "Success" {
		t.Fatalf("message=%v", body["message"])
	}
	if _, ok := body["pagination"]; ok {
		t.Fatalf("pagination should be omitted: %v", body["pagination"])
	}
	if _, ok := body["timestamp"]; ok {
		t.Fatalf("timestamp should be omitted: %v", body["timestamp"])
	}
	data, _ := body["data"].(map[string]any)
	if data["product_me"] != float64(1) {
		t.Fatalf("data=%v", body["data"])
	}
}

func TestWriteOK_WithPaginationAndTimestamp(t *testing.T) {
	w := httptest.NewRecorder()
	httpapi.WriteOK(w, []string{"a"},
		httpapi.WithPagination(&httpapi.Pagination{Page: 1, Size: 5, Total: 10, TotalPages: 2}),
		httpapi.WithTimestamp(99),
	)

	var body map[string]any
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatal(err)
	}
	if body["timestamp"] != float64(99) {
		t.Fatalf("timestamp=%v", body["timestamp"])
	}
	p, _ := body["pagination"].(map[string]any)
	if p["total"] != float64(10) {
		t.Fatalf("pagination=%v", body["pagination"])
	}
}

func TestDecodeQuery_FormTags(t *testing.T) {
	type req struct {
		Page           int32   `form:"page"`
		Size           uint32  `form:"size"`
		Text           string  `form:"text"`
		OrganizationID *uint64 `form:"organizationId"`
		Include        bool    `form:"include"`
	}
	r := httptest.NewRequest(http.MethodGet, "/?page=2&size=10&text=hello&organizationId=99&include=true", nil)
	var got req
	if err := httpapi.DecodeQuery(r, &got); err != nil {
		t.Fatal(err)
	}
	if got.Page != 2 || got.Size != 10 || got.Text != "hello" || got.OrganizationID == nil || *got.OrganizationID != 99 || !got.Include {
		t.Fatalf("got=%+v", got)
	}
}

func TestDecodeJSON(t *testing.T) {
	type body struct {
		Name string `json:"name"`
	}
	r := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`{"name":"x"}`))
	var got body
	if err := httpapi.DecodeJSON(r, &got); err != nil {
		t.Fatal(err)
	}
	if got.Name != "x" {
		t.Fatalf("got=%+v", got)
	}
}
