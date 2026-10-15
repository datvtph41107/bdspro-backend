package controller

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestBankRouteRegistered(t *testing.T) {
	mux := http.NewServeMux()
	NewBankController(nil).RegisterHTTP(mux)
	req := httptest.NewRequest(http.MethodGet, "/v2/payment/bank/all", nil)
	_, pattern := mux.Handler(req)
	if pattern == "" {
		t.Fatal("GET /v2/payment/bank/all not registered")
	}
}
