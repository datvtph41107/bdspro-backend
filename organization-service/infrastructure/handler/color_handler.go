package handler

import (
	"context"
	organizationpb "pb/types/organization"

	"organization/infrastructure/transformer"
	"organization/internal/usecase"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ColorHandler struct {
	organizationpb.UnimplementedColorServiceServer
	colorUsecase     usecase.ColorUsecase
	colorTransformer *transformer.ColorTransformer
}

func NewColorHandler(colorUsecase usecase.ColorUsecase, colorTransformer *transformer.ColorTransformer) *ColorHandler {
	return &ColorHandler{
		colorUsecase:     colorUsecase,
		colorTransformer: colorTransformer,
	}
}

func (h *ColorHandler) CreateColor(ctx context.Context, req *organizationpb.CreateColorRequest) (*organizationpb.CreateColorResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request is required")
	}

	color := h.colorTransformer.ToCreateRequest(req)
	if err := h.colorUsecase.CreateColor(ctx, color); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create color: %v", err)
	}

	return &organizationpb.CreateColorResponse{
		Id: uint32(color.ID),
	}, nil
}

func (h *ColorHandler) UpdateColor(ctx context.Context, req *organizationpb.UpdateColorRequest) (*organizationpb.UpdateColorResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request is required")
	}

	color := h.colorTransformer.ToUpdateRequest(req)
	if err := h.colorUsecase.UpdateColor(ctx, color); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update color: %v", err)
	}

	return &organizationpb.UpdateColorResponse{
		Id: req.Id,
	}, nil
}

func (h *ColorHandler) DeleteColor(ctx context.Context, req *organizationpb.DeleteColorRequest) (*organizationpb.DeleteColorResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request is required")
	}

	if err := h.colorUsecase.DeleteColor(ctx, req.Id); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to delete color: %v", err)
	}

	return &organizationpb.DeleteColorResponse{
		Id: req.Id,
	}, nil
}

func (h *ColorHandler) GetColor(ctx context.Context, req *organizationpb.GetColorRequest) (*organizationpb.GetColorResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request is required")
	}

	color, err := h.colorUsecase.GetColorByID(ctx, req.Id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get color: %v", err)
	}

	if color == nil {
		return nil, status.Error(codes.NotFound, "color not found")
	}

	return &organizationpb.GetColorResponse{
		Color: h.colorTransformer.ToProto(color),
	}, nil
}

func (h *ColorHandler) GetColors(ctx context.Context, req *organizationpb.GetColorsRequest) (*organizationpb.GetColorsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request is required")
	}
	page := 0
	size := 10
	if req.Page != nil {
		page = int(*req.Page)
	}
	if req.Size != nil {
		size = int(*req.Size)
	}

	var isActive *bool
	if req.IsActive != nil {
		isActive = req.IsActive
	}

	colors, total, err := h.colorUsecase.GetColors(ctx, page, size, isActive)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get colors: %v", err)
	}

	return &organizationpb.GetColorsResponse{
		Data:  h.colorTransformer.ToProtoList(colors),
		Total: uint32(total),
	}, nil
}
