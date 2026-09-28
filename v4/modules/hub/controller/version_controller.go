package controller

import (
	"context"
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
	httpapi.GET(mux, "/v2/hub/app-bundle/update-check/{appName}", c.checkAppUpdate)
	mux.HandleFunc("POST /v2/hub/app-bundle/version", withHubAPIKey(httpapi.BindJSON(c.createVersionFromDev)))
}

func (c *VersionController) checkAppUpdate(ctx context.Context, req dto.CheckAppUpdateRequest) (*httpapi.Result, error) {
	result, err := c.version.CheckAppUpdate(ctx, req.AppName, req)
	if err != nil {
		return nil, err
	}
	return protoDataResult(result)
}

func (c *VersionController) createVersionFromDev(ctx context.Context, req hubpb.CreateVersionRequest) (*httpapi.Result, error) {
	result, err := c.version.CreateVersionFromDev(ctx, &req)
	if err != nil {
		return nil, err
	}
	return protoDataResult(result)
}

func withHubAPIKey(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		apiKey := r.Header.Get("api-key")
		if apiKey == "" {
			apiKey = r.Header.Get("x-api-key")
		}
		next(w, r.WithContext(services.WithAPIKey(r.Context(), apiKey)))
	}
}
