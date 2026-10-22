package controller

import (
	"context"
	"net/http"

	"common/httpapi"
	"mdn4/modules/hub/dto"
	"mdn4/modules/hub/services"
)

// APIKeyController handles /v2/hub/admin/api-keys endpoints.
type APIKeyController struct {
	svc services.APIKeyService
}

func NewAPIKeyController(svc services.APIKeyService) *APIKeyController {
	return &APIKeyController{svc: svc}
}

func (c *APIKeyController) RegisterHTTP(mux *http.ServeMux) {
	mux.HandleFunc("POST /v2/hub/admin/api-keys", httpapi.BindJSON(c.create))
	httpapi.GET(mux, "/v2/hub/admin/api-keys", c.list)
	mux.HandleFunc("DELETE /v2/hub/admin/api-keys/{id}", httpapi.BindJSON(c.delete))
	mux.HandleFunc("POST /v2/hub/admin/api-keys/verify", httpapi.BindJSON(c.verify))
}

func (c *APIKeyController) create(ctx context.Context, req dto.CreateAPIKeyRequest) (*httpapi.Result, error) {
	item, err := c.svc.CreateAPIKey(ctx, req)
	if err != nil {
		return nil, err
	}
	return httpapi.OK(item), nil
}

func (c *APIKeyController) list(ctx context.Context, _ httpapi.Empty) (*httpapi.Result, error) {
	items, err := c.svc.GetAPIKeys(ctx)
	if err != nil {
		return nil, err
	}
	return httpapi.OK(items), nil
}

func (c *APIKeyController) delete(ctx context.Context, req dto.APIKeyIDPath) (*httpapi.Result, error) {
	if err := c.svc.DeleteAPIKey(ctx, req.ID); err != nil {
		return nil, err
	}
	return httpapi.OK(map[string]string{"message": "deleted"}), nil
}

func (c *APIKeyController) verify(ctx context.Context, req dto.VerifyAPIKeyRequest) (*httpapi.Result, error) {
	item, err := c.svc.VerifyAPIKey(ctx, req.APIKey)
	if err != nil {
		return nil, err
	}
	return httpapi.OK(map[string]interface{}{"valid": true, "data": item}), nil
}
