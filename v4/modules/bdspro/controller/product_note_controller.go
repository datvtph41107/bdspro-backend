package controller

import (
	"net/http"
	"strconv"

	"userv4/internal/httpapi"
	"userv4/modules/bdspro/dto"
	"userv4/modules/bdspro/services"
)

// ProductNoteController owns GET /v2/bdspro/product-notes/all/:productId.
type ProductNoteController struct {
	notes services.ProductNoteService
}

func NewProductNoteController(notes services.ProductNoteService) *ProductNoteController {
	return &ProductNoteController{notes: notes}
}

func (c *ProductNoteController) RegisterHTTP(mux *http.ServeMux) {
	mux.HandleFunc("GET /v2/bdspro/product-notes/all/{productId}", c.httpListByProduct)
}

func (c *ProductNoteController) httpListByProduct(w http.ResponseWriter, r *http.Request) {
	productID, err := strconv.ParseUint(r.PathValue("productId"), 10, 64)
	if err != nil || productID == 0 {
		writeHTTPError(w, http.StatusBadRequest, "productId is required")
		return
	}
	var req dto.SyncByIDRequest
	if err := httpapi.DecodeQuery(r, &req); err != nil {
		writeHTTPError(w, http.StatusBadRequest, err.Error())
		return
	}
	result, err := c.notes.ListByProduct(r.Context(), productID, req)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeProtoList(w, result.GetData(), result.GetTotal())
}
