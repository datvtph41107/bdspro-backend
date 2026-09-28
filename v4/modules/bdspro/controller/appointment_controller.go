package controller

import (
	"context"
	"net/http"

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
	httpapi.GET(mux, "/v4/appointment/me", c.listMe)
	httpapi.GET(mux, "/v4/appointment/{id}", c.getByID)
	httpapi.GET(mux, "/v2.1/appointment/list/by-product/{id}", c.listByProduct)
}

func (c *AppointmentController) listMe(ctx context.Context, req dto.AppointmentListRequest) (*httpapi.Result, error) {
	result, err := c.appointment.ListMe(ctx, req)
	if err != nil {
		return nil, err
	}
	return appointmentListResult(result)
}

func (c *AppointmentController) getByID(ctx context.Context, req httpapi.PathUint64) (*httpapi.Result, error) {
	result, err := c.appointment.GetByID(ctx, req.ID)
	if err != nil {
		return nil, err
	}
	return protoResult(result)
}

func (c *AppointmentController) listByProduct(ctx context.Context, req dto.SyncByIDPathRequest) (*httpapi.Result, error) {
	result, err := c.appointment.ListByProduct(ctx, req.ID, req.ToSyncByIDRequest())
	if err != nil {
		return nil, err
	}
	return appointmentListResult(result)
}
