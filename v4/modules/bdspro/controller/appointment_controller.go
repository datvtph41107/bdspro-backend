package controller

import (
	"net/http"
	"strconv"

	"userv4/internal/httpapi"
	"userv4/modules/bdspro/dto"
	"userv4/modules/bdspro/services"
)

// AppointmentController owns GET /v2.1/appointment/list/by-product/:id.
type AppointmentController struct {
	appointment services.AppointmentService
}

func NewAppointmentController(appointment services.AppointmentService) *AppointmentController {
	return &AppointmentController{appointment: appointment}
}

func (c *AppointmentController) RegisterHTTP(mux *http.ServeMux) {
	mux.HandleFunc("GET /v2.1/appointment/list/by-product/{id}", c.httpListByProduct)
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
	writeProtoList(w, result.GetData(), int64(result.GetTotal()))
}
