package controller

import (
	"net/http"

	"mdn4/modules/assistant/services"
)

// Handler is the modules/assistant HTTP entry.
type Handler struct {
	product *ProductController
}

func NewHandler(analyze services.ProductAnalyzeService) *Handler {
	return &Handler{
		product: NewProductController(analyze),
	}
}

func (h *Handler) RegisterHTTP(mux *http.ServeMux) {
	if h == nil || h.product == nil {
		return
	}
	h.product.RegisterHTTP(mux)
}
