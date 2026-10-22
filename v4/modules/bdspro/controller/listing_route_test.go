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

type stubListingService struct{}

func (stubListingService) ListMine(context.Context, dto.ListingListRequest) (*services.ListingListOutcome, error) {
	return &services.ListingListOutcome{
		List: &bdspropb.PostListResponse{
			Data:  []*bdspropb.Post{{Id: 1, Title: "L1"}},
			Total: 1,
		},
		Timestamp: 1,
	}, nil
}

func TestListingMeRouteRegistered(t *testing.T) {
	mux := http.NewServeMux()
	NewListingController(stubListingService{}).RegisterHTTP(mux)

	for _, path := range []string{"/v4/bdspro/listing/me", "/v2/bdspro/v2/listing/me"} {
		req := httptest.NewRequest(http.MethodGet, path+"?page=0&size=10", nil)
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)
		if rec.Code == http.StatusNotFound {
			t.Fatalf("%s: route not registered", path)
		}
	}
}

var _ services.ListingService = stubListingService{}
