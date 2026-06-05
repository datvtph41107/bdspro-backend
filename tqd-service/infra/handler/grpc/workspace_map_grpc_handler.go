package handler_grpc

import (
	_dto "common/domain/dto"
	_errors "common/errors"
	_provider "common/provider"
	_utils "common/utils"
	"context"

	tqdpb "pb/types/tqd"
	"tqd/infra/mapper"
	"tqd/internal/dto"
	"tqd/internal/usecase"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type MapWorkspaceGrpcHandler struct {
	tqdpb.UnimplementedMapWorkspaceServiceServer

	usecase      usecase.MapWorkspaceUsecase
	mapper       *mapper.WorkspaceMapper
	SyncProvider *_provider.SyncProvider
}

func NewMapWorkspaceGrpcHandler(
	usecase usecase.MapWorkspaceUsecase,
	mapper *mapper.WorkspaceMapper,
	syncProvider *_provider.SyncProvider,
) *MapWorkspaceGrpcHandler {
	return &MapWorkspaceGrpcHandler{
		usecase:      usecase,
		mapper:       mapper,
		SyncProvider: syncProvider,
	}
}

func (h *MapWorkspaceGrpcHandler) currentUserID(ctx context.Context) (uint64, error) {
	userID := _utils.GetOriginIdFromContext(ctx)
	if userID == 0 {
		return 0, _errors.UnauthorizedException()
	}

	return userID, nil
}

func pagableFromWorkspaceRequest(req *tqdpb.WorkspacePageRequest) _dto.IPagable {
	page := req.GetPage()
	size := req.GetSize()

	return _dto.NewPagableFromGrpc(&page, &size, nil)
}

func (h *MapWorkspaceGrpcHandler) ListFollowedParcels(
	ctx context.Context,
	req *tqdpb.WorkspacePageRequest,
) (*tqdpb.ListFollowedParcelsResponse, error) {
	userID, err := h.currentUserID(ctx)
	if err != nil {
		return nil, err
	}

	pagable := pagableFromWorkspaceRequest(req)

	result, err := h.usecase.ListFollowedParcels(ctx, userID, pagable)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &tqdpb.ListFollowedParcelsResponse{
		Items:      h.mapper.ToProtoFollowedParcels(result.Items),
		NextCursor: result.NextCursor,
		Total:      result.Total,
	}, nil
}

func (h *MapWorkspaceGrpcHandler) FollowParcel(
	ctx context.Context,
	req *tqdpb.FollowParcelRequest,
) (*tqdpb.WorkspaceMutationResponse, error) {
	userID, err := h.currentUserID(ctx)
	if err != nil {
		return nil, err
	}

	if req.GetParcelId() == 0 {
		return nil, status.Error(codes.InvalidArgument, "parcelId is required")
	}

	id, err := h.usecase.FollowParcel(ctx, dto.FollowParcelRequestDTO{
		UserID:   userID,
		ParcelID: req.GetParcelId(),
		Note:     req.GetNote(),
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &tqdpb.WorkspaceMutationResponse{
		Success: true,
		Message: "followed parcel saved",
		Id:      id,
	}, nil
}

func (h *MapWorkspaceGrpcHandler) RemoveFollowedParcel(
	ctx context.Context,
	req *tqdpb.RemoveFollowedParcelRequest,
) (*tqdpb.WorkspaceMutationResponse, error) {
	userID, err := h.currentUserID(ctx)
	if err != nil {
		return nil, err
	}

	if req.GetFollowId() == 0 && req.GetParcelId() == 0 {
		return nil, status.Error(codes.InvalidArgument, "followId or parcelId is required")
	}

	err = h.usecase.RemoveFollowedParcel(ctx, dto.RemoveFollowedParcelRequestDTO{
		UserID:   userID,
		FollowID: req.GetFollowId(),
		ParcelID: req.GetParcelId(),
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &tqdpb.WorkspaceMutationResponse{
		Success: true,
		Message: "followed parcel removed",
		Id:      req.GetFollowId(),
	}, nil
}

func (h *MapWorkspaceGrpcHandler) ListViewHistory(
	ctx context.Context,
	req *tqdpb.WorkspacePageRequest,
) (*tqdpb.ListViewHistoryResponse, error) {
	userID, err := h.currentUserID(ctx)
	if err != nil {
		return nil, err
	}

	pagable := pagableFromWorkspaceRequest(req)

	result, err := h.usecase.ListViewHistory(ctx, userID, pagable)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &tqdpb.ListViewHistoryResponse{
		Items:      h.mapper.ToProtoHistories(result.Items),
		NextCursor: result.NextCursor,
		Total:      result.Total,
	}, nil
}

func (h *MapWorkspaceGrpcHandler) AddViewHistory(
	ctx context.Context,
	req *tqdpb.AddViewHistoryRequest,
) (*tqdpb.WorkspaceMutationResponse, error) {
	userID, err := h.currentUserID(ctx)
	if err != nil {
		return nil, err
	}

	if req.GetEntityId() == 0 {
		return nil, status.Error(codes.InvalidArgument, "entityId is required")
	}

	id, err := h.usecase.AddViewHistory(ctx, dto.AddViewHistoryRequestDTO{
		UserID:     userID,
		EntityType: req.GetEntityType(),
		EntityID:   req.GetEntityId(),
		ParcelID:   req.GetParcelId(),
		Source:     req.GetSource(),
		Zoom:       req.GetZoom(),
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &tqdpb.WorkspaceMutationResponse{
		Success: true,
		Message: "view history added",
		Id:      id,
	}, nil
}

func (h *MapWorkspaceGrpcHandler) RemoveViewHistory(
	ctx context.Context,
	req *tqdpb.RemoveViewHistoryRequest,
) (*tqdpb.WorkspaceMutationResponse, error) {
	userID, err := h.currentUserID(ctx)
	if err != nil {
		return nil, err
	}

	if req.GetHistoryId() == 0 {
		return nil, status.Error(codes.InvalidArgument, "historyId is required")
	}

	err = h.usecase.RemoveViewHistory(ctx, dto.RemoveViewHistoryRequestDTO{
		UserID:    userID,
		HistoryID: req.GetHistoryId(),
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &tqdpb.WorkspaceMutationResponse{
		Success: true,
		Message: "view history removed",
		Id:      req.GetHistoryId(),
	}, nil
}

func (h *MapWorkspaceGrpcHandler) ClearViewHistory(
	ctx context.Context,
	req *tqdpb.ClearViewHistoryRequest,
) (*tqdpb.WorkspaceMutationResponse, error) {
	userID, err := h.currentUserID(ctx)
	if err != nil {
		return nil, err
	}

	if err := h.usecase.ClearViewHistory(ctx, userID); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &tqdpb.WorkspaceMutationResponse{
		Success: true,
		Message: "view history cleared",
	}, nil
}

func (h *MapWorkspaceGrpcHandler) ListGeneratedReports(
	ctx context.Context,
	req *tqdpb.WorkspacePageRequest,
) (*tqdpb.ListGeneratedReportsResponse, error) {
	userID, err := h.currentUserID(ctx)
	if err != nil {
		return nil, err
	}

	pagable := pagableFromWorkspaceRequest(req)

	result, err := h.usecase.ListGeneratedReports(ctx, userID, pagable)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &tqdpb.ListGeneratedReportsResponse{
		Items:      h.mapper.ToProtoReports(result.Items),
		NextCursor: result.NextCursor,
		Total:      result.Total,
	}, nil
}

func (h *MapWorkspaceGrpcHandler) RemoveGeneratedReport(
	ctx context.Context,
	req *tqdpb.RemoveGeneratedReportRequest,
) (*tqdpb.WorkspaceMutationResponse, error) {
	userID, err := h.currentUserID(ctx)
	if err != nil {
		return nil, err
	}

	if req.GetReportId() == 0 {
		return nil, status.Error(codes.InvalidArgument, "reportId is required")
	}

	err = h.usecase.RemoveGeneratedReport(ctx, dto.RemoveGeneratedReportRequestDTO{
		UserID:   userID,
		ReportID: req.GetReportId(),
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &tqdpb.WorkspaceMutationResponse{
		Success: true,
		Message: "report removed",
		Id:      req.GetReportId(),
	}, nil
}

func (h *MapWorkspaceGrpcHandler) RegenerateReport(
	ctx context.Context,
	req *tqdpb.RegenerateReportRequest,
) (*tqdpb.WorkspaceMutationResponse, error) {
	userID, err := h.currentUserID(ctx)
	if err != nil {
		return nil, err
	}

	if req.GetReportId() == 0 {
		return nil, status.Error(codes.InvalidArgument, "reportId is required")
	}

	err = h.usecase.RegenerateReport(ctx, userID, req.GetReportId())
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &tqdpb.WorkspaceMutationResponse{
		Success: true,
		Message: "report regeneration started",
		Id:      req.GetReportId(),
	}, nil
}

func (h *MapWorkspaceGrpcHandler) ShareReport(
	ctx context.Context,
	req *tqdpb.ShareReportRequest,
) (*tqdpb.ShareReportResponse, error) {
	userID, err := h.currentUserID(ctx)
	if err != nil {
		return nil, err
	}

	if req.GetReportId() == 0 {
		return nil, status.Error(codes.InvalidArgument, "reportId is required")
	}

	result, err := h.usecase.ShareReport(ctx, userID, req.GetReportId())
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &tqdpb.ShareReportResponse{
		Success:  result.Success,
		ShareUrl: result.ShareURL,
		Message:  result.Message,
	}, nil
}
