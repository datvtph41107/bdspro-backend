package handler_grpc

import (
	_utils "common/utils"
	"context"
	"encoding/json"
	"time"

	sharepb "pb/types/shared"
	tqdpb "pb/types/tqd"
	"tqd/internal/domain/jsonb"
	qh_domain "tqd/internal/domain/qh"
	"tqd/internal/enums"
	"tqd/internal/interface/repo"
	"tqd/internal/usecase"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type LayerGrpcHandler struct {
	tqdpb.UnimplementedLayerServiceServer
	layerUsecase usecase.LayerUsecase
}

func NewLayerGrpcHandler(layerUsecase usecase.LayerUsecase) *LayerGrpcHandler {
	return &LayerGrpcHandler{
		layerUsecase: layerUsecase,
	}
}

func (h *LayerGrpcHandler) CreateLayer(ctx context.Context, req *tqdpb.CreateLayerRequest) (*tqdpb.LayerResponse, error) {
	userID := uint64(342)

	var styleConfig jsonb.JSONB
	if req.StyleConfig != "" {
		if err := json.Unmarshal([]byte(req.StyleConfig), &styleConfig); err != nil {
			return nil, status.Error(codes.InvalidArgument, "invalid style_config JSON")
		}
	}

	// Parse dates
	var effectiveDate *time.Time
	if req.EffectiveDate != nil && *req.EffectiveDate != "" {
		t, err := time.Parse("2006-01-02", *req.EffectiveDate)
		if err != nil {
			return nil, status.Error(codes.InvalidArgument, "invalid effective_date format, expected YYYY-MM-DD")
		}
		effectiveDate = &t
	}

	var expiryDate *time.Time
	if req.ExpiryDate != nil && *req.ExpiryDate != "" {
		t, err := time.Parse("2006-01-02", *req.ExpiryDate)
		if err != nil {
			return nil, status.Error(codes.InvalidArgument, "invalid expiry_date format, expected YYYY-MM-DD")
		}
		expiryDate = &t
	}

	// Build domain
	layer := &qh_domain.QHLayer{
		Name:             req.Name,
		DisplayName:      req.DisplayName,
		Description:      req.Description,
		Type:             enums.LayerType(req.Type),
		Visible:          req.Visible,
		DisplayOrder:     int(req.DisplayOrder),
		MinZoom:          int(req.MinZoom),
		MaxZoom:          int(req.MaxZoom),
		LayerURL:         req.LayerUrl,
		Avatar:           req.Avatar,
		ImageURL:         req.ImageUrl,
		ThumbnailURL:     req.ThumbnailUrl,
		StyleConfig:      styleConfig,
		LegalDoc:         req.LegalDoc,
		LegalDocURL:      req.LegalDocUrl,
		DocumentNumber:   req.DocumentNumber,
		IssuingAuthority: req.IssuingAuthority,
		EffectiveDate:    effectiveDate,
		ExpiryDate:       expiryDate,
	}

	// Set status if provided
	if req.Status != nil {
		layer.Status = enums.LayerStatus(*req.Status)
	}

	created, err := h.layerUsecase.Create(ctx, layer, userID)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return h.toProtoLayer(created, true)
}

func (h *LayerGrpcHandler) UpdateLayer(ctx context.Context, req *tqdpb.UpdateLayerRequest) (*tqdpb.LayerResponse, error) {
	userID := uint64(1)

	layer := &qh_domain.QHLayer{}

	if req.Name != nil {
		layer.Name = *req.Name
	}
	if req.DisplayName != nil {
		layer.DisplayName = *req.DisplayName
	}
	if req.Description != nil {
		layer.Description = *req.Description
	}
	if req.Type != nil {
		layer.Type = enums.LayerType(*req.Type)
	}
	if req.Status != nil {
		layer.Status = enums.LayerStatus(*req.Status)
	}
	if req.Visible != nil {
		layer.Visible = *req.Visible
	}
	if req.DisplayOrder != nil {
		layer.DisplayOrder = int(*req.DisplayOrder)
	}
	if req.MinZoom != nil {
		layer.MinZoom = int(*req.MinZoom)
	}
	if req.MaxZoom != nil {
		layer.MaxZoom = int(*req.MaxZoom)
	}
	if req.LayerUrl != nil {
		layer.LayerURL = *req.LayerUrl
	}
	if req.Avatar != nil {
		layer.Avatar = *req.Avatar
	}
	if req.ImageUrl != nil {
		layer.ImageURL = *req.ImageUrl
	}
	if req.ThumbnailUrl != nil {
		layer.ThumbnailURL = *req.ThumbnailUrl
	}
	if req.StyleConfig != nil {
		var styleConfig jsonb.JSONB
		if err := json.Unmarshal([]byte(*req.StyleConfig), &styleConfig); err == nil {
			layer.StyleConfig = styleConfig
		}
	}
	if req.LegalDoc != nil {
		layer.LegalDoc = *req.LegalDoc
	}
	if req.LegalDocUrl != nil {
		layer.LegalDocURL = *req.LegalDocUrl
	}
	if req.DocumentNumber != nil {
		layer.DocumentNumber = *req.DocumentNumber
	}
	if req.IssuingAuthority != nil {
		layer.IssuingAuthority = *req.IssuingAuthority
	}
	if req.EffectiveDate != nil && *req.EffectiveDate != "" {
		t, err := time.Parse("2006-01-02", *req.EffectiveDate)
		if err == nil {
			layer.EffectiveDate = &t
		}
	}
	if req.ExpiryDate != nil && *req.ExpiryDate != "" {
		t, err := time.Parse("2006-01-02", *req.ExpiryDate)
		if err == nil {
			layer.ExpiryDate = &t
		}
	}

	updated, err := h.layerUsecase.Update(ctx, req.Id, layer, userID)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return h.toProtoLayer(updated, true)
}

func (h *LayerGrpcHandler) DeleteLayer(ctx context.Context, req *tqdpb.DeleteLayerRequest) (*sharepb.Empty, error) {
	// TODO: Get userID from auth context
	userID := uint64(1)

	if err := h.layerUsecase.Delete(ctx, req.Id, userID); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &sharepb.Empty{}, nil
}

func (h *LayerGrpcHandler) GetLayer(ctx context.Context, req *tqdpb.GetLayerRequest) (*tqdpb.LayerResponse, error) {
	layer, err := h.layerUsecase.GetByID(ctx, req.Id)
	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}

	return h.toProtoLayer(layer, true)
}

