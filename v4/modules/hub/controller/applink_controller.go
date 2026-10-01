package controller

import (
	"context"
	"net/http"

	"userv4/internal/httpapi"
	"userv4/modules/hub/dto"
	"userv4/modules/hub/services"
)

// ApplinkController handles POST /v2/hub/applinks and GET /v2/hub/applinks/{code}.
type ApplinkController struct {
	svc services.ApplinkService
}

func NewApplinkController(svc services.ApplinkService) *ApplinkController {
	return &ApplinkController{svc: svc}
}

func (c *ApplinkController) RegisterHTTP(mux *http.ServeMux) {
	mux.HandleFunc("POST /v2/hub/applinks", httpapi.BindJSON(c.create))
	httpapi.GET(mux, "/v2/hub/applinks/{code}", c.getByCode)
}

func (c *ApplinkController) create(ctx context.Context, req dto.CreateApplinkRequest) (*httpapi.Result, error) {
	result, err := c.svc.CreateApplink(ctx, req)
	if err != nil {
		return nil, err
	}
	return httpapi.OK(result), nil
}

func (c *ApplinkController) getByCode(ctx context.Context, req dto.GetApplinkByCodePath) (*httpapi.Result, error) {
	result, err := c.svc.GetApplinkByCode(ctx, req.Code)
	if err != nil {
		return nil, err
	}
	return httpapi.OK(result), nil
}
