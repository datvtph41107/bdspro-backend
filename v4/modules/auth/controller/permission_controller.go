package controller

import (
	"context"
	"encoding/json"
	"net/http"

	authpb "pb/types/auth"

	"userv4/modules/auth/dto"
	"userv4/modules/auth/services"
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
		mux.HandleFunc("GET "+prefix+"/auth/permission/me", c.httpGetPermissionMe)
	}
}

func (c *PermissionController) httpGetPermissionMe(w http.ResponseWriter, r *http.Request) {
	result, err := c.GetPermissionKeys(r.Context())
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writePermissionKeysHTTP(w, result)
}

func writePermissionKeysHTTP(w http.ResponseWriter, result *dto.PermissionKeysResponse) {
	keys := []string{}
	var timestamp int64
	if result != nil {
		if result.Keys != nil {
			keys = result.Keys
		}
		timestamp = result.Timestamp
	}
	// Match gateway PermissionKeysResponse marshal: timestamp at root, keys under data.
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]any{
		"code":    0,
		"message": "Success",
		"data": map[string]any{
			"keys": keys,
		},
		"timestamp": timestamp,
	})
}