func (h *LayerGrpcHandler) ListLayers(ctx context.Context, req *tqdpb.ListLayersRequest) (*tqdpb.ListLayersResponse, error) {
	filter := &repo.LayerFilter{
		Page:  int(*req.Page),
		Limit: int(*req.Limit),
	}

	if req.Status != nil {
		s := enums.LayerStatus(*req.Status)
		filter.Status = &s
	}
	if req.Type != nil {
		t := enums.LayerType(*req.Type)
		filter.Type = &t
	}
	if req.Visible != nil {
		v := *req.Visible
		filter.Visible = &v
	}
	if req.Search != nil && *req.Search != "" {
		filter.Search = req.Search
	}
	if req.OrderBy != nil {
		filter.OrderBy = *req.OrderBy
	}
	if req.OrderDir != nil {
		filter.OrderDir = *req.OrderDir
	}

	layers, total, err := h.layerUsecase.List(ctx, filter)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	protoLayers := make([]*tqdpb.LayerResponse, 0, len(layers))
	for _, layer := range layers {
		protoLayer, err := h.toProtoLayer(&layer, true)
		if err != nil {
			continue
		}
		protoLayers = append(protoLayers, protoLayer)
	}

	return &tqdpb.ListLayersResponse{
		Data:  protoLayers,
		Total: total,
		Page:  int32(filter.Page),
		Size:  int32(filter.Limit),
	}, nil
}

