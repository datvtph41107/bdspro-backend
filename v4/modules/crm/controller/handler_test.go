package controller

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func newRouteOnlyHandler() *Handler {
	return NewHandler(nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil,
		&FeedbackController{}, &AppointmentController{}, &EnumController{})
}

func TestRegisterHTTP_ProtoPathsResolve(t *testing.T) {
	mux := http.NewServeMux()
	newRouteOnlyHandler().RegisterHTTP(mux)

	cases := []struct{ method, path, pattern string }{
		{http.MethodGet, "/v2/crm/contact/by-product/638?size=999", "GET /v2/crm/contact/by-product/{id}"},
		{http.MethodGet, "/v2/crm/contact/detail/5", "GET /v2/crm/contact/detail/{id}"},
		{http.MethodGet, "/v2/crm/contact/detail/5/products", "GET /v2/crm/contact/detail/{id}/products"},
		{http.MethodGet, "/v2/crm/contact/relation-ship/5", "GET /v2/crm/contact/relation-ship/{id}"},
		{http.MethodPut, "/v2/crm/contact/edit/5", "PUT /v2/crm/contact/{first}/{second}"},
		{http.MethodPut, "/v2/crm/contact/5/products", "PUT /v2/crm/contact/{first}/{second}"},
		{http.MethodGet, "/v2/feedback/report/reasons", "GET /v2/feedback/report/reasons"},
		{http.MethodGet, "/v2/feedback/rate/detail/3", "GET /v2/feedback/rate/detail/{id}"},
		{http.MethodGet, "/v2/feedback/rate/product/stats/42", "GET /v2/feedback/rate/{registry}/stats/{ownerId}"},
		{http.MethodPost, "/v2/feedback/admin/report/3/approve", "POST /v2/feedback/admin/report/{id}/approve"},
		{http.MethodGet, "/v2/feedback/report-reason/3", "GET /v2/feedback/report-reason/{id}"},
		{http.MethodGet, "/v2/feedback/registry", "GET /v2/feedback/registry"},
		{http.MethodGet, "/v2.1/appointment/stats", "GET /v2.1/appointment/stats"},
		{http.MethodGet, "/v2.1/appointment/admin/list", "GET /v2.1/appointment/admin/list"},
		{http.MethodGet, "/v2.1/appointment/12", "GET /v2.1/appointment/{id}"},
		{http.MethodPut, "/v2.1/appointment/12/status/approve", "PUT /v2.1/appointment/{id}/status/{action}"},
		{http.MethodGet, "/v2/enums/crm/priority", "GET /v2/enums/crm/{name}"},
	}
	for _, tc := range cases {
		_, pattern := mux.Handler(httptest.NewRequest(tc.method, tc.path, nil))
		if pattern != tc.pattern {
			t.Errorf("%s %s matched %q, want %q", tc.method, tc.path, pattern, tc.pattern)
		}
	}
}

func TestContactPutDispatch(t *testing.T) {
	var hit, id string
	record := func(name string) http.Handler {
		return http.HandlerFunc(func(_ http.ResponseWriter, r *http.Request) {
			hit, id = name, r.PathValue("id")
		})
	}
	mux := http.NewServeMux()
	mux.HandleFunc("PUT /v2/crm/contact/{first}/{second}",
		(&ContactController{}).putDispatch(record("edit"), record("note"), record("products")))

	cases := []struct{ path, hit, id string }{
		{"/v2/crm/contact/edit/5", "edit", "5"},
		{"/v2/crm/contact/note/6", "note", "6"},
		{"/v2/crm/contact/7/products", "products", "7"},
	}
	for _, tc := range cases {
		hit, id = "", ""
		mux.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest(http.MethodPut, tc.path, nil))
		if hit != tc.hit || id != tc.id {
			t.Errorf("%s → %s(id=%s), want %s(id=%s)", tc.path, hit, id, tc.hit, tc.id)
		}
	}

	w := httptest.NewRecorder()
	mux.ServeHTTP(w, httptest.NewRequest(http.MethodPut, "/v2/crm/contact/foo/bar", nil))
	if w.Code != http.StatusNotFound {
		t.Fatalf("unknown PUT status=%d", w.Code)
	}
}
