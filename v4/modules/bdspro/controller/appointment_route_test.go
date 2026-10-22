package controller

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"mdn4/modules/bdspro/dto"
	crmpb "pb/types/crm"
)

type routeStubAppointment struct{}

func (routeStubAppointment) GetByID(context.Context, uint64) (*crmpb.Appointment, error) {
	return &crmpb.Appointment{}, nil
}

func (routeStubAppointment) ListMe(context.Context, dto.AppointmentListRequest) (*crmpb.ListAppointmentsResponse, error) {
	return &crmpb.ListAppointmentsResponse{}, nil
}

func (routeStubAppointment) ListByProduct(context.Context, uint64, dto.SyncByIDRequest) (*crmpb.ListAppointmentsResponse, error) {
	return &crmpb.ListAppointmentsResponse{}, nil
}

func (routeStubAppointment) CreateAppointment(context.Context, string, string, string, string, string, []int64, uint32, []int64, uint32) (uint64, error) {
	return 0, nil
}

func (routeStubAppointment) DeleteAppointment(context.Context, uint64) error { return nil }

func TestAppointmentRoutes_RegisterNoConflict(t *testing.T) {
	mux := http.NewServeMux()
	NewAppointmentController(routeStubAppointment{}).RegisterHTTP(mux)

	for _, path := range []string{
		"/v4/appointment/me",
		"/v4/appointment/65",
		"/v2.1/appointment/list/by-product/9",
	} {
		rr := httptest.NewRecorder()
		mux.ServeHTTP(rr, httptest.NewRequest(http.MethodGet, path, nil))
		if rr.Code == http.StatusNotFound {
			t.Fatalf("%s => %d", path, rr.Code)
		}
	}
}
