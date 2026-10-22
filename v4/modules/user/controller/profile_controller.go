package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	_errors "common/errors"
	"common/httpapi"
	"mdn4/modules/user/dto"
	"mdn4/modules/user/services"
	sharepb "pb/types/shared"
	userpb "pb/types/user"

	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

var protoJSON = protojson.MarshalOptions{
	UseProtoNames:   false,
	EmitUnpopulated: false,
}

// APIVersionPrefixes mirrors gateway HTTP version mounts used by auth modules.
var APIVersionPrefixes = []string{"/v2", "/v2.1", "/v3"}

// ProfileController is the modules/user entry for profile RPCs/HTTP.
type ProfileController struct {
	profile services.ProfileService
	demo    services.DemoService
	session services.TileSessionService
}

func NewProfileController(profile services.ProfileService, demo services.DemoService, session services.TileSessionService) *ProfileController {
	return &ProfileController{profile: profile, demo: demo, session: session}
}

func (c *ProfileController) GetInfoMe(ctx context.Context) (*sharepb.UserV3Proto, error) {
	return c.profile.GetInfoMe(ctx)
}

func (c *ProfileController) GetProfileInfoV3(ctx context.Context, id uint64, access string, timestamp int64) (*sharepb.TotalResponseProto, error) {
	return c.profile.GetProfileInfoV3(ctx, id, access, timestamp)
}

func (c *ProfileController) GetProfileInfo(ctx context.Context, id uint64) (*userpb.ProfilePublicInfo, error) {
	return c.profile.GetProfileInfo(ctx, id)
}

func (c *ProfileController) ListAccountMe(ctx context.Context) (*userpb.AccountInfoResponse, error) {
	return c.profile.ListAccountMe(ctx)
}

func (c *ProfileController) SetInfo(ctx context.Context, req *userpb.SetInfoRequest) error {
	return c.profile.SetInfo(ctx, req)
}

func (c *ProfileController) SetPerson(ctx context.Context, req *sharepb.UserV3Proto) error {
	return c.profile.SetPerson(ctx, req)
}

func (c *ProfileController) UpdatePrivacySetting(ctx context.Context, req *userpb.UpdatePrivacySettingRequest) error {
	return c.profile.UpdatePrivacySetting(ctx, req)
}

func (c *ProfileController) QuickSetup(ctx context.Context, req *sharepb.UserV3Proto) error {
	return c.profile.QuickSetup(ctx, req)
}

func (c *ProfileController) SearchProfilePublic(ctx context.Context, text string, page, size uint32) (*userpb.SearchProfilePublicResponse, error) {
	return c.profile.SearchProfilePublic(ctx, text, page, size)
}

func (c *ProfileController) FindByPhone(ctx context.Context, phone string) (*userpb.FindByPhoneResponse, error) {
	return c.profile.FindByPhone(ctx, phone)
}

func (c *ProfileController) GetProfileCompletion(ctx context.Context) (*userpb.ProfileCompletionResponse, error) {
	return c.profile.GetProfileCompletion(ctx)
}

func (c *ProfileController) GetStatusOnline(ctx context.Context, id uint64) (*userpb.StatusOnlineResponse, error) {
	return c.profile.GetStatusOnline(ctx, id)
}

func (c *ProfileController) SetDemoMode(ctx context.Context, enable bool) error {
	if c == nil || c.demo == nil {
		return _errors.ReturnError(500, "demo service chưa được cấu hình")
	}
	return c.demo.SetDemoMode(ctx, enable)
}

func (c *ProfileController) GenTileSessionToken(ctx context.Context) (*services.TileSessionResult, error) {
	if c == nil || c.session == nil {
		return nil, _errors.ReturnError(500, "tile session service chưa được cấu hình")
	}
	return c.session.GenTileSession(ctx)
}

// RegisterHTTP mounts profile routes under every API version prefix.
func (c *ProfileController) RegisterHTTP(mux *http.ServeMux) {
	for _, prefix := range APIVersionPrefixes {
		httpapi.GET(mux, prefix+"/user/profile/info/me", c.getInfoMe)
		httpapi.GET(mux, prefix+"/user/profile/info/{id}", c.getProfileInfo)
		httpapi.POST(mux, prefix+"/user/profile/session", c.genTileSession)
		httpapi.GET(mux, prefix+"/user/profile/accounts/me", c.listAccountMe)
		httpapi.PUT(mux, prefix+"/user/profile/set-info", c.setInfo)
		httpapi.PUT(mux, prefix+"/user/profile/data", c.setPerson)
		httpapi.PUT(mux, prefix+"/user/profile/privacy", c.updatePrivacy)
		httpapi.POST(mux, prefix+"/user/profile/quick-setup", c.quickSetup)
		httpapi.GET(mux, prefix+"/user/profile/search/public", c.searchPublic)
		httpapi.GET(mux, prefix+"/user/profile/find-by-phone", c.findByPhone)
		httpapi.GET(mux, prefix+"/user/profile/completion", c.getCompletion)
		httpapi.GET(mux, prefix+"/user/status-online/{id}", c.getStatusOnline)
		httpapi.POST(mux, prefix+"/user/demo-mode", c.setDemoMode)
	}
	httpapi.GET(mux, "/v3/user/profile/{access}/{id}", c.getProfileInfoV3)
}

func (c *ProfileController) getInfoMe(ctx context.Context, _ httpapi.Empty) (*httpapi.Result, error) {
	result, err := c.GetInfoMe(ctx)
	if err != nil {
		return nil, err
	}
	return protoMessageResult(result, true)
}

func (c *ProfileController) getProfileInfo(ctx context.Context, req dto.ProfileInfoRequest) (*httpapi.Result, error) {
	result, err := c.profile.GetProfileInfo(ctx, req.ID)
	if err != nil {
		return nil, err
	}
	return protoMessageResult(result, false)
}

