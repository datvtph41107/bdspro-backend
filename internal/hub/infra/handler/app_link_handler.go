package handler

// ─────────────────────────────────────────────────────────────────────────────
// hub/infra/handler/applink_handler.go  —  gRPC handler
// ─────────────────────────────────────────────────────────────────────────────

import (
	"context"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"hub/config"
	"hub/helpers"
	"hub/infra/mapper"
	"hub/internal/usecase"
	hubpb "pb/types/hub"
)

// ApplinkHandler implements hubpb.ApplinkServiceServer
type ApplinkHandler struct {
	hubpb.UnimplementedApplinkServiceServer
	uc usecase.IApplinkUsecase
}

func NewApplinkHandler(uc usecase.IApplinkUsecase) *ApplinkHandler {
	return &ApplinkHandler{uc: uc}
}

// ── CreateApplink ─────────────────────────────────────────────────────────────
// Public — không bị chặn bởi ApiKey interceptor (không có trong ProtectedMethods)
// QR share cá nhân gọi: truyền refId + action → nhận base64(code)
func (h *ApplinkHandler) CreateApplink(
	ctx context.Context,
	req *hubpb.CreateApplinkRequest,
) (*hubpb.CreateApplinkResponse, error) {
	if req.GetRefId() == 0 {
		return nil, status.Error(codes.InvalidArgument, "ref_id is required")
	}

	action := mapper.ProtoActionToDomain(req.GetAction())
	// IsValid check thực hiện trong usecase.CreateApplink

	applink, err := h.uc.CreateApplink(ctx, req.GetRefId(), action)
	if err != nil {
		return nil, toGRPCStatus(err)
	}

	return &hubpb.CreateApplinkResponse{
		Applink: mapper.ApplinkToProto(applink, config.AppProperties.Applink.BaseURL),
	}, nil
}

// ── GetApplinkByCode ──────────────────────────────────────────────────────────
// Public — không bị chặn bởi ApiKey interceptor
// Web landing page gọi khi user click deep link:
//
//	?code=<base64> → refId + action → FE redirect myapp://...
func (h *ApplinkHandler) GetApplinkByCode(
	ctx context.Context,
	req *hubpb.GetApplinkByCodeRequest,
) (*hubpb.GetApplinkByCodeResponse, error) {
	if req.GetCode() == "" {
		return nil, status.Error(codes.InvalidArgument, "code is required")
	}

	applink, err := h.uc.GetApplinkByCode(ctx, req.GetCode())
	if err != nil {
		return nil, toGRPCStatus(err)
	}

	return &hubpb.GetApplinkByCodeResponse{
		Applink: mapper.ApplinkToProto(applink, config.AppProperties.Applink.BaseURL),
	}, nil
}

// ── error → gRPC status ───────────────────────────────────────────────────────
func toGRPCStatus(err error) error {
	switch {
	case errors.Is(err, helpers.ErrApplinkNotFound):
		return status.Error(codes.NotFound, err.Error())
	case errors.Is(err, helpers.ErrApplinkInvalidCode):
		return status.Error(codes.InvalidArgument, err.Error())
	case errors.Is(err, helpers.ErrApplinkCodeExhausted):
		return status.Error(codes.ResourceExhausted, err.Error())
	default:
		return status.Error(codes.Internal, "internal server error")
	}
}
