package controller

import (
	"net/http"

	hubpb "pb/types/hub"
	"userv4/internal/httpapi"
	"userv4/modules/hub/dto"
	"userv4/modules/hub/services"
)

// VersionController owns app-bundle version HTTP routes under /v2/hub/app-bundle.
type VersionController struct {
	version services.VersionService
}

func NewVersionController(version services.VersionService) *VersionController {
	return &VersionController{version: version}
}

func (c *VersionController) RegisterHTTP(mux *http.ServeMux) {
	mux.HandleFunc("GET /v2/hub/app-bundle/update-check/{appName}", c.httpCheckAppUpdate)
	mux.HandleFunc("POST /v2/hub/app-bundle/version", c.httpCreateVersionFromDev)
}

func (c *VersionController) httpCheckAppUpdate(w http.ResponseWriter, r *http.Request) {
	appName := r.PathValue("appName")
	if appName == "" {
		writeHTTPError(w, http.StatusBadRequest, "appName is required")
		return
	}
	var req dto.CheckAppUpdateRequest
	if err := httpapi.DecodeQuery(r, &req); err != nil {
		writeHTTPError(w, http.StatusBadRequest, err.Error())
		return
	}
	result, err := c.version.CheckAppUpdate(r.Context(), appName, req)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeProtoData(w, result)
}

func (c *VersionController) httpCreateVersionFromDev(w http.ResponseWriter, r *http.Request) {
	var req hubpb.CreateVersionRequest
	if err := httpapi.DecodeJSON(r, &req); err != nil {
		writeHTTPError(w, http.StatusBadRequest, err.Error())
		return
	}
	apiKey := r.Header.Get("api-key")
	if apiKey == "" {
		apiKey = r.Header.Get("x-api-key")
	}
	reqCtx := services.WithAPIKey(r.Context(), apiKey)
	result, err := c.version.CreateVersionFromDev(reqCtx, &req)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeProtoData(w, result)
}