func (h *LayerGrpcHandler) UpdateLayerStatus(ctx context.Context, req *tqdpb.UpdateLayerStatusRequest) (*tqdpb.LayerResponse, error) {
	userID := uint64(1)

	if err := h.layerUsecase.UpdateStatus(ctx, req.Id, enums.LayerStatus(req.Status), userID); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	layer, err := h.layerUsecase.GetByID(ctx, req.Id)
	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}

	return h.toProtoLayer(layer, true)
}

func (h *LayerGrpcHandler) ListClientLayers(ctx context.Context, req *tqdpb.ListClientLayersRequest) (*tqdpb.ListLayersResponse, error) {
	filter := &repo.ClientLayerFilter{
		Page:  int(*req.Page),
		Limit: int(*req.Limit),
	}

	if req.Type != nil {
		t := enums.LayerType(*req.Type)
		filter.Type = &t
	}
	if req.Search != nil && *req.Search != "" {
		filter.Search = req.Search
	}

	layers, total, err := h.layerUsecase.ListClient(ctx, filter)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	protoLayers := make([]*tqdpb.LayerResponse, 0, len(layers))
	for _, layer := range layers {
		protoLayer, err := h.toProtoLayer(&layer, false) // isAdmin = false
		if err != nil {
			continue
		}
		protoLayers = append(protoLayers, protoLayer)
	}

	return &tqdpb.ListLayersResponse{
		Data:  protoLayers,
		Total: total,
		Page:  int32(filter.Page),
		Size:  int32(filter.Limit),
	}, nil
}

func (h *LayerGrpcHandler) GetClientLayer(ctx context.Context, req *tqdpb.GetClientLayerRequest) (*tqdpb.LayerResponse, error) {
	layer, err := h.layerUsecase.GetClientByID(ctx, req.Id)
	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}

	return h.toProtoLayer(layer, false)
}

func (h *LayerGrpcHandler) toProtoLayer(layer *qh_domain.QHLayer, isAdmin bool) (*tqdpb.LayerResponse, error) {
	resp := &tqdpb.LayerResponse{
		Id:               layer.ID,
		Name:             layer.Name,
		DisplayName:      layer.DisplayName,
		Description:      layer.Description,
		Type:             uint32(layer.Type),
		Visible:          layer.Visible,
		DisplayOrder:     int32(layer.DisplayOrder),
		MinZoom:          int32(layer.MinZoom),
		MaxZoom:          int32(layer.MaxZoom),
		LayerUrl:         layer.LayerURL,
		Avatar:           layer.Avatar,
		ImageUrl:         layer.ImageURL,
		ThumbnailUrl:     layer.ThumbnailURL,
		LegalDoc:         layer.LegalDoc,
		LegalDocUrl:      layer.LegalDocURL,
		DocumentNumber:   layer.DocumentNumber,
		IssuingAuthority: layer.IssuingAuthority,
		CreatedAt:        _utils.FormatTimeToString(&layer.CreatedAt),
		UpdatedAt:        _utils.FormatTimeToString(&layer.UpdatedAt),
	}

	// Status - client only sees active if visible
	if isAdmin {
		resp.Status = uint32(layer.Status)
	} else {
		// Client chỉ thấy status = active (10)
		resp.Status = uint32(enums.LayerStatusActive)
	}

	// Style config
	if layer.StyleConfig != nil {
		styleBytes, _ := json.Marshal(layer.StyleConfig)
		resp.StyleConfig = string(styleBytes)
	}

	// Dates
	if layer.EffectiveDate != nil {
		dateStr := _utils.FormatTimeToString(layer.EffectiveDate)
		resp.EffectiveDate = &dateStr
	}
	if layer.ExpiryDate != nil {
		dateStr := _utils.FormatTimeToString(layer.ExpiryDate)
		resp.ExpiryDate = &dateStr
	}

	return resp, nil
}
