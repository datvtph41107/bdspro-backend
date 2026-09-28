package controller

import (
	"net/http"
	"strconv"

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
	mux.HandleFunc("GET /v2/notification/property-history/{id}", c.httpListBySubject)
}

func (c *PropertyHistoryController) httpListBySubject(w http.ResponseWriter, r *http.Request) {
	subjectID, err := strconv.ParseUint(r.PathValue("id"), 10, 64)
	if err != nil || subjectID == 0 {
		writeHTTPError(w, http.StatusBadRequest, "id is required")
		return
	}
	var req dto.SyncByIDRequest
	if err := httpapi.DecodeQuery(r, &req); err != nil {
		writeHTTPError(w, http.StatusBadRequest, err.Error())
		return
	}
	result, err := c.history.ListBySubject(r.Context(), subjectID, req)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeProtoList(w, result.GetData(), result.GetTotal())
}
