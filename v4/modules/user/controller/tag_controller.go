package controller

import (
	"context"
	"net/http"

	userpb "pb/types/user"
	"userv4/internal/httpapi"
	"userv4/modules/user/dto"
	"userv4/modules/user/services"
)

// TagController owns tag HTTP routes under /v2/user.
type TagController struct {
	svc services.TagService
}

func NewTagController(svc services.TagService) *TagController {
	return &TagController{svc: svc}
}

func (c *TagController) RegisterHTTP(mux *http.ServeMux) {
	httpapi.POST(mux, "/v2/user/admin/tags", c.adminCreateTag)
	httpapi.PUT(mux, "/v2/user/admin/tags/{id}", c.adminUpdateTag)
	httpapi.DELETE(mux, "/v2/user/admin/tags/{id}", c.adminDeleteTag)
	httpapi.GET(mux, "/v2/user/admin/tags", c.adminListTags)
	httpapi.GET(mux, "/v2/user/profile/tags/{tagType}", c.listTagsByType)
	httpapi.GET(mux, "/v2/pub/user/profile/tags/{tagType}", c.getProfileTagsByType)
}

func (c *TagController) adminCreateTag(ctx context.Context, req dto.AdminTagWriteRequest) (*httpapi.Result, error) {
	result, err := c.svc.AdminCreateTag(ctx, &userpb.AdminTagRequest{
		Name:      req.Name,
		TagType:   req.TagType,
		IsActive:  req.IsActive,
		IsDefault: req.IsDefault,
	})
	if err != nil {
		return nil, err
	}
	return protoMessageResult(result, false)
}

func (c *TagController) adminUpdateTag(ctx context.Context, req dto.AdminTagWriteRequest) (*httpapi.Result, error) {
	result, err := c.svc.AdminUpdateTag(ctx, &userpb.AdminTagRequest{
		Id:        req.ID,
		Name:      req.Name,
		TagType:   req.TagType,
		IsActive:  req.IsActive,
		IsDefault: req.IsDefault,
	})
	if err != nil {
		return nil, err
	}
	return protoMessageResult(result, false)
}

func (c *TagController) adminDeleteTag(ctx context.Context, req dto.TagIDPath) (*httpapi.Result, error) {
	if err := c.svc.AdminDeleteTag(ctx, req.ID); err != nil {
		return nil, err
	}
	return httpapi.OK(nil), nil
}

func (c *TagController) adminListTags(ctx context.Context, req dto.AdminTagListRequest) (*httpapi.Result, error) {
	result, err := c.svc.AdminListTags(ctx, &userpb.AdminTagListRequest{
		TagTypes:  req.TagTypes,
		Page:      req.Page,
		Size:      req.Size,
		Text:      req.Text,
		Sort:      req.Sort,
		IsDefault: req.IsDefault,
	})
	if err != nil {
		return nil, err
	}
	return protoMessageResult(result, false)
}

func (c *TagController) listTagsByType(ctx context.Context, req dto.TagListByTypeRequest) (*httpapi.Result, error) {
	result, err := c.svc.ListTagsByType(ctx, &userpb.TagListRequest{
		TagType: req.TagType,
		Page:    req.Page,
		Size:    req.Size,
		Text:    req.Text,
		Sort:    req.Sort,
	})
	if err != nil {
		return nil, err
	}
	return protoMessageResult(result, false)
}

func (c *TagController) getProfileTagsByType(ctx context.Context, req dto.GetProfileTagsByTypeRequest) (*httpapi.Result, error) {
	result, err := c.svc.GetProfileTagsByType(ctx, &userpb.GetProfileTagsByTypeRequest{
		ProfileId: req.ProfileID,
		TagType:   req.TagType,
	})
	if err != nil {
		return nil, err
	}
	return protoMessageResult(result, false)
}
