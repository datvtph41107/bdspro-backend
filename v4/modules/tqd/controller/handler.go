package controller

import (
	"net/http"

	"userv4/modules/tqd/services"
)

// Handler is the modules/tqd HTTP entry for location (and later planning/map) routes.
type Handler struct {
	location *LocationController
}

func NewHandler(locationSvc services.LocationService) *Handler {
	return &Handler{
		location: NewLocationController(locationSvc),
	}
}

// RegisterHTTP mounts tqd routes with gateway path parity.
func (h *Handler) RegisterHTTP(mux *http.ServeMux) {
	h.location.RegisterHTTP(mux)
}
