package controller

import (
	"net/http"
	"strconv"

	crmpb "pb/types/crm"
	"userv4/internal/httpapi"
	"userv4/modules/bdspro/dto"
	"userv4/modules/bdspro/services"
)

// AppointmentController owns appointment HTTP routes (migrated from /v2.1/appointment*).
type AppointmentController struct {
	appointment services.AppointmentService
}

func NewAppointmentController(appointment services.AppointmentService) *AppointmentController {
	return &AppointmentController{appointment: appointment}
}

func (c *AppointmentController) RegisterHTTP(mux *http.ServeMux) {
	mux.HandleFunc("GET /v4/appointment/me", c.httpListMe)
	mux.HandleFunc("GET /v4/appointment/{id}", c.httpGetByID)
	mux.HandleFunc("GET /v2.1/appointment/list/by-product/{id}", c.httpListByProduct)
}

func (c *AppointmentController) httpListMe(w http.ResponseWriter, r *http.Request) {
	var req dto.AppointmentListRequest
	if err := httpapi.DecodeQuery(r, &req); err != nil {
		writeHTTPError(w, http.StatusBadRequest, err.Error())
		return
	}
	result, err := c.appointment.ListMe(r.Context(), req)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeAppointmentListResponse(w, result)
}

func (c *AppointmentController) httpGetByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(r.PathValue("id"), 10, 64)
	if err != nil || id == 0 {
		writeHTTPError(w, http.StatusBadRequest, "id is required")
		return
	}
	result, err := c.appointment.GetByID(r.Context(), id)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeProtoData(w, result)
}

func (c *AppointmentController) httpListByProduct(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(r.PathValue("id"), 10, 64)
	if err != nil || id == 0 {
		writeHTTPError(w, http.StatusBadRequest, "id is required")
		return
	}
	var req dto.SyncByIDRequest
	if err := httpapi.DecodeQuery(r, &req); err != nil {
		writeHTTPError(w, http.StatusBadRequest, err.Error())
		return
	}
	result, err := c.appointment.ListByProduct(r.Context(), id, req)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeAppointmentListResponse(w, result)
}

func writeAppointmentListResponse(w http.ResponseWriter, resp *crmpb.ListAppointmentsResponse) {
	if resp == nil {
		writeOK(w, []any{})
		return
	}
	data, err := protoMessages(resp.GetData())
	if err != nil {
		writeHTTPError(w, http.StatusInternalServerError, err.Error())
		return
	}
	opts := []httpapi.OKOption{
		httpapi.WithPagination(&httpapi.Pagination{Total: int64(resp.GetTotal())}),
	}
	if resp.GetTimestamp() > 0 {
		opts = append(opts, httpapi.WithTimestamp(resp.GetTimestamp()))
	}
	writeOK(w, data, opts...)
}