func (c *ProfileController) getProfileInfoV3(ctx context.Context, req dto.ProfileInfoV3Request) (*httpapi.Result, error) {
	result, err := c.GetProfileInfoV3(ctx, req.ID, req.Access, req.Timestamp)
	if err != nil {
		return nil, err
	}
	return protoMessageResult(result, false)
}

func (c *ProfileController) genTileSession(ctx context.Context, _ httpapi.Empty) (*httpapi.Result, error) {
	result, err := c.GenTileSessionToken(ctx)
	if err != nil {
		return nil, err
	}
	return httpapi.OK(dto.TileSessionResponse{
		SessionK:          result.SessionK,
		ExpiresIn:         result.ExpiresIn,
		SessionEncryptKey: result.SessionEncryptKey,
	}), nil
}

func (c *ProfileController) listAccountMe(ctx context.Context, _ httpapi.Empty) (*httpapi.Result, error) {
	result, err := c.profile.ListAccountMe(ctx)
	if err != nil {
		return nil, err
	}
	return protoMessageResult(result, false)
}

func (c *ProfileController) setInfo(ctx context.Context, req userpb.SetInfoRequest) (*httpapi.Result, error) {
	if err := c.profile.SetInfo(ctx, &req); err != nil {
		return nil, err
	}
	return httpapi.OK(nil), nil
}

func (c *ProfileController) setPerson(ctx context.Context, req sharepb.UserV3Proto) (*httpapi.Result, error) {
	if err := c.profile.SetPerson(ctx, &req); err != nil {
		return nil, err
	}
	return httpapi.OK(nil), nil
}

func (c *ProfileController) updatePrivacy(ctx context.Context, req userpb.UpdatePrivacySettingRequest) (*httpapi.Result, error) {
	if err := c.profile.UpdatePrivacySetting(ctx, &req); err != nil {
		return nil, err
	}
	return httpapi.OK(nil), nil
}

func (c *ProfileController) quickSetup(ctx context.Context, req sharepb.UserV3Proto) (*httpapi.Result, error) {
	if err := c.profile.QuickSetup(ctx, &req); err != nil {
		return nil, err
	}
	return httpapi.OK(nil), nil
}

func (c *ProfileController) searchPublic(ctx context.Context, req dto.SearchProfilePublicRequest) (*httpapi.Result, error) {
	result, err := c.profile.SearchProfilePublic(ctx, req.Text, req.Page, req.Size)
	if err != nil {
		return nil, err
	}
	return protoMessageResult(result, false)
}

func (c *ProfileController) findByPhone(ctx context.Context, req dto.FindByPhoneRequest) (*httpapi.Result, error) {
	result, err := c.profile.FindByPhone(ctx, req.Phone)
	if err != nil {
		return nil, err
	}
	return protoMessageResult(result, false)
}

func (c *ProfileController) getCompletion(ctx context.Context, _ httpapi.Empty) (*httpapi.Result, error) {
	result, err := c.profile.GetProfileCompletion(ctx)
	if err != nil {
		return nil, err
	}
	return protoMessageResult(result, false)
}

func (c *ProfileController) getStatusOnline(ctx context.Context, req dto.StatusOnlineRequest) (*httpapi.Result, error) {
	result, err := c.profile.GetStatusOnline(ctx, req.ID)
	if err != nil {
		return nil, err
	}
	return protoMessageResult(result, false)
}

func (c *ProfileController) setDemoMode(ctx context.Context, req dto.DemoModeRequest) (*httpapi.Result, error) {
	if err := c.SetDemoMode(ctx, req.Enable); err != nil {
		return nil, err
	}
	return httpapi.OK(nil), nil
}

// Handler aggregates user-module HTTP controllers.
type Handler struct {
	profile  *ProfileController
	catalog  *CatalogController
	tag      *TagController
	kyc      *KYCController
	bookmark *BookmarkController
	admin    *AdminUserController
}

func NewHandler(
	profile *ProfileController,
	catalog *CatalogController,
	tag *TagController,
	kyc *KYCController,
	bookmark *BookmarkController,
	admin *AdminUserController,
) *Handler {
	return &Handler{
		profile:  profile,
		catalog:  catalog,
		tag:      tag,
		kyc:      kyc,
		bookmark: bookmark,
		admin:    admin,
	}
}

func (h *Handler) RegisterHTTP(mux *http.ServeMux) {
	if h == nil {
		return
	}
	if h.profile != nil {
		h.profile.RegisterHTTP(mux)
	}
	if h.catalog != nil {
		h.catalog.RegisterHTTP(mux)
	}
	if h.tag != nil {
		h.tag.RegisterHTTP(mux)
	}
	if h.kyc != nil {
		h.kyc.RegisterHTTP(mux)
	}
	if h.bookmark != nil {
		h.bookmark.RegisterHTTP(mux)
	}
	if h.admin != nil {
		h.admin.RegisterHTTP(mux)
	}
}

func protoMessageResult(msg proto.Message, liftTimestamp bool) (*httpapi.Result, error) {
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
	opts := []httpapi.OKOption{}
	if liftTimestamp {
		if ts := extractTimestamp(data["timestamp"]); ts > 0 {
			delete(data, "timestamp")
			opts = append(opts, httpapi.WithTimestamp(ts))
		}
	}
	return httpapi.OK(data, opts...), nil
}

func extractTimestamp(v any) int64 {
	switch t := v.(type) {
	case float64:
		return int64(t)
	case int64:
		return t
	case json.Number:
		n, _ := t.Int64()
		return n
	case string:
		var n int64
		_, _ = fmt.Sscan(t, &n)
		return n
	default:
		return 0
	}
}
