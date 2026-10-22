package controller

import (
	"context"
	"net/http"

	"common/httpapi"
	"mdn4/modules/bdspro/dto"
	"mdn4/modules/bdspro/services"
)

// ListingController owns listing list routes on the v4 DSN.
type ListingController struct {
	listings services.ListingService
}

func NewListingController(listings services.ListingService) *ListingController {
	return &ListingController{listings: listings}
}

func (c *ListingController) RegisterHTTP(mux *http.ServeMux) {
	httpapi.GET(mux, "/v2/bdspro/v2/listing/me", c.listMine)
	httpapi.GET(mux, "/v4/bdspro/listing/me", c.listMine)
}

func (c *ListingController) listMine(ctx context.Context, req dto.ListingListRequest) (*httpapi.Result, error) {
	out, err := c.listings.ListMine(ctx, req)
	if err != nil {
		return nil, err
	}
	return listingListResult(out, req.Page, req.Size)
}
