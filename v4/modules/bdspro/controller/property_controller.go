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
	// One pattern covers both public paths — ServeMux rejects registering
	// detail/{id} and {id}/timestamp together (overlap on .../detail/timestamp).
	mux.HandleFunc("GET /v2/bdspro/v2/property/{a}/{b}", c.httpPropertyTwoSeg)
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

func (c *PropertyController) httpPropertyTwoSeg(w http.ResponseWriter, r *http.Request) {
	a, b := r.PathValue("a"), r.PathValue("b")
	switch {
	case a == "detail":
		c.httpGetDetail(w, r, b)
	case b == "detail":
		c.httpGetDetail(w, r, a)
	case b == "timestamp":
		c.httpGetTimestamps(w, r, a)
	default:
		writeHTTPError(w, http.StatusNotFound, "not found")
	}
}

func (c *PropertyController) httpGetDetail(w http.ResponseWriter, r *http.Request, idRaw string) {
	id, err := strconv.ParseUint(idRaw, 10, 64)
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

func (c *PropertyController) httpGetTimestamps(w http.ResponseWriter, r *http.Request, idRaw string) {
	id, err := strconv.ParseUint(idRaw, 10, 64)
	if err != nil || id == 0 {
		writeHTTPError(w, http.StatusBadRequest, "id is required")
		return
	}
	result, err := c.property.GetTimestamps(r.Context(), id)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	if result == nil {
		writeOK(w, map[string]int64{})
		return
	}
	writeOK(w, result.Timestamps)
}
