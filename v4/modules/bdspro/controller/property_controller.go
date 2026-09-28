package controller

import (
	"net/http"
	"strconv"

	"userv4/internal/httpapi"
	"userv4/modules/bdspro/dto"
	"userv4/modules/bdspro/services"
)

// PropertyController owns core property HTTP routes under /v2/bdspro/v2/property.
type PropertyController struct {
	property services.PropertyService
}

func NewPropertyController(property services.PropertyService) *PropertyController {
	return &PropertyController{property: property}
}

func (c *PropertyController) RegisterHTTP(mux *http.ServeMux) {
	mux.HandleFunc("GET /v2/bdspro/v2/property/me", c.httpListMe)
	mux.HandleFunc("GET /v2/bdspro/v2/property/detail/{id}", c.httpGetDetail)
}

func (c *PropertyController) httpListMe(w http.ResponseWriter, r *http.Request) {
	var req dto.PropertyListRequest
	if err := httpapi.DecodeQuery(r, &req); err != nil {
		writeHTTPError(w, http.StatusBadRequest, err.Error())
		return
	}
	result, err := c.property.ListMe(r.Context(), req)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writePropertyListResponse(w, result)
}

func (c *PropertyController) httpGetDetail(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(r.PathValue("id"), 10, 64)
	if err != nil || id == 0 {
		writeHTTPError(w, http.StatusBadRequest, "id is required")
		return
	}
	var req dto.PropertyListRequest
	if err := httpapi.DecodeQuery(r, &req); err != nil {
		writeHTTPError(w, http.StatusBadRequest, err.Error())
		return
	}
	req.ID = id
	result, err := c.property.GetDetail(r.Context(), req)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeProtoData(w, result)
}
