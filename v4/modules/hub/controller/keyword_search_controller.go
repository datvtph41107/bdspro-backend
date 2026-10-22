package controller

import (
	"context"
	"net/http"

	"common/httpapi"
	"mdn4/modules/hub/dto"
	"mdn4/modules/hub/services"
)

// KeywordSearchController handles GET /v2/hub/keyword-search-stats.
type KeywordSearchController struct {
	svc services.KeywordSearchService
}

func NewKeywordSearchController(svc services.KeywordSearchService) *KeywordSearchController {
	return &KeywordSearchController{svc: svc}
}

func (c *KeywordSearchController) RegisterHTTP(mux *http.ServeMux) {
	httpapi.GET(mux, "/v2/hub/keyword-search-stats", c.getStats)
}

func (c *KeywordSearchController) getStats(ctx context.Context, req dto.GetKeywordSearchStatsQuery) (*httpapi.Result, error) {
	resp, err := c.svc.GetStats(ctx, req)
	if err != nil {
		return nil, err
	}
	return httpapi.OK(resp), nil
}
