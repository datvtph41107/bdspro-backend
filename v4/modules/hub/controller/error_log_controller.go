package controller

import (
	"context"
	"net/http"

	"common/httpapi"
	"mdn4/modules/hub/dto"
	"mdn4/modules/hub/services"
)

// ErrorLogController handles POST /v2/hub/error/log.
type ErrorLogController struct {
	svc services.ErrorLogService
}

func NewErrorLogController(svc services.ErrorLogService) *ErrorLogController {
	return &ErrorLogController{svc: svc}
}

func (c *ErrorLogController) RegisterHTTP(mux *http.ServeMux) {
	mux.HandleFunc("POST /v2/hub/error/log", httpapi.BindJSON(c.logError))
}

func (c *ErrorLogController) logError(ctx context.Context, req dto.LogErrorRequest) (*httpapi.Result, error) {
	if err := c.svc.LogError(ctx, req); err != nil {
		return nil, err
	}
	return httpapi.OK(map[string]string{"message": "ok"}), nil
}
