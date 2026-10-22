package httpapi_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"userv4/internal/httpapi"
)

type sampleReq struct {
	DealID    uint64 `path:"id|dealId"`
	Page      uint32 `form:"page"`
	Timestamp int64  `form:"timestamp"`
}

func TestDecodeRequest_PathAndQuery(t *testing.T) {
	mux := http.NewServeMux()
	httpapi.GET(mux, "/v4/deal/{id}/x", func(ctx context.Context, req sampleReq) (*httpapi.Result, error) {
		return httpapi.OK(map[string]any{
			"dealId":    req.DealID,
			"page":      req.Page,
			"timestamp": req.Timestamp,
		}), nil
	})

	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, httptest.NewRequest(http.MethodGet, "/v4/deal/13/x?page=2&timestamp=99", nil))
	if rr.Code != http.StatusOK {
		t.Fatalf("status=%d body=%s", rr.Code, rr.Body.String())
	}
	var env httpapi.Envelope
	if err := json.Unmarshal(rr.Body.Bytes(), &env); err != nil {
		t.Fatal(err)
	}
	data, _ := env.Data.(map[string]any)
	if data["dealId"].(float64) != 13 || data["page"].(float64) != 2 || data["timestamp"].(float64) != 99 {
		t.Fatalf("unexpected data: %+v", data)
	}
}

func TestDecodeRequest_DealIdAlias(t *testing.T) {
	mux := http.NewServeMux()
	httpapi.GET(mux, "/v4/x/list/{dealId}", func(ctx context.Context, req sampleReq) (*httpapi.Result, error) {
		return httpapi.OK(req.DealID), nil
	})
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, httptest.NewRequest(http.MethodGet, "/v4/x/list/7", nil))
	if rr.Code != http.StatusOK {
		t.Fatalf("status=%d", rr.Code)
	}
}

func TestDecodeRequest_PathAllowZero(t *testing.T) {
	type allowZeroReq struct {
		OwnerID uint64 `path:"ownerId" pathAllowZero:"true"`
	}
	mux := http.NewServeMux()
	httpapi.GET(mux, "/list/user/{ownerId}", func(ctx context.Context, req allowZeroReq) (*httpapi.Result, error) {
		return httpapi.OK(req.OwnerID), nil
	})
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, httptest.NewRequest(http.MethodGet, "/list/user/0", nil))
	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d body=%s", rr.Code, rr.Body.String())
	}
}

func TestDecodeRequest_MissingPathID(t *testing.T) {
	mux := http.NewServeMux()
	httpapi.GET(mux, "/v4/deal/{id}/x", func(ctx context.Context, req sampleReq) (*httpapi.Result, error) {
		return httpapi.OK(nil), nil
	})
	// ServeMux won't match empty id; use zero id path that matches pattern.
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, httptest.NewRequest(http.MethodGet, "/v4/deal/0/x", nil))
	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d body=%s", rr.Code, rr.Body.String())
	}
}
