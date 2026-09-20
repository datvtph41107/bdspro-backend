// handler/layer_grpc_handler.go
package handler_grpc

import (
	_dto "common/domain/dto"
	_errors "common/errors"
	_utils "common/utils"
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strings"
	"time"

	sharepb "pb/types/shared"
	tqdpb "pb/types/tqd"
	"tqd/infra/mapper"
	qh_domain "tqd/internal/domain/qh"
	"tqd/internal/enums"
	"tqd/internal/interface/repo"
	"tqd/internal/usecase"

	"github.com/lib/pq"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type LayerGrpcHandler struct {
	tqdpb.UnimplementedLayerServiceServer
	layerUsecase usecase.LayerUsecase
	SyncProvider *_utils.SyncUtil
}

func NewLayerGrpcHandler(layerUsecase usecase.LayerUsecase, syncProvider *_utils.SyncUtil) *LayerGrpcHandler {
	return &LayerGrpcHandler{
		layerUsecase: layerUsecase,
		SyncProvider: syncProvider,
	}
}

func (h *LayerGrpcHandler) CreateLayer(ctx context.Context, req *tqdpb.CreateLayerRequest) (*tqdpb.LayerResponse, error) {
	userID := _utils.GetProfileIdWithContext(ctx)
	if userID == 0 {
		return nil, _errors.ReturnError(401, "unauthorized")
	}

	// Validate required fields
	// if req.Name == "" {
	// 	return nil, status.Error(codes.InvalidArgument, "name is required")
	// }
	if req.DisplayName == "" {
		return nil, status.Error(codes.InvalidArgument, "display_name is required")
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

	// Build domain - KHÔNG gán LayerURL, StyleConfig
	layer := &qh_domain.QHLayer{
		Name:          req.Name,
		DisplayName:   req.DisplayName,
		Description:   req.Description,
		Type:          enums.LayerType(req.Type),
		Visible:       req.Visible,
		DisplayOrder:  int(req.DisplayOrder),
		MinZoom:       req.MinZoom,
		MaxZoom:       req.MaxZoom,
		Avatar:        req.Avatar,
		ImageURL:      req.ImageUrl,
		ThumbnailURL:  req.ThumbnailUrl,
		SourceType:    req.SourceType,
		EffectiveDate: effectiveDate,
		ExpiryDate:    expiryDate,
		LegalStatus:   enums.LegalStatusUnknown,
		TrustValue:    0,
	}

	if req.LegalStatus != nil {
		layer.LegalStatus = enums.LegalStatus(*req.LegalStatus)
	}
	if req.TrustValue != nil {
		layer.TrustValue = float32(*req.TrustValue)
	}
	if req.FamilyId != nil && *req.FamilyId != 0 {
		fid := *req.FamilyId
		layer.FamilyID = &fid
	}
	if req.LegalDoc != "" {
		layer.LegalDoc = req.LegalDoc
	}
	if len(req.PublishScopes) > 0 {
		layer.PublishScopes = pq.Int32Array(req.PublishScopes)
	}

	// Set status if provided
	if req.Status != nil {
		layer.Status = enums.LayerStatus(*req.Status)
	} else {
		layer.Status = enums.LayerStatusDraft
	}

	created, err := h.layerUsecase.Create(ctx, layer, userID, req.ReplaceLayerIds)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	t := time.Now()
	h.SyncProvider.PutTimeRequest(ctx, h.SyncProvider.GetKey(ctx, _utils.SyncKeyTQDLayerList, 0), t.UnixMilli())
	h.SyncProvider.PutTimeRequest(ctx, h.SyncProvider.GetKeyWithoutMe(ctx, _utils.SyncKeyTQDLayerDetail, created.ID), t.UnixMilli())

	return h.toProtoLayer(created, true)
}

func (h *LayerGrpcHandler) UpdateLayer(ctx context.Context, req *tqdpb.UpdateLayerRequest) (*tqdpb.LayerResponse, error) {
	userID := _utils.GetOriginIdFromContext(ctx)

	// Lấy layer hiện tại để merge
	existing, err := h.layerUsecase.GetByID(ctx, req.Id)
	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}

	// Chỉ update các field được phép (KHÔNG update layerUrl, styleConfig)
	if req.Name != nil {
		existing.Name = *req.Name
	}
	if req.DisplayName != nil {
		existing.DisplayName = *req.DisplayName
	}
	if req.Description != nil {
		existing.Description = *req.Description
	}
	if req.Type != nil {
		existing.Type = enums.LayerType(*req.Type)
	}
	if req.Status != nil {
		existing.Status = enums.LayerStatus(*req.Status)
	}
	if req.Visible != nil {
		existing.Visible = *req.Visible
	}
	if req.DisplayOrder != nil {
		existing.DisplayOrder = int(*req.DisplayOrder)
	}
	if req.SourceType != "" {
		existing.SourceType = req.SourceType
	}
	if req.MinZoom != nil {
		existing.MinZoom = *req.MinZoom
	}
	if req.MaxZoom != nil {
		existing.MaxZoom = *req.MaxZoom
	}
	if req.Avatar != nil {
		existing.Avatar = *req.Avatar
	}
	if req.ImageUrl != nil {
		existing.ImageURL = *req.ImageUrl
	}
	if req.ThumbnailUrl != nil {
		existing.ThumbnailURL = *req.ThumbnailUrl
	}
	if req.EffectiveDate != nil && *req.EffectiveDate != "" {
		t, err := time.Parse("2006-01-02", *req.EffectiveDate)
		if err == nil {
			existing.EffectiveDate = &t
		}
	}
	if req.ExpiryDate != nil && *req.ExpiryDate != "" {
		t, err := time.Parse("2006-01-02", *req.ExpiryDate)
		if err == nil {
			existing.ExpiryDate = &t
		}
	}
	if req.LegalStatus != nil {
		existing.LegalStatus = enums.LegalStatus(*req.LegalStatus)
	}
	if req.TrustValue != nil {
		existing.TrustValue = float32(*req.TrustValue)
	}
	if req.FamilyId != nil {
		if *req.FamilyId == 0 {
			existing.FamilyID = nil
		} else {
			fid := *req.FamilyId
			existing.FamilyID = &fid
		}
	}
	if req.PublishScopes != nil {
		existing.PublishScopes = pq.Int32Array(req.PublishScopes)
	}
	if req.LegalDoc != nil {
		existing.LegalDoc = *req.LegalDoc
	}

	updated, err := h.layerUsecase.Update(ctx, req.Id, existing, userID)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	t := time.Now()
	h.SyncProvider.PutTimeRequest(ctx, h.SyncProvider.GetKey(ctx, _utils.SyncKeyTQDLayerList, 0), t.UnixMilli())
	h.SyncProvider.PutTimeRequest(ctx, h.SyncProvider.GetKeyWithoutMe(ctx, _utils.SyncKeyTQDLayerDetail, req.Id), t.UnixMilli())

	return h.toProtoLayer(updated, true)
}

// ToggleLayerBasicVisibility — đảo dấu displayOrder (âm ↔ dương) để ẩn/hiện layer cơ bản trên client.
// @bind: internal/usecase.LayerUsecase.ToggleLayerBasicVisibility
func (h *LayerGrpcHandler) ToggleLayerBasicVisibility(ctx context.Context, req *tqdpb.ToggleLayerBasicVisibilityRequest) (*tqdpb.LayerResponse, error) {
	userID := _utils.GetOriginIdFromContext(ctx)
	if userID == 0 {
		return nil, _errors.ReturnError(401, "unauthorized")
	}
	if req.Id == 0 {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}
	layer, err := h.layerUsecase.ToggleLayerBasicVisibility(ctx, req.Id, userID)
	if err != nil {
		if strings.Contains(strings.ToLower(err.Error()), "not found") {
			return nil, status.Error(codes.NotFound, err.Error())
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	t := time.Now()
	h.SyncProvider.PutTimeRequest(ctx, h.SyncProvider.GetKey(ctx, _utils.SyncKeyTQDLayerList, 0), t.UnixMilli())
	h.SyncProvider.PutTimeRequest(ctx, h.SyncProvider.GetKeyWithoutMe(ctx, _utils.SyncKeyTQDLayerDetail, req.Id), t.UnixMilli())

	return h.toProtoLayer(layer, true)
}

func (h *LayerGrpcHandler) DeleteLayer(ctx context.Context, req *tqdpb.DeleteLayerRequest) (*sharepb.Empty, error) {
	userID := _utils.GetOriginIdFromContext(ctx)
	if userID == 0 {
		return nil, _errors.ReturnError(401, "unauthorized")
	}

	if err := h.layerUsecase.Delete(ctx, req.Id, userID); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	t := time.Now()
	h.SyncProvider.PutTimeRequest(ctx, h.SyncProvider.GetKey(ctx, _utils.SyncKeyTQDLayerList, 0), t.UnixMilli())
	h.SyncProvider.PutTimeRequest(ctx, h.SyncProvider.GetKeyWithoutMe(ctx, _utils.SyncKeyTQDLayerDetail, req.Id), t.UnixMilli())

	return &sharepb.Empty{}, nil
}

func (h *LayerGrpcHandler) HardDeleteLayer(ctx context.Context, req *tqdpb.DeleteLayerRequest) (*sharepb.Empty, error) {
	userID := _utils.GetOriginIdFromContext(ctx)
	if userID == 0 {
		return nil, _errors.ReturnError(401, "unauthorized")
	}

	if err := h.layerUsecase.HardDelete(ctx, req.Id, userID); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	t := time.Now()
	h.SyncProvider.PutTimeRequest(ctx, h.SyncProvider.GetKey(ctx, _utils.SyncKeyTQDLayerList, 0), t.UnixMilli())
	h.SyncProvider.PutTimeRequest(ctx, h.SyncProvider.GetKeyWithoutMe(ctx, _utils.SyncKeyTQDLayerDetail, req.Id), t.UnixMilli())

	return &sharepb.Empty{}, nil
}

func (h *LayerGrpcHandler) GetLayer(ctx context.Context, req *tqdpb.GetLayerRequest) (*tqdpb.LayerResponse, error) {
	key := h.SyncProvider.GetKeyWithoutMe(ctx, _utils.SyncKeyTQDLayerDetail, req.Id)
	updated := h.SyncProvider.HasUpdated(ctx, key, req.Timestamp)
	h.SyncProvider.PutTimeRequest(ctx, key, req.Timestamp)
	if !updated {
		return &tqdpb.LayerResponse{}, nil
	}

	layer, err := h.layerUsecase.GetByID(ctx, req.Id)
	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}

	return h.toProtoLayer(layer, true)
}

func (h *LayerGrpcHandler) ListLayers(ctx context.Context, req *tqdpb.ListLayersRequest) (*tqdpb.ListLayersResponse, error) {
	key := h.SyncProvider.GetKey(ctx, _utils.SyncKeyTQDLayerList, 0)
	updated := h.SyncProvider.HasUpdated(ctx, key, req.Timestamp)
	h.SyncProvider.PutTimeRequest(ctx, key, req.Timestamp)
	if !updated {
		return &tqdpb.ListLayersResponse{}, nil
	}

	pagable := _dto.NewPagableFromGrpc(&req.Page, &req.Size, req.OrderBy)

	filter := &repo.LayerFilter{
		Pagable: *pagable,
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
		filter.Visible = req.Visible
	}
	if req.Search != nil && *req.Search != "" {
		filter.Search = req.Search
	}

	if req.OrderBy != nil {
		filter.OrderBy = *req.OrderBy
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
	}, nil
}

func (h *LayerGrpcHandler) ListClientLayers(ctx context.Context, req *tqdpb.ListClientLayersRequest) (*tqdpb.ListLayersResponse, error) {
	key := h.SyncProvider.GetKey(ctx, _utils.SyncKeyTQDLayerList, 0)
	updated := h.SyncProvider.HasUpdated(ctx, key, req.Timestamp)
	h.SyncProvider.PutTimeRequest(ctx, key, req.Timestamp)
	if !updated {
		return &tqdpb.ListLayersResponse{}, nil
	}

	layerFilter := repo.ClientLayerFilter{
		Pagable: _dto.Pagable{
			Page: req.GetPage(),
			Size: 1000,
		},
	}

	if req.Search != nil {
		layerFilter.Search = req.Search
	}
	if req.Type != nil {
		t := enums.LayerType(*req.Type)
		layerFilter.Type = &t
	}

	layers, total, err := h.layerUsecase.ListClient(ctx, &layerFilter)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	protoLayers := make([]*tqdpb.LayerResponse, 0, len(layers))
	for _, layer := range layers {
		protoLayer, err := mapper.ToProtoLayer(&layer)
		if err != nil {
			slog.ErrorContext(ctx, fmt.Sprintf("Failed to convert layer %d: %v", layer.ID, err))
			continue
		}
		protoLayers = append(protoLayers, protoLayer)
	}

	return &tqdpb.ListLayersResponse{
		Data:     protoLayers,
		Total:    total,
		Page:     req.GetPage(),
		PageSize: req.GetSize(),
	}, nil
}

func (h *LayerGrpcHandler) GetClientLayer(ctx context.Context, req *tqdpb.GetClientLayerRequest) (*tqdpb.LayerResponse, error) {
	key := h.SyncProvider.GetKeyWithoutMe(ctx, _utils.SyncKeyTQDLayerDetail, req.Id)
	updated := h.SyncProvider.HasUpdated(ctx, key, req.Timestamp)
	h.SyncProvider.PutTimeRequest(ctx, key, req.Timestamp)
	if !updated {
		return &tqdpb.LayerResponse{}, nil
	}

	layer, err := h.layerUsecase.GetClientByID(ctx, req.Id)
	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}

	return mapper.ToProtoLayer(layer)
}

func (h *LayerGrpcHandler) SearchTxtClientLayer(ctx context.Context, req *tqdpb.GetClientLayerRequest) (*tqdpb.ListLayersResponse, error) {
	layers, total, err := h.layerUsecase.ListClient(ctx, &repo.ClientLayerFilter{
		Pagable: _dto.Pagable{Page: 0, Size: 100},
	})
	if err != nil {
		return nil, _errors.ReturnError(500, err.Error())
	}

	protoLayers := make([]*tqdpb.LayerResponse, 0, len(layers))
	for _, layer := range layers {
		protoLayer, err := mapper.ToProtoLayer(&layer)
		if err != nil {
			continue
		}
		protoLayers = append(protoLayers, protoLayer)
	}

	return &tqdpb.ListLayersResponse{
		Data:  protoLayers,
		Total: total,
	}, nil
}

func (h *LayerGrpcHandler) BuildPMTiles(req *tqdpb.BuildPMTilesRequest, stream grpc.ServerStreamingServer[tqdpb.BuildPMTilesProgress]) error {
	ctx := stream.Context()

	if req.LayerId == 0 {
		return status.Error(codes.InvalidArgument, "layer_id is required")
	}

	minZoom := req.MinZoom
	maxZoom := req.MaxZoom
	if minZoom <= 0 {
		minZoom = 0
	}
	if maxZoom <= 0 {
		maxZoom = 14
	}
	if maxZoom > 16 {
		maxZoom = 16
	}
	if minZoom > maxZoom {
		minZoom = maxZoom
	}

	outputDir := "./files/tiles"
	if req.OutputDir != nil && *req.OutputDir != "" {
		outputDir = *req.OutputDir
	}

	progressFn := func(p usecase.BuildProgress) {
		stream.Send(&tqdpb.BuildPMTilesProgress{
			Status:       p.Status,
			Message:      p.Message,
			CurrentZoom:  p.CurrentZoom,
			TilesWritten: p.TilesWritten,
			TotalTiles:   p.TotalTiles,
			OutputPath:   p.OutputPath,
		})
	}

	err := h.layerUsecase.BuildPMTiles(ctx, req.LayerId, minZoom, maxZoom, outputDir, progressFn)
	if err == nil {
		return nil
	}

	var buildErr *usecase.ErrBuildInProgress
	if errors.As(err, &buildErr) {
		stream.Send(&tqdpb.BuildPMTilesProgress{
			Status:       buildErr.Progress.Status,
			Message:      "Another build is already in progress: " + buildErr.Progress.Message,
			CurrentZoom:  buildErr.Progress.CurrentZoom,
			TilesWritten: buildErr.Progress.TilesWritten,
			TotalTiles:   buildErr.Progress.TotalTiles,
			OutputPath:   buildErr.Progress.OutputPath,
		})
		return status.Error(codes.ResourceExhausted, err.Error())
	}

	stream.Send(&tqdpb.BuildPMTilesProgress{
		Status:  "error",
		Message: err.Error(),
	})
	return status.Error(codes.Internal, err.Error())
}

func (h *LayerGrpcHandler) toProtoLayer(layer *qh_domain.QHLayer, isAdmin bool) (*tqdpb.LayerResponse, error) {
	if layer == nil {
		return nil, nil
	}
	resp, err := mapper.ToProtoLayer(layer)
	if err != nil {
		return nil, err
	}

	if isAdmin {
		resp.Status = uint32(layer.Status)
	}

	return resp, nil
}
