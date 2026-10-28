package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	_jwt "common/jwt"
)

func TestPublicHTTPRoutes_OAuthLoginAndCallbackOnly(t *testing.T) {
	ok := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) { w.WriteHeader(http.StatusOK) })
	h := _jwt.JWTAuthHTTPMiddleware(publicHTTPRoutes, nil)(ok)

	cases := []struct {
		path string
		want int
	}{
		{"/v2/auth/oauth/google/login", http.StatusOK},
		{"/v2/auth/oauth/google/callback", http.StatusOK},
		{"/v2/auth/oauth/facebook/callback", http.StatusOK},
		{"/v2/auth/oauth/connected-accounts", http.StatusUnauthorized},
	}
	for _, tc := range cases {
		rec := httptest.NewRecorder()
		h.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, tc.path, nil))
		if rec.Code != tc.want {
			t.Errorf("%s: status = %d, want %d", tc.path, rec.Code, tc.want)
		}
	}
}

func TestPublicHTTPRoutes_FeedbackGuestAccess(t *testing.T) {
	ok := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) { w.WriteHeader(http.StatusOK) })
	h := _jwt.JWTAuthHTTPMiddleware(publicHTTPRoutes, nil)(ok)

	cases := []struct {
		method, path string
		want         int
	}{
		{http.MethodGet, "/v2/feedback/rate/product/list/42", http.StatusOK},
		{http.MethodGet, "/v2/feedback/rate/product/stats/42", http.StatusOK},
		{http.MethodGet, "/v2/feedback/report/reasons", http.StatusOK},
		{http.MethodPost, "/v2/feedback/report-reason", http.StatusUnauthorized},
		{http.MethodGet, "/v2/feedback/admin/report/product", http.StatusUnauthorized},
		{http.MethodGet, "/v2/enums/crm/priority", http.StatusUnauthorized},
		{http.MethodGet, "/v2.1/appointment", http.StatusUnauthorized},
		{http.MethodGet, "/v2/internal/seo/static-paths", http.StatusUnauthorized},
		{http.MethodGet, "/v2/crm/contact/42/relation-ship", http.StatusOK},
		{http.MethodGet, "/v2/crm/contact/42/detail", http.StatusUnauthorized},
	}
	for _, tc := range cases {
		rec := httptest.NewRecorder()
		h.ServeHTTP(rec, httptest.NewRequest(tc.method, tc.path, nil))
		if rec.Code != tc.want {
			t.Errorf("%s %s: status = %d, want %d", tc.method, tc.path, rec.Code, tc.want)
		}
	}
}
