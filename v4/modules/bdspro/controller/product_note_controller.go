package controller

import (
	"context"
	"net/http"

	"common/httpapi"
	"mdn4/modules/bdspro/dto"
	"mdn4/modules/bdspro/services"
)

// ProductNoteController owns GET /v2/bdspro/product-notes/all/:productId.
type ProductNoteController struct {
	notes services.ProductNoteService
}

func NewProductNoteController(notes services.ProductNoteService) *ProductNoteController {
	return &ProductNoteController{notes: notes}
}

func (c *ProductNoteController) RegisterHTTP(mux *http.ServeMux) {
	httpapi.GET(mux, "/v2/bdspro/product-notes/all/{productId}", c.listByProduct)
}

func (c *ProductNoteController) listByProduct(ctx context.Context, req dto.SyncByIDPathRequest) (*httpapi.Result, error) {
	result, err := c.notes.ListByProduct(ctx, req.ID, req.ToSyncByIDRequest())
	if err != nil {
		return nil, err
	}
	return protoListResult(result.GetData(), result.GetTotal())
}
