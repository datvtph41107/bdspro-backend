package controller

import (
	"net/http"

	"userv4/internal/httpapi"
	"userv4/modules/bdspro/dto"
	"userv4/modules/bdspro/services"
)

// DealController owns GET /v2.1/org/deals/user.
type DealController struct {
	deal services.DealService
}

func NewDealController(deal services.DealService) *DealController {
	return &DealController{deal: deal}
}

func (c *DealController) RegisterHTTP(mux *http.ServeMux) {
	mux.HandleFunc("GET /v2.1/org/deals/user", c.httpGetDealsUser)
}

func (c *DealController) httpGetDealsUser(w http.ResponseWriter, r *http.Request) {
	var req dto.DealSearchRequest
	if err := httpapi.DecodeQuery(r, &req); err != nil {
		writeHTTPError(w, http.StatusBadRequest, err.Error())
		return
	}
	result, err := c.deal.GetDealsUser(r.Context(), req)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeGetDealsResponse(w, result)
}
