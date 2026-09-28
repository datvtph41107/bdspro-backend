package controller

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	bdspropb "pb/types/bdspro"
	"userv4/modules/bdspro/dto"
)

type routeStubProperty struct{}

func (routeStubProperty) ListMe(context.Context, dto.PropertyListRequest) (*bdspropb.PropertyListResponse, error) {
	return &bdspropb.PropertyListResponse{}, nil
}

func (routeStubProperty) GetDetail(context.Context, dto.PropertyListRequest) (*bdspropb.PropertyDetailResponse, error) {
	return &bdspropb.PropertyDetailResponse{Id: 30}, nil
}

func (routeStubProperty) GetTimestamps(context.Context, uint64) (*bdspropb.PropertyTimestampsResponse, error) {
	return &bdspropb.PropertyTimestampsResponse{Timestamps: map[string]int64{"property_detail": 1}}, nil
}

func TestPropertyRoutes_RegisterNoConflict(t *testing.T) {
	mux := http.NewServeMux()
	NewPropertyController(routeStubProperty{}).RegisterHTTP(mux)

	for _, path := range []string{
		"/v2/bdspro/v2/property/me",
		"/v2/bdspro/v2/property/30/detail",
		"/v2/bdspro/v2/property/30/timestamp",
	} {
		rr := httptest.NewRecorder()
		mux.ServeHTTP(rr, httptest.NewRequest(http.MethodGet, path, nil))
		if rr.Code == http.StatusNotFound {
			t.Fatalf("%s => %d", path, rr.Code)
		}
	}
}
