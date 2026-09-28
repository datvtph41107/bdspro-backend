package controller

import (
	"context"
	"net/http"

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
	httpapi.GET(mux, "/v2/hub/user-guides/groups", c.getGroups)
	httpapi.GET(mux, "/v2/hub/user-guides/by-key", c.getAllByKey)
	httpapi.GET(mux, "/v2/hub/user-guides/by-key/{key}", c.getByKey)
	httpapi.GET(mux, "/v2/hub/user-guides", c.listSimple)
	httpapi.GET(mux, "/v2/hub/user-guides/{id}", c.getByID)
}

func (c *UserGuideController) listSimple(ctx context.Context, req dto.GetUserGuideSimpleRequest) (*httpapi.Result, error) {
	result, err := c.userGuide.GetUserGuideSimple(ctx, req)
	if err != nil {
		return nil, err
	}
	if result == nil {
		return emptyProtoListResult()
	}
	return protoListWithTotalResult(result.GetData(), result.GetTotal())
}

func (c *UserGuideController) getByID(ctx context.Context, req httpapi.PathUint64) (*httpapi.Result, error) {
	result, err := c.userGuide.GetUserGuide(ctx, req.ID)
	if err != nil {
		return nil, err
	}
	return wrappedDataResult(result)
}

func (c *UserGuideController) getByKey(ctx context.Context, req dto.UserGuideByKeyPath) (*httpapi.Result, error) {
	result, err := c.userGuide.GetUserGuideByKey(ctx, req.Key)
	if err != nil {
		return nil, err
	}
	return wrappedDataResult(result)
}

func (c *UserGuideController) getAllByKey(ctx context.Context, req dto.GetAllUserGuidesByKeyRequest) (*httpapi.Result, error) {
	result, err := c.userGuide.GetAllUserGuidesByKey(ctx, req)
	if err != nil {
		return nil, err
	}
	return allUserGuidesByKeyResult(result)
}

func (c *UserGuideController) getGroups(ctx context.Context, _ httpapi.Empty) (*httpapi.Result, error) {
	result, err := c.userGuide.GetUserGuideGroups(ctx)
	if err != nil {
		return nil, err
	}
	if result == nil {
		return httpapi.OK([]any{}), nil
	}
	data, err := protoMessages(result.GetData())
	if err != nil {
		return nil, err
	}
	return httpapi.OK(data), nil
}
