package controller

import (
	"context"
	"encoding/json"
	"net/http"

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
		httpapi.GET(mux, prefix+"/user/profile/info/me", c.getInfoMe)
		httpapi.POST(mux, prefix+"/user/demo-mode", c.setDemoMode)
	}
	httpapi.GET(mux, "/v3/user/profile/{access}/{id}", c.getProfileInfoV3)
}

func (c *ProfileController) getInfoMe(ctx context.Context, _ httpapi.Empty) (*httpapi.Result, error) {
	result, err := c.GetInfoMe(ctx)
	if err != nil {
		return nil, err
	}
	return userV3Result(result)
}

func (c *ProfileController) getProfileInfoV3(ctx context.Context, req dto.ProfileInfoV3Request) (*httpapi.Result, error) {
	result, err := c.GetProfileInfoV3(ctx, req.ID, req.Access, req.Timestamp)
	if err != nil {
		return nil, err
	}
	return totalResponseResult(result)
}

func (c *ProfileController) setDemoMode(ctx context.Context, req dto.DemoModeRequest) (*httpapi.Result, error) {
	if err := c.SetDemoMode(ctx, req.Enable); err != nil {
		return nil, err
	}
	return httpapi.OK(nil), nil
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

func userV3Result(msg *sharepb.UserV3Proto) (*httpapi.Result, error) {
	if msg == nil {
		return nil, _errors.ReturnError(500, "empty profile")
	}
	timestamp := msg.GetTimestamp()
	body, err := protoJSON.Marshal(msg)
	if err != nil {
		return nil, err
	}
	var data map[string]any
	if err := json.Unmarshal(body, &data); err != nil {
		return nil, err
	}
	delete(data, "timestamp")
	opts := []httpapi.OKOption{}
	if timestamp > 0 {
		opts = append(opts, httpapi.WithTimestamp(timestamp))
	}
	return httpapi.OK(data, opts...), nil
}

func totalResponseResult(msg *sharepb.TotalResponseProto) (*httpapi.Result, error) {
	if msg == nil {
		return httpapi.OK(map[string]any{}), nil
	}
	body, err := protoJSON.Marshal(msg)
	if err != nil {
		return nil, err
	}
	var data map[string]any
	if err := json.Unmarshal(body, &data); err != nil {
		return nil, err
	}
	if data == nil {
		data = map[string]any{}
	}
	return httpapi.OK(data), nil
}
