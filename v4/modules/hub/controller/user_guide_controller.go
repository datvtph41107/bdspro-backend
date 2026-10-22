package controller

import (
	"context"
	"net/http"

	"common/httpapi"
	"mdn4/modules/hub/dto"
	"mdn4/modules/hub/services"
)

// UserGuideController owns user-facing user-guide HTTP routes under /v2/hub/user-guides.
type UserGuideController struct {
	userGuide services.UserGuideService
}

func NewUserGuideController(userGuide services.UserGuideService) *UserGuideController {
	return &UserGuideController{userGuide: userGuide}
}

func (c *UserGuideController) RegisterHTTP(mux *http.ServeMux) {
	// Public
	httpapi.GET(mux, "/v2/hub/user-guides/groups", c.getGroups)
	httpapi.GET(mux, "/v2/hub/user-guides/by-key", c.getAllByKey)
	httpapi.GET(mux, "/v2/hub/user-guides/by-key/{key}", c.getByKey)
	httpapi.GET(mux, "/v2/hub/user-guides", c.listSimple)
	httpapi.GET(mux, "/v2/hub/user-guides/{id}", c.getByID)
	// Admin
	httpapi.GET(mux, "/v2/hub/admin/user-guides", c.adminList)
	mux.HandleFunc("POST /v2/hub/user-guides", httpapi.BindJSON(c.adminCreate))
	mux.HandleFunc("PUT /v2/hub/user-guides/{id}", httpapi.BindJSON(c.adminUpdate))
	mux.HandleFunc("DELETE /v2/hub/user-guides/{id}", httpapi.BindJSON(c.adminDelete))
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

func (c *UserGuideController) adminList(ctx context.Context, req dto.GetUserGuideSimpleRequest) (*httpapi.Result, error) {
	result, err := c.userGuide.AdminList(ctx, req)
	if err != nil {
		return nil, err
	}
	if result == nil {
		return emptyProtoListResult()
	}
	return protoListWithTotalResult(result.GetData(), result.GetTotal())
}

func (c *UserGuideController) adminCreate(ctx context.Context, req dto.CreateUserGuideRequest) (*httpapi.Result, error) {
	result, err := c.userGuide.AdminCreate(ctx, req)
	if err != nil {
		return nil, err
	}
	data, err := protoToAny(result)
	if err != nil {
		return nil, err
	}
	return httpapi.OK(data), nil
}

func (c *UserGuideController) adminUpdate(ctx context.Context, req dto.UpdateUserGuideRequest) (*httpapi.Result, error) {
	createReq := dto.CreateUserGuideRequest{Title: req.Title, Description: req.Description, GroupKey: req.GroupKey, Key: req.Key, Mode: req.Mode}
	result, err := c.userGuide.AdminUpdate(ctx, req.ID, createReq)
	if err != nil {
		return nil, err
	}
	data, err := protoToAny(result)
	if err != nil {
		return nil, err
	}
	return httpapi.OK(data), nil
}

func (c *UserGuideController) adminDelete(ctx context.Context, req dto.DeleteUserGuidePath) (*httpapi.Result, error) {
	if err := c.userGuide.AdminDelete(ctx, req.ID); err != nil {
		return nil, err
	}
	return httpapi.OK(map[string]string{"message": "deleted"}), nil
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
