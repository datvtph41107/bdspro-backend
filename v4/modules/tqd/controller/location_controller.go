package controller

import (
	"net/http"

	tqdpb "pb/types/tqd"
	"userv4/internal/httpapi"
	"userv4/modules/tqd/dto"
	"userv4/modules/tqd/services"
)

// LocationController owns TQD client location HTTP routes under /v2/tqd/client.
type LocationController struct {
	location services.LocationService
}

func NewLocationController(location services.LocationService) *LocationController {
	return &LocationController{location: location}
}

func (c *LocationController) RegisterHTTP(mux *http.ServeMux) {
	mux.HandleFunc("GET /v2/tqd/client/provinces", c.httpListProvinces)
	mux.HandleFunc("GET /v2/tqd/client/wards", c.httpListWards)
}

func (c *LocationController) httpListProvinces(w http.ResponseWriter, r *http.Request) {
	var req dto.ListProvincesRequest
	if err := httpapi.DecodeQuery(r, &req); err != nil {
		writeHTTPError(w, http.StatusBadRequest, err.Error())
		return
	}
	result, err := c.location.ListProvinces(r.Context(), req)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeListProvincesResponse(w, result)
}

func (c *LocationController) httpListWards(w http.ResponseWriter, r *http.Request) {
	var req dto.ListWardsRequest
	if err := httpapi.DecodeQuery(r, &req); err != nil {
		writeHTTPError(w, http.StatusBadRequest, err.Error())
		return
	}
	if req.Page != nil && *req.Page < 1 {
		writeHTTPError(w, http.StatusBadRequest, "page must be >= 1")
		return
	}
	result, err := c.location.ListWards(r.Context(), req)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeListWardsResponse(w, result)
}

func writeListProvincesResponse(w http.ResponseWriter, resp *tqdpb.ListProvincesResponse) {
	if resp == nil {
		writeOK(w, []any{})
		return
	}
	data, err := protoMessages(resp.GetProvinces())
	if err != nil {
		writeHTTPError(w, http.StatusInternalServerError, err.Error())
		return
	}
	var opts []httpapi.OKOption
	if resp.GetTimestamp() > 0 {
		opts = append(opts, httpapi.WithTimestamp(resp.GetTimestamp()))
	}
	writeOK(w, data, opts...)
}

func writeListWardsResponse(w http.ResponseWriter, resp *tqdpb.ClientListWardsResponse) {
	if resp == nil {
		writeOK(w, []any{})
		return
	}
	data, err := protoMessages(resp.GetData())
	if err != nil {
		writeHTTPError(w, http.StatusInternalServerError, err.Error())
		return
	}
	var opts []httpapi.OKOption
	if p := resp.GetPagination(); p != nil {
		opts = append(opts, httpapi.WithPagination(&httpapi.Pagination{
			Page:  p.GetPage(),
			Size:  p.GetSize(),
			Total: p.GetTotal(),
		}))
	}
	if resp.GetTimestamp() > 0 {
		opts = append(opts, httpapi.WithTimestamp(resp.GetTimestamp()))
	}
	writeOK(w, data, opts...)
}
