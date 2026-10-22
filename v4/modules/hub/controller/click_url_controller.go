package controller

import (
	"context"
	"net/http"

	_enums "common/domain/enum"
	"common/httpapi"
	"mdn4/modules/hub/dto"
	"mdn4/modules/hub/services"
)

// ClickURLController handles POST /v2/hub/v2/click-url.
type ClickURLController struct {
	svc services.ClickURLService
}

func NewClickURLController(svc services.ClickURLService) *ClickURLController {
	return &ClickURLController{svc: svc}
}

func (c *ClickURLController) RegisterHTTP(mux *http.ServeMux) {
	mux.HandleFunc("POST /v2/hub/v2/click-url", httpapi.BindJSON(c.track))
}

func (c *ClickURLController) track(ctx context.Context, req dto.TrackClickURLBatchRequest) (*httpapi.Result, error) {
	var profileID *uint64
	if pid, ok := ctx.Value(_enums.ProfileIDKey).(uint64); ok && pid > 0 {
		profileID = &pid
	}
	if err := c.svc.TrackBatch(ctx, req, profileID); err != nil {
		return nil, err
	}
	return httpapi.OK(map[string]string{"message": "ok"}), nil
}
