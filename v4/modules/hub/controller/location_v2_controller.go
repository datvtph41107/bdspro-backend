package controller

import (
	"net/http"
	"strconv"

	"userv4/internal/httpapi"
	"userv4/modules/hub/dto"
	"userv4/modules/hub/services"
)

// LocationV2Controller owns v2 location HTTP routes under /v2/hub/v2.
type LocationV2Controller struct {
	location services.LocationV2Service
}

func NewLocationV2Controller(location services.LocationV2Service) *LocationV2Controller {
	return &LocationV2Controller{location: location}
}

func (c *LocationV2Controller) RegisterHTTP(mux *http.ServeMux) {
	mux.HandleFunc("GET /v2/hub/v2/provinces", c.httpGetProvinces)
	mux.HandleFunc("GET /v2/hub/v2/province/{code}", c.httpGetProvinceByCode)
	mux.HandleFunc("GET /v2/hub/v2/wards", c.httpGetWards)
	mux.HandleFunc("GET /v2/hub/v2/ward/{code}", c.httpGetWardByCode)
	mux.HandleFunc("GET /v2/hub/v2/location/search", c.httpSearchLocation)
}

func (c *LocationV2Controller) httpGetProvinces(w http.ResponseWriter, r *http.Request) {
	var req dto.GetProvincesV2Request
	if err := httpapi.DecodeQuery(r, &req); err != nil {
		writeHTTPError(w, http.StatusBadRequest, err.Error())
		return
	}
	result, err := c.location.GetProvincesV2(r.Context(), req)
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

func (c *LocationV2Controller) httpGetProvinceByCode(w http.ResponseWriter, r *http.Request) {
	code, err := strconv.ParseInt(r.PathValue("code"), 10, 32)
	if err != nil || code == 0 {
		writeHTTPError(w, http.StatusBadRequest, "code is required")
		return
	}
	var req dto.GetProvinceV2ByCodeRequest
	if err := httpapi.DecodeQuery(r, &req); err != nil {
		writeHTTPError(w, http.StatusBadRequest, err.Error())
		return
	}
	result, err := c.location.GetProvinceV2ByCode(r.Context(), int32(code), req)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeWrappedData(w, result)
}

func (c *LocationV2Controller) httpGetWards(w http.ResponseWriter, r *http.Request) {
	var req dto.GetWardsV2Request
	if err := httpapi.DecodeQuery(r, &req); err != nil {
		writeHTTPError(w, http.StatusBadRequest, err.Error())
		return
	}
	result, err := c.location.GetWardsV2(r.Context(), req)
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

func (c *LocationV2Controller) httpGetWardByCode(w http.ResponseWriter, r *http.Request) {
	code, err := strconv.ParseInt(r.PathValue("code"), 10, 32)
	if err != nil || code == 0 {
		writeHTTPError(w, http.StatusBadRequest, "code is required")
		return
	}
	result, err := c.location.GetWardV2ByCode(r.Context(), int32(code))
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeWrappedData(w, result)
}

func (c *LocationV2Controller) httpSearchLocation(w http.ResponseWriter, r *http.Request) {
	var req dto.SearchLocationV2Request
	if err := httpapi.DecodeQuery(r, &req); err != nil {
		writeHTTPError(w, http.StatusBadRequest, err.Error())
		return
	}
	result, err := c.location.SearchLocationV2(r.Context(), req)
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
