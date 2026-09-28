package controller

import (
	"net/http"

	"userv4/modules/crm/services"
)

// Handler is the modules/crm HTTP entry for contact routes.
type Handler struct {
	contact *ContactController
}

func NewHandler(contactSvc services.ContactService) *Handler {
	return &Handler{
		contact: NewContactController(contactSvc),
	}
}

// RegisterHTTP mounts CRM routes with gateway path parity.
func (h *Handler) RegisterHTTP(mux *http.ServeMux) {
	h.contact.RegisterHTTP(mux)
}
