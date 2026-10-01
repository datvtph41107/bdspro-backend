package controller

import (
	"context"
	"net/http"

	_errors "common/errors"
	"userv4/internal/httpapi"
	"userv4/modules/tqd/dto"
)

// AdminCommercialController owns /v2/tqd/admin/commercial/* (stubs until quota wiring).
type AdminCommercialController struct{}

func NewAdminCommercialController() *AdminCommercialController {
	return &AdminCommercialController{}
}

func (c *AdminCommercialController) RegisterHTTP(mux *http.ServeMux) {
	httpapi.GET(mux, "/v2/tqd/admin/commercial/usage", c.listUsage)
	httpapi.POST(mux, "/v2/tqd/admin/commercial/usage:reconcile", c.reconcile)
	httpapi.GET(mux, "/v2/tqd/admin/commercial/generated-reports", c.listReports)
}

func (c *AdminCommercialController) listUsage(ctx context.Context, _ dto.EmptyJSON) (*httpapi.Result, error) {
	return nil, _errors.ReturnError(501, "admin commercial usage not configured")
}

func (c *AdminCommercialController) reconcile(ctx context.Context, _ dto.EmptyJSON) (*httpapi.Result, error) {
	return nil, _errors.ReturnError(501, "admin commercial reconcile not configured")
}

func (c *AdminCommercialController) listReports(ctx context.Context, _ dto.EmptyJSON) (*httpapi.Result, error) {
	return nil, _errors.ReturnError(501, "admin generated-reports not configured")
}
