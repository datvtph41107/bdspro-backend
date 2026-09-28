package controller

import (
	"net/http"
	"strconv"

	"userv4/internal/httpapi"
	"userv4/modules/crm/dto"
	"userv4/modules/crm/services"
)

// ContactController owns GET /v2/crm/contact/list/user/{ownerId}.
type ContactController struct {
	contact services.ContactService
}

func NewContactController(contact services.ContactService) *ContactController {
	return &ContactController{contact: contact}
}

func (c *ContactController) RegisterHTTP(mux *http.ServeMux) {
	mux.HandleFunc("GET /v2/crm/contact/list/user/{ownerId}", c.httpListUser)
}

func (c *ContactController) httpListUser(w http.ResponseWriter, r *http.Request) {
	ownerID, err := strconv.ParseUint(r.PathValue("ownerId"), 10, 64)
	if err != nil {
		writeHTTPError(w, http.StatusBadRequest, "ownerId is required")
		return
	}
	var req dto.ContactSearchRequest
	if err := httpapi.DecodeQuery(r, &req); err != nil {
		writeHTTPError(w, http.StatusBadRequest, err.Error())
		return
	}
	result, err := c.contact.ListUser(r.Context(), ownerID, req)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeContactListResponse(w, result)
}
