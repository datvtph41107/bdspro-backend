package controller

import (
	"context"
	"net/http"

	"userv4/internal/httpapi"
	"userv4/modules/bdspro/dto"
	"userv4/modules/bdspro/services"
)

// DealInternalNoteController owns POST /v4/deal/{id}/internal-note.
type DealInternalNoteController struct {
	internalNote services.DealInternalNoteService
}

func NewDealInternalNoteController(internalNote services.DealInternalNoteService) *DealInternalNoteController {
	return &DealInternalNoteController{internalNote: internalNote}
}

func (c *DealInternalNoteController) RegisterHTTP(mux *http.ServeMux) {
	httpapi.POST(mux, "/v4/deal/{id}/internal-note", c.add)
}

func (c *DealInternalNoteController) add(ctx context.Context, req dto.AddDealInternalNoteRequest) (*httpapi.Result, error) {
	result, err := c.internalNote.Add(ctx, req)
	if err != nil {
		return nil, err
	}
	return protoResult(result)
}
