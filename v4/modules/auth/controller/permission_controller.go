package controller

import (
	"context"
	"net/http"

	authpb "pb/types/auth"

	"common/httpapi"
	"mdn4/modules/auth/dto"
	"mdn4/modules/auth/services"
)

// PermissionController is the modules/auth entry for permission RPCs/HTTP.
// Legacy PermissionHandler delegates GetPermissionKeys here.
type PermissionController struct {
	permission services.PermissionService
}

func NewPermissionController(permission services.PermissionService) *PermissionController {
	return &PermissionController{permission: permission}
}

func (c *PermissionController) GetPermissionKeys(ctx context.Context) (*dto.PermissionKeysResponse, error) {
	return c.permission.GetMyPermissionKeys(ctx)
}

func (c *PermissionController) GetPermissionKeysPB(ctx context.Context) (*authpb.PermissionKeysResponse, error) {
	result, err := c.GetPermissionKeys(ctx)
	if err != nil {
		return nil, err
	}
	if result == nil {
		return &authpb.PermissionKeysResponse{}, nil
	}
	return &authpb.PermissionKeysResponse{
		Keys:      result.Keys,
		Timestamp: result.Timestamp,
	}, nil
}

// RegisterHTTP mounts permission/me under every API version prefix, same as gateway.
//
//	GET {prefix}/auth/permission/me
func (c *PermissionController) RegisterHTTP(mux *http.ServeMux) {
	for _, prefix := range APIVersionPrefixes {
		httpapi.GET(mux, prefix+"/auth/permission/me", c.getPermissionMe)
	}
}

func (c *PermissionController) getPermissionMe(ctx context.Context, _ httpapi.Empty) (*httpapi.Result, error) {
	result, err := c.GetPermissionKeys(ctx)
	if err != nil {
		return nil, err
	}
	return permissionKeysResult(result), nil
}

func permissionKeysResult(result *dto.PermissionKeysResponse) *httpapi.Result {
	keys := []string{}
	var timestamp int64
	if result != nil {
		if result.Keys != nil {
			keys = result.Keys
		}
		timestamp = result.Timestamp
	}
	opts := []httpapi.OKOption{}
	if timestamp > 0 {
		opts = append(opts, httpapi.WithTimestamp(timestamp))
	}
	return httpapi.OK(map[string]any{"keys": keys}, opts...)
}
