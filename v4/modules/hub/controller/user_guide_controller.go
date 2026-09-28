package controller

import (
	"net/http"
	"strconv"

	"userv4/internal/httpapi"
	"userv4/modules/hub/dto"
	"userv4/modules/hub/services"
)

// UserGuideController owns user-facing user-guide HTTP routes under /v2/hub/user-guides.
type UserGuideController struct {
	userGuide services.UserGuideService
}

func NewUserGuideController(userGuide services.UserGuideService) *UserGuideController {
	return &UserGuideController{userGuide: userGuide}
}

func (c *UserGuideController) RegisterHTTP(mux *http.ServeMux) {
	// Static / multi-segment paths before {id}.
	mux.HandleFunc("GET /v2/hub/user-guides/groups", c.httpGetGroups)
	mux.HandleFunc("GET /v2/hub/user-guides/by-key", c.httpGetAllByKey)
	mux.HandleFunc("GET /v2/hub/user-guides/by-key/{key}", c.httpGetByKey)
	mux.HandleFunc("GET /v2/hub/user-guides", c.httpListSimple)
	mux.HandleFunc("GET /v2/hub/user-guides/{id}", c.httpGetByID)
}

func (c *UserGuideController) httpListSimple(w http.ResponseWriter, r *http.Request) {
	var req dto.GetUserGuideSimpleRequest
	if err := httpapi.DecodeQuery(r, &req); err != nil {
		writeHTTPError(w, http.StatusBadRequest, err.Error())
		return
	}
	result, err := c.userGuide.GetUserGuideSimple(r.Context(), req)
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

func (c *UserGuideController) httpGetByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(r.PathValue("id"), 10, 64)
	if err != nil || id == 0 {
		writeHTTPError(w, http.StatusBadRequest, "id is required")
		return
	}
	result, err := c.userGuide.GetUserGuide(r.Context(), id)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeWrappedData(w, result)
}

func (c *UserGuideController) httpGetByKey(w http.ResponseWriter, r *http.Request) {
	key := r.PathValue("key")
	if key == "" {
		writeHTTPError(w, http.StatusBadRequest, "key is required")
		return
	}
	result, err := c.userGuide.GetUserGuideByKey(r.Context(), key)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeWrappedData(w, result)
}

func (c *UserGuideController) httpGetAllByKey(w http.ResponseWriter, r *http.Request) {
	var req dto.GetAllUserGuidesByKeyRequest
	if err := httpapi.DecodeQuery(r, &req); err != nil {
		writeHTTPError(w, http.StatusBadRequest, err.Error())
		return
	}
	result, err := c.userGuide.GetAllUserGuidesByKey(r.Context(), req)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeAllUserGuidesByKey(w, result)
}

func (c *UserGuideController) httpGetGroups(w http.ResponseWriter, r *http.Request) {
	result, err := c.userGuide.GetUserGuideGroups(r.Context())
	if err != nil {
		writeServiceError(w, err)
		return
	}
	if result == nil {
		writeOK(w, []any{})
		return
	}
	data, err := protoMessages(result.GetData())
	if err != nil {
		writeHTTPError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeOK(w, data)
}
