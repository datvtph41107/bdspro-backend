package controller

import (
	"context"
	"net/http"

	_utils "common/utils"
	"userv4/internal/httpapi"
	"userv4/modules/hub/dto"
	"userv4/modules/hub/services"
)

// InteractiveEventController handles POST /v2/hub/v2/interactive-event.
type InteractiveEventController struct {
	svc services.InteractiveEventService
}

func NewInteractiveEventController(svc services.InteractiveEventService) *InteractiveEventController {
	return &InteractiveEventController{svc: svc}
}

func (c *InteractiveEventController) RegisterHTTP(mux *http.ServeMux) {
	mux.HandleFunc("POST /v2/hub/v2/interactive-event", httpapi.BindJSON(c.track))
}

func (c *InteractiveEventController) track(ctx context.Context, req dto.TrackInteractiveEventBatchRequest) (*httpapi.Result, error) {
	var profileID *uint64
	if pid := _utils.GetProfileIdWithContext(ctx); pid > 0 {
		profileID = &pid
	}
	if err := c.svc.TrackBatch(ctx, req, profileID); err != nil {
		return nil, err
	}
	return httpapi.OK(map[string]string{"message": "ok"}), nil
}
