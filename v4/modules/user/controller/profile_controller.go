package controller

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	_errors "common/errors"
	sharepb "pb/types/shared"
	"userv4/internal/httpapi"
	"userv4/modules/user/dto"
	"userv4/modules/user/services"

	"google.golang.org/protobuf/encoding/protojson"
)

var protoJSON = protojson.MarshalOptions{
	UseProtoNames:   false,
	EmitUnpopulated: false,
}

// APIVersionPrefixes mirrors gateway HTTP version mounts used by auth modules.
var APIVersionPrefixes = []string{"/v2", "/v2.1", "/v3"}

// ProfileController is the modules/user entry for profile info + demo-mode RPCs/HTTP.
// Legacy GrpcProfileService delegates GetProfileInfoMe / GetProfileInfoV3 / SetDemoMode here.
type ProfileController struct {
	profile services.ProfileService
	demo    services.DemoService
}

func NewProfileController(profile services.ProfileService, demo services.DemoService) *ProfileController {
	return &ProfileController{profile: profile, demo: demo}
}

func (c *ProfileController) GetInfoMe(ctx context.Context) (*sharepb.UserV3Proto, error) {
	return c.profile.GetInfoMe(ctx)
}

func (c *ProfileController) GetProfileInfoV3(ctx context.Context, id uint64, access string, timestamp int64) (*sharepb.TotalResponseProto, error) {
	return c.profile.GetProfileInfoV3(ctx, id, access, timestamp)
}

func (c *ProfileController) SetDemoMode(ctx context.Context, enable bool) error {
	if c == nil || c.demo == nil {
		return _errors.ReturnError(500, "demo service chưa được cấu hình")
	}
	return c.demo.SetDemoMode(ctx, enable)
}

// RegisterHTTP mounts profile routes under every API version prefix.
//
//	GET  {prefix}/user/profile/info/me
//	POST {prefix}/user/demo-mode
//	GET  /v3/user/profile/{access}/{id}
func (c *ProfileController) RegisterHTTP(mux *http.ServeMux) {
	for _, prefix := range APIVersionPrefixes {
		mux.HandleFunc("GET "+prefix+"/user/profile/info/me", c.httpGetInfoMe)
		mux.HandleFunc("POST "+prefix+"/user/demo-mode", c.httpSetDemoMode)
	}
	mux.HandleFunc("GET /v3/user/profile/{access}/{id}", c.httpGetProfileInfoV3)
}

func (c *ProfileController) httpGetInfoMe(w http.ResponseWriter, r *http.Request) {
	result, err := c.GetInfoMe(r.Context())
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeUserV3HTTP(w, result)
}

func (c *ProfileController) httpGetProfileInfoV3(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(r.PathValue("id"), 10, 64)
	if err != nil {
		writeHTTPError(w, http.StatusBadRequest, "id is required")
		return
	}
	access := r.PathValue("access")
	var req dto.ProfileInfoV3Request
	if err := httpapi.DecodeQuery(r, &req); err != nil {
		writeHTTPError(w, http.StatusBadRequest, err.Error())
		return
	}
	result, err := c.GetProfileInfoV3(r.Context(), id, access, req.Timestamp)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeTotalResponseHTTP(w, result)
}

func (c *ProfileController) httpSetDemoMode(w http.ResponseWriter, r *http.Request) {
	var req dto.DemoModeRequest
	if err := httpapi.DecodeJSON(r, &req); err != nil {
		writeHTTPError(w, http.StatusBadRequest, err.Error())
		return
	}
	if err := c.SetDemoMode(r.Context(), req.Enable); err != nil {
		writeServiceError(w, err)
		return
	}
	httpapi.WriteOK(w, nil)
}

// Handler aggregates user-module HTTP controllers.
type Handler struct {
	profile *ProfileController
}

func NewHandler(profile *ProfileController) *Handler {
	return &Handler{profile: profile}
}

func (h *Handler) RegisterHTTP(mux *http.ServeMux) {
	if h == nil || h.profile == nil {
		return
	}
	h.profile.RegisterHTTP(mux)
}

func writeUserV3HTTP(w http.ResponseWriter, msg *sharepb.UserV3Proto) {
	if msg == nil {
		writeHTTPError(w, http.StatusInternalServerError, "empty profile")
		return
	}
	timestamp := msg.GetTimestamp()
	body, err := protoJSON.Marshal(msg)
	if err != nil {
		writeHTTPError(w, http.StatusInternalServerError, err.Error())
		return
	}
	var data map[string]any
	if err := json.Unmarshal(body, &data); err != nil {
		writeHTTPError(w, http.StatusInternalServerError, err.Error())
		return
	}
	delete(data, "timestamp")
	// Match gateway UserV3Proto marshal: timestamp at root, remaining fields under data.
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]any{
		"code":      0,
		"message":   "Success",
		"data":      data,
		"timestamp": timestamp,
	})
}

func writeTotalResponseHTTP(w http.ResponseWriter, msg *sharepb.TotalResponseProto) {
	if msg == nil {
		httpapi.WriteOK(w, map[string]any{})
		return
	}
	body, err := protoJSON.Marshal(msg)
	if err != nil {
		writeHTTPError(w, http.StatusInternalServerError, err.Error())
		return
	}
	var data map[string]any
	if err := json.Unmarshal(body, &data); err != nil {
		writeHTTPError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if data == nil {
		data = map[string]any{}
	}
	httpapi.WriteOK(w, data)
}

func writeHTTPError(w http.ResponseWriter, httpStatus int, message string) {
	httpapi.WriteError(w, httpStatus, message)
}

func writeServiceError(w http.ResponseWriter, err error) {
	httpapi.WriteServiceError(w, err)
}
