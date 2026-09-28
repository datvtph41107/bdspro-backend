package controller

import (
	"net/http"
	"strconv"

	"userv4/internal/httpapi"
	"userv4/modules/hub/dto"
	"userv4/modules/hub/services"
)

// FAQController owns user-facing FAQ HTTP routes under /v2/hub/faqs.
type FAQController struct {
	faq services.FAQService
}

func NewFAQController(faq services.FAQService) *FAQController {
	return &FAQController{faq: faq}
}

func (c *FAQController) RegisterHTTP(mux *http.ServeMux) {
	mux.HandleFunc("GET /v2/hub/faqs", c.httpListSimple)
	mux.HandleFunc("GET /v2/hub/faqs/{id}", c.httpGetDetail)
}

func (c *FAQController) httpListSimple(w http.ResponseWriter, r *http.Request) {
	var req dto.GetFAQSimpleRequest
	if err := httpapi.DecodeQuery(r, &req); err != nil {
		writeHTTPError(w, http.StatusBadRequest, err.Error())
		return
	}
	result, err := c.faq.GetFAQSimple(r.Context(), req)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	if result == nil {
		writeListWithTotal(w, []any{}, 0)
		return
	}
	writeProtoListWithTotal(w, result.GetData(), result.GetTotal())
}

func (c *FAQController) httpGetDetail(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(r.PathValue("id"), 10, 64)
	if err != nil || id == 0 {
		writeHTTPError(w, http.StatusBadRequest, "id is required")
		return
	}
	result, err := c.faq.GetFAQSimpleDetail(r.Context(), id)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeProtoData(w, result)
}
