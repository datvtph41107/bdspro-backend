package controller

import (
	"context"
	"net/http"

	"userv4/internal/httpapi"
	"userv4/modules/bdspro/dto"
	"userv4/modules/bdspro/services"
)

// PropertyHistoryController owns GET /v2/notification/property-history/:id.
type PropertyHistoryController struct {
	history services.PropertyHistoryService
}

func NewPropertyHistoryController(history services.PropertyHistoryService) *PropertyHistoryController {
	return &PropertyHistoryController{history: history}
}

func (c *PropertyHistoryController) RegisterHTTP(mux *http.ServeMux) {
	httpapi.GET(mux, "/v2/notification/property-history/{id}", c.listBySubject)
}

func (c *PropertyHistoryController) listBySubject(ctx context.Context, req dto.SyncByIDPathRequest) (*httpapi.Result, error) {
	result, err := c.history.ListBySubject(ctx, req.ID, req.ToSyncByIDRequest())
	if err != nil {
		return nil, err
	}
	return protoListResult(result.GetData(), result.GetTotal())
}
