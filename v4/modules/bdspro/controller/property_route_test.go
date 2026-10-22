package controller

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"mdn4/modules/bdspro/dto"
	"mdn4/modules/bdspro/services"
	bdspropb "pb/types/bdspro"
)

type routeStubProperty struct{}

type routeStubPropertyType struct{}

func (routeStubPropertyType) List(context.Context, dto.PropertyTypeListRequest) (*services.PropertyTypeListOutcome, error) {
	return &services.PropertyTypeListOutcome{}, nil
}

func (routeStubProperty) ListMe(context.Context, dto.PropertyListRequest) (*bdspropb.PropertyListResponse, error) {
	return &bdspropb.PropertyListResponse{}, nil
}

func (routeStubProperty) GetDetail(context.Context, dto.PropertyListRequest) (*bdspropb.PropertyDetailResponse, int64, error) {
	return &bdspropb.PropertyDetailResponse{Id: 30}, 1, nil
}

func (routeStubProperty) GetTimestamps(context.Context, uint64) (*bdspropb.PropertyTimestampsResponse, error) {
	return &bdspropb.PropertyTimestampsResponse{Timestamps: map[string]int64{"property_detail": 1}}, nil
}

func (routeStubProperty) CreateProperty(context.Context, *bdspropb.CreatePropertyRequest) (*bdspropb.Response, error) {
	return &bdspropb.Response{Id: 1, Message: "successfully"}, nil
}

func TestPropertyRoutes_RegisterNoConflict(t *testing.T) {
	mux := http.NewServeMux()
	NewPropertyController(routeStubProperty{}, routeStubPropertyType{}).RegisterHTTP(mux)

	for _, path := range []string{
		"/v2/bdspro/v2/property/me",
		"/v2/bdspro/v2/property/detail/30",
		"/v4/bdspro/property/detail/30",
		"/v2/bdspro/v2/property/timestamp/30",
	} {
		rr := httptest.NewRecorder()
		mux.ServeHTTP(rr, httptest.NewRequest(http.MethodGet, path, nil))
		if rr.Code == http.StatusNotFound {
			t.Fatalf("%s => %d", path, rr.Code)
		}
	}

	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/v2/bdspro/property/new", http.NoBody)
	req.Header.Set("Content-Type", "application/json")
	mux.ServeHTTP(rr, req)
	if rr.Code == http.StatusNotFound {
		t.Fatalf("POST /v2/bdspro/property/new => %d", rr.Code)
	}
}
