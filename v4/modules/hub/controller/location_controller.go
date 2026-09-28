package controller

import (
	"net/http"

	"userv4/internal/httpapi"
	"userv4/modules/hub/dto"
	"userv4/modules/hub/services"
)

// LocationController owns v1 location HTTP routes under /v2/hub.
type LocationController struct {
	location services.LocationService
}

func NewLocationController(location services.LocationService) *LocationController {
	return &LocationController{location: location}
}

func (c *LocationController) RegisterHTTP(mux *http.ServeMux) {
	mux.HandleFunc("GET /v2/hub/provinces", c.httpGetProvinces)
	mux.HandleFunc("GET /v2/hub/districts", c.httpGetDistricts)
	mux.HandleFunc("GET /v2/hub/wards", c.httpGetWards)
	mux.HandleFunc("GET /v2/hub/location/search", c.httpSearchLocation)
}

func (c *LocationController) httpGetProvinces(w http.ResponseWriter, r *http.Request) {
	var req dto.GetProvincesRequest
	if err := httpapi.DecodeQuery(r, &req); err != nil {
		writeHTTPError(w, http.StatusBadRequest, err.Error())
		return
	}
	result, err := c.location.GetProvinces(r.Context(), req)
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

func (c *LocationController) httpGetDistricts(w http.ResponseWriter, r *http.Request) {
	var req dto.GetDistrictsRequest
	if err := httpapi.DecodeQuery(r, &req); err != nil {
		writeHTTPError(w, http.StatusBadRequest, err.Error())
		return
	}
	result, err := c.location.GetDistricts(r.Context(), req)
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

func (c *LocationController) httpGetWards(w http.ResponseWriter, r *http.Request) {
	var req dto.GetWardsRequest
	if err := httpapi.DecodeQuery(r, &req); err != nil {
		writeHTTPError(w, http.StatusBadRequest, err.Error())
		return
	}
	result, err := c.location.GetWards(r.Context(), req)
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

func (c *LocationController) httpSearchLocation(w http.ResponseWriter, r *http.Request) {
	var req dto.SearchLocationRequest
	if err := httpapi.DecodeQuery(r, &req); err != nil {
		writeHTTPError(w, http.StatusBadRequest, err.Error())
		return
	}
	result, err := c.location.SearchLocation(r.Context(), req)
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
