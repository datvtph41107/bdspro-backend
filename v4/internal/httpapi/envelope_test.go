package httpapi_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"userv4/internal/httpapi"
)

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
