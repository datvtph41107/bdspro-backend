package handler_grpc

import (
	"context"

	_dto "common/domain/dto"
	_errors "common/errors"
	_utils "common/utils"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	tqdpb "pb/types/tqd"
	"tqd/infra/mapper"
	"tqd/internal/usecase"
)

// RegionLabelGrpcHandler xử lý gRPC cho RegionLabelService
type RegionLabelGrpcHandler struct {
	tqdpb.UnimplementedRegionLabelServiceServer
	regionLabelUsecase usecase.RegionLabelUsecase
	regionUsecase      usecase.RegionUsecase
	labelUsecase       usecase.QHLabelUsecase
	regionMapper       *mapper.RegionMapper
	regionLabelMapper  *mapper.RegionLabelMapper
}

func NewRegionLabelGrpcHandler(
	regionLabelUsecase usecase.RegionLabelUsecase,
	regionUsecase usecase.RegionUsecase,
	labelUsecase usecase.QHLabelUsecase,
) *RegionLabelGrpcHandler {
	return &RegionLabelGrpcHandler{
		regionLabelUsecase: regionLabelUsecase,
		regionUsecase:      regionUsecase,
		labelUsecase:       labelUsecase,
		regionMapper:       mapper.NewRegionMapper(),
		regionLabelMapper:  mapper.NewRegionLabelMapper(),
	}
}

// =====================================================
// ASSIGN LABELS
// =====================================================

// AssignLabelsToRegion - Gán label(s) cho region (thay thế toàn bộ)
func (h *RegionLabelGrpcHandler) AssignLabelsToRegion(ctx context.Context, req *tqdpb.AssignLabelsToRegionRequest) (*tqdpb.AssignLabelsToRegionResponse, error) {
	if req.RegionId == 0 {
		return nil, status.Error(codes.InvalidArgument, "regionId is required")
	}

	userID := _utils.GetOriginIdFromContext(ctx)
	if userID == 0 {
		return nil, _errors.ReturnError(401, "unauthorized")
	}
	notes := ""
	if req.Notes != nil {
		notes = *req.Notes
	}

	err := h.regionLabelUsecase.AssignLabels(ctx, req.RegionId, req.LabelIds, userID, "api", notes)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return h.regionLabelMapper.ToProtoAssignResponse(req.RegionId, req.LabelIds, "Labels assigned successfully"), nil
}

// AddLabelsToRegion - Thêm label(s) cho region (giữ nguyên label cũ)
func (h *RegionLabelGrpcHandler) AddLabelsToRegion(ctx context.Context, req *tqdpb.AddLabelsToRegionRequest) (*tqdpb.AddLabelsToRegionResponse, error) {
	if req.RegionId == 0 {
		return nil, status.Error(codes.InvalidArgument, "regionId is required")
	}
	if len(req.LabelIds) == 0 {
		return nil, status.Error(codes.InvalidArgument, "labelIds is required")
	}

	userID := _utils.GetOriginIdFromContext(ctx)
	if userID == 0 {
		return nil, _errors.ReturnError(401, "unauthorized")
	}
	notes := ""
	if req.Notes != nil {
		notes = *req.Notes
	}

	err := h.regionLabelUsecase.AddLabels(ctx, req.RegionId, req.LabelIds, userID, "api", notes)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return h.regionLabelMapper.ToProtoAddResponse(req.RegionId, req.LabelIds, "Labels added successfully"), nil
}

// RemoveLabelFromRegion - Gỡ một label khỏi region
func (h *RegionLabelGrpcHandler) RemoveLabelFromRegion(ctx context.Context, req *tqdpb.RemoveLabelFromRegionRequest) (*tqdpb.RemoveLabelFromRegionResponse, error) {
	if req.RegionId == 0 {
		return nil, status.Error(codes.InvalidArgument, "regionId is required")
	}
	if req.LabelId == 0 {
		return nil, status.Error(codes.InvalidArgument, "labelId is required")
	}

	err := h.regionLabelUsecase.RemoveLabel(ctx, req.RegionId, req.LabelId)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return h.regionLabelMapper.ToProtoRemoveResponse(req.RegionId, req.LabelId, "Label removed successfully"), nil
}

// =====================================================
// QUERY METHODS
// =====================================================

// GetLabelsByRegion - Lấy danh sách labels của region
func (h *RegionLabelGrpcHandler) GetLabelsByRegion(ctx context.Context, req *tqdpb.GetLabelsByRegionRequest) (*tqdpb.GetLabelsByRegionResponse, error) {
	if req.RegionId == 0 {
		return nil, status.Error(codes.InvalidArgument, "regionId is required")
	}

	labels, err := h.regionLabelUsecase.GetLabelsByRegion(ctx, req.RegionId)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return h.regionLabelMapper.ToProtoGetLabelsResponse(labels), nil
}

// GetRegionsByLabel - Lấy danh sách regions theo label
func (h *RegionLabelGrpcHandler) GetRegionsByLabel(ctx context.Context, req *tqdpb.GetRegionsByLabelRequest) (*tqdpb.GetRegionsByLabelResponse, error) {
	if req.LabelId == 0 {
		return nil, status.Error(codes.InvalidArgument, "labelId is required")
	}

	pagable := _dto.NewPagableFromGrpc(&req.Page, &req.Size, nil)

	regions, total, err := h.regionLabelUsecase.GetRegionsByLabel(ctx, req.LabelId, req.LayerId, pagable)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	data := make([]*tqdpb.RegionResponse, len(regions))
	for i, r := range regions {
		data[i] = h.regionMapper.ToProtoClientResponse(&r)
	}

	return &tqdpb.GetRegionsByLabelResponse{
		Data:     data,
		Total:    total,
		Page:     int32(pagable.GetPage()),
		PageSize: int32(pagable.GetSize()),
	}, nil
}

// =====================================================
// BATCH OPERATIONS
// =====================================================

// BatchAssignLabels - Gán label hàng loạt cho nhiều regions
func (h *RegionLabelGrpcHandler) BatchAssignLabels(ctx context.Context, req *tqdpb.BatchAssignLabelsRequest) (*tqdpb.BatchAssignLabelsResponse, error) {
	if len(req.RegionIds) == 0 || len(req.LabelIds) == 0 {
		return nil, status.Error(codes.InvalidArgument, "regionIds and labelIds are required")
	}

	userID := _utils.GetOriginIdFromContext(ctx)
	if userID == 0 {
		return nil, _errors.ReturnError(401, "unauthorized")
	}
	notes := ""
	if req.Notes != nil {
		notes = *req.Notes
	}

	successCount, failureCount, failures, err := h.regionLabelUsecase.BatchAssignLabels(
		ctx, req.RegionIds, req.LabelIds, req.Operation, userID, "api", notes,
	)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	message := "Batch assign completed"
	if failureCount > 0 {
		message = "Batch assign completed with some failures"
	}

	return h.regionLabelMapper.ToProtoBatchAssignResponse(successCount, failureCount, failures, message), nil
}
