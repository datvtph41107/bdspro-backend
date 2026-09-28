package controller

import (
	"context"
	"encoding/json"
	"net/http"

	sharepb "pb/types/shared"
	"userv4/internal/httpapi"
	"userv4/modules/user/services"

	"google.golang.org/protobuf/encoding/protojson"
)

var protoJSON = protojson.MarshalOptions{
	UseProtoNames:   false,
	EmitUnpopulated: false,
}

// APIVersionPrefixes mirrors gateway HTTP version mounts used by auth modules.
var APIVersionPrefixes = []string{"/v2", "/v2.1", "/v3"}

// ProfileController is the modules/user entry for profile info RPCs/HTTP.
// Legacy GrpcProfileService delegates GetProfileInfoMe here.
type ProfileController struct {
	profile services.ProfileService
}

func NewProfileController(profile services.ProfileService) *ProfileController {
	return &ProfileController{profile: profile}
}

func (c *ProfileController) GetInfoMe(ctx context.Context) (*sharepb.UserV3Proto, error) {
	return c.profile.GetInfoMe(ctx)
}

// RegisterHTTP mounts profile/info/me under every API version prefix.
//
//	GET {prefix}/user/profile/info/me
func (c *ProfileController) RegisterHTTP(mux *http.ServeMux) {
	for _, prefix := range APIVersionPrefixes {
		mux.HandleFunc("GET "+prefix+"/user/profile/info/me", c.httpGetInfoMe)
	}
}

func (c *ProfileController) httpGetInfoMe(w http.ResponseWriter, r *http.Request) {
	result, err := c.GetInfoMe(r.Context())
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeUserV3HTTP(w, result)
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

func writeHTTPError(w http.ResponseWriter, httpStatus int, message string) {
	httpapi.WriteError(w, httpStatus, message)
}

func writeServiceError(w http.ResponseWriter, err error) {
	httpapi.WriteServiceError(w, err)
}
